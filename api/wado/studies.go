package wado

import (
	"bytes"
	"dicom-store-api/fs"
	"dicom-store-api/models"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/suyashkumar/dicom"
	"io/ioutil"
	"net/http"
)

// ErrStudyValidation defines the list of error types returned from study resource.
var (
	ErrStudiesValidation = errors.New("studies validation error")
)

// StudiesStore defines database operations for a study.
type StudiesStore interface {
	Get(accountID int) (*models.Study, error)
	Update(s *models.Study) error
}

// StudiesResource implements study management handler.
type StudiesResource struct {
	Store StudiesStore
}

// NewStudiesResource creates and returns a study resource.
func NewStudiesResource(store StudiesStore) *StudiesResource {
	return &StudiesResource{
		Store: store,
	}
}

func (rs *StudiesResource) router() *chi.Mux {
	r := chi.NewRouter()
	//r.Use(rs.studiesCtx)
	r.Post("/", rs.save)
	return r
}

//func (rs *StudiesResource) studiesCtx(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		ctx := context.WithValue(r.Context(), ctxStudy, nil)
//		next.ServeHTTP(w, r.WithContext(ctx))
//	})
//}

type studiesSaveRequest struct {
	*models.Study
	ProtectedID int `json:"id"`
}

func (d *studiesSaveRequest) Bind(r *http.Request) error {
	return nil
}

type studiesSaveResponse struct {
	Study *models.Study
}

func newStudiesSaveResponse(s *models.Study) *studiesSaveResponse {
	return &studiesSaveResponse{
		Study: s,
	}
}

func (rs *StudiesResource) save(w http.ResponseWriter, r *http.Request) {
	//s := r.Context().Value(ctxStudy).(*models.Study)
	//data := &studiesSaveRequest{}
	//if err := render.Bind(r, data); err != nil {
	//	render.Render(w, r, ErrInvalidRequest(err))
	//}
	//file, _, err := r.FormFile("receipt") // r is *http.Request
	//switch err {
	//case nil:
	//case http.ErrMissingFile:
	//	fmt.Println("no file")
	//	render.Render(w, r, ErrInternalServerError)
	//	return
	//default:
	//	fmt.Println(err)
	//	render.Render(w, r, ErrInternalServerError)
	//	return
	//}
	//var buff bytes.Buffer
	//fileSize, err := buff.ReadFrom(file)
	//if err != nil {
	//	//log(r).WithField("profileCtx", claims.Sub).Error(err)
	//	render.Render(w, r, ErrInternalServerError)
	//	return
	//}
	//fmt.Println(fileSize) // this will return you a file size.

	const MaxUploadSize = 10 << 20 // 10MB
	//if r.ContentLength > MaxUploadSize {
	//	http.Error(w, "The uploaded image is too big. Please use an image less than 10MB in size", http.StatusBadRequest)
	//	return
	//}
	bodyReader := http.MaxBytesReader(w, r.Body, MaxUploadSize)

	defer bodyReader.Close()

	body, err := ioutil.ReadAll(bodyReader)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dataset, _ := dicom.Parse(bytes.NewReader(body), MaxUploadSize, nil)

	fmt.Println(dataset)

	//
	//// Dataset is also JSON serializable out of the box.
	//j, _ := json.Marshal(dataset)
	//fmt.Println(j)

	// todo create dicom objects from dataset
	// todo save dicom objects to database
	// todo generate filepath for dicom objects
	study := models.Study{}
	//studies, err := rs.Store.FindByPatient("patient1")
	//if err != nil {
	//	render.Render(w, r, ErrInternalServerError)
	//	return
	//}

	// todo use filepath to save dicom objects to filesystem
	fs.Save("./test.dcm", body)

	render.Respond(w, r, newStudiesSaveResponse(&study))
}
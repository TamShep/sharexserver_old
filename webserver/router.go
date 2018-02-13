package webserver

import (
	"net/http"
	"github.com/gorilla/mux"
	"github.com/mmichaelb/sharexserver/storage"
	"log"
	"fmt"
)

const contentTypeHeader = "Content-Type"

// ShareXRouter represents the main router which serves/accepts files.
type ShareXRouter struct {
	// Storage is an implementation of the Storage interface which is used by the ShareX router.
	Storage storage.FileStorage
}

// BindToRouter binds the ShareX router to the given super-router.
func (shareXRouter *ShareXRouter) BindToRouter(router *mux.Router) {
	// register endpoints
	router.Path("/upload").Methods(http.MethodPost).HandlerFunc(shareXRouter.handleUpload)
	router.Path(fmt.Sprintf("/{%v}", callReferenceVar)).HandlerFunc(shareXRouter.handleRequest)
}

// sendInternalError generalizes the internal error method.
func (shareXRouter *ShareXRouter) sendInternalError(writer http.ResponseWriter, action string, err error) {
	http.Error(writer, "500 an internal error occurred", http.StatusInternalServerError)
	log.Printf("An error occurred while doing the action \"%v\", %T: %+v\n", action, err, err)
}

// Close stops and closes the ShareX router. It returns an error if something goes wrong.
func (shareXRouter *ShareXRouter) Close() error {
	return shareXRouter.Storage.Close()
}

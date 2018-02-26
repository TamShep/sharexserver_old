package router

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/mmichaelb/sharexserver/pkg/storage"
	"log"
	"net/http"
)

const contentTypeHeader = "Content-Type"

// ShareXRouter represents the main router which serves/accepts files.
type ShareXRouter struct {
	// Storage is an implementation of the Storage interface which is used by the ShareX router.
	Storage storage.FileStorage
}

// GetHandler returns an instance of the http.Handler to afford easy dependency access.
func (shareXRouter *ShareXRouter) GetHandler() http.Handler {
	// create new router instance
	router := mux.NewRouter()
	// register endpoints
	router.Path("/upload").Methods(http.MethodPost).HandlerFunc(shareXRouter.handleUpload)
	router.Path(fmt.Sprintf("/{%v}", callReferenceVar)).HandlerFunc(shareXRouter.handleRequest)
	return router
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

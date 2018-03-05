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
	// WhitelistedContentTypes is a slice of content types which will be displayed embed in the browser.
	WhitelistedContentTypes []string
}

// WrapHandler wraps the endpoints to the given mux.Router. At the moment this is bound to the usage of gorilla/mux in
// your dependency but in the future this should be generalized. //TODO
func (shareXRouter *ShareXRouter) WrapHandler(router *mux.Router) {
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

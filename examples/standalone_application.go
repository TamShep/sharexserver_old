package main

import (
	"github.com/mmichaelb/sharexserver/pkg/router"
	"github.com/mmichaelb/sharexserver/pkg/storage/storages"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
	"log"
	"net/http"
)

func main() {
	// initialize main gorilla/mux router
	mainRouter := mux.NewRouter()
	mainRouter.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte("Hello there, this is my custom ShareX server application."))
	})
	// use MongoDB storage type
	fileStorage := &storages.MongoStorage{
		// MongoDB database name
		DatabaseName: "sharex-db",
		// MongoDB collection name
		CollectionName: "sharex-collection",
		// set file folder
		DataFolder: "./files/",
		// default DialInfo to connect unauthorized to the local MongoDB server
		DialInfo: &mgo.DialInfo{
			Addrs: []string{"localhost"},
		},
	}
	if err := fileStorage.Initialize(); err != nil {
		log.Println("Could not initialize file storage!")
		panic(err)
	}
	// setup ShareX router
	shareXRouter := router.ShareXRouter{
		Storage:                 fileStorage,
		WhitelistedContentTypes: []string{"image/png", "image/jpeg"},
	}
	// add ShareX handler to main router
	shareXRouter.WrapHandler(mainRouter.PathPrefix("/sharex/").Subrouter())
	httpServer := http.Server{
		Handler: mainRouter,        // use the gorilla/mux router as the http handler
		Addr:    "localhost:10711", // bind to local loop-back interface on port 8080
	}
	// run server and log occurring errors
	log.Fatal(httpServer.ListenAndServe())
}

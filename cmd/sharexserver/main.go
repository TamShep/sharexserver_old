package main

import (
	"bufio"
	"github.com/gorilla/mux"
	"github.com/mmichaelb/sharexserver/pkg/storage/storages"
	"github.com/mmichaelb/sharexserver/pkg/webserver"
	"gopkg.in/mgo.v2"
	"log"
	"net/http"
	"os"
	"time"
)

// Address determines the address, the web server should listen to
const Address = "localhost:10711"

// general information about the application
var applicationName = "{application_name}"
var version = "{version}"
var branch = "{branch}"
var commit = "{commit}"
var author = "mmichaelb"

func main() {
	log.Printf("Starting %v %v (%v/%v) by %v...\n", applicationName, version, branch, commit, author)
	router := mux.NewRouter()
	storage := &storages.MongoStorage{
		DialInfo: &mgo.DialInfo{
			Addrs:   []string{"localhost"},
			Timeout: time.Second * 4,
		},
		DataFolder:     "files/",
		DatabaseName:   "sharexserver",
		CollectionName: "uploads",
	}
	log.Printf("Initializing file storage (%T)...\n", storage)
	if err := storage.Initialize(); err != nil {
		log.Println("There was an error while initializing the storage.")
		panic(err)
	}
	log.Println("Done with storage initialization! Continuing with the binding of the ShareX router...")
	shareXRouter := &webserver.ShareXRouter{
		Storage: storage,
	}
	shareXRouter.BindToRouter(router)
	httpServer := http.Server{
		Addr:    Address,
		Handler: router,
	}
	log.Println("Running ShareX server in background. Enter \"close\" or \"stop\" to shutdown the ShareX server!")
	var closed bool
	go func() {
		if err := httpServer.ListenAndServe(); err != nil && !closed {
			panic(err)
		}
	}()
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		if scanner.Text() == "close" || scanner.Text() == "stop" {
			break
		}
	}
	log.Println("Shutting down ShareX server...")
	if err := httpServer.Close(); err != nil {
		log.Printf("There was an error while closing the ShareX server, %T: %v\n", err, err)
	}
	if err := storage.Close(); err != nil {
		log.Printf("There was an error while closing the ShareX file storage, %T: %v\n", err, err)
	}
	log.Println("Thank you for using the ShareX server. Bye!")
}

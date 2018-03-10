package main

import (
	"bufio"
	"flag"
	"github.com/gorilla/mux"
	"github.com/mmichaelb/sharexserver/internal/sharexserver"
	"github.com/mmichaelb/sharexserver/internal/sharexserver/config"
	"github.com/mmichaelb/sharexserver/pkg/router"
	"github.com/mmichaelb/sharexserver/pkg/storage"
	"log"
	"net/http"
	"os"
	"strconv"
)

// general information about the application
var applicationName = "{application_name}"
var version = "{version}"
var branch = "{branch}"
var commit = "{commit}"
var author = "mmichaelb"

var configFilepath = flag.String(
	"config", "./config.toml", "The filepath to the configuration file used by the ShareX server.")

func main() {
	// parse flags
	flag.Parse()
	// main start process
	log.Printf("Starting %v %v (%v/%v) by %v...\n", applicationName, version, branch, commit, author)
	// load main configuration
	log.Printf("Loading configuration file from %s...\n", strconv.Quote(*configFilepath))
	if err := config.LoadMainConfig(*configFilepath); err != nil {
		log.Fatalf("Could not load configuration from file, %T: %v\n", err, err)
	}
	log.Printf("Successfully loaded %d configuration keys.\n", len(config.Cfg.AllKeys()))
	// setup default mux router
	muxRouter := mux.NewRouter()
	var fileStorage storage.FileStorage
	var err error
	// determine from configuration value which file storage system should be used
	storageEngine := config.Cfg.GetString("storage_engine")
	switch storageEngine {
	case "MongoDB+file":
		// default storage system (MongoDB + standard system files)
		fileStorage, err = config.ParseMongoStorageFromConfig(config.Cfg.GetString("storage_engine_config"))
		break
	default:
		log.Fatalf("Unknown storage engine: %s\n", strconv.Quote(storageEngine))
	}
	// an error occurred while creating the file storage system
	if err != nil {
		log.Fatalf("Could not read parse %s storage system from configuration file, %T: %v.\n",
			strconv.Quote(storageEngine), err, err)
	}
	// initialization via interface method Initialize of the file storage instance
	log.Printf("Initializing file storage (%s)...\n", strconv.Quote(storageEngine))
	if err := fileStorage.Initialize(); err != nil {
		log.Fatalf("There was an error while initializing the storage (%s), %T: %v\n",
			strconv.Quote(storageEngine), err, err)
	}
	log.Println("Done with storage initialization! Continuing with the binding of the ShareX muxRouter...")
	// bind ShareXRouter to previously initialized mux muxRouter
	shareXRouter := &router.ShareXRouter{
		Storage:                 fileStorage,
		WhitelistedContentTypes: config.Cfg.GetStringSlice("whitelisted_content_types"),
	}
	// bind ShareX server handler to existing mux muxRouter
	shareXRouter.WrapHandler(muxRouter.PathPrefix("/").Subrouter())
	var handler http.Handler
	// check if a reverse proxy is used
	if reverseProxyHeader := config.Cfg.GetString("reverse_proxy_header"); reverseProxyHeader != "" {
		handler = sharexserver.WrapRouterToReverseProxyRouter(muxRouter, reverseProxyHeader)
	} else {
		handler = muxRouter
	}
	webserverAddress := config.Cfg.GetString("webserver_address")
	httpServer := http.Server{
		Addr:    webserverAddress,
		Handler: handler,
	}
	log.Printf("Running ShareX server in background and listening for connections on %s. "+
		"Enter \"close\" or \"stop\" to shutdown the ShareX server!\n", strconv.Quote(webserverAddress))
	var closed bool
	go func() {
		// run http server in background
		if err := httpServer.ListenAndServe(); err != nil && !closed {
			panic(err)
		}
	}()
	// scan for interruption message
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
	if err := fileStorage.Close(); err != nil {
		log.Printf("There was an error while closing the ShareX file storage, %T: %v\n", err, err)
	}
	log.Println("Thank you for using the ShareX server. Bye!")
}
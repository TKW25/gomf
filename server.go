package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	if len(os.Args) < 2 {
		// TODO: Is the below usage string correctly formatted for Go?
		log.Fatalln("ERROR: No config path provided\nUSAGE: go run server.go configPath")
	}

	LoadConfig(os.Args[1])
	Setup()

	// Check if directory exists, if it doesn't create it
	if _, err := os.Stat(Config.UploadDir); os.IsNotExist(err) {
		log.Println("Upload folder doesn't exist, creating it...")
		if err = os.Mkdir(Config.UploadDir, 0777); err != nil {
			log.Fatalln(err)
		}
	} else if err != nil {
		log.Fatalln(err)
	}

	router := mux.NewRouter()
	//TODO: Might want to grab some of this info from environment variables rather than the config file
	router.HandleFunc(fmt.Sprintf("%v", Config.UploadEndpoint), ReceiveFile).Methods("POST")
	router.Use(LoggingMiddleware, PanicMiddleware)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", Config.UploadPort), router))
}

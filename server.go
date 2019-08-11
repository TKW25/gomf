package main

import (
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
	router.HandleFunc("/upload", ReceiveFile).Methods("POST") //TODO: add end point to config file and listen here
	log.Fatal(http.ListenAndServe(":8000", router))           //TODO: add port to config file nad load it here
}

package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {

	LoadConfig("config/gomf.config") //TODO: should get path as a variable

	// Check if directory exists, if it doesn't create it
	if _, err := os.Stat(Config.UploadDir); os.IsNotExist(err) {
		log.Println("Upload folder doesn't exist, creating it...")
		if err = os.Mkdir(Config.UploadDir, 0777); err != nil {
			log.Println(3)
			log.Fatalln(err)
		}
	} else if err != nil {
		log.Println(4)
		log.Fatalln(err)
	}

	router := mux.NewRouter()
	router.HandleFunc("/upload", ReceiveFile).Methods("POST")
	log.Fatal(http.ListenAndServe(":8000", router))
}

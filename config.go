package main

import (
	"encoding/json"
	"log"
	"os"
)

type configuration struct {
	FileNameLength int      `json:"fileNameLength"` // The length of filenames
	UploadDir      string   `json:"uploadDir"`      // The directory to upload files to
	MaxFileSize    int64    `json:"maxFileSize"`    // The maximum file size to be accepted
	Database       database `json:"database"`       // The database connection information
	UploadEndpoint string   `json:"uploadEndpoint"` // The endpoint to receive uploads on
	UploadPort     int      `json:"uploadPort"`     // The port to listen for requests on e.g. 80
}

// databse contains the config information for the database
type database struct {
	Host      string `json:"host"`      // Database host address e.g. 127.0.0.1
	User      string `json:"user"`      // Database user name e.g. postgres
	Password  string `json:"password"`  // Password for the database user
	DBname    string `json:"dbName"`    // The name of the database to use
	Port      int    `json:"port"`      // The port to connect to the database on
	TableName string `json:"tableName"` // The name of the database's table name
}

// Global singleton holding the configuration information
var Config *configuration

// LoadConfig sets Config to the values provided by the provided configuration file.
// If no configuration file is provided, it will set Config to the default values.
// If an error occurs while opening or decoding the configuration file, it will close the program.
func LoadConfig(path string) {
	if Config != nil {
		log.Println("INFO: Config is a singleton")
		return
	}

	file, err := os.Open(path)
	if err != nil {
		log.Println(err)
		log.Fatalf("ERROR: Could not open path %v, exiting...\n", path)
	}

	var configuration *configuration = new(configuration)
	err = json.NewDecoder(file).Decode(configuration)
	if err != nil {
		log.Println(err)
		log.Fatalln("ERROR: Failed to decode the config file, exiting...")
	}

	Config = configuration
}

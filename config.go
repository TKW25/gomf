package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/alecthomas/units"
)

type configuration struct {
	FileNameLength int    // The length of filenames
	UploadDir      string // The directory to upload files to
	MaxFileSize    int64  // The maximum file size to be accepted
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

	if path == "" {
		log.Println("WARNING: No path provided, setting Config to the default...")
		path = "files"
		Config = &configuration{20, "files", int64(units.Megabyte) * 512}
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

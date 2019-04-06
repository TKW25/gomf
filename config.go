package main

import (
	"encoding/json"
	"log"
	"os"
)

type configuration struct {
	Size int
}

// Global singleton holding the configuration information
var Config *configuration

// LoadConfig sets Config to the values provided by the provided configuration file.
// If no configuration file is provided, it will st Config to the default values.
// If an error occurs while opening or decoding the configuration file, it will close the program.
func LoadConfig(path string) {
	if Config != nil {
		log.Println("INFO: Config is a singleton")
		return
	}

	if path == "" {
		log.Println("WARNING: No path provided, setting Config to the default...")
		Config = &configuration{20}
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

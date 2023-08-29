package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

// yamlLoader - function of loading configuration from yaml file
func (config *ServerConfig) yamlLoader(path string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err2 := file.Close(); err2 != nil {
			log.Println(err2)
		}
	}()
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatal(err)
	}
}

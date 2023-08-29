package config

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

// yamlLoader - функция загрузки конфига из yaml файла
//
// yamlLoader - function for loading config from yaml file
func (config *AgentConfig) yamlLoader() {
	// читаем конфиг из файла
	// read config from yamlFile
	yamlFile, err := os.Open("./config/agent.yaml")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err2 := yamlFile.Close(); err2 != nil {
			log.Println(err2)
		}
	}()
	decoder := yaml.NewDecoder(yamlFile)
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Agent config loaded from yaml file")
}

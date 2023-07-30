// Package config содержит конфигурацию агента. Инициализация конфига происхоит из файла config/agent.yaml.
// Затем конфигурация перезаписывается флагами командной строки и переменными окружения если они заданы.
//
// Package config contains the agent configuration. The config is initialized from the config/agent.yaml file.
// Then the configuration is overwritten by command line flags and environment variables if they are set.
package config

import (
	"flag"
	"log"
	"os"
	"strconv"
	"time"
	"unicode"

	"gopkg.in/yaml.v3"
)

// AgentConfig - структура, описывающая конфигурацию агента.
//
// AgentConfig - a structure that describes the agent configuration.
type AgentConfig struct {
	ServerAddress  string        `yaml:"server"`
	Key            string        `yaml:"key"`
	RateLimit      int           `yaml:"rate_limit"`
	ReportInterval time.Duration `yaml:"report"`
	PollInterval   time.Duration `yaml:"poll"`
}

// GetConfig - функция, возвращающая конфигурацию агента.
//
// GetConfig is a function that returns the agent configuration.
func GetConfig() *AgentConfig {
	var config *AgentConfig
	// читаем конфиг из файла
	// read config from file
	file, err := os.Open("./config/agent.yaml")
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
	// инициализируем флаги
	// initialize flags
	fs := flag.NewFlagSet("agent", flag.ContinueOnError)
	fs.DurationVar(&config.ReportInterval, "r", config.ReportInterval, "Report interval")
	fs.DurationVar(&config.PollInterval, "p", config.PollInterval, "Poll interval")
	fs.StringVar(&config.ServerAddress, "a", config.ServerAddress, "Server address")
	fs.StringVar(&config.Key, "k", config.Key, "Key")
	fs.IntVar(&config.RateLimit, "l", config.RateLimit, "Rate limit")
	// парсим флаги
	// parse flags
	err = fs.Parse(os.Args[1:]) //nolint:errcheck
	if err != nil {
		log.Println(err)
	}
	// если ключ не задан - обнуляем его
	// if the key is not set - zero it
	if !isSet(fs, "k") {
		config.Key = ""
	}
	// если адрес сервера задан в переменной окружения - перезаписываем
	// if the server address is set in the environment variable - rewrite
	if envServerAddress := os.Getenv("ADDRESS"); envServerAddress != "" {
		config.ServerAddress = envServerAddress
	}
	// если интервал отчета задан в переменной окружения - перезаписываем
	// if the report interval is set in the environment variable - rewrite
	if envReportInterval := os.Getenv("REPORT_INTERVAL"); envReportInterval != "" {
		if isNumeric(envReportInterval) {
			envReportInterval += "s"
		}
		config.ReportInterval, _ = time.ParseDuration(envReportInterval)
	}
	// если интервал опроса задан в переменной окружения - перезаписываем
	// if the poll interval is set in the environment variable - rewrite
	if envPollInterval := os.Getenv("POLL_INTERVAL"); envPollInterval != "" {
		if isNumeric(envPollInterval) {
			envPollInterval += "s"
		}
		config.PollInterval, _ = time.ParseDuration(envPollInterval)
	}
	// если ключ задан в переменной окружения - перезаписываем
	// if the key is set in the environment variable - rewrite
	if envKey := os.Getenv("KEY"); envKey != "" {
		config.Key = envKey
	}
	// если лимит задан в переменной окружения - перезаписываем
	// if the limit is set in the environment variable - rewrite
	if envRateLimit := os.Getenv("RATE_LIMIT"); envRateLimit != "" {
		if isNumeric(envRateLimit) {
			config.RateLimit, _ = strconv.Atoi(envRateLimit)
		}
	}
	// возвращаем конфиг
	// return config
	return config
}

// isNumeric - функция, проверяющая, является ли строка числом.
//
// isNumeric is a function that checks whether a string is a number.
func isNumeric(s string) bool {
	for _, c := range s {
		if !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}

// isSet - функция, проверяющая, установлен ли флаг.
//
// isSet is a function that checks if the flag is set.
func isSet(fs *flag.FlagSet, name string) bool {
	set := false
	fs.Visit(func(f *flag.Flag) {
		if f.Name == name {
			set = true
		}
	})
	return set
}

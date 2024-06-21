package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type RequestConfig struct {
	Count          string `json:"count"`
	Countries      string `json:"countries"`
	SearchType     string `json:"search_type"`
	BaseURL        string `json:"base_url"`
	AcceptLanguage string `json:"accept_language"`
	ContentType    string `json:"content_type"`
	Method         string `json:"method"`
	Accept         string `json:"accept"`
}

func LoadRequestConfig() (RequestConfig, error) {
	// Open the configuration file
	filename := "requests-config.json"

	file, err := os.Open("./config/" + filename)

	if err != nil {
		return RequestConfig{}, fmt.Errorf("error opening config file '%s': %w", filename, err)
	}
	defer file.Close()

	// Read the entire file contents
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return RequestConfig{}, fmt.Errorf("error reading config file '%s': %w", filename, err)
	}

	// Decode the JSON data
	var config RequestConfig
	err = json.Unmarshal(data, &config)
	if err != nil {
		return RequestConfig{}, fmt.Errorf("error unmarshaling config JSON: %w", err)
	}

	return config, nil
}

package handler

import (
	"encoding/json"
	"github.com/softwareplace/mock-server/pkg/env"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func LoadResponses() {
	loadMockResponses()
	go func() {
		watchAndReload()
	}()
	time.Sleep(256 * time.Millisecond)
	ConfigLoaded = true
}

func loadMockResponses() {
	mockJsonFilesBasePath := env.GetAppEnv().MockPath

	var newResponses []MockConfigResponse

	err := filepath.Walk(mockJsonFilesBasePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && (strings.HasSuffix(info.Name(), ".json") || strings.HasSuffix(info.Name(), ".yaml") || strings.HasSuffix(info.Name(), ".yml")) {
			data, err := os.ReadFile(path)
			if err != nil {
				log.Printf("Failed to read file %s: %v", path, err)
				return nil
			}

			var response MockConfigResponse
			if strings.HasSuffix(info.Name(), ".json") {
				if err := json.Unmarshal(data, &response); err != nil {
					log.Printf("Failed to parse JSON in file %s: %v", path, err)
					return nil
				}
			} else {
				if err := yaml.Unmarshal(data, &response); err != nil {
					log.Printf("Failed to parse YAML in file %s: %v", path, err)
					return nil
				}
			}

			newResponses = append(newResponses, response)
		}
		return nil
	})

	if err != nil {
		log.Printf("Failed to read directory: %v", err)
	}

	if len(newResponses) == 0 {
		panic("No mock responses found")
	}

	MockConfigResponses = newResponses
}

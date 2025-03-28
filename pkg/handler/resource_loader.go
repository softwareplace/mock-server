package handler

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	errohandler "github.com/softwareplace/goserve/error"
	"github.com/softwareplace/mock-server/pkg/env"
	"github.com/softwareplace/mock-server/pkg/model"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func LoadResponses(onFileChangeDetected OnFileChangDetected) {
	loadMockResponses()
	go func() {
		watchAndReload(onFileChangeDetected)
	}()
	time.Sleep(256 * time.Millisecond)
	onFileChangeDetected(false)
}

func loadMockResponses() {
	mockJsonFilesBasePath := env.GetAppEnv().MockPath

	var newResponses []model.MockConfigResponse

	errohandler.Handler(func() {
		err := filepath.Walk(mockJsonFilesBasePath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && (strings.HasSuffix(info.Name(), ".json") || strings.HasSuffix(info.Name(), ".yaml") || strings.HasSuffix(info.Name(), ".yml")) {
				data, err := os.ReadFile(path)
				if err != nil {
					log.Errorf("Failed to read file %s: %v", path, err)
					return nil
				}

				var response model.MockConfigResponse
				if strings.HasSuffix(info.Name(), ".json") {
					if err := json.Unmarshal(data, &response); err != nil {
						log.Errorf("Failed to parse JSON in file %s: %v", path, err)
						return nil
					}
				} else {
					if err := yaml.Unmarshal(data, &response); err != nil {
						log.Errorf("Failed to parse YAML in file %s: %v", path, err)
						return nil
					}
				}

				response.Redirect.StoreResponsesDir = env.UserHomePathFix(response.Redirect.StoreResponsesDir)
				response.MockFilePath = path
				newResponses = append(newResponses, response)
			}
			return nil
		})

		if err != nil {
			log.Errorf("Failed to read directory: %v", err)
		}
	}, func(err error) {
		log.Errorf("Failed to load mock files: %v", err)
	})

	model.MockConfigResponses = newResponses
}

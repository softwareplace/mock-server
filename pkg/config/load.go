package config

import (
	log "github.com/sirupsen/logrus"
	"github.com/softwareplace/mock-server/pkg/file"
	"github.com/softwareplace/mock-server/pkg/model"
)

func Load(configFilePath string) {
	if configFilePath != "" {

		config, err := file.FromYaml(configFilePath, model.MockServerConfig{})
		if err != nil {
			log.Errorf("Failed to load config file: %v", err)
			return
		}

		model.Config = config
	}
}

func HasAValidRedirectConfig() bool {
	return model.Config != nil && model.Config.RedirectConfig != nil && model.Config.RedirectConfig.Url != ""
}

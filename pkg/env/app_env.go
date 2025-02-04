package env

import (
	"flag"
	"fmt"
	"github.com/softwareplace/mock-server/pkg/config"
	"github.com/softwareplace/mock-server/pkg/model"
	"os"
	"strings"
)

type AppEnv struct {
	Port         string
	MockPath     string
	ContextPath  string
	ServerConfig string
}

var env *AppEnv

func SetAppEnv(appEnv *AppEnv) {
	env = appEnv
}

func GetAppEnv() *AppEnv {
	if env == nil {
		serverConfig := flag.String("config", "", "The configuration file to use for the mock server")
		mockPath := flag.String("mock", "", "Directory path containing JSON files")
		contextPath := flag.String("context-path", "/", "The context path to use for the mock server")
		portFlag := flag.String("port", "8080", "Port to run the mock server on")

		flag.Parse()

		config.Load(*serverConfig)

		if model.Config != nil {
			if model.Config.MockPath != "" {
				*mockPath = model.Config.MockPath
			}
			if model.Config.ContextPath != "" {
				*contextPath = model.Config.ContextPath
			}
			if model.Config.Port != "" {
				*portFlag = model.Config.Port
			}
		}

		if *mockPath == "" {
			flag.Usage()
			fmt.Println("Error: The 'mock' flag is required and cannot be empty.")
			os.Exit(1)
		}

		env = &AppEnv{
			Port:         *portFlag,
			MockPath:     *mockPath,
			ContextPath:  strings.TrimSuffix(*contextPath, "/") + "/",
			ServerConfig: *serverConfig,
		}
	}
	return env
}

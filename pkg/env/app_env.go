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

// UserHomePathFix resolves a path that starts with '~' to the user's home directory.
// If the home directory cannot be determined, it logs an error and exits the application.
// Parameters:
//   - path: The file path string, potentially starting with '~'.
//
// Returns:
//   - The resolved file path with '~' replaced by the user's home directory, or the original path if no substitution is needed.
func UserHomePathFix(path string) string {
	if strings.HasPrefix(path, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Printf("Error: Unable to resolve user home directory for path: %v\n", err)
			os.Exit(1)
		}
		path = strings.Replace(path, "~", homeDir, 1)
	}
	return path
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

			if model.Config.RedirectConfig != nil {
				model.Config.RedirectConfig.StoreResponsesDir = UserHomePathFix(model.Config.RedirectConfig.StoreResponsesDir)
			}
		}

		if *mockPath == "" {
			flag.Usage()
			fmt.Println("Error: The 'mock' flag is required and cannot be empty.")
			os.Exit(1)
		}

		*mockPath = UserHomePathFix(*mockPath)
		*serverConfig = UserHomePathFix(*serverConfig)

		fmt.Printf("Using server configuration file at: %s\n", *serverConfig)
		fmt.Printf("Using mock data path at: %s\n", *mockPath)

		env = &AppEnv{
			Port:         *portFlag,
			MockPath:     *mockPath,
			ContextPath:  strings.TrimSuffix(*contextPath, "/") + "/",
			ServerConfig: *serverConfig,
		}
	}
	return env
}

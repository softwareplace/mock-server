package env

import (
	"flag"
	"fmt"
	"os"
)

type AppEnv struct {
	Port         string
	MockPath     string
	RedirectPath string
}

var env *AppEnv

func GetAppEnv() *AppEnv {
	if env == nil {
		mockPath := flag.String("mock", "", "Directory path containing JSON files")
		redirectConfig := flag.String("redirect", "", "Directory path containing YAML files with redirect rules")
		portFlag := flag.String("port", "8080", "Port to run the mock server on")

		flag.Parse()

		if *mockPath == "" {
			flag.Usage()
			fmt.Println("Error: The 'mock' flag is required and cannot be empty.")
			os.Exit(1)
		}

		env = &AppEnv{
			Port:         *portFlag,
			MockPath:     *mockPath,
			RedirectPath: *redirectConfig,
		}
	}
	return env
}

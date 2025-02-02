package env

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

type AppEnv struct {
	Port        string
	MockPath    string
	ContextPath string
}

var env *AppEnv

func SetAppEnv(appEnv *AppEnv) {
	env = appEnv
}

func GetAppEnv() *AppEnv {
	if env == nil {
		mockPath := flag.String("mock", "", "Directory path containing JSON files")
		contextPath := flag.String("context-path", "/", "The context path to use for the mock server")
		portFlag := flag.String("port", "8080", "Port to run the mock server on")

		flag.Parse()

		if *mockPath == "" {
			flag.Usage()
			fmt.Println("Error: The 'mock' flag is required and cannot be empty.")
			os.Exit(1)
		}

		env = &AppEnv{
			Port:        *portFlag,
			MockPath:    *mockPath,
			ContextPath: strings.TrimSuffix(*contextPath, "/") + "/",
		}
	}
	return env
}

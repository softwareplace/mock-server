package main

import (
	"github.com/softwareplace/http-utils/server"
	"github.com/softwareplace/mock-server/pkg/env"
	"github.com/softwareplace/mock-server/pkg/handler"
	"time"
)

func main() {
	appEnv := env.GetAppEnv()
	handler.LoadResponses()

	appServer := server.Default().
		WithContextPath(appEnv.ContextPath)

	if !handler.ConfigLoaded {
		for !handler.ConfigLoaded {
			time.Sleep(256 * time.Millisecond)
		}
	}

	handler.Register(appServer)

	appServer.
		WithPort(appEnv.Port).
		NotFoundHandler().
		StartServer()
}

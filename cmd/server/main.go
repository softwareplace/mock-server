package main

import (
	"github.com/softwareplace/http-utils/server"
	"github.com/softwareplace/mock-server/pkg/env"
	"github.com/softwareplace/mock-server/pkg/handler"
)

func main() {
	appEnv := env.GetAppEnv()
	server.Default().
		WithPort(appEnv.Port).
		CustomNotFoundHandler(handler.RequestHandler()).
		StartServer()
}

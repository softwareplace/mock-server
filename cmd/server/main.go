package main

import (
	apicontext "github.com/softwareplace/http-utils/context"
	"github.com/softwareplace/http-utils/logger"
	"github.com/softwareplace/http-utils/server"
	"github.com/softwareplace/mock-server/pkg/env"
	"github.com/softwareplace/mock-server/pkg/handler"
	"log"
)

var appEnv = env.GetAppEnv()

var appServer server.Api[*apicontext.DefaultContext]

func init() {
	logger.LogSetup()
}

func main() {
	handler.LoadResponses(onFileChangeDetected)
	select {}
}

func onFileChangeDetected(restartServer bool) {
	if restartServer {
		if appServer != nil {
			err := appServer.StopServer()
			if err != nil {
				log.Fatalf("Failed to stop server: %v", err)
			}
		}
	}

	appServer = server.Default().
		Port(appEnv.Port).
		ContextPath(appEnv.ContextPath).
		EmbeddedServer(handler.Register).
		CustomNotFoundHandler(handler.NotFound).
		StartServerInGoroutine()
}

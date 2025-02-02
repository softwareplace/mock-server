package main

import (
	"github.com/softwareplace/http-utils/api_context"
	"github.com/softwareplace/http-utils/server"
	"github.com/softwareplace/mock-server/pkg/env"
	"github.com/softwareplace/mock-server/pkg/handler"
	"log"
)

var appEnv = env.GetAppEnv()

var appServer server.ApiRouterHandler[*api_context.DefaultContext]

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
		WithContextPath(appEnv.ContextPath).
		EmbeddedServer(handler.Register).
		WithPort(appEnv.Port).
		NotFoundHandler().
		StartServerInGoroutine()
}

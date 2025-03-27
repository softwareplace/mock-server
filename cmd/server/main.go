package main

import (
	log "github.com/sirupsen/logrus"
	apicontext "github.com/softwareplace/goserve/context"
	"github.com/softwareplace/goserve/logger"
	"github.com/softwareplace/goserve/server"
	"github.com/softwareplace/mock-server/pkg/env"
	"github.com/softwareplace/mock-server/pkg/handler"
)

var (
	appEnv    *env.AppEnv
	appServer server.Api[*apicontext.DefaultContext]
)

func init() {
	logger.LogSetup()
	appEnv = env.GetAppEnv()
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

package handler

import (
	"github.com/softwareplace/http-utils/api_context"
	"github.com/softwareplace/mock-server/pkg/config"
	"github.com/softwareplace/mock-server/pkg/model"
	"net/http"
)

func NotFound(w http.ResponseWriter, r *http.Request) {
	ctx := api_context.Of[*api_context.DefaultContext](w, r, "MOCK/NOT/FOUND/HANDLER")
	if config.HasAValidRedirectConfig() {
		redirectConfig := model.Config.RedirectConfig
		requestRedirectHandler(ctx, *redirectConfig)
	} else {
		ctx.NotFount("Resource not found")
	}
}

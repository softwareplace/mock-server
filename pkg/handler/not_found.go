package handler

import (
	apicontext "github.com/softwareplace/http-utils/context"
	"github.com/softwareplace/mock-server/pkg/config"
	"github.com/softwareplace/mock-server/pkg/model"
	"net/http"
)

func NotFound(w http.ResponseWriter, r *http.Request) {
	ctx := apicontext.Of[*apicontext.DefaultContext](w, r, "MOCK/NOT/FOUND/HANDLER")
	if config.HasAValidRedirectConfig() {
		redirectConfig := model.Config.RedirectConfig
		requestRedirectHandler(ctx, *redirectConfig)
	} else {
		ctx.NotFount("Resource not found")
	}
}

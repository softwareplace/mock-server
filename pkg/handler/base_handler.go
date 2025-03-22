package handler

import (
	"fmt"
	"github.com/softwareplace/http-utils/api_context"
	"github.com/softwareplace/http-utils/server"
	"github.com/softwareplace/mock-server/pkg/env"
	"github.com/softwareplace/mock-server/pkg/model"
	"log"
	"net/http"
	"strings"
	"time"
)

func Register(appServer server.ApiRouterHandler[*api_context.DefaultContext]) {
	for _, config := range model.MockConfigResponses {
		if config.Request.Method != "" && config.Request.Path != "" {
			if config.Redirect.Url != "" || config.Response.Bodies != nil {
				contextPath := env.GetAppEnv().ContextPath
				path := strings.TrimPrefix(config.Request.Path, "/")
				log.Printf("Registering handler for %s::%s%s\n", config.Request.Method, contextPath, path)

				appServer.Add(func(ctx *api_context.ApiRequestContext[*api_context.DefaultContext]) {
					url := ctx.Request.RequestURI
					log.Printf("Request %s::%s\n", config.Request.Method, url)
					if !redirectHandler(ctx, config) {
						requestHandler(ctx, config)
					}

				}, config.Request.Path, config.Request.Method)
			} else {
				log.Printf("Invalid definition on %s. No response body or redirect URL found for %s::%s\n", config.MockFilePath, config.Request.Method, config.Request.Path)
			}

		}
	}
}

func redirectHandler(ctx *api_context.ApiRequestContext[*api_context.DefaultContext], config model.MockConfigResponse) (redirected bool) {
	if config.Redirect.Url != "" {
		return requestRedirectHandler(ctx, config.Redirect)

	}
	return false
}

func requestHandler(
	ctx *api_context.ApiRequestContext[*api_context.DefaultContext],
	config model.MockConfigResponse,
) {
	bodies := config.Response.Bodies
	var matchedBody *model.ResponseBody

	matchedBody = findMatchingBody(ctx, bodies)

	writer := *ctx.Writer

	// If no matching body is found, use the first body as a default
	if matchedBody == nil && len(bodies) > 0 {
		writer.WriteHeader(http.StatusNotFound)
		_, _ = writer.Write([]byte("Resource not found"))
		return
	}

	// If a matching body is found, return it as the response
	if matchedBody != nil {
		for key, value := range *matchedBody.Headers {
			writer.Header().Set(key, fmt.Sprintf("%v", value))
		}
		if config.Response.Delay > 0 {
			time.Sleep(time.Duration(config.Response.Delay) * time.Millisecond)
		}

		ctx.Response(matchedBody.Body, config.Response.StatusCode)
		return
	}

	ctx.Error("Resource not found", http.StatusNotFound)
}

func findMatchingBody(
	ctx *api_context.ApiRequestContext[*api_context.DefaultContext],
	bodies []model.ResponseBody,
) *model.ResponseBody {
	var matchedBody *model.ResponseBody
	// Extract query and path parameters from the incoming request

	// Iterate through the bodies to find a match
	for _, body := range bodies {
		if containsExpectedPathsAndQueries(ctx, body) {
			matchedBody = &body
			break
		}
	}
	return matchedBody
}

func containsExpectedPathsAndQueries(
	ctx *api_context.ApiRequestContext[*api_context.DefaultContext],
	body model.ResponseBody,
) bool {
	if body.Matching == nil {
		return true
	}

	return containsExpectedPaths(ctx, body) &&
		containsExpectedQueries(ctx, body) &&
		containsExpectedHeaders(ctx, body)
}

func containsExpectedPaths(
	ctx *api_context.ApiRequestContext[*api_context.DefaultContext],
	body model.ResponseBody,
) bool {
	requestedPaths := ctx.PathValues
	// Check if the paths match
	pathsMatch := len(requestedPaths) == len(body.Matching.Paths)
	for key, value := range body.Matching.Paths {
		if requestedPaths[key] != fmt.Sprintf("%v", value) {
			pathsMatch = false
			break
		}
	}
	if pathsMatch {
		log.Printf("Paths match for request %s\n", ctx.Request.URL.RequestURI())
	}
	return pathsMatch
}

func containsExpectedQueries(ctx *api_context.ApiRequestContext[*api_context.DefaultContext], body model.ResponseBody) bool {
	requestedQueries := ctx.QueryValues
	var queriesMatch = len(requestedQueries) == len(body.Matching.Queries)
	for key, value := range body.Matching.Queries {
		if len(requestedQueries[key]) == 0 {
			queriesMatch = false
			break
		}

		if requestedQueries[key][0] != fmt.Sprintf("%v", value) {
			queriesMatch = false
			break
		}
	}
	if queriesMatch {
		log.Printf("Queries match for request %s\n", ctx.Request.URL.RequestURI())
	}
	return queriesMatch
}

func containsExpectedHeaders(ctx *api_context.ApiRequestContext[*api_context.DefaultContext], body model.ResponseBody) bool {
	var requestHeaders = make(map[string][]string)

	if ctx.Request.Header != nil {
		for key, values := range ctx.Headers {
			lowerKey := strings.ToLower(key)
			requestHeaders[lowerKey] = values

		}
	}

	var headersMatch = true

	for key, value := range body.Matching.Headers {
		if len(requestHeaders[key]) == 0 {
			headersMatch = false
			break
		}
		if requestHeaders[key][0] != fmt.Sprintf("%v", value) {
			headersMatch = false
			break
		}
	}
	if headersMatch {
		log.Printf("Headers match for request %s\n", ctx.Request.URL.RequestURI())
	}
	return headersMatch
}

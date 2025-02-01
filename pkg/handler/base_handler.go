package handler

import (
	"fmt"
	"github.com/softwareplace/http-utils/api_context"
	"github.com/softwareplace/http-utils/server"
	"log"
	"net/http"
	"time"
)

func Register(appServer server.ApiRouterHandler[*api_context.DefaultContext]) {
	for _, config := range MockConfigResponses {
		if config.Request.Method != "" && config.Request.Path != "" && config.Response.Bodies != nil {
			log.Printf("Registering handler for %s::%s\n", config.Request.Method, config.Request.Path)
			appServer.Add(func(ctx *api_context.ApiRequestContext[*api_context.DefaultContext]) {
				requestHandler(ctx, config)
			}, config.Request.Path, config.Request.Method)
		}
	}
}

func requestHandler(
	ctx *api_context.ApiRequestContext[*api_context.DefaultContext],
	config MockConfigResponse,
) {
	bodies := config.Response.Bodies
	var matchedBody *ResponseBody

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
		for key, value := range matchedBody.Headers {
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
	bodies []ResponseBody,
) *ResponseBody {
	var matchedBody *ResponseBody
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
	body ResponseBody,
) bool {
	return containsExpectedPaths(ctx, body) && containsExpectedQueries(ctx, body)
}

func containsExpectedPaths(
	ctx *api_context.ApiRequestContext[*api_context.DefaultContext],
	body ResponseBody,
) bool {
	requestedPaths := ctx.PathValues
	// Check if the paths match
	pathsMatch := true
	for key, value := range body.Paths {
		if requestedPaths[key] != fmt.Sprintf("%v", value) {
			pathsMatch = false
			break
		}
	}
	return pathsMatch
}

func containsExpectedQueries(ctx *api_context.ApiRequestContext[*api_context.DefaultContext], body ResponseBody) bool {
	requestedQueries := ctx.QueryValues
	var queriesMatch = true
	for key, value := range body.Queries {
		if requestedQueries[key][0] != fmt.Sprintf("%v", value) {
			queriesMatch = false
			break
		}
	}
	return queriesMatch
}

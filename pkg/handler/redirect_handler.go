package handler

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	apicontext "github.com/softwareplace/goserve/context"
	errohandler "github.com/softwareplace/goserve/error"
	"github.com/softwareplace/mock-server/pkg/file"
	"github.com/softwareplace/mock-server/pkg/model"
	"io"
	"net/http"
	"strings"
	"time"
)

func requestRedirectHandler(
	ctx *apicontext.Request[*apicontext.DefaultContext],
	redirect model.RedirectConfig,
) bool {
	var request http.Request
	request = *ctx.Request

	for key, value := range redirect.Headers {
		request.Header.Set(key, fmt.Sprintf("%v", value))
	}

	requestedUri := request.URL.RequestURI()
	targetUri := requestedUri

	replacement := redirect.Replacement
	if len(replacement) > 0 {
		for _, replace := range replacement {
			targetUri = strings.ReplaceAll(targetUri, replace.Old, replace.New)
		}
	}
	targetUri = strings.ReplaceAll(targetUri, "//", "/")

	targetURL := strings.TrimSuffix(redirect.Url, "/") + "/" +
		strings.TrimPrefix(targetUri, "/")

	req, err := http.NewRequest(request.Method, targetURL, request.Body)

	if err != nil {
		ctx.Error("Failed to complete the request", http.StatusInternalServerError)
		return true
	}
	for key, value := range redirect.Headers {
		req.Header.Set(key, fmt.Sprintf("%v", value))
	}

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		ctx.Error(err.Error(), http.StatusInternalServerError)
		return true
	}

	// Read the response body
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Errorf("Failed to close response body: %v", err)
		}
	}(resp.Body)

	bodyBytes, err := io.ReadAll(resp.Body)

	if err != nil {
		ctx.Error(err.Error(), http.StatusInternalServerError)
		return true
	}

	writer := *ctx.Writer
	responseContentType := resp.Header.Get("Content-Type")

	if responseContentType != "" {
		writer.Header().Set("Content-Type", responseContentType)
	}

	writer.WriteHeader(resp.StatusCode)

	_, err = writer.Write(bodyBytes)

	if redirect.LogEnabled {
		log.Infof("%s -> %s returned: %s", requestedUri, targetUri, string(bodyBytes))
	}

	if redirect.StoreResponsesDir != "" {
		data := map[string]interface{}{
			"headers":   ctx.Request.Header,
			"uri":       requestedUri,
			"targetURL": targetURL,
			"body":      string(bodyBytes),
		}

		errohandler.Handler(func() {
			if strings.Contains(responseContentType, "application/json") {
				data["body"] = json.RawMessage(bodyBytes)
			}
		}, func(err error) {
			log.Errorf("Failed to store response: %v", err)
		})
		storeFile(data, redirect, requestedUri)

		return true
	}

	if err != nil {
		log.Errorf("Failed to write response body: %v", err)
		return false
	}

	return true
}

func storeFile(data map[string]interface{}, redirect model.RedirectConfig, requestedUri string) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("Failed to marshal data to JSON: %v", err)
		return
	}

	filePath := fmt.Sprintf("%s/%s_%d.json",
		strings.TrimSuffix(redirect.StoreResponsesDir, "/"),
		strings.ReplaceAll(requestedUri, "/", "_"),
		time.Now().Unix(),
	)

	err = file.SaveToFile(jsonData, filePath)
	if err != nil {
		log.Errorf("Failed to save data to file: %v", err)
		return
	}

}

package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"time"
)

type MockConfigResponse struct {
	Request  RequestConfig  `json:"request" yaml:"request"`
	Response ResponseConfig `json:"response" yaml:"response"`
}

type RequestConfig struct {
	Path        string         `json:"path" yaml:"path"`
	Method      string         `json:"method" yaml:"method"`
	ContentType string         `json:"contentType" yaml:"content-type" yaml:"contentType"`
	Queries     map[string]any `json:"queries" yaml:"queries"`
	Paths       map[string]any `json:"paths" yaml:"paths"`
}

type ResponseConfig struct {
	ContentType string                 `json:"contentType" yaml:"content-type" yaml:"contentType"`
	StatusCode  int                    `json:"statusCode" yaml:"status-code" yaml:"statusCode"`
	Delay       int                    `json:"delay" yaml:"delay"`
	Headers     map[string]any         `json:"headers" yaml:"headers"`
	Body        map[string]interface{} `json:"body" yaml:"body"`
}

var (
	Responses []MockConfigResponse
)

func RequestHandler() func(w http.ResponseWriter, r *http.Request) {
	LoadResponses()
	return func(w http.ResponseWriter, r *http.Request) {
		responses := Responses

		matchedResponse := (*MockConfigResponse)(nil)
		for _, response := range responses {
			// Match path using regex
			matched, err := regexp.MatchString(response.Request.Path, r.URL.Path)
			if err != nil {
				http.Error(w, "Invalid path configuration", http.StatusInternalServerError)
				return
			}

			if matched && response.Request.Method == r.Method {
				// Validate queries if they are present in the response
				if len(response.Request.Queries) > 0 {
					valid := true
					queryValues := r.URL.Query()

					for key, value := range response.Request.Queries {
						if queryValues.Get(key) != fmt.Sprintf("%v", value) {
							valid = false
							break
						}
					}

					if valid {
						matchedResponse = &response
						break
					}
				} else {
					matchedResponse = &response
					break
				}
			}
		}

		if matchedResponse == nil {
			http.Error(w, "Unavailable service", http.StatusServiceUnavailable)
			return
		}

		// Write the matched response
		w.Header().Set("Content-Type", matchedResponse.Response.ContentType)
		// Write headers if present in the matched response
		if len(matchedResponse.Response.Headers) > 0 {
			for key, value := range matchedResponse.Response.Headers {
				w.Header().Set(key, fmt.Sprintf("%v", value))
			}
		}

		w.WriteHeader(matchedResponse.Response.StatusCode)

		if matchedResponse.Response.Delay > 0 {
			time.Sleep(time.Duration(matchedResponse.Response.Delay) * time.Second)
		}
		if matchedResponse.Response.Body != nil {
			body, err := json.Marshal(matchedResponse.Response.Body)
			if err != nil {
				http.Error(w, "Error encoding response body", http.StatusInternalServerError)
				return
			}
			_, _ = w.Write(body)
		}
	}
}

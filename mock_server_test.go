package mock_server

import (
	"encoding/json"
	"github.com/softwareplace/http-utils/api_context"
	"github.com/softwareplace/mock-server/pkg/env"
	"github.com/softwareplace/mock-server/pkg/handler"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/softwareplace/http-utils/server"
)

var appEnv = &env.AppEnv{
	Port:        "18888",
	MockPath:    "./dev",
	ContextPath: "/",
}

func TestMockServer(t *testing.T) {
	env.SetAppEnv(appEnv)

	var appServer server.ApiRouterHandler[*api_context.DefaultContext]

	// Load mock responses
	handler.LoadResponses(func(restartServer bool) {
		// Create a test server
		appServer = server.Default().
			WithContextPath(appEnv.ContextPath).
			EmbeddedServer(handler.Register)
	})

	// Test cases
	tests := []struct {
		name           string
		method         string
		path           string
		queryParams    map[string]string
		headers        map[string]string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Test GET /api/products/1",
			method:         "GET",
			path:           "/api/products/1",
			expectedStatus: http.StatusOK,
			expectedBody:   `{"amount":2500.75,"id":1,"name":"Product"}`,
		},
		{
			name:           "Test GET /v1/products with query id=1",
			method:         "GET",
			path:           "/v1/products",
			queryParams:    map[string]string{"id": "1"},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"id":1,"name":"Product 1","amount":2500.75}`,
		},
		{
			name:           "Test GET /v1/products with query id=2",
			method:         "GET",
			path:           "/v1/products",
			queryParams:    map[string]string{"id": "2"},
			headers:        map[string]string{"id": "2", "name": "Product 2"},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"id":2,"name":"Product 2","amount":2500.75}`,
		},
		{
			name:           "Test GET /v1/redirect/products with query id=1",
			method:         "GET",
			path:           "/v1/redirect/products",
			queryParams:    map[string]string{"id": "2"},
			headers:        map[string]string{"id": "2", "name": "Product 2"},
			expectedStatus: http.StatusTemporaryRedirect,
		},

		{
			name:           "Test GET /api/user/view for id=2",
			method:         "GET",
			path:           "/api/user/view",
			queryParams:    map[string]string{"id": "2", "name": "User 2"},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"id":2,"name":"User For Queries request","email":"john.doe+2@email.com"}`,
		},
		{
			name:           "Test GET /api/user/view for id=3",
			method:         "GET",
			path:           "/api/user/view",
			queryParams:    map[string]string{"id": "3", "name": "User 3"},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"id":3,"name":"User For Queries request","email":"john.doe+3@email.com"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new request
			req, err := http.NewRequest(tt.method, tt.path, nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			// Add query parameters
			q := req.URL.Query()
			for key, value := range tt.queryParams {
				q.Add(key, value)
			}
			req.URL.RawQuery = q.Encode()

			// Add headers
			for key, value := range tt.headers {
				req.Header.Add(key, value)
			}

			// Create a response recorder
			rr := httptest.NewRecorder()

			// Serve the request
			appServer.Router().ServeHTTP(rr, req)

			// Check the status code
			if rr.Code != tt.expectedStatus {
				t.Errorf("[%s] Expected status code %d, got %d", tt.name, tt.expectedStatus, rr.Code)
			}

			// Check the response body
			if tt.expectedBody != "" {
				responseBody, err := io.ReadAll(rr.Body)
				if err != nil {
					t.Fatalf("Failed to read response body: %v", err)
				}
				if !jsonDeepEqual(responseBody, []byte(tt.expectedBody)) {
					t.Errorf("[%s] Expected body %s, got %s", tt.name, tt.expectedBody, string(responseBody))
				}

			}
		})
	}
}

func jsonDeepEqual(actual, expected []byte) bool {
	var actualJSON, expectedJSON interface{}
	if err := json.Unmarshal(actual, &actualJSON); err != nil {
		return false
	}
	if err := json.Unmarshal(expected, &expectedJSON); err != nil {
		return false
	}
	return reflect.DeepEqual(actualJSON, expectedJSON)
}

func TestDelaySimulation(t *testing.T) {
	env.SetAppEnv(appEnv)

	var appServer server.ApiRouterHandler[*api_context.DefaultContext]

	// Load mock responses
	handler.LoadResponses(func(restartServer bool) {
		// Create a test server
		appServer = server.Default().
			WithContextPath(appEnv.ContextPath).
			EmbeddedServer(handler.Register)
	})

	// Test delay simulation
	req, err := http.NewRequest("GET", "/v1/products?id=1", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()

	start := time.Now()
	appServer.Router().ServeHTTP(rr, req)
	elapsed := time.Since(start)

	// Check if the delay is approximately 256ms
	expectedDelay := 256 * time.Millisecond
	if elapsed < expectedDelay || elapsed > expectedDelay+50*time.Millisecond {
		t.Errorf("Expected delay of approximately %v, got %v", expectedDelay, elapsed)
	}

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
	}
}

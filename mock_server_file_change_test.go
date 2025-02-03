package mock_server

import (
	"encoding/base64"
	"fmt"
	"github.com/softwareplace/http-utils/api_context"
	"github.com/softwareplace/http-utils/request"
	"github.com/softwareplace/http-utils/server"
	"github.com/softwareplace/mock-server/pkg/env"
	"github.com/softwareplace/mock-server/pkg/handler"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestServerReloading(t *testing.T) {
	mockFilePath := "./dev/mock/.temp"
	mockFileName := "product-test-mode.yaml"

	mockFileFullPath := filepath.Join(mockFilePath, mockFileName)

	_removeMockTestFile(t, mockFileFullPath)

	config := request.Build("http://localhost:18888").
		WithPath("v1/products/server/reload/1000")

	env.SetAppEnv(appEnv)
	var appServer server.ApiRouterHandler[*api_context.DefaultContext]
	serverWasReloaded := false

	handler.LoadResponses(func(restartServer bool) {
		serverWasReloaded = restartServer
		appServer = createServer(appServer, restartServer)
	})

	t.Run("expects that return resource not found before add new mock file", func(t *testing.T) {
		requestService := request.NewService()

		response, err := requestService.
			Get(config)

		if response.StatusCode != 404 {
			t.Fatalf("Expected status code 404, but got: %d", response.StatusCode)
		}

		responseJson, err := requestService.ToString()

		if err != nil {
			t.Fatalf("Failed to parse to string: %v", err)
		}

		if strings.Contains(fmt.Sprintf("%v", responseJson), "404 page not found") {
			t.Log("Response error contains '404 page not found'")
		} else {
			t.Errorf("Expected response error '404 page not found', but got: %s", err)
		}
	})

	t.Run("expects that return the expected json to the matching product", func(t *testing.T) {
		err := createProjectMockConfig(t, mockFilePath, mockFileFullPath)

		if err != nil {
			t.Fatalf("Failed to write to file %s: %v", mockFileFullPath, err)
		}

		log.Println("Waiting for server to reload...")

		time.Sleep(1 * time.Second)

		if serverWasReloaded {
			log.Println("Server was reloaded successfully.")
		} else {
			t.Fatalf("Expected server to be reloaded, but got: %v", serverWasReloaded)
		}

		requestService := request.NewService()
		response, err := requestService.
			Get(config)

		if response.StatusCode != 200 {
			t.Fatalf("Expected status code 200, but got: %d", response.StatusCode)
		}

		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}

		expectedJson := `{"id": 1000,"name": "Product","description": "This is a mock product description","amount": 2500.75}`

		responseJson, err := requestService.ToString()
		if err != nil {
			t.Fatalf("Failed to marshal response to JSON: %v", err)
		}

		if !jsonDeepEqual([]byte(responseJson), []byte(expectedJson)) {
			t.Errorf("Expected JSON response: %s, but got: %s", expectedJson, responseJson)
		} else {
			t.Log("Response matches expected JSON.")
		}
	})

	_removeMockTestFile(t, mockFileFullPath)
}

func createProjectMockConfig(t *testing.T, mockFilePath string, mockFileFullPath string) error {
	addMockProductConfig := "cmVxdWVzdDoKICBwYXRoOiAiL3YxL3Byb2R1Y3RzL3NlcnZlci9yZWxvYWQve2lkfSIKICBtZXRob2Q6ICJHRVQiCnJlc3BvbnNlOgogIGNvbnRlbnQtdHlwZTogImFwcGxpY2F0aW9uL2pzb24iCiAgc3RhdHVzLWNvZGU6IDIwMAogIGJvZGllczoKICAgIC0gYm9keToKICAgICAgICBpZDogMTAwMAogICAgICAgIG5hbWU6ICJQcm9kdWN0IgogICAgICAgIGRlc2NyaXB0aW9uOiBUaGlzIGlzIGEgbW9jayBwcm9kdWN0IGRlc2NyaXB0aW9uCiAgICAgICAgYW1vdW50OiAyNTAwLjc1CiAgICAgIG1hdGNoaW5nOgogICAgICAgIHBhdGhzOgogICAgICAgICAgaWQ6IDEwMDAK"

	decodedConfig, err := base64.StdEncoding.DecodeString(addMockProductConfig)
	if err != nil {
		t.Fatalf("Failed to decode base64 string: %v", err)
	}
	// Validate if the directory path exists
	if _, err := os.Stat(mockFilePath); os.IsNotExist(err) {
		err = os.MkdirAll(mockFilePath, os.ModePerm)
		if err != nil {
			t.Fatalf("Failed to create directory %s: %v", mockFilePath, err)
		}
	}

	err = os.WriteFile(mockFileFullPath, decodedConfig, 0644)
	return err
}

func _removeMockTestFile(t *testing.T, mockFileFullPath string) {
	// Check if the file already exists; if so, remove it
	if _, err := os.Stat(mockFileFullPath); err == nil {
		err = os.Remove(mockFileFullPath)
		if err != nil {
			t.Fatalf("Failed to remove existing file %s: %v", mockFileFullPath, err)
		}
	}
}

func createServer(
	appServer server.ApiRouterHandler[*api_context.DefaultContext],
	restartServer bool,
) server.ApiRouterHandler[*api_context.DefaultContext] {
	if restartServer {
		if appServer != nil {
			err := appServer.StopServer()
			if err != nil {
				log.Fatalf("Failed to stop server: %v", err)
			}
		}
	}

	return server.Default().
		WithContextPath(appEnv.ContextPath).
		EmbeddedServer(handler.Register).
		WithPort(appEnv.Port).
		NotFoundHandler().
		StartServerInGoroutine()
}

# Mock Server

Mock Server is a lightweight HTTP server designed to mock API responses for testing and development purposes. It allows
you to define mock responses in JSON or YAML files and automatically reloads them when changes are detected. The server
supports dynamic response matching based on query parameters, headers, and path variables. Additionally, it provides
advanced features like redirection, custom headers, delay simulation, and file-based configuration.

## Features

- **Dynamic Mock Responses**: Define mock responses in JSON or YAML files.
- **Automatic Reloading**: Automatically reloads mock responses when files are changed.
- **Response Matching**: Match responses based on query parameters, headers, and path variables.
- **Redirection Support**: Redirect requests to another URL with optional string replacements.
- **Custom Headers**: Add custom headers to responses.
- **Delay Simulation**: Simulate network latency by adding delays to responses.
- **File Watching**: Watches for changes in mock files and reloads the server dynamically.
- **Configuration File**: Supports configuration via a YAML file for server settings and redirection rules.
- **Debounced Reloading**: Prevents excessive reloads with a debouncing mechanism.

## Installation

To use the Mock Server, you need to have Go installed on your machine. You can install the server by cloning the
repository and building it:

```bash
git clone https://github.com/softwareplace/mock-server.git
cd mock-server
go build -o bin/$(uname -m)/mock-server cmd/server/main.go
```

## Usage

### Environment Variables

| Variable Name | Required | Default Value | Description                                        |
|---------------|----------|---------------|----------------------------------------------------|
| `LOG_PATH`    | No       | `./logs/`     | The directory path where log files will be stored. |

### Running the Server

To start the server, run the following command:

```bash
./bin/$(uname -m)/mock-server --mock /path/to/mock/files --port 8080 --context-path /api
```

- `--mock`: Path to the directory containing your mock JSON or YAML files.
- `--port`: Port to run the server on (default: `8080`).
- `--context-path`: Base path for the API endpoints (default: `/`).
- `--config`: Alternatively, you can use a configuration file (`config.yaml`) to specify these settings:

```yaml
# The port on which the mock server will run.
port: 8080

# The directory containing the mock JSON or YAML files.
mock: /path/to/mock/files

# The base path for the API endpoints.
context-path: /api
# In case the requested URL is not found in the mock configuration, the server will redirect 
# the request to the specified URL. This feature helps handle fallback API requests gracefully, 
# such as forwarding to an upstream server or logging unknown requests for debugging.
redirect:
  # The URL to redirect incoming requests to.
  url: http://localhost:8888/

  # Enable or disable logging for redirection events.
  log-enabled: true

  # Directory to store redirected response files.
  store-responses-dir: ./.temp/

  # List of string replacements to modify the request URI.
  replacement:
    - old: mock-server  # String to be replaced in the request URI.
      new: ""           # Replacement for the specified string.
    - old: api          # String to be replaced in the request URI.
      new: ""           # Replacement for the specified string.
```

Run the server with the configuration file:

```bash
./bin/$(uname -m)/mock-server --config /path/to/config.yaml
```

### Defining Mock Responses

Mock responses are defined in JSON or YAML files. Each file should contain a `MockConfigResponse` object with the
following structure:

#### JSON Example

```json
{
  "redirect": {
    "url": "http://localhost:8888/",
    "logEnabled": true,
    "storeResponsesDir": "./.temp/",
    "replacements": [
      {
        "old": "mock-server",
        "new": ""
      },
      {
        "old": "api",
        "new": ""
      }
    ]
  },
  "request": {
    "path": "/api/user/view",
    "method": "GET",
    "contentType": "application/json"
  },
  "response": {
    "delay": 0,
    "contentType": "application/json",
    "statusCode": 200,
    "bodies": [
      {
        "matching": {
          "queries": {
            "id": 2,
            "name": "User 2"
          }
        },
        "headers": {
          "id": 2,
          "name": "User 2"
        },
        "body": {
          "id": 2,
          "name": "User For Queries request",
          "email": "john.doe+2@email.com"
        }
      },
      {
        "matching": {
          "queries": {
            "id": 3,
            "name": "User 3"
          }
        },
        "headers": {
          "id": 3,
          "name": "User 3"
        },
        "body": {
          "id": 3,
          "name": "User For Queries request",
          "email": "john.doe+3@email.com"
        }
      }
    ]
  }
}
```

#### YAML Example

```yaml
# In case requests need to be redirected to another API host
redirect:
  # The URL to redirect incoming requests to.
  url: http://localhost:8888/

  # Enable or disable logging for redirection events. When enabled, all redirection events will be logged.
  log-enabled: true

  # Directory to store responses from redirected requests. Useful for debugging and tracking.
  store-responses-dir: ./.temp/responses/

  # List of string replacements to modify the request URI before redirecting.
  replacement:
    - old: "mock-server"  # The string to be replaced in the request URI.
      new: ""             # The new value to replace the specified string.
    - old: "api"          # Another string to be replaced in the request URI.
      new: ""             # The new value to replace this string.
    - old: "v1"           # Example for replacing version segments in the URI.
      new: "v2"           # The new version string to replace the old one.
request:
  path: "/v1/products"
  method: "GET"
response:
  content-type: "application/json"
  status-code: 200
  delay: 256
  bodies:
    - body:
        id: 1
        name: "Product 1"
        amount: 2500.75
      matching:
        queries:
          id: 1
    - body:
        id: 2
        name: "Product 2"
        amount: 2500.75
      headers:
        is: 2
        name: Product 2
      matching:
        queries:
          id: 2
        headers:
          is: 2
          name: Product 2
```

### Response Matching

The server supports response matching based on query parameters, headers, and path variables. If a request matches the
criteria defined in the `matching` section of a response body, that response will be returned.

### Redirection

You can configure the server to redirect requests to another URL. The `redirect` section allows you to specify the
target URL and perform string replacements on the request URI.

```yaml
redirect:
  url: http://localhost:8888/
  log-enabled: true
  store-responses-dir: ./.temp/
  replacement:
    - old: mock-server
      new: ""
    - old: api
      new: ""
```

### Custom Headers

You can add custom headers to the response by specifying them in the `headers` section of the response body.

```yaml
headers:
  is: 2
  name: Product 2
```

### Delay Simulation

To simulate network latency, you can add a delay to the response by specifying the `delay` field in milliseconds.

```yaml
response:
  delay: 256
```

### File Watching and Automatic Reloading

The server watches for changes in the mock files directory and automatically reloads the server when changes are
detected. This feature uses a debouncing mechanism to prevent excessive reloads.

### Advanced Configuration

The server supports advanced configurations via a YAML file. You can specify the server port, mock files directory,
context path, and redirection rules in the configuration file.

```yaml
# The port on which the mock server will run.
port: 8080

# The directory containing the mock JSON or YAML files.
mock: /path/to/mock/files

# The base path for the API endpoints.
context-path: /api
# In case the requested URL is not found in the mock configuration, the server will redirect 
# the request to the specified URL. This feature helps handle fallback API requests gracefully, 
# such as forwarding to an upstream server or logging unknown requests for debugging.
redirect:
  # The URL to redirect incoming requests to.
  url: http://localhost:8888/

  # Enable or disable logging for redirection events.
  log-enabled: true

  # Directory to store redirected response files.
  store-responses-dir: ./.temp/

  # List of string replacements to modify the request URI.
  replacement:
    - old: mock-server  # String to be replaced in the request URI.
      new: ""           # Replacement for the specified string.
    - old: api          # String to be replaced in the request URI.
      new: ""           # Replacement for the specified string.
```

## Contributing

Contributions are welcome! Please open an issue or submit a pull request if you have any improvements or bug fixes.

---


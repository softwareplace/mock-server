# Mock Server

Mock Server is a lightweight HTTP server designed to mock API responses for testing and development purposes. It allows
you to define mock responses in JSON or YAML files and automatically reloads them when changes are detected. The server
supports dynamic response matching based on query parameters, headers, and path variables.

## Features

- **Dynamic Mock Responses**: Define mock responses in JSON or YAML files.
- **Automatic Reloading**: Automatically reloads mock responses when files are changed.
- **Response Matching**: Match responses based on query parameters, headers, and path variables.
- **Redirection Support**: Redirect requests to another URL with optional string replacements.
- **Custom Headers**: Add custom headers to responses.
- **Delay Simulation**: Simulate network latency by adding delays to responses.

## Installation

To use the Mock Server, you need to have Go installed on your machine. You can install the server by cloning the
repository and building it:

```bash
git clone https://github.com/softwareplace/mock-server.git
cd mock-server
go build -o bin/$(uname -m)/mock-server cmd/server/main.go
```

## Usage

### Running the Server

To start the server, run the following command:

```bash
./bin/$(uname -m)/mock-server --mock /path/to/mock/files --port 8080 --context-path /api
```

- `--mock`: Path to the directory containing your mock JSON or YAML files.
- `--port`: Port to run the server on (default: `8080`).
- `--context-path`: Base path for the API endpoints (default: `/`).

### Defining Mock Responses

Mock responses are defined in JSON or YAML files. Each file should contain a `MockConfigResponse` object with the
following structure:

```yaml
request:
  path: "/v1/products"
  method: "GET"
redirect:
  url: http://localhost:8888/
  replacement:
    - old: mock-server
      new: ""
    - old: api
      new: ""
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

### Example Mock Files

#### JSON Example

```json
{
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

## Contributing

Contributions are welcome! Please open an issue or submit a pull request if you have any improvements or bug fixes.


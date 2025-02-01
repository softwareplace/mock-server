# Mock Server

A lightweight mock server written in Go, designed to simulate API responses based on predefined configurations. This
server is useful for testing and development purposes, allowing you to mock API endpoints without needing a real
backend.

## Features

- **Dynamic Response Handling**: Define mock responses in JSON or YAML files.
- **Path and Query Matching**: Match requests based on path parameters and query strings.
- **Custom Headers**: Set custom headers for each response.
- **Response Delay**: Simulate network latency by adding delays to responses.
- **Redirection**: Redirect requests to another URL with optional path replacements.
- **Auto-Reload**: Automatically reload configurations when files change.

## Installation

1. **Clone the Repository**:
   ```bash
   git clone git@github.com:softwareplace/mock-server.git
   cd mock-server
   ```

2. **Build the Project**:
   ```bash
   go build -o bin/$(uname -m)/mock-server cmd/server/main.go
   ```
3. **Run the Server**:
   ```bash
   ./bin/$(uname -m)/mock-server --mock=/path/to/mock/files --port=8080 --context-path=/api
   ```

## Configuration

### Mock Configuration Files

Mock responses are defined in JSON or YAML files. These files should be placed in the directory specified by the
`--mock` flag.

#### Example Configuration

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
      queries:
        id: 1
    - body:
        id: 2
        name: "Product 2"
        amount: 2500.75
      queries:
        id: 2
      headers:
        is: 2
        name: Product 2
```

### Configuration Options

- **Request**:
    - `path`: The URL path to match.
    - `method`: The HTTP method to match (e.g., GET, POST).
    - `contentType`: The content type of the request (optional).

- **Response**:
    - `contentType`: The content type of the response.
    - `statusCode`: The HTTP status code to return.
    - `delay`: The delay in milliseconds before sending the response.
    - `bodies`: A list of response bodies, each with optional query and header matching.

- **Redirect**:
    - `url`: The URL to redirect to.
    - `replacement`: A list of string replacements to apply to the request URI before redirection.

## Usage

### Running the Server

To start the server, run the following command:

```bash
./mock-server --mock=/path/to/mock/files --port=8080 --context-path=/api
```

- `--mock`: Path to the directory containing mock configuration files.
- `--port`: Port to run the server on (default: 8080).
- `--context-path`: Base path for all endpoints (default: `/`).

### Example Requests

1. **Simple GET Request**:
   ```bash
   curl -X GET http://localhost:8080/api/v1/products?id=1
   ```

   **Response**:
   ```json
   {
     "id": 1,
     "name": "Product 1",
     "amount": 2500.75
   }
   ```

2. **Request with Headers**:
   ```bash
   curl -X GET http://localhost:8080/api/v1/products?id=2 -H "is: 2" -H "name: Product 2"
   ```

   **Response**:
   ```json
   {
     "id": 2,
     "name": "Product 2",
     "amount": 2500.75
   }
   ```

3. **Redirect Request**:
   ```bash
   curl -X GET http://localhost:8080/api/v1/products
   ```

   **Response**:
   The request will be redirected to `http://localhost:8888/v1/products`.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request for any improvements or bug fixes.


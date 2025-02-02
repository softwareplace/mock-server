package handler

type MockConfigResponse struct {
	Request  RequestConfig  `json:"request" yaml:"request"`   // Request contains the configuration details for the HTTP request.
	Response ResponseConfig `json:"response" yaml:"response"` // Response holds the specifications for the HTTP response configuration.
	Redirect RedirectConfig `json:"redirect" yaml:"redirect"` // Redirect defines the settings for HTTP redirection if applicable.
}

type Replacement struct {
	Old string `json:"old" yaml:"old"` // Old specifies the string to be replaced during the redirection process.
	New string `json:"new" yaml:"new"` // New specifies the replacement string for the redirection process.
}

type RedirectConfig struct {
	Url         string         `json:"url" yaml:"url"`                 // Url specifies the target URL for the redirection.
	Headers     map[string]any `json:"headers" yaml:"headers"`         // Headers to provide custom headers when redirect
	Replacement []Replacement  `json:"replacement" yaml:"replacement"` // Replacement specifies a list of string replacements to perform in the redirection process.
}

type RequestConfig struct {
	Path        string `json:"path" yaml:"path"`                                   // Path specifies the endpoint or resource location for the request in the RequestConfig struct.
	Method      string `json:"method" yaml:"method"`                               // Method specifies the HTTP method for the request in the RequestConfig struct.
	ContentType string `json:"contentType" yaml:"content-type" yaml:"contentType"` // ContentType specifies the media type of the request payload as defined in the RequestConfig struct.
}

type Matching struct {
	Queries map[string]any `json:"queries" yaml:"queries"` // Queries is a map of key-value pairs used for defining matching query parameters in requests.
	Paths   map[string]any `json:"paths" yaml:"paths"`     // Paths is a map of key-value pairs used for defining matching path parameters in requests.
	Headers map[string]any `json:"headers" yaml:"headers"` // Headers is a map of key-value pairs used for defining matching header parameters in requests.
}
type ResponseBody struct {
	Body     interface{}    `json:"body" yaml:"body"`         // Body represents the dynamic content of the response, serialized based on the provided JSON or YAML format.
	Matching *Matching      `json:"matching" yaml:"matching"` // Matching handles product retrieval. Filters the body with matching queries, headers, and path parameters if provided.
	Headers  map[string]any `json:"headers" yaml:"headers"`   // Headers in case that need to add headers to the response
}

type ResponseConfig struct {
	ContentType string         `json:"contentType" yaml:"content-type" yaml:"contentType"`
	StatusCode  int            `json:"statusCode" yaml:"status-code" yaml:"statusCode"` // StatusCode represents the HTTP status code to return in the response.
	Delay       int            `json:"delay" yaml:"delay"`                              // Delay specifies the time delay (in milliseconds) before the response is sent.
	Bodies      []ResponseBody `json:"bodies" yaml:"bodies"`                            // Bodies contains multiple response bodies to choose from. If no matching filter is set for the body, the first body will be returned.
}

var (
	MockConfigResponses []MockConfigResponse
)

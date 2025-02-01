package handler

type MockConfigResponse struct {
	Request  RequestConfig  `json:"request" yaml:"request"`
	Response ResponseConfig `json:"response" yaml:"response"`
}

type RequestConfig struct {
	Path        string `json:"path" yaml:"path"`
	Method      string `json:"method" yaml:"method"`
	ContentType string `json:"contentType" yaml:"content-type" yaml:"contentType"`
}

type ResponseBody struct {
	Body    interface{}    `json:"body" yaml:"body"`
	Queries map[string]any `json:"queries" yaml:"queries"`
	Paths   map[string]any `json:"paths" yaml:"paths"`
	Headers map[string]any `json:"headers" yaml:"headers"`
}

type ResponseConfig struct {
	ContentType string         `json:"contentType" yaml:"content-type" yaml:"contentType"`
	StatusCode  int            `json:"statusCode" yaml:"status-code" yaml:"statusCode"`
	Delay       int            `json:"delay" yaml:"delay"`
	Bodies      []ResponseBody `json:"bodies" yaml:"bodies"`
}

var (
	MockConfigResponses []MockConfigResponse
	ConfigLoaded        bool
)

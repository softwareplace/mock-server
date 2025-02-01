package handler

type MockConfigResponse struct {
	Request  RequestConfig  `json:"request" yaml:"request"`
	Response ResponseConfig `json:"response" yaml:"response"`
	Redirect RedirectConfig `json:"redirect" yaml:"redirect"`
}

type Replacement struct {
	Old string `json:"old" yaml:"old"`
	New string `json:"new" yaml:"new"`
}

type RedirectConfig struct {
	Url         string         `json:"url" yaml:"url"`
	Headers     map[string]any `json:"headers" yaml:"headers"`
	Replacement []Replacement  `json:"replacement" yaml:"replacement"`
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

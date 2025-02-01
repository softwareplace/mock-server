package redirect

type Redirect struct {
	Match   string            `yaml:"match"`
	Target  string            `yaml:"target"`
	Headers map[string]string `yaml:"headers"`
}

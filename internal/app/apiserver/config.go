package apiserver

// Config for API server
type Config struct {
	Hostname       string `yaml:"hostname"`
	AllowedOrigins string `yaml:"allowed_origins"`
}

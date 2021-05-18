package config

// Jaeger tracer
type Jaeger struct {
	HostPort string `yaml:"jaeger.hostPort"`
	LogSpans bool   `yaml:"jaeger.logSpans"`
}

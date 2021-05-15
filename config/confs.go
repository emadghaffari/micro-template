package config

var (
	// Global config
	Global GlobalConfig
)

// GlobalConfig is base of configs we need for project
type GlobalConfig struct {
	Environment string        `yaml:"environment"`
	Service     service       `yaml:"service"`
	Jaeger      jaeger        `yaml:"jaeger"`
	Log         loggingConfig `yaml:"loggingConfig"`
	ETCD        etcd          `yaml:"etcd"`
	Redis       redis         `yaml:"redis"`
	POSTGRES   database      `yaml:"database"`
	Nats        nats          `yaml:"nats"`
	JWT         JWT           `yaml:"jwt"`
}

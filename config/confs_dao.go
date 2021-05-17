package config

var (
	// Global config
	Confs cnfs = &Config{}
)

type cnfs interface {
	Set(key string, query []byte) error
	Get() Config
	GetService() interface{}
	Debug() bool
	Load(path string) error
	file(path string) error
}

// Config is base of configs we need for project
type Config struct {
	Environment string        `yaml:"environment"`
	Service     service       `yaml:"service"`
	Jaeger      jaeger        `yaml:"jaeger"`
	Log         loggingConfig `yaml:"loggingConfig"`
	ETCD        etcd          `yaml:"etcd"`
	Redis       redis         `yaml:"redis"`
	POSTGRES    database      `yaml:"database"`
	Nats        nats          `yaml:"nats"`
	JWT         JWT           `yaml:"jwt"`
}

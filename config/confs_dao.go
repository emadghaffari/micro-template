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
	Service     Service       `yaml:"service"`
	Jaeger      Jaeger        `yaml:"jaeger"`
	Log         loggingConfig `yaml:"loggingConfig"`
	ETCD        ETCD          `yaml:"etcd"`
	Redis       Redis         `yaml:"redis"`
	POSTGRES    Database      `yaml:"database"`
	Nats        NATS          `yaml:"nats"`
	JWT         JWT           `yaml:"jwt"`
}

package config

type NATS struct {
	Username       string   `yaml:"nats.username"`
	Password       string   `yaml:"nats.password"`
	Encoder        string   `yaml:"nats.encoder"`
	Auth           bool     `yaml:"nats.auth"`
	Endpoints      []string `yaml:"nats.endpoints"`
	AllowReconnect bool     `yaml:"nats.allowReconnect"`
	MaxReconnect   int      `yaml:"nats.maxReconnect"`
	ReconnectWait  int      `yaml:"nats.reconnectWait"`
	Timeout        int      `yaml:"nats.timeout"`
}

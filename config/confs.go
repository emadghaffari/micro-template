package config

import "time"

var (
	// Global config
	Global global
)

type global struct {
	Environment string `yaml:"environment"`
	GRPC        struct {
		Host     string `yaml:"grpc.host"`
		Port     string `yaml:"grpc.port"`
		Endpoint string `yaml:"grpc.endpoint"`
	}
	HTTP struct {
		Host     string `yaml:"http.host"`
		Port     string `yaml:"http.port"`
		Endpoint string `yaml:"http.endpoint"`
	}
	DEBUG struct {
		Host     string `yaml:"debug.host"`
		Port     string `yaml:"debug.port"`
		Endpoint string `yaml:"debug.endpoint"`
	}
	Service service
	Jaeger  jaeger
	Log     loggingConfig
	ETCD    etcd
	Redis   redis
}

// Service details
type service struct {
	Name  string `yaml:"service.name"`
	Redis struct {
		SMSDuration         time.Duration `yaml:"service.redis.smsDuration"`
		SMSCodeVerification time.Duration `yaml:"service.redis.smsCodeVerification"`
		UserDuration        time.Duration `yaml:"service.redis.userDuration"`
	}
}

// Jaeger tracer
type jaeger struct {
	HostPort string `yaml:"jaeger.hostPort"`
	LogSpans bool   `yaml:"jaeger.logSpans"`
}

// LoggingConfig struct
type loggingConfig struct {
	DisableColors    bool `json:"disable_colors" yaml:"log.disableColors"`
	QuoteEmptyFields bool `json:"quote_empty_fields" yaml:"log.quoteEmptyFields"`
}

type etcd struct {
	Endpoints []string `json:"endpoints" yaml:"etcd.endpoints"`
	WatchList []string `json:"watch_list" yaml:"etcd.watchList"`
}

// redis struct
type redis struct {
	Address string `json:"address" yaml:"redis.address"`
}

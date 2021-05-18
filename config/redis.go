package config

import "time"

// redis struct
type Redis struct {
	Username     string        `yaml:"redis.username"`
	Password     string        `yaml:"redis.password"`
	DB           int           `yaml:"redis.db"`
	Host         string        `yaml:"redis.host"`
	Logger       bool          `yaml:"redis.logger"`
	UserDuration time.Duration `yaml:"redis.userDuration"`
}

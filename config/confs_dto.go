package config

import (
	"encoding/json"
	"fmt"
	zapLogger "micro/pkg/logger"
	"os"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// Set method
// you can set new key in switch for manage config with config server
func (g *Config) Set(key string, query []byte) error {
	logger := zapLogger.GetZapLogger(g.GetDebug())
	if err := json.Unmarshal(query, &Confs); err != nil {
		zapLogger.Prepare(logger).
			Append(zap.Any("key", key)).
			Append(zap.Any("value", string(query))).
			Development().
			Level(zap.ErrorLevel).
			Commit(err.Error())
		return err
	}

	return nil
}

func (g Config) Get() Config {
	return g
}

// GetService is method for get a service struct with default vaules by config file
func (g *Config) GetService() interface{} {
	service := struct {
		Name string
		GRPC struct {
			Port string
			Host string
		}
	}{
		Name: g.Get().Service.Name,
		GRPC: struct {
			Port string
			Host string
		}{
			Port: g.Get().Service.GRPC.Port,
			Host: g.Get().Service.GRPC.Host,
		},
	}

	return service
}

func (g *Config) GetDebug() bool {
	return g.Get().Debug
}

func (g *Config) SetDebug(debug bool) {
	g.Debug = debug
}

// Load returns configs
func (g *Config) Load(path string) error {

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		if err := g.file(path); err != nil {
			return err
		}
	}

	return nil
}

// file func
func (g *Config) file(path string) error {
	log := zapLogger.GetZapLogger(false)

	// name of config file (without extension)
	// REQUIRED if the config file does not have the extension in the name
	// path to look for the config file in
	viper.SetConfigFile(path)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			zapLogger.Prepare(log).
				Append(zap.Any("error", fmt.Sprintf("Config file not found; ignore error if desired: %s", err))).
				Level(zap.ErrorLevel).
				Development().
				Commit("env")
		} else {
			// Config file was found but another error was produced
			zapLogger.Prepare(log).
				Append(zap.Any("error", fmt.Sprintf("Config file was found but another error was produced: %s", err))).
				Level(zap.ErrorLevel).
				Development().
				Commit("env")
		}
		return err
	}

	if err := viper.Unmarshal(&Confs); err != nil {
		// Config file can not unmarshal to struct
		zapLogger.Prepare(log).
			Append(zap.Any("error", fmt.Sprintf("Config file can not unmarshal to struct: %s", err))).
			Level(zap.ErrorLevel).
			Development().
			Commit("env")

		return err
	}

	return nil
}

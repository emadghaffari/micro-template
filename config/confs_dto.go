package config

import (
	"encoding/json"
	zapLogger "micro/pkg/logger"

	"go.uber.org/zap"
)

// Set method
// you can set new key in switch for manage config with config server
func (g *GlobalConfig) Set(key string, query []byte) error {
	logger := zapLogger.GetZapLogger(Global.Debug())
	if err := json.Unmarshal(query, &Global); err != nil {
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

// GetService is method for get a service struct with default vaules by config file
func (g *GlobalConfig) GetService() interface{} {
	service := struct {
		Name string
		GRPC struct {
			Port string
			Host string
		}
	}{
		Name: Global.Service.Name,
		GRPC: struct {
			Port string
			Host string
		}{
			Port: Global.Service.GRPC.Port,
			Host: Global.Service.GRPC.Host,
		},
	}

	return service
}

func (g *GlobalConfig) Debug() bool {
	return g.Environment != "production"
}

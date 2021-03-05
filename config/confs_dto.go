package config

import (
	"encoding/json"
	"fmt"
	zapLogger "micro/pkg/logger"

	"go.uber.org/zap"
)

func (g *global) Set(key string, query []byte) error {
	logger := zapLogger.GetZapLogger(false)

	switch key {
	case "redis":
		var r redis
		if err := json.Unmarshal(query, &r); err != nil {
			zapLogger.Prepare(logger).
				Append(zap.Any("key", key)).
				Append(zap.Any("value", string(query))).
				Development().
				Level(zap.ErrorLevel).
				Commit(err.Error())
			return err
		}
		Global.Redis = r
	default:
		zapLogger.Prepare(logger).
			Append(zap.Any("key", key)).
			Append(zap.Any("value", string(query))).
			Development().
			Level(zap.ErrorLevel).
			Commit("key not found in service")
		return fmt.Errorf("key not found in configs")
	}

	return nil
}

package config

import (
	"fmt"
	"micro/pkg/logger"
	"os"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// LoadGlobalConfiguration returns configs
func LoadGlobalConfiguration(path string) error {

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		if err := localEnvironment(path); err != nil {
			return err
		}
	}

	// READ FROM PRODUCTION CONFIG MANAGER
	if viper.GetString("environment") == "production" {
	}

	return nil
}

func localEnvironment(path string) error {
	log := logger.GetZapLogger(false)

	// name of config file (without extension)
	// REQUIRED if the config file does not have the extension in the name
	// path to look for the config file in
	viper.SetConfigFile(path)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			logger.Prepare(log).
				Append(zap.Any("error", fmt.Sprintf("Config file not found; ignore error if desired: %s", err))).
				Level(zap.PanicLevel).
				Development().
				Commit("env")
		} else {
			// Config file was found but another error was produced
			logger.Prepare(log).
				Append(zap.Any("error", fmt.Sprintf("Config file was found but another error was produced: %s", err))).
				Level(zap.PanicLevel).
				Development().
				Commit("env")
		}
		return err
	}

	if err := viper.Unmarshal(&Global); err != nil {
		// Config file can not unmarshal to struct
		logger.Prepare(log).
			Append(zap.Any("error", fmt.Sprintf("Config file can not unmarshal to struct: %s", err))).
			Level(zap.PanicLevel).
			Development().
			Commit("env")

		return err
	}

	return nil
}

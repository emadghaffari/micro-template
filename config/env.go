package config

import (
	"encoding/json"
	"fmt"
	"micro/pkg/logger"
	"os"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// Load returns configs
func Load(path string) error {

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		if err := file(path); err != nil {
			return err
		}
	}

	return nil
}

func file(path string) error {
	log := logger.GetZapLogger(false)

	// name of config file (without extension)
	// REQUIRED if the config file does not have the extension in the name
	// path to look for the config file in
	viper.SetConfigFile(path)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			logger.Prepare(log).
				Append(zap.Any("error", fmt.Sprintf("Config file not found; ignore error if desired: %s", err))).
				Level(zap.ErrorLevel).
				Development().
				Commit("env")
		} else {
			// Config file was found but another error was produced
			logger.Prepare(log).
				Append(zap.Any("error", fmt.Sprintf("Config file was found but another error was produced: %s", err))).
				Level(zap.ErrorLevel).
				Development().
				Commit("env")
		}
		return err
	}

	if err := viper.Unmarshal(&Global); err != nil {
		// Config file can not unmarshal to struct
		logger.Prepare(log).
			Append(zap.Any("error", fmt.Sprintf("Config file can not unmarshal to struct: %s", err))).
			Level(zap.ErrorLevel).
			Development().
			Commit("env")

		return err
	}

	// FIXME DELETE the comments
	b, _ := json.Marshal(Global)
	fmt.Println(string(b))

	return nil
}

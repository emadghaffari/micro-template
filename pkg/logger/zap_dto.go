package logger

import (
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

var (
	logPath string
)

// SetLogPath func
func SetLogPath(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, os.ModePerm)
	}
	logPath = path
}

// GetZapLogger func
// create new zap logger
func GetZapLogger(debug bool) *zap.Logger {
	w := zapcore.AddSync(
		&lumberjack.Logger{
			Filename:   fmt.Sprintf("%s/%s.log", logPath, time.Now().Local().Format("2006-01-02")),
			MaxSize:    500, // megabytes
			MaxBackups: 3,
			MaxAge:     28,   //days
			Compress:   true, // disabled by default
		},
	)

	prodEC := zap.NewProductionEncoderConfig()
	prodEC.EncodeTime = zapcore.RFC3339TimeEncoder
	if debug {
		core := zapcore.NewTee(
			zapcore.NewCore(
				zapcore.NewJSONEncoder(prodEC),
				w,
				zap.InfoLevel,
			),
			zapcore.NewCore(
				zapcore.NewConsoleEncoder(prodEC),
				zapcore.AddSync(os.Stdout),
				zapcore.DebugLevel,
			),
		)
		return zap.New(core)
	}

	core := zapcore.NewTee(
		zapcore.NewCore(
			zapcore.NewJSONEncoder(prodEC),
			w,
			zap.InfoLevel,
		),
	)

	return zap.New(core)
}

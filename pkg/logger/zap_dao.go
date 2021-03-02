package logger

import (
	"runtime"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Log struct {
	logger *zap.Logger
	fields []zapcore.Field
	level  zapcore.Level
}

// Add new log
func (l *Log) Add(key string, field interface{}) *Log {
	l.fields = append(l.fields, zap.Any(key, field))
	return l
}

// Append new log
func (l *Log) Append(fields ...zapcore.Field) *Log {
	l.fields = append(l.fields, fields...)
	return l
}

// Commit meth
func (l *Log) Commit(message string) {
	defer func() {
		l.logger.Sync()
		l.fields = nil
	}()

	switch l.level {
	case zapcore.InfoLevel:
		l.logger.Info(message, l.fields...)
	case zapcore.WarnLevel:
		l.logger.Warn(message, l.fields...)
	case zapcore.DebugLevel:
		l.logger.Debug(message, l.fields...)
	case zapcore.FatalLevel:
		l.logger.Fatal(message, l.fields...)
	case zapcore.PanicLevel:
		l.logger.Panic(message, l.fields...)
	default:
		l.logger.Warn(message, l.fields...)
	}
}

// Level of log
func (l *Log) Level(level zapcore.Level) *Log {
	l.level = level
	return l
}

// Development method
func (l *Log) Development() *Log {
	var caller string

	pc, _, line, ok := runtime.Caller(1)
	details := runtime.FuncForPC(pc)
	if ok && details != nil {
		caller = details.Name()
	}
	runInfo := make(map[string]interface{})

	runInfo["line"] = line
	runInfo["caller"] = caller

	l.Add("developer", runInfo)

	return l
}

// Prepare new logger
func Prepare(zap *zap.Logger) *Log {
	return &Log{logger: zap}
}

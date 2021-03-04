package logger

import (
	"runtime"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Log interface
type Log interface {
	Add(key string, field interface{}) *Log
	Append(fields ...zapcore.Field) *Log
	Level(level zapcore.Level) *Log
	Development() *Log
}

type log struct {
	logger *zap.Logger
	fields []zapcore.Field
	level  zapcore.Level
}

// Add new log
func (l *log) Add(key string, field interface{}) *log {
	l.fields = append(l.fields, zap.Any(key, field))
	return l
}

// Append new log
func (l *log) Append(fields ...zapcore.Field) *log {
	l.fields = append(l.fields, fields...)
	return l
}

// Commit meth
func (l *log) Commit(message string) {
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
	default:
		l.logger.Warn(message, l.fields...)
	}
}

// Level of log
func (l log) Level(level zapcore.Level) *log {
	l.level = level
	return &l
}

// Development method
func (l *log) Development() *log {
	var caller string = ""

	pc, _, line, ok := runtime.Caller(1)
	details := runtime.FuncForPC(pc)
	if ok && details != nil {
		caller = details.Name()
	}

	l.Add("line", line)
	l.Add("caller", caller)

	return l
}

// Prepare new logger
func Prepare(zap *zap.Logger) *log {
	return &log{logger: zap}
}

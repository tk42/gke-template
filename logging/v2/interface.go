package logging

import "go.uber.org/zap/zapcore"

type ILogger interface {
	Debug(string, ...zapcore.Field)
	Info(string, ...zapcore.Field)
	Warn(string, ...zapcore.Field)
	Error(string, ...zapcore.Field)
	Panic(string, ...zapcore.Field)
	Fatal(string, ...zapcore.Field)
}

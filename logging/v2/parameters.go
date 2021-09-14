package logging

import "go.uber.org/zap/zapcore"

type LoggingOption func(*LoggingParameters)

type LoggingParameters struct {
	LogLevel        zapcore.Level
	CountThredshold uint
	WindowThreshold int
	InitialFields   map[string]interface{}
}

func LogLevel(level zapcore.Level) LoggingOption {
	return func(op *LoggingParameters) {
		op.LogLevel = level
	}
}

func CountThredshold(thres uint) LoggingOption {
	return func(op *LoggingParameters) {
		op.CountThredshold = thres
	}
}

func WindowThreshold(window int) LoggingOption {
	return func(op *LoggingParameters) {
		op.WindowThreshold = window
	}
}

func InitialFields(fields map[string]interface{}) LoggingOption {
	return func(op *LoggingParameters) {
		op.InitialFields = fields
	}
}

func (p *LoggingParameters) validation() {
	if p.LogLevel < zapcore.DebugLevel || p.LogLevel > zapcore.FatalLevel {
		panic("Unknown LogLevel")
	}
}

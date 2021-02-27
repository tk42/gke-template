package logging

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/tk42/victolinux/env"
	"github.com/tk42/victolinux/throttle"
	stackdriver "github.com/tommy351/zap-stackdriver"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.Logger
	isDebugEnabled bool
	isInfoEnabled  bool
	throttler      throttle.Throttle
}

const (
	TIME_FORMAT = "2006-01-02_15:04:05"
)

func GetLogger(name string, cfg throttle.ThrottleConfig) *Logger {
	setLevel := env.GetString(strings.ToUpper(name)+"_LOGLEVEL", env.GetString("LOGLEVEL", ""))

	level := zap.NewAtomicLevel()
	isDebugEnabled, isInfoEnabled := false, false
	switch setLevel {
	case "DEBUG":
		isDebugEnabled = true
		isInfoEnabled = true
		level.SetLevel(zapcore.DebugLevel)
	case "INFO":
		isDebugEnabled = false
		isInfoEnabled = true
		level.SetLevel(zapcore.InfoLevel)
	case "WARN":
		level.SetLevel(zapcore.WarnLevel)
	default:
		err := fmt.Errorf("Invalid LOGLEVEL(%v)", setLevel)
		panic(err)
	}

	// Create ProcessName with CREATE_TIME
	createTime, ok := os.LookupEnv("CREATE_TIME")
	if !ok {
		createTime = time.Now().Format(TIME_FORMAT)
		os.Setenv("CREATE_TIME", createTime)
	}

	initialFields := map[string]interface{}{
		"ProcessName": env.GetString("PROCESS_NAME", "") + "_" + createTime,
	}

	myConfig := zap.Config{
		Level:            level,
		Encoding:         "json",
		EncoderConfig:    stackdriver.EncoderConfig,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		InitialFields:    initialFields,
	}
	zaplogger, err := myConfig.Build()
	if err != nil {
		panic(err)
	}
	logger := &Logger{
		zaplogger.Named(name), isDebugEnabled, isInfoEnabled, throttle.NewThrottle(cfg),
	}

	logger.Info("Successfully created Config and Logger",
		zap.String("name", name), zap.String("loglevel", setLevel),
		zap.Bool("IsDebugEnabled", isDebugEnabled), zap.Bool("IsInfoEnabled", isInfoEnabled),
	)
	return logger
}

func (l *Logger) Debug(msg string, fields ...zapcore.Field) {
	if !l.isDebugEnabled {
		return
	}
	l.throttler.Trigger()
	l.Debug(msg, fields...)
}

func (l *Logger) Info(msg string, fields ...zapcore.Field) {
	if !l.isInfoEnabled {
		return
	}
	l.throttler.Trigger()
	l.Info(msg, fields...)
}

func (l *Logger) Warn(msg string, fields ...zapcore.Field) {
	l.throttler.Trigger()
	l.Warn(msg, fields...)
}

func (l *Logger) Panic(msg string, fields ...zapcore.Field) {
	l.throttler.Trigger()
	l.Panic(msg, fields...)
}

func (l *Logger) Fatal(msg string, fields ...zapcore.Field) {
	l.throttler.Trigger()
	l.Fatal(msg, fields...)
}

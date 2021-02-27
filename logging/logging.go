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
	throttleConfig throttle.ThrottleConfig
	throttler      throttle.Throttle
}

const (
	TIME_FORMAT = "2006-01-02_15:04:05"
)

func GetLogger(name string) *Logger {
	setLevel := env.GetString(strings.ToUpper(name)+"_LOGLEVEL", env.GetString("LOGLEVEL", ""))
	setCount := uint(env.GetInt(strings.ToUpper(name)+"_COUNT_THRESHOLD", env.GetInt("COUNT_THRESHOLD", 30)))
	setWindow := env.GetInt(strings.ToUpper(name)+"_WINDOW_MSEC_THRESHOLD", env.GetInt("WINDOW_MSEC_THRESHOLD", 1000))

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

	defaultPolicy := throttle.ThrottleParameter(
		setCount, time.Duration(setWindow)*time.Millisecond,
		throttle.Reached(
			func() {
				panic(fmt.Sprintf("DETECTED THROTTLE CHECK: %v COUNT WITHIN %v MSEC", setCount, setWindow))
			},
		),
	)
	logger := &Logger{
		zaplogger.Named(name), isDebugEnabled, isInfoEnabled, defaultPolicy, throttle.NewThrottle(defaultPolicy),
	}

	logger.Info("Successfully created Config and Logger",
		zap.String("name", name), zap.String("loglevel", setLevel),
		zap.Bool("IsDebugEnabled", isDebugEnabled), zap.Bool("IsInfoEnabled", isInfoEnabled),
		zap.Uint("ThrottleCount", setCount), zap.Int("ThrottleWindow", setWindow),
	)
	return logger
}

func (l *Logger) GetThrottleConfig() throttle.ThrottleConfig {
	return l.throttleConfig
}

func (l *Logger) SetThrottleConfig(cfg throttle.ThrottleConfig) {
	l.throttleConfig = cfg
	l.throttler = throttle.NewThrottle(cfg)
}

func (l *Logger) Debug(msg string, fields ...zapcore.Field) {
	if !l.isDebugEnabled {
		return
	}
	l.Logger.Debug(msg, fields...)
}

func (l *Logger) Info(msg string, fields ...zapcore.Field) {
	if !l.isInfoEnabled {
		return
	}
	l.Logger.Info(msg, fields...)
}

func (l *Logger) Warn(msg string, fields ...zapcore.Field) {
	l.throttler.Trigger()
	if l.throttler.IsFreeze() {
		return
	}
	l.Logger.Warn(msg, fields...)
}

func (l *Logger) Error(msg string, fields ...zapcore.Field) {
	l.throttler.Trigger()
	if l.throttler.IsFreeze() {
		return
	}
	l.Logger.Error(msg, fields...)
}

func (l *Logger) Panic(msg string, fields ...zapcore.Field) {
	l.throttler.Trigger()
	if l.throttler.IsFreeze() {
		return
	}
	l.Logger.Panic(msg, fields...)
}

func (l *Logger) Fatal(msg string, fields ...zapcore.Field) {
	l.throttler.Trigger()
	if l.throttler.IsFreeze() {
		return
	}
	l.Logger.Fatal(msg, fields...)
}

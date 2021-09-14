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
	TIME_FORMAT = "20060102_150405"
)

func defaultLoggingParameters(name string) LoggingParameters {
	setLevel := strings.ToUpper(env.GetString(strings.ToUpper(name)+"_LOGLEVEL", env.GetString("LOGLEVEL", "")))
	setCount := uint(env.GetInt(strings.ToUpper(name)+"_COUNT_THRESHOLD", env.GetInt("COUNT_THRESHOLD", 30)))
	setWindow := env.GetInt(strings.ToUpper(name)+"_WINDOW_MSEC_THRESHOLD", env.GetInt("WINDOW_MSEC_THRESHOLD", 1000))

	var level zapcore.Level
	switch setLevel {
	case "DEBUG":
		level = zapcore.DebugLevel
	case "INFO":
		level = zapcore.InfoLevel
	case "WARN":
		level = zapcore.WarnLevel
	}

	return LoggingParameters{
		LogLevel:        level,
		CountThredshold: setCount,
		WindowThreshold: setWindow,
	}
}

func GetLogger(name string, ops ...LoggingOption) *Logger {
	params := defaultLoggingParameters(name)
	for _, opt := range ops {
		opt(&params)
	}

	params.validation()

	isDebugEnabled, isInfoEnabled := false, false
	switch params.LogLevel {
	case zapcore.DebugLevel:
		isDebugEnabled = true
		isInfoEnabled = true
	case zapcore.InfoLevel:
		isInfoEnabled = true
	case zapcore.WarnLevel:
		// pass
	default:
		err := fmt.Errorf("Invalid LOGLEVEL(%v)", params.LogLevel.String())
		panic(err)
	}

	// Create ProcessName with CREATE_TIME
	createTime, ok := os.LookupEnv("CREATE_TIME")
	if !ok {
		createTime = time.Now().Format(TIME_FORMAT)
		os.Setenv("CREATE_TIME", createTime)
	}

	if processName, ok := os.LookupEnv("PROCESS_NAME"); ok {
		createTime = processName + "_" + createTime
	}

	initialFields := params.InitialFields
	initialFields["ProcessName"] = createTime

	myConfig := zap.Config{
		Level:            zap.NewAtomicLevelAt(params.LogLevel),
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
		params.CountThredshold, time.Duration(params.WindowThreshold)*time.Millisecond,
		throttle.Reached(
			func() {
				panic(fmt.Sprintf("DETECTED THROTTLE CHECK: %v COUNT WITHIN %v MSEC", params.CountThredshold, params.WindowThreshold))
			},
		),
	)
	logger := &Logger{
		zaplogger.Named(name), isDebugEnabled, isInfoEnabled, defaultPolicy, throttle.NewThrottle(defaultPolicy),
	}

	logger.Info("Successfully created Config and Logger",
		zap.String("name", name), zap.String("loglevel", params.LogLevel.String()),
		zap.Bool("IsDebugEnabled", isDebugEnabled), zap.Bool("IsInfoEnabled", isInfoEnabled),
		zap.Uint("ThrottleCount", params.CountThredshold), zap.Int("ThrottleWindow", params.WindowThreshold),
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

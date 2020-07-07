package logging

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jimako1989/gke-template/env"
	stackdriver "github.com/tommy351/zap-stackdriver"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	TIME_FORMAT = "2006-01-02_15:04:05"
)

var (
	IsDebugEnabled = false
	IsInfoEnabled  = false
)

func GetLogger(name string) zap.Logger {
	setLevel := env.GetString(strings.ToUpper(name)+"_LOGLEVEL", env.GetString("LOGLEVEL", ""))

	level := zap.NewAtomicLevel()
	switch setLevel {
	case "DEBUG":
		IsDebugEnabled = true
		IsInfoEnabled = true
		level.SetLevel(zapcore.DebugLevel)
	case "INFO":
		IsDebugEnabled = false
		IsInfoEnabled = true
		level.SetLevel(zapcore.InfoLevel)
	case "WARN":
		level.SetLevel(zapcore.WarnLevel)
	default:
		err := fmt.Errorf("Invalid LOGLEVEL(%v)", setLevel)
		panic(err)
	}

	myConfig := zap.Config{
		Level:            level,
		Encoding:         "json",
		EncoderConfig:    stackdriver.EncoderConfig,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}
	logger, _ := myConfig.Build()
	logger = logger.Named(name)

	// Create ProcessName with CREATE_TIME
	createTime, ok := os.LookupEnv("CREATE_TIME")
	if !ok {
		createTime = time.Now().Format(TIME_FORMAT)
		os.Setenv("CREATE_TIME", createTime)
	}
	// ex) feedlogic1_2019-10-23_23:22:10
	logger = logger.With(zap.String("ProcessName", env.GetString("PROCESS_NAME", "")+"_"+createTime))

	logger.Info("Successfully created Config and Logger",
		zap.String("name", name), zap.String("loglevel", setLevel),
		zap.Bool("IsDebugEnabled", IsDebugEnabled), zap.Bool("IsInfoEnabled", IsInfoEnabled),
	)
	return *logger
}

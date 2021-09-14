package logging

import (
	"os"
	"testing"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestLogger(t *testing.T) {
	logger := GetLogger("Logger", LogLevel(zapcore.InfoLevel), InitialFields(map[string]interface{}{
		"my_name": "test",
	}))
	logger.Debug("debug", zap.String("Key", "String"), zap.Ints("ints", []int{10, 20}))
	logger.Info("info", zap.String("Key", "String"), zap.Ints("ints", []int{10, 20}))
	logger.Warn("warn", zap.String("Key", "String"), zap.Ints("ints", []int{10, 20}))
	logger.Error("error", zap.String("Key", "String"), zap.Ints("ints", []int{10, 20}))
}

func CheckThrottleLogger(period time.Duration, count int) {
	os.Setenv("LOGLEVEL", "DEBUG")
	logger := GetLogger("Logger")
	go func() {
		var i int
		for i = 0; i < count; i++ {
			logger.Error("error", zap.String("Key", "String"), zap.Ints("ints", []int{10, 20}))
			time.Sleep(period / time.Duration(count))
		}
	}()

	time.Sleep(period)
}

func TestErrorThrottleLogger(t *testing.T) {
	CheckThrottleLogger(time.Second, 29) // ok

	CheckThrottleLogger(time.Second, 30) // panic: DETECTED THROTTLE CHECK: 30 COUNT WITHIN 1000 MSEC
}

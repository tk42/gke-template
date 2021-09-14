package logging

import (
	"testing"

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

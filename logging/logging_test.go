package logging

import (
	"os"
	"testing"

	"go.uber.org/zap"
)

func TestLogger(t *testing.T) {
	os.Setenv("LOGLEVEL", "DEBUG")
	logger := GetLogger("Logger")
	logger.Debug("debug", zap.String("Key", "String"), zap.Ints("ints", []int{10, 20}))
	logger.Info("info", zap.String("Key", "String"), zap.Ints("ints", []int{10, 20}))
	logger.Warn("warn", zap.String("Key", "String"), zap.Ints("ints", []int{10, 20}))
	logger.Error("error", zap.String("Key", "String"), zap.Ints("ints", []int{10, 20}))
	logger.Fatal("fatal")
}

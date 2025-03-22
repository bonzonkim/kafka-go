package logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// NewLogger initializes new Logger with log rotation setting
func NewLogger(processID string, rotationSize int, rotationCount int) *zap.Logger {
	// config for zapcore
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename: fmt.Sprintf("./logs/%s.log", processID),
		MaxAge:	rotationCount,
		MaxSize: rotationSize,
	})

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		w,
		zapcore.InfoLevel,
	)

	return zap.New(core, zap.AddCaller())
}

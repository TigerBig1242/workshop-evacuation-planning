package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Log *zap.Logger

// InitLogger initializes the logger with different configurations for console and file output.
// It uses lumberjack for file rotation and zap for structured logging.
func InitLogger() error {

	stdOut := zapcore.AddSync(os.Stdout)

	// Configure lumberjack for log file rotation
	file := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "logs/evacuation.log",
		MaxSize:    10,
		MaxBackups: 3,
		MaxAge:     7,
	})

	// Set the log level to InfoLevel
	level := zap.NewAtomicLevelAt(zap.InfoLevel)

	// Configure the encoder for file output (JSON)
	productionCfg := zap.NewProductionEncoderConfig()
	productionCfg.TimeKey = "timestamp"
	productionCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	// Configure the encoder for console output (text)
	// This is used for development and debugging
	developmentCfg := zap.NewDevelopmentEncoderConfig()
	developmentCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

	// Create the console encoders (text format)
	consoleEncoder := zapcore.NewConsoleEncoder(developmentCfg)
	// Create the file encoder (JSON format)
	fileEncoder := zapcore.NewJSONEncoder(productionCfg)

	// Create the core for the logger output (console and file) used for both by NewTee
	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, stdOut, level),
		zapcore.NewCore(fileEncoder, file, level),
	)

	Log = zap.New(core)
	return nil
}

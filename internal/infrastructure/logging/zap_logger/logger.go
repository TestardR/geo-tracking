package zap_logger

import (
	"syscall"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const JsonRepack = "json-repack"

var contextNamespace = zap.Namespace("context")

func NewLogger(logLevel, outputPath, appVersion, appName, env string) (*zap.Logger, error) {
	output := []string{outputPath}
	if outputPath == "" {
		output = []string{"stderr"}
	}

	initialFields := make(map[string]interface{}, 0)
	encoding := JsonRepack
	disableStackTrace := false

	if env != "dev" {
		initialFields = map[string]interface{}{
			"appName": appName,
			"env":     env,
			"version": appVersion,
			"context": "",
		}
	} else {
		disableStackTrace = true
		encoding = "console"
	}

	zcfg := zap.Config{
		Level:             parseLevel(logLevel),
		Development:       false,
		DisableCaller:     false,
		DisableStacktrace: disableStackTrace,
		Encoding:          encoding,
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:    "brut",
			StacktraceKey: "stacktrace",
			LevelKey:      "level",
			EncodeLevel:   zapcore.CapitalLevelEncoder,
			TimeKey:       "@timestamp",
			EncodeTime:    zapcore.ISO8601TimeEncoder,
			CallerKey:     "context.caller",
			EncodeCaller:  zapcore.ShortCallerEncoder,
		},
		OutputPaths:   output,
		InitialFields: initialFields, // Complementary fields that are expected by LOM format
	}

	logger, err := zcfg.Build()
	if err != nil {
		return nil, err
	}

	return logger.With(contextNamespace), nil
}

func parseLevel(level string) zap.AtomicLevel {
	atomicLevel := zap.NewAtomicLevel()
	err := atomicLevel.UnmarshalText([]byte(level))
	if err != nil {
		syscall.Exit(1)
	}

	return atomicLevel
}

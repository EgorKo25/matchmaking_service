package logger

import (
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// SetupProductionLogger function for production env with JSON formatting
func SetupProductionLogger() (*zap.Logger, error) {
	logDir := "logs/production"
	if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
		return nil, err
	}

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	errorFile, err := os.Create(filepath.Join(logDir, "error.log"))
	if err != nil {
		return nil, err
	}
	infoFile, err := os.Create(filepath.Join(logDir, "info.log"))
	if err != nil {
		return nil, err
	}

	errorWS := zapcore.AddSync(errorFile)
	infoWS := zapcore.AddSync(infoFile)

	errorLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zap.ErrorLevel
	})

	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl <= zap.InfoLevel
	})

	consoleDebugging := zapcore.Lock(os.Stderr)

	encoder := zapcore.NewJSONEncoder(encoderCfg)

	core := zapcore.NewTee(
		zapcore.NewCore(encoder, consoleDebugging, zapcore.InfoLevel),
		zapcore.NewCore(encoder, errorWS, errorLevel),
		zapcore.NewCore(encoder, infoWS, infoLevel),
	)

	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	return logger, nil
}

// SetupLocalLogger function for production env with console(humanized) formatting
func SetupLocalLogger() (*zap.Logger, error) {
	logDir := "logs/local"
	err := os.MkdirAll(logDir, os.ModePerm)
	if err != nil {
		return nil, err
	}

	encoderCfg := zap.NewDevelopmentEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	debugFile, err := os.Create(filepath.Join(logDir, "debug.log"))
	if err != nil {
		return nil, err
	}
	errorFile, err := os.Create(filepath.Join(logDir, "error.log"))
	if err != nil {
		return nil, err
	}
	infoFile, err := os.Create(filepath.Join(logDir, "info.log"))
	if err != nil {
		return nil, err
	}

	debugWS := zapcore.AddSync(debugFile)
	errorWS := zapcore.AddSync(errorFile)
	infoWS := zapcore.AddSync(infoFile)

	debugLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zap.DebugLevel
	})

	errorLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zap.ErrorLevel
	})

	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl <= zap.InfoLevel
	})

	consoleDebugging := zapcore.Lock(os.Stdout)

	encoder := zapcore.NewConsoleEncoder(encoderCfg)

	core := zapcore.NewTee(
		zapcore.NewCore(encoder, consoleDebugging, zapcore.DebugLevel),
		zapcore.NewCore(encoder, debugWS, debugLevel),
		zapcore.NewCore(encoder, errorWS, errorLevel),
		zapcore.NewCore(encoder, infoWS, infoLevel),
	)

	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	return logger, nil
}

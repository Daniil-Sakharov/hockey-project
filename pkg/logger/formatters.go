package logger

import (
	"go.uber.org/zap/zapcore"
)

// EnvironmentFormatter создает форматер в зависимости от окружения
func EnvironmentFormatter(environment string, asJSON bool) zapcore.Encoder {
	switch environment {
	case "development", "dev":
		return createDevelopmentFormatter()
	case "staging":
		return createStagingFormatter()
	case "production", "prod":
		return createProductionFormatter()
	default:
		if asJSON {
			return zapcore.NewJSONEncoder(buildProductionEncoderConfig())
		}
		return zapcore.NewConsoleEncoder(buildProductionEncoderConfig())
	}
}

// createDevelopmentFormatter создает форматер для development с цветами
func createDevelopmentFormatter() zapcore.Encoder {
	config := zapcore.EncoderConfig{
		MessageKey:    "message",
		LevelKey:      "level",
		TimeKey:       "timestamp",
		CallerKey:     "caller",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.CapitalColorLevelEncoder, // Цвета для уровней
		EncodeTime:    zapcore.TimeEncoderOfLayout("15:04:05.000"),
		EncodeCaller:  zapcore.ShortCallerEncoder,
		EncodeName:    zapcore.FullNameEncoder,
	}
	return zapcore.NewConsoleEncoder(config)
}

// createStagingFormatter создает hybrid форматер для staging
func createStagingFormatter() zapcore.Encoder {
	config := zapcore.EncoderConfig{
		MessageKey:    "message",
		LevelKey:      "level",
		TimeKey:       "timestamp",
		CallerKey:     "caller",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.CapitalLevelEncoder,
		EncodeTime:    zapcore.ISO8601TimeEncoder,
		EncodeCaller:  zapcore.ShortCallerEncoder,
		EncodeName:    zapcore.FullNameEncoder,
	}
	return zapcore.NewJSONEncoder(config) // JSON для staging
}

// createProductionFormatter создает форматер для production
func createProductionFormatter() zapcore.Encoder {
	return zapcore.NewJSONEncoder(buildProductionEncoderConfig())
}

package yetzap

import (
	"strings"

	"github.com/pvormste/yetenv"
	"github.com/pvormste/yetlog"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ConfigureSugaredFunc defines a function which can be used to configure a sugared logger. See function NewCustomSugaredLogger().
type ConfigureSugaredFunc func() (*zap.SugaredLogger, error)

// SugaredLogger is the wrapper for the actual sugared logger.
type SugaredLogger struct {
	zapLogger *zap.SugaredLogger
}

// NewDefaultSugaredLogger creates a new sugared logger with some default configurations for different environments.
func NewDefaultSugaredLogger(environment yetenv.Environment, rawMinLevel string) (yetlog.Logger, error) {
	return NewCustomSugaredLogger(func() (*zap.SugaredLogger, error) {
		minLevel := zapcore.InfoLevel
		if err := minLevel.Set(strings.ToLower(rawMinLevel)); err != nil {
			return nil, err
		}

		var loggerConf zap.Config

		switch environment {
		case yetenv.Production:
			loggerConf = DefaultProductionConfig(minLevel)
		default:
			loggerConf = DefaultDevelopmentConfig(minLevel)
		}

		loggerConf.DisableStacktrace = true
		logger, err := loggerConf.Build(zap.AddCallerSkip(1))
		if err != nil {
			return nil, err
		}

		return logger.Sugar(), nil
	})
}

// NewCustomSugaredLogger can be used to create a custom sugared logger by providing a ConfigureSugaredFunc function.
func NewCustomSugaredLogger(zapConfigureFunc ConfigureSugaredFunc) (yetlog.Logger, error) {
	zapSugaredLogger, err := zapConfigureFunc()
	if err != nil {
		return nil, err
	}

	return SugaredLogger{
		zapLogger: zapSugaredLogger,
	}, nil
}

// WrapSugaredLogger wraps an existent sugared logger without needing to touch any configuration.
func WrapSugaredLogger(sugaredLogger *zap.SugaredLogger) yetlog.Logger {
	return SugaredLogger{
		zapLogger: sugaredLogger,
	}
}

// DefaultProductionConfig returns the default production config which is used to create a default sugared logger.
func DefaultProductionConfig(minLevel zapcore.Level) zap.Config {
	loggerConf := zap.NewProductionConfig()
	loggerConf.Level.SetLevel(minLevel)

	return loggerConf
}

// DefaultDevelopmentConfig returns the default development config which is used to create a default sugared logger.
func DefaultDevelopmentConfig(minLevel zapcore.Level) zap.Config {
	loggerConf := zap.NewDevelopmentConfig()
	loggerConf.Level.SetLevel(minLevel)
	loggerConf.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	return loggerConf
}

// Reconfigure is currently not implemented and logs a warning.
func (s SugaredLogger) Reconfigure(_ interface{}) {
	s.Warn("reconfigure is not implemented", "logger", "zap")
}

// NewNamedLogger creates a new named logger.
func (s SugaredLogger) NewNamedLogger(name string) yetlog.Logger {
	namedLogger := s.zapLogger.Named(name)
	return SugaredLogger{
		zapLogger: namedLogger,
	}
}

// Debug logs a debug message with parameters.
func (s SugaredLogger) Debug(message string, fields ...interface{}) {
	s.zapLogger.Debugw(message, fields...)
}

// Info logs a info message with parameters.
func (s SugaredLogger) Info(message string, fields ...interface{}) {
	s.zapLogger.Infow(message, fields...)
}

// Warn logs a warning message with parameters.
func (s SugaredLogger) Warn(message string, fields ...interface{}) {
	s.zapLogger.Warnw(message, fields...)
}

// Error logs a error message with paramters.
func (s SugaredLogger) Error(message string, fields ...interface{}) {
	s.zapLogger.Errorw(message, fields...)
}

// Fatal logs a fatal message with paramters.
func (s SugaredLogger) Fatal(message string, fields ...interface{}) {
	s.zapLogger.Fatalw(message, fields...)
}

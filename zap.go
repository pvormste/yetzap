package yetzap

import (
	"strings"

	"github.com/pvormste/yetwebutils/yetenv"
	"github.com/pvormste/yetwebutils/yetlog"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ConfigureSugaredFunc func() (*zap.SugaredLogger, error)

type SugaredLogger struct {
	zapLogger *zap.SugaredLogger
}

func NewDefaultSugaredLogger(environment yetenv.Environment, rawMinLevel string) (yetlog.Logger, error) {
	return NewCustomSugaredLogger(func() (*zap.SugaredLogger, error) {
		minLevel := zapcore.InfoLevel
		if err := minLevel.Set(strings.ToLower(rawMinLevel)); err != nil {
			return nil, err
		}

		var logger *zap.Logger
		var err error

		switch environment {
		case yetenv.Production:
			loggerConf := DefaultProductionConfig(minLevel)
			logger, err = loggerConf.Build()
		default:
			loggerConf := DefaultDevelopmentConfig(minLevel)
			logger, err = loggerConf.Build()
		}

		if err != nil {
			return nil, err
		}

		return logger.Sugar(), nil
	})
}

func NewCustomSugaredLogger(zapConfigureFunc ConfigureSugaredFunc) (yetlog.Logger, error) {
	zapSugaredLogger, err := zapConfigureFunc()
	if err != nil {
		return nil, err
	}

	return SugaredLogger{
		zapLogger: zapSugaredLogger,
	}, nil
}

func DefaultProductionConfig(minLevel zapcore.Level) zap.Config {
	loggerConf := zap.NewProductionConfig()
	loggerConf.Level.SetLevel(minLevel)

	return loggerConf
}

func DefaultDevelopmentConfig(minLevel zapcore.Level) zap.Config {
	loggerConf := zap.NewDevelopmentConfig()
	loggerConf.Level.SetLevel(minLevel)
	loggerConf.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	return loggerConf
}

func (s SugaredLogger) Reconfigure(_ interface{}) {
	s.Warn("reconfigure is not implemented", "logger", "zap")
}

func (s SugaredLogger) NewNamedLogger(name string) yetlog.Logger {
	namedLogger := s.zapLogger.Named(name)
	return SugaredLogger{
		zapLogger: namedLogger,
	}
}

func (s SugaredLogger) Debug(message string, fields ...interface{}) {
	s.zapLogger.Debugw(message, fields...)
}

func (s SugaredLogger) Info(message string, fields ...interface{}) {
	s.zapLogger.Infow(message, fields...)
}

func (s SugaredLogger) Warn(message string, fields ...interface{}) {
	s.zapLogger.Warnw(message, fields...)
}

func (s SugaredLogger) Error(message string, fields ...interface{}) {
	s.zapLogger.Errorw(message, fields...)
}

func (s SugaredLogger) Fatal(message string, fields ...interface{}) {
	s.zapLogger.Fatalw(message, fields...)
}

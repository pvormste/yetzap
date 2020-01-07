package yetzap

import (
	"github.com/pvormste/yetwebutils/yetenv"
	"github.com/pvormste/yetwebutils/yetlog"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ConfigureSugaredFunc func(*zap.Logger) (*zap.SugaredLogger, error)

type SugaredLogger struct {
	zapLogger *zap.SugaredLogger
}

func NewDefaultSugaredLogger(environment yetenv.Environment, rawMinLevel string) (yetlog.Logger, error) {
	return NewCustomSugaredLogger(func(logger *zap.Logger) (sugaredLogger *zap.SugaredLogger, err error) {
		minLevel := zapcore.InfoLevel
		if err := minLevel.Set(rawMinLevel); err != nil {
			return nil, err
		}

		switch environment {
		case yetenv.Production:
			loggerConf := zap.NewProductionConfig()
			loggerConf.Level.SetLevel(minLevel)
			logger, err = loggerConf.Build()
		default:
			loggerConf := zap.NewDevelopmentConfig()
			loggerConf.Level.SetLevel(minLevel)
			loggerConf.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
			logger, err = loggerConf.Build()
		}

		if err != nil {
			return nil, err
		}

		return logger.Sugar(), nil
	})
}

func NewCustomSugaredLogger(zapConfigureFunc ConfigureSugaredFunc) (yetlog.Logger, error) {
	zapSugaredLogger, err := zapConfigureFunc(&zap.Logger{})
	if err != nil {
		return nil, err
	}

	return SugaredLogger{
		zapLogger: zapSugaredLogger,
	}, nil
}

func (s SugaredLogger) Reconfigure(options interface{}) {
	panic("implement me")
}

func (s SugaredLogger) NewNamedLogger(name string) yetlog.Logger {
	panic("implement me")
}

func (s SugaredLogger) Debug(message string, fields ...interface{}) {
	panic("implement me")
}

func (s SugaredLogger) Info(message string, fields ...interface{}) {
	panic("implement me")
}

func (s SugaredLogger) Warn(message string, fields ...interface{}) {
	panic("implement me")
}

func (s SugaredLogger) Error(message string, fields ...interface{}) {
	panic("implement me")
}

func (s SugaredLogger) Fatal(message string, fields ...interface{}) {
	panic("implement me")
}

[![GoDoc](https://godoc.org/github.com/pvormste/yetzap?status.svg)](https://godoc.org/github.com/pvormste/yetzap) ![](https://github.com/pvormste/yetzap/workflows/lint/badge.svg?branch=master)

# yetzap

`yetzap` is a wrapper package for [uber's zap logger](https://github.com/uber-go/zap) which implements the [yetlog interface](https://github.com/pvormste/yetlog). 
It only supports a small subset of the zap logger but it should be enough for most cases.

## Usage

```go
env := yetenv.Develop
zaplogger, err := yetzap.NewDefaultSugaredLogger(env, "info")

if err != nil {
    // handle error
}

zaplogger.Info("started server", "port", 8080)
```

## Custom zap logger instance

Provide a `ConfigureSugaredFunc` to the `NewCustomSugaredLogger()` function.

Example:
```go
func MyLoggerConstructor(rawMinLevel string) (yetlog.Logger, error) {
	return NewCustomSugaredLogger(func() (*zap.SugaredLogger, error) {
		minLevel := zapcore.InfoLevel
		if err := minLevel.Set(strings.ToLower(rawMinLevel)); err != nil {
			return nil, err
		}

        	loggerConf := DefaultDevelopmentConfig(minLevel)
        	logger, err := loggerConf.Build()

		if err != nil {
			return nil, err
		}

		return logger.Sugar(), nil
	})
}
```

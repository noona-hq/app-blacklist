package logger

import (
	"github.com/noona-hq/blacklist/config"
	"go.uber.org/zap"
)

type Logger struct {
	*zap.SugaredLogger
}

func FromConfig(cfg Config) (*Logger, error) {
	c := zap.NewProductionConfig()
	if !cfg.Structured {
		c.Encoding = "console"
		c.EncoderConfig = zap.NewDevelopmentEncoderConfig()
	}

	_ = c.Level.UnmarshalText([]byte(cfg.Level))
	c.DisableStacktrace = !cfg.Stacktrace

	skip := zap.AddCallerSkip(1)

	l, err := c.Build(skip, zap.WithCaller(false))
	if err != nil {
		return nil, err
	}

	return &Logger{
		SugaredLogger: l.Sugar(),
	}, nil
}

func New() *Logger {
	var cfg = new(Config)
	_ = config.Process(cfg)
	logger, _ := FromConfig(*cfg)
	return logger
}

func NoOp() *Logger {
	return &Logger{
		SugaredLogger: zap.NewNop().Sugar(),
	}
}

func NewWithContext(context ...interface{}) *Logger {
	return New().With(context)
}

func (l *Logger) With(context ...interface{}) *Logger {
	return &Logger{
		SugaredLogger: l.SugaredLogger.With(context...),
	}
}

func (l *Logger) Info(args ...interface{}) {
	l.SugaredLogger.Info(args...)
}

func (l *Logger) Infof(template string, args ...interface{}) {
	l.SugaredLogger.Infof(template, args...)
}

func (l *Logger) Infow(msg string, keysAndValues ...interface{}) {
	l.SugaredLogger.Infow(msg, keysAndValues...)
}

func (l *Logger) Error(args ...interface{}) {
	l.SugaredLogger.Error(args...)
}

func (l *Logger) Errorf(template string, args ...interface{}) {
	l.SugaredLogger.Errorf(template, args...)
}

func (l *Logger) Errorw(msg string, keysAndValues ...interface{}) {

	l.SugaredLogger.Errorw(msg, keysAndValues...)
}

func (l *Logger) Warn(args ...interface{}) {
	l.SugaredLogger.Warn(args...)
}

func (l *Logger) Warnf(template string, args ...interface{}) {
	l.SugaredLogger.Warnf(template, args...)
}

func (l *Logger) Warnw(msg string, keysAndValues ...interface{}) {
	l.SugaredLogger.Warnw(msg, keysAndValues...)
}

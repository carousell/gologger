package logger

import (
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var NoLogFunc = func(_ string, _ ...interface{}) {}

type LogFunc func(template string, args ...interface{})

// Logger wrapper the base zap SugaredLogger for some extra features
type Logger struct {
	prefix string
	logger *zap.SugaredLogger
	core   zapcore.Core
}

// MustNew is like New but panics if the logger can not be built
func MustNew() *Logger {
	return MustNamed("")
}

// MustNamed is like Named but panics if the logger can not be built
func MustNamed(name string) *Logger {
	log, err := Named(name)
	if err != nil {
		panic(err)
	}

	return log
}

// New register a wrapped zap logger with no name
func New() (*Logger, error) {
	return Named("")
}

// Named build a simple wrapped zap logger
func Named(name string) (*Logger, error) {
	logger, err := NewZap(name)
	if err != nil {
		return nil, err
	}

	log := wrap(logger, "").Skip(1)
	return log, nil
}

// Wrap wrap zap logger with pkg logger
func Wrap(logger *zap.SugaredLogger) *Logger {
	log := wrap(logger, "")
	return log.Skip(1)
}

func wrap(logger *zap.SugaredLogger, prefix string) *Logger {
	if prefix != "" {
		prefix = strings.TrimSpace(prefix)
		prefix += " "
	}
	log := &Logger{
		prefix: prefix,
		logger: logger,
		core:   logger.Desugar().Core(),
	}

	return log
}

// Unwrap returns core logger instance, so you can use it directly.
// you can also re-wrap this logger with Wrap function.
func (c *Logger) Unwrap() *zap.SugaredLogger {
	return c.Skip(-1).logger
}

// Skip return new instances that increases the number of callers skipped by caller annotation
func (c *Logger) Skip(skip int) *Logger {
	logger := c.logger.Desugar().WithOptions(zap.AddCallerSkip(skip)).Sugar()
	return wrap(logger, c.prefix)
}

// AddHook add new hook to the core logger
func (c *Logger) AddHook(fn func(e zapcore.Entry) error) *Logger {
	logger := c.logger.Desugar().WithOptions(zap.Hooks(fn)).Sugar()
	return wrap(logger, c.prefix)
}

// With uses fmt.Sprint to construct new logger with prefix
func (c *Logger) With(args ...interface{}) *Logger {
	return wrap(c.logger, c.makeMsg("", args))
}

// Withf makes new instance and uses fmt.Sprintf to make prefix for later logs
func (c *Logger) Withf(template string, args ...interface{}) *Logger {
	return wrap(c.logger, c.makeMsg(template, args))
}

// Withw makes new instance included additional context keys and values for later logs (prefer for json logs).
func (c *Logger) Withw(keysAndValues ...interface{}) *Logger {
	return wrap(c.logger.With(keysAndValues), c.prefix)
}

// Sync flushes any buffered log entries.
func (c *Logger) Sync() error {
	return c.logger.Sync()
}

// Debug uses fmt.Sprint to construct and log a message with prefix.
func (c *Logger) Debug(args ...interface{}) {
	if !c.core.Enabled(zap.DebugLevel) {
		return
	}
	c.logger.Debugf(c.makeMsg("", args))
}

// Debugf uses fmt.Sprintf to log a templated message with prefix.
func (c *Logger) Debugf(template string, args ...interface{}) {
	if !c.core.Enabled(zap.DebugLevel) {
		return
	}
	c.logger.Debugf(c.makeMsg(template, args))
}

// Debugw logs a message with some additional context with keys and values (prefer for json logs).
func (c *Logger) Debugw(msg string, keysAndValues ...interface{}) {
	if !c.core.Enabled(zap.DebugLevel) {
		return
	}
	c.logger.Debugw(c.makeMsg(msg, nil), keysAndValues...)
}

// Info uses fmt.Sprint to construct and log a message with prefix.
func (c *Logger) Info(args ...interface{}) {
	if !c.core.Enabled(zap.InfoLevel) {
		return
	}
	c.logger.Infof(c.makeMsg("", args))
}

// Infof uses fmt.Sprintf to log a templated message with prefix.
func (c *Logger) Infof(template string, args ...interface{}) {
	if !c.core.Enabled(zap.InfoLevel) {
		return
	}
	c.logger.Infof(c.makeMsg(template, args))
}

// Infow logs a message with some additional context with keys and values (prefer for json logs).
func (c *Logger) Infow(msg string, keysAndValues ...interface{}) {
	if !c.core.Enabled(zap.InfoLevel) {
		return
	}
	c.logger.Infow(c.makeMsg(msg, nil), keysAndValues...)
}

// Warn uses fmt.Sprint to construct and log a message with prefix.
func (c *Logger) Warn(args ...interface{}) {
	if !c.core.Enabled(zap.WarnLevel) {
		return
	}
	c.logger.Warnf(c.makeMsg("", args))
}

// Warnf uses fmt.Sprintf to log a templated message with prefix.
func (c *Logger) Warnf(template string, args ...interface{}) {
	if !c.core.Enabled(zap.WarnLevel) {
		return
	}
	c.logger.Warnf(c.makeMsg(template, args))
}

// Warnw logs a message with some additional context with keys and values (prefer for json logs).
func (c *Logger) Warnw(msg string, keysAndValues ...interface{}) {
	if !c.core.Enabled(zap.WarnLevel) {
		return
	}
	c.logger.Warnw(c.makeMsg(msg, nil), keysAndValues...)
}

// Error uses fmt.Sprint to construct and log a message with prefix.
func (c *Logger) Error(args ...interface{}) {
	if !c.core.Enabled(zap.ErrorLevel) {
		return
	}
	c.logger.Errorf(c.makeMsg("", args))
}

// Errorf uses fmt.Sprintf to log a templated message with prefix.
func (c *Logger) Errorf(template string, args ...interface{}) {
	if !c.core.Enabled(zap.ErrorLevel) {
		return
	}
	c.logger.Errorf(c.makeMsg(template, args))
}

// Errorw logs a message with some additional context with keys and values (prefer for json logs).
func (c *Logger) Errorw(msg string, keysAndValues ...interface{}) {
	if !c.core.Enabled(zap.ErrorLevel) {
		return
	}
	c.logger.Errorw(c.makeMsg(msg, nil), keysAndValues...)
}

// Panic uses fmt.Sprint to construct and log a message with prefix.
func (c *Logger) Panic(args ...interface{}) {
	if !c.core.Enabled(zap.PanicLevel) {
		return
	}
	c.logger.Panicf(c.makeMsg("", args))
}

// Panicf uses fmt.Sprintf to log a templated message with prefix.
func (c *Logger) Panicf(template string, args ...interface{}) {
	if !c.core.Enabled(zap.PanicLevel) {
		return
	}
	c.logger.Panicf(c.makeMsg(template, args))
}

// Panicw logs a message with some additional context with keys and values (prefer for json logs).
func (c *Logger) Panicw(msg string, keysAndValues ...interface{}) {
	if !c.core.Enabled(zap.PanicLevel) {
		return
	}
	c.logger.Panicw(c.makeMsg(msg, nil), keysAndValues...)
}

// Fatal uses fmt.Sprint to construct and log a message with prefix.
func (c *Logger) Fatal(args ...interface{}) {
	if !c.core.Enabled(zap.FatalLevel) {
		return
	}
	c.logger.Fatalf(c.makeMsg("", args))
}

// Fatalf uses fmt.Sprintf to log a templated message with prefix.
func (c *Logger) Fatalf(template string, args ...interface{}) {
	if !c.core.Enabled(zap.FatalLevel) {
		return
	}
	c.logger.Fatalf(c.makeMsg(template, args))
}

// Fatalw logs a message with some additional context with keys and values (prefer for json logs).
func (c *Logger) Fatalw(msg string, keysAndValues ...interface{}) {
	if !c.core.Enabled(zap.FatalLevel) {
		return
	}
	c.logger.Fatalw(c.makeMsg(msg, nil), keysAndValues...)
}

// Infod is like Infof but return a function to log message with duration.
func (c *Logger) Infod() LogFunc {
	if !c.core.Enabled(zap.InfoLevel) {
		return NoLogFunc
	}
	ts := time.Now()
	return func(template string, args ...interface{}) {
		c.logger.Infof(c.makeDurationMsg(template, args, time.Since(ts)))
	}
}

// Warnd is like Warnf but return a function to log message with duration only if the duration is longer than expected
func (c *Logger) Warnd(dur time.Duration) LogFunc {
	if !c.core.Enabled(zap.WarnLevel) {
		return NoLogFunc
	}
	ts := time.Now()
	return func(template string, args ...interface{}) {
		el := time.Since(ts)
		if el > dur {
			c.logger.Warnf(c.makeDurationMsg(template, args, time.Since(ts)))
		}
	}
}

// Debugd is like Debugf but return a function to log message with duration.
func (c *Logger) Debugd() LogFunc {
	if !c.core.Enabled(zap.DebugLevel) {
		return NoLogFunc
	}
	ts := time.Now()
	return func(template string, args ...interface{}) {
		c.logger.Debugf(c.makeDurationMsg(template, args, time.Since(ts)))
	}
}

// Autod return a function to log message with elapsed duration, log level could be debug or warn depended on duration time
func (c *Logger) Autod(dur time.Duration) LogFunc {
	if !c.core.Enabled(zap.WarnLevel) {
		return NoLogFunc
	}
	ts := time.Now()
	return func(template string, args ...interface{}) {
		el := time.Since(ts)
		if el > dur {
			c.logger.Warnf(c.makeDurationMsg(template, args, el))
		} else {
			c.logger.Debugf(c.makeDurationMsg(template, args, el))
		}
	}
}

func (c *Logger) makeMsg(template string, args []interface{}) string {
	return c.prefix + getMessage(template, args)
}

func (c *Logger) makeDurationMsg(template string, args []interface{}, el time.Duration) string {
	return fmt.Sprintf("%s: %s", c.makeMsg(template, args), el)
}

// getMessage format with Sprint, Sprintf, or neither.
func getMessage(template string, args []interface{}) string {
	if len(args) == 0 {
		return template
	}

	if template != "" {
		return fmt.Sprintf(template, args...)
	}

	if len(args) == 1 {
		if str, ok := args[0].(string); ok {
			return str
		}
	}

	return fmt.Sprint(args...)
}

package logger

import (
	"os"
	"strconv"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	zapLog    *zap.SugaredLogger
	preFields []zap.Field
	once      sync.Once
	// undo      func()
)

type LoggerOption func(*LoggerConfig)

// LoggerConfig defines the configuration options for the logger.
type LoggerConfig struct {
	Debug          bool
	CommonLogPath  string
	ErrorLogPath   string
	MaxSize        int // in megabytes
	MaxBackups     int
	MaxAge         int // in days
	Compress       bool
	FileEncoder    zapcore.Encoder
	ConsoleEncoder zapcore.Encoder
}

func InitLogger(opts ...LoggerOption) {
	once.Do(func() {

		pe := zap.NewProductionEncoderConfig()
		pe.EncodeTime = zapcore.ISO8601TimeEncoder
		defaultFileEncoder := zapcore.NewJSONEncoder(pe)

		pe.EncodeLevel = zapcore.CapitalColorLevelEncoder // colorize log level
		defaultConsoleEncoder := zapcore.NewConsoleEncoder(pe)

		config := LoggerConfig{
			Debug:         false,
			CommonLogPath: "logs/common.log",
			ErrorLogPath:  "logs/error.log",
			MaxSize:       500, // in megabytes
			MaxBackups:    5,
			MaxAge:        28,   // in days
			Compress:      true, // disabled by default

			FileEncoder:    defaultFileEncoder,
			ConsoleEncoder: defaultConsoleEncoder,
		}

		for _, opt := range opts {
			opt(&config)
		}

		commonWriter := zapcore.AddSync(&lumberjack.Logger{
			Filename:   config.CommonLogPath,
			MaxSize:    config.MaxSize, // megabytes
			MaxBackups: config.MaxBackups,
			MaxAge:     config.MaxAge, // days
			Compress:   config.Compress,
		})

		errorWriter := zapcore.AddSync(&lumberjack.Logger{
			Filename:   config.ErrorLogPath,
			MaxSize:    config.MaxSize, // megabytes
			MaxBackups: config.MaxBackups,
			MaxAge:     config.MaxAge, // days
			Compress:   config.Compress,
		})

		logLevel, err := strconv.Atoi(os.Getenv("LOG_LEVEL"))
		if err != nil {
			logLevel = 0
		}
		level := zapcore.Level(logLevel)

		zapCore := zapcore.NewTee(
			zapcore.NewCore(config.FileEncoder, commonWriter, zap.InfoLevel),
			zapcore.NewCore(config.FileEncoder, errorWriter, zap.ErrorLevel),
			zapcore.NewCore(config.ConsoleEncoder, zapcore.AddSync(os.Stdout), level),
		)

		zapLog = zap.New(zapCore, zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zap.ErrorLevel)).Sugar()
	})
}

func WithCommonLogPath(path string) LoggerOption {
	return func(c *LoggerConfig) {
		c.CommonLogPath = path
	}
}

func WithErrorLogPath(path string) LoggerOption {
	return func(c *LoggerConfig) {
		c.ErrorLogPath = path
	}
}

func WithMaxSize(size int) LoggerOption {
	return func(c *LoggerConfig) {
		c.MaxSize = size
	}
}

func WithMaxBackups(backups int) LoggerOption {
	return func(c *LoggerConfig) {
		c.MaxBackups = backups
	}
}

func WithMaxAge(age int) LoggerOption {
	return func(c *LoggerConfig) {
		c.MaxAge = age
	}
}

func WithCompress(compress bool) LoggerOption {
	return func(c *LoggerConfig) {
		c.Compress = compress
	}
}

func WithDebug(debug bool) LoggerOption {
	return func(c *LoggerConfig) {
		c.Debug = debug
	}
}

func WithFileEncoder(encoder zapcore.Encoder) LoggerOption {
	return func(c *LoggerConfig) {
		c.FileEncoder = encoder
	}
}

func WithConsoleEncoder(encoder zapcore.Encoder) LoggerOption {
	return func(c *LoggerConfig) {
		c.ConsoleEncoder = encoder
	}
}

func With(fileds ...zap.Field) *zap.SugaredLogger {
	preFields = append(preFields, fileds...)
	return zapLog
}

func WithPair(key string, val interface{}) *zap.SugaredLogger {
	preFields = append(preFields, zap.Any(key, val))
	return zapLog
}

func WithString(key string, val string) *zap.SugaredLogger {
	preFields = append(preFields, zap.String(key, val))
	return zapLog
}

func WithStringPrefix(prefix string) *zap.SugaredLogger {
	return zapLog.With(zap.String("prefix", prefix))
}

func Sync() error {
	return zapLog.Sync()
}

// Debug uses fmt.Sprint to construct and log a message.
func Debug(args ...interface{}) {
	zapLog.Debug(args...)
}

// Info uses fmt.Sprint to construct and log a message.
func Info(args ...interface{}) {
	zapLog.Info(args...)
}

// Warn uses fmt.Sprint to construct and log a message.
func Warn(args ...interface{}) {
	zapLog.Warn(args...)
}

// Error uses fmt.Sprint to construct and log a message.
func Error(args ...interface{}) {
	zapLog.Error(args...)
}

// DPanic uses fmt.Sprint to construct and log a message. In development, the
// logger then panics. (See DPanicLevel for details.)
func DPanic(args ...interface{}) {
	zapLog.DPanic(args...)
}

// Panic uses fmt.Sprint to construct and log a message, then panics.
func Panic(args ...interface{}) {
	zapLog.Panic(args...)
}

// Fatal uses fmt.Sprint to construct and log a message, then calls os.Exit.
func Fatal(args ...interface{}) {
	zapLog.Fatal(args...)
}

// Debugf uses fmt.Sprintf to log a templated message.
func Debugf(template string, args ...interface{}) {
	zapLog.Debugf(template, args...)
}

// Infof uses fmt.Sprintf to log a templated message.
func Infof(template string, args ...interface{}) {
	zapLog.Infof(template, args...)
}

func Printf(template string, args ...interface{}) {
	Infof(template, args...)
}

// Warnf uses fmt.Sprintf to log a templated message.
func Warnf(template string, args ...interface{}) {
	zapLog.Warnf(template, args...)
}

// Errorf uses fmt.Sprintf to log a templated message.
func Errorf(template string, args ...interface{}) {
	zapLog.Errorf(template, args...)
}

// DPanicf uses fmt.Sprintf to log a templated message. In development, the
// logger then panics. (See DPanicLevel for details.)
func DPanicf(template string, args ...interface{}) {
	zapLog.DPanicf(template, args...)
}

// Panicf uses fmt.Sprintf to log a templated message, then panics.
func Panicf(template string, args ...interface{}) {
	zapLog.Panicf(template, args...)
}

// Fatalf uses fmt.Sprintf to log a templated message, then calls os.Exit.
func Fatalf(template string, args ...interface{}) {
	zapLog.Fatalf(template, args...)
}

// Debugw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
//
// # When debug-level logging is disabled, this is much faster than
//
// ./(args...)zapLog
func Debugw(msg string, keysAndValues ...interface{}) {
	zapLog.Debugw(msg, keysAndValues...)
}

// Infow logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func Infow(msg string, keysAndValues ...interface{}) {
	zapLog.Infow(msg, keysAndValues...)
}

// Warnw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func Warnw(msg string, keysAndValues ...interface{}) {
	zapLog.Warnw(msg, keysAndValues...)
}

// Errorw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func Errorw(msg string, keysAndValues ...interface{}) {
	zapLog.Errorw(msg, keysAndValues...)
}

// DPanicw logs a message with some additional context. In development, the
// logger then panics. (See DPanicLevel for details.) The variadic key-value
// pairs are treated as they are in With.
func DPanicw(msg string, keysAndValues ...interface{}) {
	zapLog.DPanicw(msg, keysAndValues...)
}

// Panicw logs a message with some additional context, then panics. The
// variadic key-value pairs are treated as they are in With.
func Panicw(msg string, keysAndValues ...interface{}) {
	zapLog.Panicw(msg, keysAndValues...)
}

// Fatalw logs a message with some additional context, then calls os.Exit. The
// variadic key-value pairs are treated as they are in With.
func Fatalw(msg string, keysAndValues ...interface{}) {
	zapLog.Fatalw(msg, keysAndValues...)
}

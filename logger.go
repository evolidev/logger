package logger

import (
	"context"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/evolidev/console/color"
	"github.com/evolidev/filesystem"
)

const logFormat = "%s"
const textColor = 245
const timeColor = 240
const debugColor = 3
const successColor = 2
const errorColor = 1
const logColor = 61

const (
	LevelFatal   = slog.Level(12)
	LevelSuccess = slog.Level(2)
	LevelLog     = slog.Level(1)
)

var LevelNames = map[slog.Leveler]string{
	LevelFatal:   "FATAL",
	LevelSuccess: "SUCCESS",
	LevelLog:     "LOG",
}

type Config struct {
	EnableColors bool
	Name         string
	Stdout       io.Writer
	Path         string
	PrefixColor  int
	OutputJSON   bool
}

type Logger struct {
	log      *slog.Logger
	plainLog *slog.Logger
	config   *Config
}

var Verbose = 0

func NewLogger(c *Config) *Logger {
	if c == nil {
		c = &Config{
			Name: "app",
		}
	}

	var colorfulWriters []io.Writer
	var plainWriters []io.Writer

	var output io.Writer
	if c.Stdout != nil {
		output = c.Stdout
	} else {
		output = os.Stdout
	}

	if c.EnableColors {
		colorfulWriters = append(colorfulWriters, output)
	} else {
		plainWriters = append(plainWriters, output)
	}

	if c.Path != "" {
		if !filesystem.Exists(c.Path) {
			filesystem.Write(c.Path, "")
		}

		f, err := os.OpenFile(c.Path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}

		plainWriters = append(plainWriters, f)
	}

	opts := slog.HandlerOptions{
		Level: slog.Level(-10),
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.LevelKey {
				level := a.Value.Any().(slog.Level)
				levelLabel, exists := LevelNames[level]
				if !exists {
					levelLabel = level.String()
				}

				a.Value = slog.StringValue(levelLabel)
			}

			return a
		},
	}

	if c.OutputJSON {
		return &Logger{
			log:      slog.New(slog.NewJSONHandler(io.MultiWriter(colorfulWriters...), &opts)),
			plainLog: slog.New(slog.NewJSONHandler(io.MultiWriter(plainWriters...), &opts)),
			config:   c,
		}
	}

	return &Logger{
		log:      slog.New(slog.NewTextHandler(io.MultiWriter(colorfulWriters...), &opts)),
		plainLog: slog.New(slog.NewTextHandler(io.MultiWriter(plainWriters...), &opts)),
		config:   c,
	}
}

func NewLoggerByName(name string, colorCode int) *Logger {
	return NewLogger(&Config{
		Name:        name,
		PrefixColor: colorCode,
	})
}

func (l *Logger) getPrefix() string {
	var prefixColor = l.config.PrefixColor
	prefix := ""

	if l.config.Name != "" {
		prefix = color.Text(prefixColor, "["+l.config.Name+"]") + " "
	}

	return prefix
}

func (l *Logger) getPlainPrefix() string {
	prefix := ""

	if l.config.Name != "" {
		prefix = "[" + l.config.Name + "]" + " "
	}

	return prefix
}

func (l *Logger) Log(msg interface{}, args ...interface{}) {
	var ctx = context.Background()
	l.log.Log(
		ctx,
		LevelLog,
		fmt.Sprint(l.getPrefix(), color.Text(textColor, msg)),
		args...,
	)

	l.plainLog.Log(
		ctx,
		LevelLog,
		fmt.Sprint(l.getPlainPrefix(), msg),
		args...,
	)
}

func (l *Logger) Info(msg interface{}, args ...interface{}) {
	l.log.Info(
		fmt.Sprint(l.getPrefix(), color.Text(textColor, msg)),
		args...,
	)

	l.plainLog.Info(
		fmt.Sprint(l.getPlainPrefix(), msg),
		args...,
	)
}

func (l *Logger) Success(msg interface{}, args ...interface{}) {
	var ctx = context.Background()
	l.log.Log(
		ctx,
		LevelSuccess,
		fmt.Sprint(l.getPrefix(), color.Text(textColor, msg)),
		args...,
	)

	l.plainLog.Log(
		ctx,
		LevelSuccess,
		fmt.Sprint(l.getPlainPrefix(), msg),
		args...,
	)
}

func (l *Logger) Error(msg interface{}, args ...interface{}) {
	l.log.Error(
		fmt.Sprint(l.getPrefix(), color.Text(errorColor, "Error"), color.Text(textColor, msg)),
		args...,
	)

	l.plainLog.Error(
		fmt.Sprint(l.getPlainPrefix(), msg),
		args...,
	)
}

func (l *Logger) Debug(msg interface{}, args ...interface{}) {
	l.log.Debug(
		fmt.Sprint(l.getPrefix(), color.Text(textColor, msg)),
		args...,
	)

	l.plainLog.Debug(
		fmt.Sprint(l.getPlainPrefix(), msg),
		args...,
	)
}

func (l *Logger) Print(msg interface{}, args ...interface{}) {
	l.log.Debug(color.Text(textColor, msg), args...)

	l.plainLog.Debug(fmt.Sprintf(logFormat, msg), args...)
}

func (l *Logger) Fatal(msg interface{}, args ...interface{}) {
	ctx := context.Background()
	l.log.Log(ctx, LevelFatal, color.Text(textColor, msg), args...)

	l.plainLog.Log(ctx, LevelFatal, fmt.Sprintf(logFormat, msg), args...)

	time.Sleep(1 * time.Second)
	os.Exit(1)
}

var appLogger = NewLogger(nil)

func GetAppLogger() *Logger {
	return appLogger
}

func SetAppLogger(l *Logger) {
	appLogger = l
}

func Debug(msg interface{}, args ...interface{}) {
	appLogger.Debug(msg, args...)
}

func Info(msg interface{}, args ...interface{}) {
	appLogger.Info(msg, args...)
}

func Error(msg interface{}, args ...interface{}) {
	appLogger.Error(msg, args...)
}

func Fatal(msg interface{}, args ...interface{}) {
	appLogger.Fatal(msg, args...)
}

func Log(msg interface{}, args ...interface{}) {
	appLogger.Log(msg, args...)
}

func Success(msg interface{}, args ...interface{}) {
	appLogger.Success(msg, args...)
}

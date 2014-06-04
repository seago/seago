package logger

import (
	"os"
	"path/filepath"

	l4g "github.com/seateam/log4go"
)

var levelMap = map[string]string{
	"ERROR":    "EROR",
	"CRITICAL": "CRIT",
	"DEBUG":    "DEBG",
	"FINE":     "FINE",
	"INFO":     "INFO",
}

var logger l4g.Logger

func SetLogger(logLevel, logDirName, logName string) {
	var level l4g.Level
	switch logLevel {
	case "debug":
		level = l4g.DEBUG
	case "info":
		level = l4g.INFO
	case "error":
		level = l4g.ERROR
	case "all":
		level = l4g.FINEST
	default:
		level = l4g.DEBUG
	}
	logger = make(l4g.Logger)
	logDir, _ := filepath.Abs(logDirName)
	if _, err := os.Stat(logDir); err != nil {
		os.Mkdir(logDir, os.ModePerm)
	}
	fileName := filepath.Join(logDir, logName)
	fileLog := l4g.NewFileLogWriter(fileName, true)
	fileLog.SetFormat("[%D %T] [%L] (%S) %M")
	fileLog.SetRotateDaily(true)
	logger.AddFilter("file", level, fileLog)
	consoleLog := l4g.NewConsoleLogWriter()
	logger.AddFilter("stdout", level, consoleLog)
}

func Close() {
	logger.Info("Close logger")
	for _, filter := range logger {
		filter.Close()
	}
}

func Debug(args0 interface{}, args ...interface{}) {
	logger.Debug(args0, args...)
}
func Fatal(args0 interface{}, args ...interface{}) {
	logger.Critical(args0, args...)
}
func Error(args0 interface{}, args ...interface{}) {
	logger.Error(args0, args...)
}
func Info(args0 interface{}, args ...interface{}) {
	logger.Info(args0, args...)
}
func Warn(args0 interface{}, args ...interface{}) {
	logger.Warn(args0, args...)
}

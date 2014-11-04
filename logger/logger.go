package logger

import (
	l4g "github.com/seateam/log4go"
)

var levelMap = map[string]string{
	"ERROR":    "EROR",
	"CRITICAL": "CRIT",
	"DEBUG":    "DEBG",
	"FINE":     "FINE",
	"INFO":     "INFO",
}

type Logger struct {
	logger l4g.Logger
	level  l4g.Level
}

func New(l string) *Logger {
	var level l4g.Level
	switch l {
	case "debug":
		level = l4g.DEBUG
	case "info":
		level = l4g.INFO
	case "error":
		level = l4g.ERROR
	case "all":
		level = l4g.FINEST
	case "criticals":
		level = l4g.CRITICAL
	case "fine":
		level = l4g.FINE
	case "trace":
		level = l4g.TRACE
	case "warning":
		level = l4g.WARNING
	default:
		level = l4g.DEBUG
	}
	return &Logger{logger: make(l4g.Logger), level: level}
}

func (l *Logger) SetFileLogWriter(fname string, rotate bool) {
	fileLog := l4g.NewFileLogWriter(fname, rotate)
	fileLog.SetFormat("[%D %T] [%L] (%S) %M")
	if rotate == true {
		fileLog.SetRotateDaily(true)
	}
	l.logger.AddFilter("file", l.level, fileLog)
}

func (l *Logger) SetConsoleLogWriter() {
	consoleLog := l4g.NewConsoleLogWriter()
	l.logger.AddFilter("stdout", l.level, consoleLog)
}

func (l *Logger) SetSocketLogWriter(proto, hostport string) {
	socketLog := l4g.NewSocketLogWriter(proto, hostport)
	l.logger.AddFilter("socket", l.level, socketLog)
}

func (l *Logger) Info(arg0 interface{}, args ...interface{}) {
	l.logger.Info(arg0, args...)
}

func (l *Logger) Debug(arg0 interface{}, args ...interface{}) {
	l.logger.Debug(arg0, args...)
}
func (l *Logger) Trace(arg0 interface{}, args ...interface{}) {
	l.logger.Trace(arg0, args...)
}

func (l *Logger) Fatal(arg0 interface{}, args ...interface{}) {
	l.logger.Critical(arg0, args...)
}
func (l *Logger) Error(arg0 interface{}, args ...interface{}) {
	l.logger.Error(arg0, args...)
}
func (l *Logger) Warning(arg0 interface{}, args ...interface{}) {
	l.logger.Warn(arg0, args...)
}
func (l *Logger) Fine(arg0 interface{}, args ...interface{}) {
	l.logger.Fine(arg0, args...)
}

var StdLogger = New("debug")

func (l *Logger) Close() {
	l.Info("Close logger")
	for _, filter := range l.logger {
		filter.Close()
	}
}

func SetLevel(l string) {
	var level l4g.Level
	switch l {
	case "debug":
		level = l4g.DEBUG
	case "info":
		level = l4g.INFO
	case "error":
		level = l4g.ERROR
	case "all":
		level = l4g.FINEST
	case "criticals":
		level = l4g.CRITICAL
	case "fine":
		level = l4g.FINE
	case "trace":
		level = l4g.TRACE
	case "warning":
		level = l4g.WARNING
	default:
		level = l4g.DEBUG
	}
	StdLogger.level = level
}

func SetFileLogWriter(fname string, rotate bool) {
	StdLogger.SetFileLogWriter(fname, rotate)
}

func SetConsoleLogWriter() {
	StdLogger.SetConsoleLogWriter()
}

func Debug(arg0 interface{}, args ...interface{}) {
	StdLogger.Debug(arg0, args...)
}
func Fatal(arg0 interface{}, args ...interface{}) {
	StdLogger.Fatal(arg0, args...)
}
func Error(arg0 interface{}, args ...interface{}) {
	StdLogger.Error(arg0, args...)
}
func Info(arg0 interface{}, args ...interface{}) {
	StdLogger.Info(arg0, args...)
}
func Warning(arg0 interface{}, args ...interface{}) {
	StdLogger.Warning(arg0, args...)
}

func Fine(arg0 interface{}, args ...interface{}) {
	StdLogger.Fine(arg0, args...)
}

func Trace(arg0 interface{}, args ...interface{}) {
	StdLogger.Trace(arg0, args...)
}

func Close() {
	StdLogger.Close()
}

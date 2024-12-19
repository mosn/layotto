/*
 * Copyright 2021 Layotto Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package logger

import (
	"fmt"
	"strings"
	"sync"

	"mosn.io/pkg/log"
)

const (
	// TraceLevel is for logging verbose message with a set of methods and properties to help track code execution.
	TraceLevel LogLevel = "trace"
	// DebugLevel has verbose message.
	DebugLevel LogLevel = "debug"
	// InfoLevel is default log level.
	InfoLevel LogLevel = "info"
	// WarnLevel is for logging messages about possible issues.
	WarnLevel LogLevel = "warn"
	// ErrorLevel is for logging errors.
	ErrorLevel LogLevel = "error"
	// FatalLevel is for logging fatal messages.
	FatalLevel LogLevel = "fatal"

	// UndefinedLevel is for undefined log level.
	UndefinedLevel LogLevel = "undefined"

	logKeyDebug    = "debug"
	logKeyAccess   = "access"
	logKeyError    = "error"
	fileNameDebug  = "layotto.debug.log"
	fileNameAccess = "layotto.access.log"
	fileNameError  = "layotto.error.log"
)

var (
	loggerListeners    sync.Map
	defaultLoggerLevel = InfoLevel
	defaultLogFilePath = "./"
)

// LogLevel is Logger Level type.
type LogLevel string

// ComponentLoggerListener is the interface for setting log config.
type ComponentLoggerListener interface {
	OnLogLevelChanged(outputLevel LogLevel)
}

// RegisterComponentLoggerListener registers a logger for a component logger listener.
func RegisterComponentLoggerListener(componentName string, logger ComponentLoggerListener) {
	loggerListeners.Store(componentName, logger)
}

// SetComponentLoggerLevel sets the log level for a component.
func SetComponentLoggerLevel(componentName string, level string) {
	logLevel := toLogLevel(level)
	logger, ok := loggerListeners.Load(componentName)
	if !ok {
		log.DefaultLogger.Warnf("component logger for %s not found", componentName)
	} else {
		componentLoggerListener, ok := logger.(ComponentLoggerListener)
		if !ok {
			log.DefaultLogger.Warnf("component logger for %s is not ComponentLoggerListener", componentName)
		} else {
			componentLoggerListener.OnLogLevelChanged(logLevel)
		}
	}
}

// SetDefaultLoggerLevel sets the default log output level.
func SetDefaultLoggerLevel(level string) {
	if level != "" {
		defaultLoggerLevel = toLogLevel(level)
	}
}

// SetDefaultLoggerFilePath sets the default log file path.
func SetDefaultLoggerFilePath(filePath string) {
	defaultLogFilePath = filePath
}

// layottoLogger is the implementation for layotto.
type layottoLogger struct {
	// name is the name of logger that is published to log as a component.
	name string

	logLevel LogLevel

	loggers map[string]log.ErrorLogger
}

// Logger api for logging.
type Logger interface {
	// Trace logs a message at level Trace.
	Trace(args ...interface{})
	// Tracef logs a message at level Trace.
	Tracef(format string, args ...interface{})
	// Debug logs a message at level Debug.
	Debug(args ...interface{})
	// Debugf logs a message at level Debug.
	Debugf(format string, args ...interface{})
	// Info logs a message at level Info.
	Info(args ...interface{})
	// Infof logs a message at level Info.
	Infof(format string, args ...interface{})
	// Warn logs a message at level Warn.
	Warn(args ...interface{})
	// Warnf logs a message at level Warn.
	Warnf(format string, args ...interface{})
	// Error logs a message at level Error.
	Error(args ...interface{})
	// Errorf logs a message at level Error.
	Errorf(format string, args ...interface{})
	// Fatal logs a message at level Fatal.
	Fatal(args ...interface{})
	// Fatalf logs a message at level Fatal.
	Fatalf(format string, args ...interface{})
	// SetLogLevel sets the log output level
	SetLogLevel(outputLevel LogLevel)
	// GetLogLevel get the log output level
	GetLogLevel() LogLevel
}

// toLogLevel converts to LogLevel.
func toLogLevel(level string) LogLevel {
	switch strings.ToLower(level) {
	case "trace":
		return TraceLevel
	case "debug":
		return DebugLevel
	case "info":
		return InfoLevel
	case "warn":
		return WarnLevel
	case "error":
		return ErrorLevel
	case "fatal":
		return FatalLevel
	}
	// unsupported log level
	return UndefinedLevel
}

// ToLogPriority converts to Logger priority.
func ToLogPriority(level LogLevel) int {
	switch level {
	case TraceLevel:
		return 1
	case DebugLevel:
		return 2
	case InfoLevel:
		return 3
	case WarnLevel:
		return 4
	case ErrorLevel:
		return 5
	case FatalLevel:
		return 6
	}
	return 0
}

// NewLayottoLogger creates new Logger instance.
func NewLayottoLogger(name string) Logger {
	ll := &layottoLogger{
		name:     name,
		logLevel: defaultLoggerLevel,
		loggers:  make(map[string]log.ErrorLogger),
	}

	dMosnLogger, err := log.GetOrCreateLogger(defaultLogFilePath+fileNameDebug, nil)

	dLogger := &log.SimpleErrorLog{
		Logger: dMosnLogger,
		Level:  log.DEBUG,
	}
	if err != nil {
		ll.loggers[logKeyDebug] = log.DefaultLogger
		log.DefaultLogger.Errorf("Failed to create mosn logger: %v", err)
	} else {
		dLogger.SetLogLevel(toMosnLoggerLevel(defaultLoggerLevel))
		ll.loggers[logKeyDebug] = dLogger
	}

	aMosnLogger, err := log.GetOrCreateLogger(defaultLogFilePath+fileNameAccess, nil)

	aLogger := &log.SimpleErrorLog{
		Logger: aMosnLogger,
		Level:  log.INFO,
	}
	if err != nil {
		ll.loggers[logKeyAccess] = log.DefaultLogger
		log.DefaultLogger.Errorf("Failed to create mosn logger: %v", err)
	} else {
		aLogger.SetLogLevel(toMosnLoggerLevel(defaultLoggerLevel))
		ll.loggers[logKeyAccess] = aLogger
	}

	eMosnLogger, err := log.GetOrCreateLogger(defaultLogFilePath+fileNameError, nil)

	eLogger := &log.SimpleErrorLog{
		Logger: eMosnLogger,
		Level:  log.ERROR,
	}
	if err != nil {
		ll.loggers[logKeyError] = log.DefaultLogger
		log.DefaultLogger.Errorf("Failed to create mosn logger: %v", err)
	} else {
		eLogger.SetLogLevel(toMosnLoggerLevel(defaultLoggerLevel))
		ll.loggers[logKeyError] = eLogger
	}
	return ll
}

// Tracef logs a message at level Trace.
func (l *layottoLogger) Tracef(format string, args ...interface{}) {
	l.loggers[logKeyDebug].Tracef("[%s] %s", l.name, fmt.Sprintf(format, args...))
}

// Trace logs a message at level Trace.
func (l *layottoLogger) Trace(args ...interface{}) {
	l.loggers[logKeyDebug].Tracef("%s", args...)
}

// Debugf logs a message at level Debug.
func (l *layottoLogger) Debugf(format string, args ...interface{}) {
	l.loggers[logKeyDebug].Debugf("[%s] %s", l.name, fmt.Sprintf(format, args...))
}

// Debug logs a message at level Debug.
func (l *layottoLogger) Debug(args ...interface{}) {
	l.loggers[logKeyDebug].Debugf("%s", args...)
}

// Infof logs a message at level Info.
func (l *layottoLogger) Infof(format string, args ...interface{}) {
	l.loggers[logKeyAccess].Infof("[%s] %s", l.name, fmt.Sprintf(format, args...))
}

// Info logs a message at level Info.
func (l *layottoLogger) Info(args ...interface{}) {
	l.loggers[logKeyAccess].Infof("%s", args...)
}

// Warnf logs a message at level Warn.
func (l *layottoLogger) Warnf(format string, args ...interface{}) {
	l.loggers[logKeyAccess].Warnf("[%s] %s", l.name, fmt.Sprintf(format, args...))
}

// Warn logs a message at level Warn.
func (l *layottoLogger) Warn(args ...interface{}) {
	l.loggers[logKeyAccess].Warnf("%s", args...)
}

// Errorf logs a message at level Error.
func (l *layottoLogger) Errorf(format string, args ...interface{}) {
	l.loggers[logKeyError].Errorf("[%s] %s", l.name, fmt.Sprintf(format, args...))
}

// Error logs a message at level Error.
func (l *layottoLogger) Error(args ...interface{}) {
	l.loggers[logKeyError].Errorf("%s", args...)
}

// Fatalf logs a message at level Fatal.
func (l *layottoLogger) Fatalf(format string, args ...interface{}) {
	l.loggers[logKeyError].Fatalf("[%s] %s", l.name, fmt.Sprintf(format, args...))
}

// Fatal logs a message at level Fatal.
func (l *layottoLogger) Fatal(args ...interface{}) {
	l.loggers[logKeyError].Fatalf("%s", args...)
}

// GetLogLevel gets the log output level.
func (l *layottoLogger) GetLogLevel() LogLevel {
	return l.logLevel
}

// toMosnLoggerLevel converts to logrus.Level.
func toMosnLoggerLevel(lvl LogLevel) log.Level {
	// ignore error because it will never happen
	l, _ := parseLevel(string(lvl))
	return l
}

// parseLevel takes a string level and returns the Mosn logger level constant.
func parseLevel(lvl string) (log.Level, error) {
	switch strings.ToLower(lvl) {
	case "fatal":
		return log.FATAL, nil
	case "error":
		return log.ERROR, nil
	case "warn", "warning":
		return log.WARN, nil
	case "info":
		return log.INFO, nil
	case "debug":
		return log.DEBUG, nil
	case "trace":
		return log.TRACE, nil
	}

	var l log.Level
	return l, fmt.Errorf("not a valid mosn Level: %q", lvl)
}

// SetLogLevel sets log output level.
func (l *layottoLogger) SetLogLevel(outputLevel LogLevel) {
	l.logLevel = outputLevel
	l.loggers[logKeyDebug].SetLogLevel(toMosnLoggerLevel(outputLevel))
	l.loggers[logKeyAccess].SetLogLevel(toMosnLoggerLevel(outputLevel))
	l.loggers[logKeyError].SetLogLevel(toMosnLoggerLevel(outputLevel))
}

//
// Copyright (C) 2021 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/logger
//

/*

Package logger outputs log message in the well-defined format,
using UTC date and times.

  2020/12/01 20:30:40 main.go:11: [level] message {"json": "context"}
*/
package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

/*

Logger interface defines finer grained log levels.
*/
type Logger interface {
	// build new logger with context
	With(Note) Logger

	/*
		system is unusable, panic execution of current routine/application,
		it is notpossible to gracefully terminate it.
	*/
	Emergency(format string, v ...interface{}) error

	/*
		system is failed, response actions must be taken immediately,
		the application is not able to execute correctly but still
		able to gracefully exit.
	*/
	Critical(format string, v ...interface{}) error

	/*
		system is failed, unable to recover from error.
		The failure do not have global catastrophic impacts but
		local functionality is impaired, incorrect result is returned.
	*/
	Error(format string, v ...interface{}) error

	/*
		system is failed, unable to recover, degraded functionality.
		The failure is ignored and application still capable to deliver
		incomplete but correct results.
	*/
	Warning(format string, v ...interface{}) error

	/*
		system is failed, error is recovered, no impact
	*/
	Notice(format string, v ...interface{}) error

	/*
		output informative status about system
	*/
	Info(format string, v ...interface{}) error

	/*
		output debug status about system
	*/
	Debug(format string, v ...interface{}) error
}

/*

LogLevel constants
*/
const (
	lEMERGENCY = iota
	lCRITICAL
	lERROR
	lWARNING
	lNOTICE
	lINFO
	lDEBUG
)

/*

LogLevel constants
*/
const (
	EMERGENCY = "emergency"
	CRITICAL  = "critical"
	ERROR     = "error"
	WARNING   = "warning"
	NOTICE    = "notice"
	INFO      = "info"
	DEBUG     = "debug"
)

// logger instance is ensemble of multiple standard loggers, each per log level
type logger struct {
	emergency,
	critical,
	err,
	warning,
	notice,
	info,
	debug *log.Logger
}

const flags = log.Lmsgprefix | log.LUTC | log.Ldate | log.Ltime | log.Lshortfile

var (
	global *logger
)

func init() {
	global = newLogger(logLevelFromEnv(), os.Stderr)
}

func logLevelFromEnv() int {
	loglevel, defined := os.LookupEnv("CONFIG_LOGGER_LEVEL")
	if !defined {
		return lDEBUG
	}

	seq := [...]string{EMERGENCY, CRITICAL, ERROR, WARNING, NOTICE, INFO, DEBUG}
	for level, name := range seq {
		if strings.ToLower(loglevel) == name {
			return level
		}
	}

	return lDEBUG
}

/*

Config logger, the configuration support
* definition of log level using string constant
* output destination

  logger.Config("WARNING", os.Stderr)
	logger.Config(os.Getenv("CONFIG_LOG_LEVEL"), os.Stderr)
*/
func Config(loglevel string, out io.Writer) {
	seq := [...]string{EMERGENCY, CRITICAL, ERROR, WARNING, NOTICE, INFO, DEBUG}
	for level, name := range seq {
		if strings.ToLower(loglevel) == name {
			global = newLogger(level, out)
			return
		}
	}
}

func newLogger(loglevel int, out io.Writer) *logger {
	l := &logger{}

	switch loglevel {
	case lDEBUG:
		l.debug = log.New(out, "["+DEBUG+"] ", flags)
		fallthrough
	case lINFO:
		l.info = log.New(out, "["+INFO+"] ", flags)
		fallthrough
	case lNOTICE:
		l.notice = log.New(out, "["+NOTICE+"] ", flags)
		fallthrough
	case lWARNING:
		l.warning = log.New(out, "["+WARNING+"] ", flags)
		fallthrough
	case lERROR:
		l.err = log.New(out, "["+ERROR+"] ", flags)
		fallthrough
	case lCRITICAL:
		l.critical = log.New(out, "["+CRITICAL+"] ", flags)
		fallthrough
	case lEMERGENCY:
		l.emergency = log.New(out, "["+EMERGENCY+"] ", flags)
	}
	return l
}

/*

Emergency logging:
system is unusable, panic execution of current routine/application,
it is notpossible to gracefully terminate it.
*/
func Emergency(format string, v ...interface{}) error {
	if global.emergency != nil {
		if err := global.emergency.Output(2, fmt.Sprintf(format, v...)); err != nil {
			return err
		}
		panic(fmt.Errorf(format, v...))
	}

	return nil
}

/*

Critical logging:
system is failed, response actions must be taken immediately,
the application is not able to execute correctly but still
able to gracefully exit.
*/
func Critical(format string, v ...interface{}) error {
	if global.critical != nil {
		return global.critical.Output(2, fmt.Sprintf(format, v...))
	}

	return nil
}

/*

Error logging:
system is failed, unable to recover from error.
The failure do not have global catastrophic impacts but
local functionality is impaired, incorrect result is returned.
*/
func Error(format string, v ...interface{}) error {
	if global.err != nil {
		return global.err.Output(2, fmt.Sprintf(format, v...))
	}

	return nil
}

/*

Warning logging:
system is failed, unable to recover, degraded functionality.
The failure is ignored and application still capable to deliver
incomplete but correct results.
*/
func Warning(format string, v ...interface{}) error {
	if global.warning != nil {
		return global.warning.Output(2, fmt.Sprintf(format, v...))
	}

	return nil
}

/*

Notice logging:
system is failed, error is recovered, no impact
*/
func Notice(format string, v ...interface{}) error {
	if global.notice != nil {
		return global.notice.Output(2, fmt.Sprintf(format, v...))
	}

	return nil
}

/*

Info logging:
output informative status about system
*/
func Info(format string, v ...interface{}) error {
	if global.info != nil {
		return global.info.Output(2, fmt.Sprintf(format, v...))
	}

	return nil
}

/*

Debug logging:
output debug status about system
*/
func Debug(format string, v ...interface{}) error {
	if global.debug != nil {
		return global.debug.Output(2, fmt.Sprintf(format, v...))
	}

	return nil
}

/*

Note type annotates log message with semi-structured data
*/
type Note map[string]interface{}

/*

With returns contextual logger
*/
func With(note Note) Logger {
	return newContext(note)
}

type context struct {
	note string
}

func newContext(note Note) Logger {
	return &context{note: toJSON(note)}
}

func toJSON(note Note) string {
	var out strings.Builder

	for key, val := range note {
		out.WriteString(", ")
		k, _ := json.Marshal(key)
		out.Write(k)
		out.WriteString(": ")
		v, _ := json.Marshal(val)
		out.Write(v)
	}

	return out.String()[2:]
}

func (log *context) With(note Note) Logger {
	if len(log.note) != 0 {
		return &context{note: log.note + ", " + toJSON(note)}
	}

	return &context{note: toJSON(note)}
}

func (log *context) sprintf(format string, v ...interface{}) string {
	if len(log.note) == 0 {
		return fmt.Sprintf(format, v...)
	} else {
		var out strings.Builder
		out.WriteString("{ ")
		out.WriteString(log.note)
		out.WriteString(", \"message\": \"")
		out.WriteString(fmt.Sprintf(format, v...))
		out.WriteString("\" }")
		return out.String()
	}
}

/*

Emergency logging:
system is unusable, panic execution of current routine/application,
it is notpossible to gracefully terminate it.
*/
func (log *context) Emergency(format string, v ...interface{}) error {
	if global.emergency != nil {
		if err := global.emergency.Output(2, log.sprintf(format, v...)); err != nil {
			return nil
		}
		panic(fmt.Errorf(format, v...))
	}

	return nil
}

/*

Critical logging:
system is failed, response actions must be taken immediately,
the application is not able to execute correctly but still
able to gracefully exit.
*/
func (log *context) Critical(format string, v ...interface{}) error {
	if global.critical != nil {
		return global.critical.Output(2, log.sprintf(format, v...))
	}

	return nil
}

/*

Error logging:
system is failed, unable to recover from error.
The failure do not have global catastrophic impacts but
local functionality is impaired, incorrect result is returned.
*/
func (log *context) Error(format string, v ...interface{}) error {
	if global.err != nil {
		return global.err.Output(2, log.sprintf(format, v...))
	}

	return nil
}

/*

Warning logging:
system is failed, unable to recover, degraded functionality.
The failure is ignored and application still capable to deliver
incomplete but correct results.
*/
func (log *context) Warning(format string, v ...interface{}) error {
	if global.warning != nil {
		return global.warning.Output(2, log.sprintf(format, v...))
	}

	return nil
}

/*

Notice logging:
system is failed, error is recovered, no impact.
*/
func (log *context) Notice(format string, v ...interface{}) error {
	if global.notice != nil {
		return global.notice.Output(2, log.sprintf(format, v...))
	}

	return nil
}

/*

Info logging:
output informative status about system
*/
func (log *context) Info(format string, v ...interface{}) error {
	if global.info != nil {
		return global.info.Output(2, log.sprintf(format, v...))
	}

	return nil
}

/*

Debug logging:
output debug status about system
*/
func (log *context) Debug(format string, v ...interface{}) error {
	if global.debug != nil {
		return global.debug.Output(2, log.sprintf(format, v...))
	}

	return nil
}

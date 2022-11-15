//
// Copyright (C) 2021 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/logger
//

// Package logger outputs log message as JSON in the well-defined format,
// using UTC date and times.
//
//	2020/12/01 20:30:40 main.go:11: [level] {"lvl": "level", "json": "context", "msg": "some text"}
package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

// LogLevel constants
const (
	lEMERGENCY = iota
	lCRITICAL
	lERROR
	lWARNING
	lNOTICE
	lINFO
	lDEBUG
)

// LogLevel constants
const (
	EMERGENCY = "emergency" // EMR
	CRITICAL  = "critical"  // CRT
	ERROR     = "error"     // ERR
	WARNING   = "warning"   // WRN
	NOTICE    = "notice"    // NTC
	INFO      = "info"      // INF
	DEBUG     = "debug"     // DEB
)

const (
	EMR = "emr"
	CRT = "crt"
	ERR = "err"
	WRN = "wrn"
	NTC = "ntc"
	INF = "inf"
	DEB = "deb"
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

// Config logger, the configuration support
// - definition of log level using string constant
// - output destination
//
//	  logger.Config("WARNING", os.Stderr)
//		logger.Config(os.Getenv("CONFIG_LOG_LEVEL"), os.Stderr)
func Config(loglevel string, out io.Writer) {
	long := [...]string{EMERGENCY, CRITICAL, ERROR, WARNING, NOTICE, INFO, DEBUG}
	for level, name := range long {
		if strings.ToLower(loglevel) == name {
			global = newLogger(level, out)
			return
		}
	}

	short := [...]string{EMR, CRT, ERR, WRN, NTC, INF, DEB}
	for level, name := range short {
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
		l.debug = log.New(out, "["+DEB+"] ", flags)
		fallthrough
	case lINFO:
		l.info = log.New(out, "["+INF+"] ", flags)
		fallthrough
	case lNOTICE:
		l.notice = log.New(out, "["+NTC+"] ", flags)
		fallthrough
	case lWARNING:
		l.warning = log.New(out, "["+WRN+"] ", flags)
		fallthrough
	case lERROR:
		l.err = log.New(out, "["+ERR+"] ", flags)
		fallthrough
	case lCRITICAL:
		l.critical = log.New(out, "["+CRT+"] ", flags)
		fallthrough
	case lEMERGENCY:
		l.emergency = log.New(out, "["+EMR+"] ", flags)
	}
	return l
}

// Emergency logging:
//
// system is unusable, panic execution of current routine/application,
// it is notpossible to gracefully terminate it.
func Emergency(seq ...Message) error {
	if global.emergency != nil {
		msg := Messages(seq).toString(EMERGENCY)
		if err := global.emergency.Output(2, msg); err != nil {
			return err
		}
		panic(msg)
	}

	return nil
}

// Critical logging:
//
// system is failed, response actions must be taken immediately,
// the application is not able to execute correctly but still
// able to gracefully exit.
func Critical(seq ...Message) error {
	if global.critical != nil {
		return global.critical.Output(2, Messages(seq).toString(CRT))
	}

	return nil
}

// Error logging:
//
// system is failed, unable to recover from error.
// The failure do not have global catastrophic impacts but
// local functionality is impaired, incorrect result is returned.
func Error(seq ...Message) error {
	if global.err != nil {
		return global.err.Output(2, Messages(seq).toString(ERR))
	}

	return nil
}

// Warning logging:
//
// system is failed, unable to recover, degraded functionality.
// The failure is ignored and application still capable to deliver
// incomplete but correct results.
func Warning(seq ...Message) error {
	if global.warning != nil {
		return global.warning.Output(2, Messages(seq).toString(WRN))
	}

	return nil
}

// Notice logging:
//
// system is failed, error is recovered, no impact
func Notice(seq ...Message) error {
	if global.notice != nil {
		return global.notice.Output(2, Messages(seq).toString(NTC))
	}

	return nil
}

// Info logging:
//
// output informative status about system
func Info(seq ...Message) error {
	if global.info != nil {
		return global.info.Output(2, Messages(seq).toString(INF))
	}

	return nil
}

// Debug logging:
//
// output debug status about system
func Debug(seq ...Message) error {
	if global.debug != nil {
		return global.debug.Output(2, Messages(seq).toString(DEB))
	}

	return nil
}

// Groups sequence of messages as re-usable context
func Context(seq ...Message) (msg Message) {
	buf := strings.Builder{}
	Messages(seq).toJSON(&buf)
	msg.val = buf.String()
	return
}

type Types interface {
	string | []byte | int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64 | bool
}

// Builds loggable Key/Val pair
func Val[T Types](key string, val T) (msg Message) {
	msg.key = escape(key)
	switch v := any(val).(type) {
	case string:
		msg.val = escape(v)
		msg.ext = `"`
	case []byte:
		msg.val = escape(string(v))
		msg.ext = `"`
	case int:
		msg.val = strconv.FormatInt(int64(v), 10)
	case int8:
		msg.val = strconv.FormatInt(int64(v), 10)
	case int16:
		msg.val = strconv.FormatInt(int64(v), 10)
	case int32:
		msg.val = strconv.FormatInt(int64(v), 10)
	case int64:
		msg.val = strconv.FormatInt(v, 10)
	case uint:
		msg.val = strconv.FormatUint(uint64(v), 10)
	case uint8:
		msg.val = strconv.FormatUint(uint64(v), 10)
	case uint16:
		msg.val = strconv.FormatUint(uint64(v), 10)
	case uint32:
		msg.val = strconv.FormatUint(uint64(v), 10)
	case uint64:
		msg.val = strconv.FormatUint(uint64(v), 10)
	case float32:
		msg.val = strconv.FormatFloat(float64(v), 'g', -1, 32)
	case float64:
		msg.val = strconv.FormatFloat(v, 'g', -1, 64)
	case bool:
		msg.val = strconv.FormatBool(v)
	}
	return
}

func Str(key string, val fmt.Stringer) (msg Message) {
	msg.key = escape(key)
	msg.val = escape(val.String())
	msg.ext = `"`
	return
}

func Fmt(key string, format string, v ...interface{}) (msg Message) {
	msg.key = escape(key)
	msg.ext = `"`

	if len(v) == 0 {
		msg.val = escape(format)
	} else {
		msg.val = escape(fmt.Sprintf(format, v...))
	}

	return
}

func Msg(val string, v ...interface{}) (msg Message) {
	msg.key = "msg"
	msg.ext = `"`

	if len(v) == 0 {
		msg.val = escape(val)
	} else {
		msg.val = escape(fmt.Sprintf(val, v...))
	}

	return
}

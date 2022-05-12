//
// Copyright (C) 2021 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/logger
//

package logger_test

import (
	"bytes"
	"regexp"
	"testing"

	"github.com/fogfish/it"
	"github.com/fogfish/logger"
)

func testLogger(loglevel string, logf func(format string, v ...interface{}) error) bool {
	var buf bytes.Buffer

	logger.Config(loglevel, &buf)
	logf("foobar %d", 1)

	matched, _ := regexp.Match("\\d{4}/\\d{2}/\\d{2} \\d{2}:\\d{2}:\\d{2} logger_test.go:24: ."+loglevel+". foobar 1", buf.Bytes())

	return matched
}

func TestGlobalLogger(t *testing.T) {
	it.Ok(t).
		IfTrue(testLogger(logger.CRITICAL, logger.Critical)).
		IfTrue(testLogger(logger.ERROR, logger.Error)).
		IfTrue(testLogger(logger.WARNING, logger.Warning)).
		IfTrue(testLogger(logger.NOTICE, logger.Notice)).
		IfTrue(testLogger(logger.INFO, logger.Info)).
		IfTrue(testLogger(logger.DEBUG, logger.Debug))
}

func testLoggerSilent(loglevel string, logf func(format string, v ...interface{}) error) bool {
	var buf bytes.Buffer

	logger.Config(loglevel, &buf)
	logf("foobar %d", 1)

	return buf.String() == ""
}

func TestGlobalLoggerSilentDebug(t *testing.T) {
	it.Ok(t).
		IfFalse(testLoggerSilent(logger.INFO, logger.Critical)).
		IfFalse(testLoggerSilent(logger.INFO, logger.Error)).
		IfFalse(testLoggerSilent(logger.INFO, logger.Warning)).
		IfFalse(testLoggerSilent(logger.INFO, logger.Notice)).
		IfFalse(testLoggerSilent(logger.INFO, logger.Info)).
		IfTrue(testLoggerSilent(logger.INFO, logger.Debug))
}

func TestGlobalLoggerSilentInfo(t *testing.T) {
	it.Ok(t).
		IfFalse(testLoggerSilent(logger.NOTICE, logger.Critical)).
		IfFalse(testLoggerSilent(logger.NOTICE, logger.Error)).
		IfFalse(testLoggerSilent(logger.NOTICE, logger.Warning)).
		IfFalse(testLoggerSilent(logger.NOTICE, logger.Notice)).
		IfTrue(testLoggerSilent(logger.NOTICE, logger.Info)).
		IfTrue(testLoggerSilent(logger.NOTICE, logger.Debug))
}

func TestGlobalLoggerSilentNotice(t *testing.T) {
	it.Ok(t).
		IfFalse(testLoggerSilent(logger.WARNING, logger.Critical)).
		IfFalse(testLoggerSilent(logger.WARNING, logger.Error)).
		IfFalse(testLoggerSilent(logger.WARNING, logger.Warning)).
		IfTrue(testLoggerSilent(logger.WARNING, logger.Notice)).
		IfTrue(testLoggerSilent(logger.WARNING, logger.Info)).
		IfTrue(testLoggerSilent(logger.WARNING, logger.Debug))
}

func TestGlobalLoggerSilentWarning(t *testing.T) {
	it.Ok(t).
		IfFalse(testLoggerSilent(logger.ERROR, logger.Critical)).
		IfFalse(testLoggerSilent(logger.ERROR, logger.Error)).
		IfTrue(testLoggerSilent(logger.ERROR, logger.Warning)).
		IfTrue(testLoggerSilent(logger.ERROR, logger.Notice)).
		IfTrue(testLoggerSilent(logger.ERROR, logger.Info)).
		IfTrue(testLoggerSilent(logger.ERROR, logger.Debug))
}

func TestGlobalLoggerSilentError(t *testing.T) {
	it.Ok(t).
		IfFalse(testLoggerSilent(logger.CRITICAL, logger.Critical)).
		IfTrue(testLoggerSilent(logger.CRITICAL, logger.Error)).
		IfTrue(testLoggerSilent(logger.CRITICAL, logger.Warning)).
		IfTrue(testLoggerSilent(logger.CRITICAL, logger.Notice)).
		IfTrue(testLoggerSilent(logger.CRITICAL, logger.Info)).
		IfTrue(testLoggerSilent(logger.CRITICAL, logger.Debug))
}

func TestGlobalLoggerSilentCritical(t *testing.T) {
	it.Ok(t).
		IfTrue(testLoggerSilent(logger.EMERGENCY, logger.Critical)).
		IfTrue(testLoggerSilent(logger.EMERGENCY, logger.Error)).
		IfTrue(testLoggerSilent(logger.EMERGENCY, logger.Warning)).
		IfTrue(testLoggerSilent(logger.EMERGENCY, logger.Notice)).
		IfTrue(testLoggerSilent(logger.EMERGENCY, logger.Info)).
		IfTrue(testLoggerSilent(logger.EMERGENCY, logger.Debug))
}

func testLoggerWith(loglevel string, logf func(format string, v ...interface{}) error) bool {
	var buf bytes.Buffer

	logger.Config(loglevel, &buf)
	logf("foobar %d", 1)

	matched, _ := regexp.Match("\\d{4}/\\d{2}/\\d{2} \\d{2}:\\d{2}:\\d{2} logger_test.go:114: ."+loglevel+". \\{[^}]+\"message\": \"foobar 1\" \\}", buf.Bytes())

	return matched
}

func TestLoggerWith(t *testing.T) {
	l := logger.With(logger.Note{
		"foo": "bar",
		"bar": 1,
	})

	it.Ok(t).
		IfTrue(testLoggerWith(logger.CRITICAL, l.Critical)).
		IfTrue(testLoggerWith(logger.ERROR, l.Error)).
		IfTrue(testLoggerWith(logger.WARNING, l.Warning)).
		IfTrue(testLoggerWith(logger.NOTICE, l.Notice)).
		IfTrue(testLoggerWith(logger.INFO, l.Info)).
		IfTrue(testLoggerWith(logger.DEBUG, l.Debug))
}

func TestLoggerWithMultiple(t *testing.T) {
	a := logger.With(logger.Note{"foo": "bar"})
	b := a.With(logger.Note{"bar": 1})

	var buf bytes.Buffer
	logger.Config(logger.DEBUG, &buf)

	b.Debug("some text")
	matched, _ := regexp.Match("\\d{4}/\\d{2}/\\d{2} \\d{2}:\\d{2}:\\d{2} logger_test.go:143: .debug. \\{ \"foo\": \"bar\", \"bar\": 1, \"message\": \"some text\" \\}", buf.Bytes())
	it.Ok(t).IfTrue(matched)

	a.Debug("some text")
	matched, _ = regexp.Match("\\d{4}/\\d{2}/\\d{2} \\d{2}:\\d{2}:\\d{2} logger_test.go:147: .debug. \\{ \"foo\": \"bar\", \"message\": \"some text\" \\}", buf.Bytes())
	it.Ok(t).IfTrue(matched)
}

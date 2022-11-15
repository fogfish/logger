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
	"github.com/fogfish/logger/v2"
)

func testLogger(loglevel string, logf func(seq ...logger.Message) error) bool {
	var buf bytes.Buffer

	logger.Config(loglevel, &buf)
	if err := logf(logger.Val("foobar", 1)); err != nil {
		panic(err)
	}

	matched, _ := regexp.Match("\\d{4}/\\d{2}/\\d{2} \\d{2}:\\d{2}:\\d{2} logger_test.go:24: ."+loglevel+". \\{\"lvl\": \""+loglevel+"\", \"foobar\": 1\\}", buf.Bytes())

	return matched
}

func TestGlobalLogger(t *testing.T) {
	it.Ok(t).
		IfTrue(testLogger(logger.CRT, logger.Critical)).
		IfTrue(testLogger(logger.ERR, logger.Error)).
		IfTrue(testLogger(logger.WRN, logger.Warning)).
		IfTrue(testLogger(logger.NTC, logger.Notice)).
		IfTrue(testLogger(logger.INF, logger.Info)).
		IfTrue(testLogger(logger.DEB, logger.Debug))
}

func testLoggerSilent(loglevel string, logf func(seq ...logger.Message) error) bool {
	var buf bytes.Buffer

	logger.Config(loglevel, &buf)
	if err := logf(logger.Val("foobar", 1)); err != nil {
		panic(err)
	}

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

func testLoggerWith(loglevel string, logf func(seq ...logger.Message) error) bool {
	var buf bytes.Buffer

	logger.Config(loglevel, &buf)
	if err := logf(logger.Val("foobar", 1), logger.Msg("baz")); err != nil {
		panic(err)
	}

	matched, _ := regexp.Match("\\d{4}/\\d{2}/\\d{2} \\d{2}:\\d{2}:\\d{2} logger_test.go:118: ."+loglevel+". \\{\"lvl\": \""+loglevel+"\", \"foobar\": 1, \"msg\": \"baz\"\\}", buf.Bytes())
	return matched
}

func TestLoggerWith(t *testing.T) {
	it.Ok(t).
		IfTrue(testLoggerWith(logger.CRT, logger.Critical)).
		IfTrue(testLoggerWith(logger.ERR, logger.Error)).
		IfTrue(testLoggerWith(logger.WRN, logger.Warning)).
		IfTrue(testLoggerWith(logger.NTC, logger.Notice)).
		IfTrue(testLoggerWith(logger.INF, logger.Info)).
		IfTrue(testLoggerWith(logger.DEB, logger.Debug))
}

type null struct{}

func (null) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func BenchmarkLogger(b *testing.B) {
	logger.Config(logger.DEBUG, &null{})
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		logger.Debug(
			logger.Val("key", string([]byte{0x00, 0x00})),
			logger.Val("key", "xxx"),
		)
	}
}

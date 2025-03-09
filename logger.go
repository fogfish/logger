//
// Copyright (C) 2021 - 2025 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/logger
//

// Package logger configures slog for AWS CloudWatch
package logger

import (
	"log/slog"
	"os"
)

// Create New Logger
func New(opts ...Option) *slog.Logger {
	if _, has := os.LookupEnv("AWS_LAMBDA_FUNCTION_NAME"); has {
		return slog.New(NewJSONHandler(opts...))
	}

	if preset, has := os.LookupEnv("CONFIG_LOG_PROFILE"); has {
		switch preset {
		case "CloudWatch":
			return slog.New(NewJSONHandler(opts...))
		default:
			return slog.New(NewStdioHandler(opts...))
		}
	}

	return slog.New(NewStdioHandler(opts...))
}

const (
	// EMERGENCY
	// system is unusable, panic execution of current routine/application,
	// it is notpossible to gracefully terminate it.
	EMERGENCY = slog.Level(100)
	EMR       = EMERGENCY

	// CRITICAL
	// system is failed, response actions must be taken immediately,
	// the application is not able to execute correctly but still
	// able to gracefully exit.
	CRITICAL = slog.Level(50)
	CRT      = CRITICAL

	// ERROR
	// system is failed, unable to recover from error.
	// The failure do not have global catastrophic impacts but
	// local functionality is impaired, incorrect result is returned.
	ERROR = slog.LevelError
	ERR   = ERROR

	// WARN
	// system is failed, unable to recover, degraded functionality.
	// The failure is ignored and application still capable to deliver
	// incomplete but correct results.
	WARN = slog.LevelWarn
	WRN  = WARN

	// NOTICE
	// system is failed, error is recovered, no impact
	NOTICE = slog.Level(2)
	NTC    = NOTICE

	// INFO
	// output informative status about system
	INFO = slog.LevelInfo
	INF  = INFO

	// DEBUG
	// output debug status about system
	DEBUG = slog.LevelDebug
	DEB   = DEBUG
)

var (
	levelShortName = map[slog.Level]string{
		EMERGENCY: "EMR",
		CRITICAL:  "CRT",
		ERROR:     "ERR",
		WARN:      "WRN",
		NOTICE:    "NTC",
		INFO:      "INF",
		DEBUG:     "DEB",
	}

	levelLongName = map[slog.Level]string{
		EMERGENCY: "EMERGENCY",
		CRITICAL:  "CRITICAL",
		ERROR:     "ERROR",
		WARN:      "WARN",
		NOTICE:    "NOTICE",
		INFO:      "INFO",
		DEBUG:     "DEBUG",
	}

	colorReset = "\x1b[0m"

	levelColorForName = map[slog.Level]string{
		EMERGENCY: "\x1b[1;48;5;198m", // Red (intensive)
		CRITICAL:  "\x1b[48;5;160m",   // Red (mid-intensive)
		ERROR:     "\x1b[1;31m",       // Red
		WARN:      "\x1b[1;38;5;178m", // Orange
		NOTICE:    "\x1b[0;38;5;111m", // Blue
		INFO:      "\x1b[38;5;255m",   // White
		DEBUG:     "\x1b[38;5;248m",   // Grey
	}

	levelColorForText = map[slog.Level]string{
		EMERGENCY: "\x1b[1;38;5;198m", // Red
		CRITICAL:  "\x1b[0;38;5;160m", // Red
		ERROR:     "\x1b[38;5;15m",    // White
		WARN:      "\x1b[38;5;15m",    // White
		NOTICE:    "\x1b[38;5;15m",    // White
		INFO:      "\x1b[38;5;15m",    // White
		DEBUG:     "\x1b[0m",          // default
	}

	levelColorForAttr = map[slog.Level]string{
		EMERGENCY: "\x1b[38;5;255m", // White
		CRITICAL:  "\x1b[38;5;255m", // White
		ERROR:     "\x1b[38;5;250m", // Grey
		WARN:      "\x1b[38;5;244m", // Drak Gray
		NOTICE:    "\x1b[38;5;244m", // Dark Gray
		INFO:      "\x1b[90m",       // Darker Gray
		DEBUG:     "\x1b[90m",       // Darker Gray
	}

	longNames = map[string]slog.Level{
		"EMERGENCY": EMERGENCY,
		"CRITICAL":  CRITICAL,
		"ERROR":     ERROR,
		"WARN":      WARN,
		"NOTICE":    NOTICE,
		"INFO":      INFO,
		"DEBUG":     DEBUG,
	}
)

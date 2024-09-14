//
// Copyright (C) 2021 - 2024 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/logger
//

package xlog

import (
	"context"
	"log/slog"
)

// Note: intentionally copied from logger for dependencies
const (
	// EMERGENCY
	// system is unusable, panic execution of current routine/application,
	// it is notpossible to gracefully terminate it.
	xEMERGENCY = slog.Level(100)

	// CRITICAL
	// system is failed, response actions must be taken immediately,
	// the application is not able to execute correctly but still
	// able to gracefully exit.
	xCRITICAL = slog.Level(50)

	// ERROR
	// system is failed, unable to recover from error.
	// The failure do not have global catastrophic impacts but
	// local functionality is impaired, incorrect result is returned.
	xERROR = slog.LevelError

	// WARN
	// system is failed, unable to recover, degraded functionality.
	// The failure is ignored and application still capable to deliver
	// incomplete but correct results.
	xWARN = slog.LevelWarn

	// NOTICE
	// system is failed, error is recovered, no impact
	xNOTICE = slog.Level(2)

	// INFO
	// output informative status about system
	xINFO = slog.LevelInfo

	// DEBUG
	// output debug status about system
	xDEBUG = slog.LevelDebug
)

// EMERGENCY
// system is unusable, panic execution of current routine/application,
// it is notpossible to gracefully terminate it.
func Emergency(err error, args ...any) {
	slog.Log(context.Background(), xEMERGENCY, err.Error(), args...)
	panic(err)
}

// CRITICAL
// system is failed, response actions must be taken immediately,
// the application is not able to execute correctly but still
// able to gracefully exit.
func Critical(err error, args ...any) {
	slog.Log(context.Background(), xCRITICAL, err.Error(), args...)
}

// ERROR
// system is failed, unable to recover from error.
// The failure do not have global catastrophic impacts but
// local functionality is impaired, incorrect result is returned.
func Error(err error, args ...any) {
	slog.Log(context.Background(), xERROR, err.Error(), args...)
}

// WARN
// system is failed, unable to recover, degraded functionality.
// The failure is ignored and application still capable to deliver
// incomplete but correct results.
func Warn(msg string, args ...any) {
	slog.Log(context.Background(), xWARN, msg, args...)
}

// NOTICE
// system is failed, error is recovered, no impact
func Notice(msg string, args ...any) {
	slog.Log(context.Background(), xNOTICE, msg, args...)
}

// INFO
// output informative status about system
func Info(msg string, args ...any) {
	slog.Log(context.Background(), xINFO, msg, args...)
}

// DEBUG
// output debug status about system
func Debug(msg string, args ...any) {
	slog.Log(context.Background(), xDEBUG, msg, args...)
}

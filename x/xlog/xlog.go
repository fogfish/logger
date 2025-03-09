//
// Copyright (C) 2021 - 2025 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/logger
//

package xlog

import (
	"context"
	"log/slog"
	"runtime"
	"time"

	"github.com/fogfish/logger/v3"
)

func log(
	ctx context.Context,
	logger *slog.Logger,
	level slog.Level,
	msg string,
	args ...any,
) {
	if !logger.Enabled(ctx, level) {
		return
	}
	var pc uintptr
	var pcs [1]uintptr
	// skip [runtime.Callers, this function, this function's caller]
	runtime.Callers(3, pcs[:])
	pc = pcs[0]
	r := slog.NewRecord(time.Now(), level, msg, pc)
	r.Add(args...)
	if ctx == nil {
		ctx = context.Background()
	}
	_ = logger.Handler().Handle(ctx, r)
}

// EMERGENCY
// system is unusable, panic execution of current routine/application,
// it is notpossible to gracefully terminate it.
func Emergency(msg string, err error, args ...any) {
	if err != nil {
		args = append(args, slog.Any("err", err))
	}
	log(context.Background(), slog.Default(), logger.EMERGENCY, msg, args...)
	panic(err)
}

// CRITICAL
// system is failed, response actions must be taken immediately,
// the application is not able to execute correctly but still
// able to gracefully exit.
func Critical(msg string, err error, args ...any) {
	if err != nil {
		args = append(args, slog.Any("err", err))
	}
	log(context.Background(), slog.Default(), logger.CRITICAL, msg, args...)
}

// ERROR
// system is failed, unable to recover from error.
// The failure do not have global catastrophic impacts but
// local functionality is impaired, incorrect result is returned.
func Error(msg string, err error, args ...any) {
	if err != nil {
		args = append(args, slog.Any("err", err))
	}
	log(context.Background(), slog.Default(), logger.ERROR, msg, args...)
}

// WARN
// system is failed, unable to recover, degraded functionality.
// The failure is ignored and application still capable to deliver
// incomplete but correct results.
func Warn(msg string, err error, args ...any) {
	if err != nil {
		args = append(args, slog.Any("err", err))
	}
	log(context.Background(), slog.Default(), logger.WARN, msg, args...)
}

// NOTICE
// system is failed, error is recovered, no impact
func Notice(msg string, args ...any) {
	log(context.Background(), slog.Default(), logger.NOTICE, msg, args...)
}

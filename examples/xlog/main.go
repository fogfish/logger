//
// Copyright (C) 2021 - 2025 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/logger
//

package main

import (
	"io"
	"log/slog"

	log "github.com/fogfish/logger/v3"
	"github.com/fogfish/logger/x/xlog"
)

func main() {
	slog.SetDefault(log.New(log.WithLogLevel(log.DEBUG)))

	obj := slog.Group("obj",
		slog.Any("key", "val"),
	)

	slog.Debug("debug status about system.", obj)
	slog.Info("informative status about system.", obj)
	xlog.Notice("system is failed, error is recovered, no impact.", obj)
	xlog.Warn("system is failed, unable to recover, degraded functionality.", io.EOF, obj)
	xlog.Error("system is failed, unable to recover from error.", io.EOF, obj)
	xlog.Critical("system is failed, response actions must be taken immediately. Doing graceful exit.", io.EOF, obj)
	xlog.Emergency("system is unusable, panic execution of app.", io.EOF, obj)
}

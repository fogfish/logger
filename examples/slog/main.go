//
// Copyright (C) 2021 - 2025 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/logger
//

package main

import (
	"context"
	"log/slog"

	log "github.com/fogfish/logger/v3"
)

func main() {
	slog.SetDefault(log.New(log.WithLogLevel(log.DEBUG)))

	obj := slog.Group("obj",
		slog.Any("key", "val"),
	)

	slog.Debug("debug status about system.", obj)
	slog.Info("informative status about system.", obj)
	slog.Log(context.Background(), log.NOTICE, "system is failed, error is recovered, no impact.", obj)
	slog.Warn("system is failed, unable to recover, degraded functionality.", obj)
	slog.Error("system is failed, unable to recover from error.", obj)
	slog.Log(context.Background(), log.CRITICAL, "system is failed, response actions must be taken immediately. Doing graceful exit.", obj)
	slog.Log(context.Background(), log.EMERGENCY, "system is unusable, panic execution of app.", obj)
}

//
// Copyright (C) 2021 Dmitry Kolesnikov
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
	slog.Debug("debug message", "key", "val")
	slog.Info("info message", "key", "val")

	slog.Log(context.Background(), log.EMERGENCY, "system emergency")
}

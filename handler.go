//
// Copyright (C) 2021 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/logger
//

package logger

import (
	"context"
	"log/slog"
	"runtime"
	"strings"

	"github.com/fogfish/logger/v3/internal/trie"
)

type handler struct {
	slog.Handler
	trie *trie.Node
}

func (h *handler) Enabled(context.Context, slog.Level) bool { return true }

func (h *handler) Handle(ctx context.Context, r slog.Record) error {
	if r.PC == 0 {
		return h.Handler.Handle(ctx, r)
	}

	fs := runtime.CallersFrames([]uintptr{r.PC})
	f, _ := fs.Next()

	parts := strings.Split(f.File, "go/src/")
	path := parts[0]
	if len(parts) > 1 {
		path = parts[1]
	}

	_, n := h.trie.Lookup(path)

	if len(n.Path) != 0 && n.Level <= r.Level {
		return h.Handler.Handle(ctx, r)
	}

	return nil
}

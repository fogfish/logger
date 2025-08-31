//
// Copyright (C) 2021 - 2025 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/logger
//

package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"runtime"
	"strings"
	"sync"

	"github.com/fogfish/logger/v3/internal/trie"
)

// JSON logger handler
func NewJSONHandler(opts ...Option) slog.Handler {
	config := defaultOpts(CloudWatch...)
	for _, opt := range opts {
		opt(config)
	}

	h := slog.NewJSONHandler(config.writer,
		&slog.HandlerOptions{
			AddSource:   config.addSource,
			Level:       config.level,
			ReplaceAttr: config.attributes.handle,
		},
	)

	if config.trie == nil {
		return h
	}

	return &modTrieHandler{
		Handler: h,
		trie:    config.trie,
	}
}

//------------------------------------------------------------------------------

// The handler perform module-based logging
type modTrieHandler struct {
	slog.Handler
	trie *trie.Node
}

func (h *modTrieHandler) Enabled(ctx context.Context, level slog.Level) bool {
	// For module-based filtering, we need to allow the handler to process the record
	// and make the filtering decision in Handle() where we have access to the source path
	return true
}

func (h *modTrieHandler) Handle(ctx context.Context, r slog.Record) error {
	if r.PC == 0 {
		return h.Handler.Handle(ctx, r)
	}

	fs := runtime.CallersFrames([]uintptr{r.PC})
	f, _ := fs.Next()

	// TODO: "go/pkg/mod"
	parts := strings.Split(f.File, "go/src/")
	path := parts[0]
	if len(parts) > 1 {
		path = parts[1]
	}

	_, n := h.trie.Lookup(path)

	// If a specific module rule is found, use that rule
	if len(n.Path) != 0 {
		if n.Level <= r.Level {
			return h.Handler.Handle(ctx, r)
		}
		// Module rule found but level doesn't match, don't log
		return nil
	}

	// No specific module rule found, fall back to default handler behavior
	// But we need to check if the underlying handler would accept this level
	if h.Handler.Enabled(ctx, r.Level) {
		return h.Handler.Handle(ctx, r)
	}

	return nil
}

//------------------------------------------------------------------------------

type stdioHandler struct {
	w io.Writer
	h slog.Handler
	b *bytes.Buffer
	m sync.Mutex
}

// Standard I/O handler
func NewStdioHandler(opts ...Option) slog.Handler {
	config := defaultOpts(Console...)
	for _, opt := range opts {
		opt(config)
	}

	b := &bytes.Buffer{}
	h := slog.NewJSONHandler(b,
		&slog.HandlerOptions{
			AddSource:   config.addSource,
			Level:       config.level,
			ReplaceAttr: config.attributes.handle,
		},
	)

	if config.trie == nil {
		return &stdioHandler{w: config.writer, b: b, h: h}
	}

	return &stdioHandler{w: config.writer, b: b,
		h: &modTrieHandler{
			Handler: h,
			trie:    config.trie,
		},
	}
}

func (h *stdioHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.h.Enabled(ctx, level)
}

func (h *stdioHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &stdioHandler{h: h.h.WithAttrs(attrs), b: h.b}
}

func (h *stdioHandler) WithGroup(name string) slog.Handler {
	return &stdioHandler{h: h.h.WithGroup(name), b: h.b}
}

func (h *stdioHandler) Handle(ctx context.Context, r slog.Record) error {
	attrs, err := h.computeAttrs(ctx, r)
	if err != nil {
		return err
	}

	// If attrs is nil, it means the message was filtered out
	if attrs == nil {
		return nil
	}

	time := attrs["time"]

	level := attrs["level"]
	cl := levelColorForText[r.Level]
	msg := cl + r.Message + colorReset

	delete(attrs, "level")
	delete(attrs, "time")
	delete(attrs, "msg")

	if len(attrs) == 0 {
		fmt.Fprintln(h.w, time, level, msg)
		return nil
	}

	bytes, err := json.MarshalIndent(attrs, "", "  ")
	if err != nil {
		return err
	}

	ca := levelColorForAttr[r.Level]
	obj := ca + string(bytes) + colorReset

	fmt.Fprintln(h.w, time, level, msg, obj)
	return nil
}

func (h *stdioHandler) computeAttrs(
	ctx context.Context,
	r slog.Record,
) (map[string]any, error) {
	h.m.Lock()
	defer func() {
		h.b.Reset()
		h.m.Unlock()
	}()

	if err := h.h.Handle(ctx, r); err != nil {
		return nil, err
	}

	// If buffer is empty, it means the message was filtered out by modTrieHandler
	if h.b.Len() == 0 {
		return nil, nil
	}

	var attrs map[string]any
	err := json.Unmarshal(h.b.Bytes(), &attrs)
	if err != nil {
		return nil, err
	}
	return attrs, nil
}

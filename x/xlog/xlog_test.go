//
// Copyright (C) 2021 - 2025 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/logger
//

package xlog_test

import (
	"bytes"
	"io"
	"log/slog"
	"strings"
	"testing"

	"github.com/fogfish/logger/v3"
	"github.com/fogfish/logger/x/xlog"
)

func TestStdioLogger(t *testing.T) {
	b := &bytes.Buffer{}
	slog.SetDefault(slog.New(logger.NewStdioHandler(logger.WithWriter(b))))

	t.Run("Notice", func(t *testing.T) {
		defer b.Reset()

		xlog.Notice("test")
		txt := b.String()
		if !strings.Contains(txt, "NTC") ||
			!strings.Contains(txt, "test") ||
			!strings.Contains(txt, "source") ||
			!strings.Contains(txt, "lggr") {
			t.Errorf("unexpected log line %s", txt)
		}
	})

	t.Run("Info", func(t *testing.T) {
		defer b.Reset()

		xlog.Warn("test", io.EOF)
		txt := b.String()
		if !strings.Contains(txt, "WRN") ||
			!strings.Contains(txt, "EOF") ||
			!strings.Contains(txt, "test") ||
			!strings.Contains(txt, "source") ||
			!strings.Contains(txt, "lggr") {
			t.Errorf("unexpected log line %s", txt)
		}
	})

	t.Run("Error", func(t *testing.T) {
		defer b.Reset()

		xlog.Error("test", io.EOF)
		txt := b.String()
		if !strings.Contains(txt, "ERR") ||
			!strings.Contains(txt, "EOF") ||
			!strings.Contains(txt, "test") ||
			!strings.Contains(txt, "source") ||
			!strings.Contains(txt, "lggr") {
			t.Errorf("unexpected log line %s", txt)
		}
	})

	t.Run("Critical", func(t *testing.T) {
		defer b.Reset()

		xlog.Critical("test", io.EOF)
		txt := b.String()
		if !strings.Contains(txt, "CRT") ||
			!strings.Contains(txt, "EOF") ||
			!strings.Contains(txt, "test") ||
			!strings.Contains(txt, "source") ||
			!strings.Contains(txt, "lggr") {
			t.Errorf("unexpected log line %s", txt)
		}
	})
}

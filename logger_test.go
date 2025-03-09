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
	"log/slog"
	"strings"
	"testing"
)

func TestStdioLogger(t *testing.T) {
	b := &bytes.Buffer{}
	log := slog.New(NewStdioHandler(WithWriter(b), WithLogLevel(DEBUG)))

	t.Run("Debug", func(t *testing.T) {
		defer b.Reset()

		log.Debug("test")
		txt := b.String()
		if !strings.Contains(txt, "DEB") ||
			!strings.Contains(txt, "test") ||
			!strings.Contains(txt, "source") ||
			!strings.Contains(txt, "lggr") {
			t.Errorf("unexpected log line %s", txt)
		}
	})

	t.Run("Info", func(t *testing.T) {
		defer b.Reset()

		log.Info("test")
		txt := b.String()
		if !strings.Contains(txt, "INF") ||
			!strings.Contains(txt, "test") ||
			!strings.Contains(txt, "source") ||
			!strings.Contains(txt, "lggr") {
			t.Errorf("unexpected log line %s", txt)
		}
	})

	t.Run("Warn", func(t *testing.T) {
		defer b.Reset()

		log.Warn("test")
		txt := b.String()
		if !strings.Contains(txt, "WRN") ||
			!strings.Contains(txt, "test") ||
			!strings.Contains(txt, "source") ||
			!strings.Contains(txt, "lggr") {
			t.Errorf("unexpected log line %s", txt)
		}
	})

	t.Run("Error", func(t *testing.T) {
		defer b.Reset()

		log.Error("test")
		txt := b.String()
		if !strings.Contains(txt, "ERR") ||
			!strings.Contains(txt, "test") ||
			!strings.Contains(txt, "source") ||
			!strings.Contains(txt, "lggr") {
			t.Errorf("unexpected log line %s", txt)
		}
	})
}

func TestJSONLogger(t *testing.T) {
	b := &bytes.Buffer{}
	log := slog.New(NewJSONHandler(WithWriter(b), WithLogLevel(DEBUG)))

	t.Run("Debug", func(t *testing.T) {
		defer b.Reset()

		log.Debug("test")
		txt := b.String()
		if !strings.Contains(txt, "DEBUG") ||
			!strings.Contains(txt, "test") ||
			!strings.Contains(txt, "source") ||
			!strings.Contains(txt, "lggr") {
			t.Errorf("unexpected log line %s", txt)
		}
	})

	t.Run("Info", func(t *testing.T) {
		defer b.Reset()

		log.Info("test")
		txt := b.String()
		if !strings.Contains(txt, "INFO") ||
			!strings.Contains(txt, "test") ||
			!strings.Contains(txt, "source") ||
			!strings.Contains(txt, "lggr") {
			t.Errorf("unexpected log line %s", txt)
		}
	})

	t.Run("Warn", func(t *testing.T) {
		defer b.Reset()

		log.Warn("test")
		txt := b.String()
		if !strings.Contains(txt, "WARN") ||
			!strings.Contains(txt, "test") ||
			!strings.Contains(txt, "source") ||
			!strings.Contains(txt, "lggr") {
			t.Errorf("unexpected log line %s", txt)
		}
	})

	t.Run("Error", func(t *testing.T) {
		defer b.Reset()

		log.Error("test")
		txt := b.String()
		if !strings.Contains(txt, "ERROR") ||
			!strings.Contains(txt, "test") ||
			!strings.Contains(txt, "source") ||
			!strings.Contains(txt, "lggr") {
			t.Errorf("unexpected log line %s", txt)
		}
	})
}

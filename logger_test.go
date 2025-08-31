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

func TestStdioLoggerWithTrie(t *testing.T) {
	b := &bytes.Buffer{}
	log := slog.New(NewStdioHandler(
		WithWriter(b),
		WithLogLevel(INFO), // Default level is INFO
		WithLogLevelForMod(map[string]slog.Level{
			"github.com/fogfish/logger": ERROR, // Only allow ERROR and above for this specific module
		}),
	))

	t.Run("FilteredByModuleRule", func(t *testing.T) {
		defer b.Reset()

		// This should be filtered out because logger module has ERROR level requirement
		log.Info("this info message should be filtered by module rule")
		txt := b.String()
		// Buffer should be empty since message was filtered
		if txt != "" {
			t.Errorf("expected empty buffer, got: %s", txt)
		}
	})

	t.Run("AllowedByModuleRule", func(t *testing.T) {
		defer b.Reset()

		// This should pass through because it meets the module's ERROR requirement
		log.Error("this error message should appear due to module rule")
		txt := b.String()
		if !strings.Contains(txt, "ERR") ||
			!strings.Contains(txt, "this error message should appear due to module rule") {
			t.Errorf("unexpected log line %s", txt)
		}
	})

	t.Run("DefaultBehaviorForUnspecifiedModule", func(t *testing.T) {
		defer b.Reset()

		// Test needs to simulate logging from a different module path
		// Since we can't easily change the source path in tests, let's create
		// a logger that doesn't match any trie rules
		logOther := slog.New(NewStdioHandler(
			WithWriter(b),
			WithLogLevel(INFO), // Default level is INFO
			WithLogLevelForMod(map[string]slog.Level{
				"some/other/module": ERROR, // Rule for a different module
			}),
		))

		// This should pass through because no specific rule exists and it meets default INFO level
		logOther.Info("this info message should appear due to default level")
		txt := b.String()
		if !strings.Contains(txt, "INF") ||
			!strings.Contains(txt, "this info message should appear due to default level") {
			t.Errorf("expected info message to appear with default behavior, got: %s", txt)
		}
	})

	t.Run("FilteredByDefaultLevel", func(t *testing.T) {
		defer b.Reset()

		// Test with a logger that has higher default level
		logStrict := slog.New(NewStdioHandler(
			WithWriter(b),
			WithLogLevel(ERROR), // Default level is ERROR
			WithLogLevelForMod(map[string]slog.Level{
				"some/other/module": DEBUG, // Rule for a different module (not this one)
			}),
		))

		// This should be filtered out because it doesn't meet default ERROR level
		logStrict.Info("this info message should be filtered by default level")
		txt := b.String()
		// Buffer should be empty since message was filtered by default level
		if txt != "" {
			t.Errorf("expected empty buffer due to default level filtering, got: %s", txt)
		}
	})
}

//
// Copyright (C) 2021 - 2025 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/logger
//

package logger

import (
	"log/slog"
	"path/filepath"
	"testing"
	"time"
)

func TestAttrLogTimeFormat(t *testing.T) {
	format := "2006-01-02 15:04:05"
	attr := attrLogTimeFormat(format)
	timestamp := time.Now()
	a := slog.Attr{Key: slog.TimeKey, Value: slog.AnyValue(timestamp)}

	result := attr(nil, a)
	expected := timestamp.Format(format)

	if result.Value.String() != expected {
		t.Errorf("expected %s, got %s", expected, result.Value.String())
	}
}

func TestAttrNoTimestamp(t *testing.T) {
	attr := attrNoTimestamp
	timestamp := time.Now()
	a := slog.Attr{Key: slog.TimeKey, Value: slog.AnyValue(timestamp)}

	result := attr(nil, a)
	if result.Key != "" {
		t.Errorf("expected empty key, got %s", result.Key)
	}

	result = attr([]string{"group"}, a)
	if result.Key != slog.TimeKey {
		t.Errorf("expected %s, got %s", slog.TimeKey, result.Key)
	}
}

func TestAttrLogLevel7(t *testing.T) {
	attr := attrLogLevel7(false)
	level := slog.LevelInfo
	a := slog.Attr{Key: slog.LevelKey, Value: slog.AnyValue(level)}

	result := attr(nil, a)
	expected := levelLongName[level]

	if result.Value.String() != expected {
		t.Errorf("expected %s, got %s", expected, result.Value.String())
	}
}

func TestAttrLogLevel7Shorten(t *testing.T) {
	attr := attrLogLevel7Shorten(false)
	level := slog.LevelInfo
	a := slog.Attr{Key: slog.LevelKey, Value: slog.AnyValue(level)}

	result := attr(nil, a)
	expected := levelShortName[level]

	if result.Value.String() != expected {
		t.Errorf("expected %s, got %s", expected, result.Value.String())
	}
}

func TestAttrSourceFileName(t *testing.T) {
	attr := attrSourceFileName
	source := &slog.Source{File: "/path/to/file.go"}
	a := slog.Attr{Key: slog.SourceKey, Value: slog.AnyValue(source)}

	result := attr(nil, a)
	expected := filepath.Base(source.File)

	if result.Value.Any().(*slog.Source).File != expected {
		t.Errorf("expected %s, got %s", expected, result.Value.Any().(*slog.Source).File)
	}
}

func TestAttrSourceShorten(t *testing.T) {
	attr := attrSourceShorten
	source := &slog.Source{File: "/path/to/go/src/github.com/fogfish/logger/attributes.go", Function: "github.com/fogfish/logger.TestFunc"}
	a := slog.Attr{Key: slog.SourceKey, Value: slog.AnyValue(source)}

	result := attr(nil, a)
	expectedFile := "gthb.fgfs.lggr/attributes.go"
	expectedFunc := "gthb.fgfs/logger.TestFunc"

	if result.Value.Any().(*slog.Source).File != expectedFile {
		t.Errorf("expected %s, got %s", expectedFile, result.Value.Any().(*slog.Source).File)
	}

	if result.Value.Any().(*slog.Source).Function != expectedFunc {
		t.Errorf("expected %s, got %s", expectedFunc, result.Value.Any().(*slog.Source).Function)
	}
}

func TestShorten(t *testing.T) {
	path := "github.com/fogfish/logger/attributes.go"
	expected := "gthb.fgfs.lggr/attributes.go"

	result := shorten(path)
	if result != expected {
		t.Errorf("expected %s, got %s", expected, result)
	}
}

func TestShortenSegment(t *testing.T) {
	segment := "logger"
	expected := "lggr"

	result := shortenSegment(segment)
	if result != expected {
		t.Errorf("expected %s, got %s", expected, result)
	}
}

func TestRemoveVowels(t *testing.T) {
	segment := "logger"
	expected := "lggr"

	result := removeVowels(segment)
	if result != expected {
		t.Errorf("expected %s, got %s", expected, result)
	}
}

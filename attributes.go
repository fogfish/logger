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
	"strings"
	"time"
)

// Log attributes, forammting combinator
type Attributes []func(groups []string, a slog.Attr) slog.Attr

func (attrs Attributes) handle(groups []string, a slog.Attr) slog.Attr {
	for _, f := range attrs {
		a = f(groups, a)
	}
	return a
}

//------------------------------------------------------------------------------

// Logs timestamp using the time format string
func attrLogTimeFormat(format string) func([]string, slog.Attr) slog.Attr {
	return func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.TimeKey {
			t := a.Value.Any().(time.Time)
			return slog.String(slog.TimeKey, t.Format(format))
		}

		return a
	}
}

// Skip timestamp from log
func attrNoTimestamp(groups []string, a slog.Attr) slog.Attr {
	if a.Key == slog.TimeKey && len(groups) == 0 {
		return slog.Attr{}
	}

	return a
}

//------------------------------------------------------------------------------

// Logs longs level names using 7-levels and colors
func attrLogLevel7(color bool) func([]string, slog.Attr) slog.Attr {
	return func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.LevelKey {
			lvl := a.Value.Any().(slog.Level)

			name, has := levelLongName[lvl]
			if !has {
				return a
			}

			if !color {
				return slog.String(slog.LevelKey, name)
			}

			color := levelColorForName[lvl]
			return slog.String(slog.LevelKey, color+name+colorReset)
		}

		return a
	}
}

// Logs short level names (3 chars) using 7-levels and colors
func attrLogLevel7Shorten(color bool) func([]string, slog.Attr) slog.Attr {
	return func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.LevelKey {
			lvl, ok := a.Value.Any().(slog.Level)
			if !ok {
				return a
			}

			name, has := levelShortName[lvl]
			if !has {
				return a
			}

			if !color {
				return slog.String(slog.LevelKey, name)
			}

			color := levelColorForName[lvl]
			return slog.String(slog.LevelKey, color+name+colorReset)
		}

		return a
	}
}

//------------------------------------------------------------------------------

// Logs only name of the file instead of full path
func attrSourceFileName(groups []string, a slog.Attr) slog.Attr {
	if a.Key == slog.SourceKey {
		source, _ := a.Value.Any().(*slog.Source)
		if source != nil {
			source.File = filepath.Base(source.File)
		}
	}

	return a
}

// Logs full path to the source file, applying shortening heuristic
func attrSourceShorten(groups []string, a slog.Attr) slog.Attr {
	if a.Key == slog.SourceKey {
		source, _ := a.Value.Any().(*slog.Source)
		if source != nil {
			parts := strings.Split(source.File, "go/src/")
			path := parts[0]
			if len(parts) > 1 {
				path = parts[1]
			}
			source.File = shorten(path)
			source.Function = shorten(source.Function)
		}
	}

	return a
}

func shorten(path string) string {
	seq := strings.Split(path, string(filepath.Separator))
	if len(seq) == 1 {
		return path
	}

	for i := 0; i < len(seq)-1; i++ {
		seq[i] = shortenSegment(seq[i])
	}
	return strings.Join(seq[0:len(seq)-1], ".") + "/" + seq[len(seq)-1]
}

// shortenSegment heuristically shortens a path segment
func shortenSegment(segment string) string {
	length := len(segment)
	switch {
	case length <= 3:
		return segment // Keep short names unchanged
	case length <= 6:
		return removeVowels(segment) // Remove vowels from medium-length names
	default:
		short := removeVowels(segment)
		if len(short) > 4 {
			return short[:4]
		}
		return short
	}
}

func removeVowels(segment string) string {
	vowels := "aeiouAEIOU"
	keep := []rune{}
	for i, r := range segment {
		if i == 0 || !strings.ContainsRune(vowels, r) {
			keep = append(keep, r)
		}
	}
	if len(keep) < 2 { // Ensure at least 2 characters remain
		return segment
	}
	return string(keep)
}

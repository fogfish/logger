//
// Copyright (C) 2021 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/logger
//

// Package logger configures slog for AWS CloudWatch
package logger

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

func init() {
	slog.SetDefault(
		Config(
			WithLogLevelFromEnv(),
		),
	)
}

const (
	// EMERGENCY
	// system is unusable, panic execution of current routine/application,
	// it is notpossible to gracefully terminate it.
	EMERGENCY = slog.Level(100)

	// CRITICAL
	// system is failed, response actions must be taken immediately,
	// the application is not able to execute correctly but still
	// able to gracefully exit.
	CRITICAL = slog.Level(50)

	// ERROR
	// system is failed, unable to recover from error.
	// The failure do not have global catastrophic impacts but
	// local functionality is impaired, incorrect result is returned.
	ERROR = slog.LevelError

	// WARN
	// system is failed, unable to recover, degraded functionality.
	// The failure is ignored and application still capable to deliver
	// incomplete but correct results.
	WARN = slog.LevelWarn

	// NOTICE
	// system is failed, error is recovered, no impact
	NOTICE = slog.Level(2)

	// INFO
	// output informative status about system
	INFO = slog.LevelInfo

	// DEBUG
	// output debug status about system
	DEBUG = slog.LevelDebug
)

var (
	shortLevel = map[slog.Leveler]string{
		EMERGENCY: "EMR",
		CRITICAL:  "CRT",
		ERROR:     "ERR",
		WARN:      "WRN",
		NOTICE:    "NTC",
		INFO:      "INF",
		DEBUG:     "DEB",
	}

	longLevel = map[slog.Leveler]string{
		EMERGENCY: "EMERGENCY",
		CRITICAL:  "CRITICAL",
		ERROR:     "ERROR",
		WARN:      "WARN",
		NOTICE:    "NOTICE",
		INFO:      "INFO",
		DEBUG:     "DEBUG",
	}

	longNames = map[string]slog.Level{
		"EMERGENCY": EMERGENCY,
		"CRITICAL":  CRITICAL,
		"ERROR":     ERROR,
		"WARN":      WARN,
		"NOTICE":    NOTICE,
		"INFO":      INFO,
		"DEBUG":     DEBUG,
	}
)

type opts struct {
	writer       io.Writer
	level        slog.Leveler
	shortLevel   bool
	noTimestamp  bool
	filenameOnly bool
	shortestPath bool
}

func defaultOpts() *opts {
	return &opts{
		writer:       os.Stdout,
		level:        INFO,
		shortLevel:   false,
		noTimestamp:  true,
		filenameOnly: false,
		shortestPath: true,
	}
}

func (opts opts) fmt(groups []string, a slog.Attr) slog.Attr {
	if a.Key == slog.LevelKey && opts.shortLevel {
		lvl := a.Value.Any().(slog.Level)
		if name, has := shortLevel[lvl]; has {
			return slog.String(slog.LevelKey, name)
		}
	}

	if a.Key == slog.LevelKey {
		switch a.Value.Any() {
		case EMERGENCY:
			return slog.String(slog.LevelKey, longLevel[EMERGENCY])
		case CRITICAL:
			return slog.String(slog.LevelKey, longLevel[CRITICAL])
		case NOTICE:
			return slog.String(slog.LevelKey, longLevel[NOTICE])
		default:
			return a
		}
	}

	if a.Key == slog.TimeKey && len(groups) == 0 && opts.noTimestamp {
		return slog.Attr{}
	}

	if a.Key == slog.SourceKey && opts.filenameOnly {
		source, _ := a.Value.Any().(*slog.Source)
		if source != nil {
			source.File = filepath.Base(source.File)
		}
		return a
	}

	if a.Key == slog.SourceKey && opts.shortestPath {
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
		return a
	}

	return a
}

func shorten(path string) string {
	seq := strings.Split(path, string(filepath.Separator))
	if len(seq) == 1 {
		return path
	}

	for i := 0; i < len(seq)-1; i++ {
		if len(seq[i]) > 0 {
			seq[i] = strings.ToLower(seq[i][0:1])
		}
	}
	return strings.Join(seq[0:len(seq)-1], ".") + "/" + seq[len(seq)-1]
}

// Config options for slog
type Option func(*opts)

// Config Log Level, default INFO
func WithLogLevel(level slog.Leveler) Option {
	return func(o *opts) {
		o.level = level
	}
}

// Config Log Level from env CONFIG_LOG_LEVEL, default INFO
func WithLogLevelFromEnv() Option {
	return func(o *opts) {
		level, defined := os.LookupEnv("CONFIG_LOG_LEVEL")
		if !defined {
			return
		}

		if lvl, has := longNames[level]; has {
			o.level = lvl
		}
	}
}

// Config Log writer, default os.Stdout
func WithWriter(w io.Writer) Option {
	return func(o *opts) {
		o.writer = w
	}
}

// Config Log Level to be 3 letters only
func WithShortLogLevel() Option {
	return func(o *opts) {
		o.shortLevel = true
	}
}

// Enables timestamp in the log messages
func WithTimestamp() Option {
	return func(o *opts) {
		o.noTimestamp = false
	}
}

// Logs absolute path for the source file
func WithFilePath() Option {
	return func(o *opts) {
		o.filenameOnly = false
		o.shortestPath = false
	}
}

// Logs file name of the source file
func WithFileName() Option {
	return func(o *opts) {
		o.filenameOnly = true
		o.shortestPath = false
	}
}

// Config logger
func Config(opts ...Option) *slog.Logger {
	config := defaultOpts()
	for _, opt := range opts {
		opt(config)
	}

	return slog.New(
		slog.NewJSONHandler(os.Stdout,
			&slog.HandlerOptions{
				AddSource:   true,
				Level:       slog.LevelDebug,
				ReplaceAttr: config.fmt,
			},
		),
	)
}

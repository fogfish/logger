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

	"github.com/fogfish/logger/v3/internal/trie"
)

func init() {
	// Config Optimal for logging of serverless with AWS Cloud Watch
	slog.SetDefault(
		New(
			WithLogLevelFromEnv(),
			WithLogLevel7(),
			WithSourceShorten(),
			WithoutTimestamp(),
			WithLogLevelForModFromEnv(),
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
	shortLevel = map[slog.Level]string{
		EMERGENCY: "EMR",
		CRITICAL:  "CRT",
		ERROR:     "ERR",
		WARN:      "WRN",
		NOTICE:    "NTC",
		INFO:      "INF",
		DEBUG:     "DEB",
	}

	longLevel = map[slog.Level]string{
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

// combinator of attribute formatting
type Attributes []func(groups []string, a slog.Attr) slog.Attr

func (attrs Attributes) handle(groups []string, a slog.Attr) slog.Attr {
	for _, f := range attrs {
		a = f(groups, a)
	}
	return a
}

type opts struct {
	writer     io.Writer
	level      slog.Leveler
	attributes Attributes
	addSource  bool
	trie       *trie.Node
}

func defaultOpts() *opts {
	return &opts{
		writer:     os.Stdout,
		level:      INFO,
		attributes: Attributes{},
		addSource:  false,
	}
}

// Config options for slog
type Option func(*opts)

// Config Log writer, default os.Stdout
func WithWriter(w io.Writer) Option {
	return func(o *opts) {
		o.writer = w
	}
}

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

// WithLevel7 enables from DEBUG to EMERGENCY levels
func WithLogLevel7() Option {
	return func(o *opts) {
		o.attributes = append(o.attributes, attrLogLevel7)
	}
}

func attrLogLevel7(groups []string, a slog.Attr) slog.Attr {
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

	return a
}

// Config Log Level to be 3 letters only
func WithLogLevelShorten() Option {
	return func(o *opts) {
		o.attributes = append(o.attributes, attrLogLevelShorten)
	}
}

func attrLogLevelShorten(groups []string, a slog.Attr) slog.Attr {
	if a.Key == slog.LevelKey {
		lvl := a.Value.Any().(slog.Level)
		if name, has := shortLevel[lvl]; has {
			return slog.String(slog.LevelKey, name)
		}
	}
	return a
}

// Config Log Levels per module
func WithLogLevelForMod(mods map[string]slog.Level) Option {
	return func(o *opts) {
		o.trie = trie.New()
		for mod, lvl := range mods {
			o.trie.Append(mod, lvl)
		}
	}
}

// Config Log Levels per module from env variables CONFIG_LOG_LEVEL_{NAME}
//
// CONFIG_LOG_LEVEL_DEBUG=github.com/fogfish/logger/*:github.com/your/app
func WithLogLevelForModFromEnv() Option {
	return func(o *opts) {
		root := trie.New()
		for lvl, name := range longLevel {
			for _, mod := range fromEnvMods("CONFIG_LOG_LEVEL_" + name) {
				root.Append(mod, lvl)
			}
		}

		if len(root.Heir) != 0 {
			o.trie = root
		}
	}
}

func fromEnvMods(key string) []string {
	value, defined := os.LookupEnv(key)
	if !defined {
		return nil
	}

	return strings.Split(value, ":")
}

// Exclude timestamp, required by CloudWatch
func WithoutTimestamp() Option {
	return func(o *opts) {
		o.attributes = append(o.attributes, attrNoTimestamp)
	}
}

func attrNoTimestamp(groups []string, a slog.Attr) slog.Attr {
	if a.Key == slog.TimeKey && len(groups) == 0 {
		return slog.Attr{}
	}

	return a
}

// Logs file name of the source file
func WithSourceFileName() Option {
	return func(o *opts) {
		o.addSource = true
		o.attributes = append(o.attributes, attrSourceFileName)
	}
}

func attrSourceFileName(groups []string, a slog.Attr) slog.Attr {
	if a.Key == slog.SourceKey {
		source, _ := a.Value.Any().(*slog.Source)
		if source != nil {
			source.File = filepath.Base(source.File)
		}
	}

	return a
}

// Shorten Source file to letters only
func WithSourceShorten() Option {
	return func(o *opts) {
		o.addSource = true
		o.attributes = append(o.attributes, attrSourceShorten)
	}
}

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
		if len(seq[i]) > 0 {
			seq[i] = strings.ToLower(seq[i][0:1])
		}
	}
	return strings.Join(seq[0:len(seq)-1], ".") + "/" + seq[len(seq)-1]
}

// Enable logging of source file
func WithSource() Option {
	return func(o *opts) {
		o.addSource = true
	}
}

// Create New Logger
func New(opts ...Option) *slog.Logger {
	return slog.New(NewJSONHandler(opts...))
}

// Create's new handler
func NewJSONHandler(opts ...Option) slog.Handler {
	config := defaultOpts()
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

	return &handler{
		Handler: h,
		trie:    config.trie,
	}
}

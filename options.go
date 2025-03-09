//
// Copyright (C) 2021 - 2025 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/logger
//

package logger

import (
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/fogfish/logger/v3/internal/trie"
)

var (
	// Preset for console logging
	Console = []Option{
		func(o *opts) { o.attributes = append(o.attributes, attrLogLevel7Shorten(true)) },
		WithTimeFormat("[15:04:05.000]"),
		WithLogLevel(INFO),
		WithLogLevelFromEnv(),
		WithSourceShorten(),
		WithLogLevelForModFromEnv(),
	}

	// Preset for CloudWatch logging
	CloudWatch = []Option{
		WithLogLevelFromEnv(),
		WithLogLevel7(),
		WithSourceShorten(),
		WithoutTimestamp(),
		WithLogLevelForModFromEnv(),
	}
)

// The logger config option
type Option func(*opts)

type opts struct {
	writer     io.Writer
	level      slog.Leveler
	attributes Attributes
	addSource  bool
	trie       *trie.Node
}

func defaultOpts(preset ...Option) *opts {
	opt := &opts{
		writer:     os.Stdout,
		level:      INFO,
		attributes: Attributes{},
		addSource:  false,
	}
	for _, f := range preset {
		f(opt)
	}

	return opt
}

// Config Log writer, default os.Stdout
func WithWriter(w io.Writer) Option {
	return func(o *opts) {
		o.writer = w
	}
}

// Exclude timestamp, required by CloudWatch
func WithoutTimestamp() Option {
	return func(o *opts) {
		o.attributes = append(o.attributes, attrNoTimestamp)
	}
}

// Configure the time format
func WithTimeFormat(format string) Option {
	return func(o *opts) {
		o.attributes = append(o.attributes, attrLogTimeFormat(format))
	}
}

// Config Log Level, default INFO
func WithLogLevel(level slog.Leveler) Option {
	return func(o *opts) {
		o.level = level
	}
}

// Config Log Level from env CONFIG_LOG_LEVEL, default INFO
//
//	export CONFIG_LOG_LEVEL=DEBUG
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
		o.attributes = append(o.attributes, attrLogLevel7(false))
	}
}

// Config Log Level to be 3 letters only from DEB to EMR
func WithLogLevelShorten() Option {
	return func(o *opts) {
		o.attributes = append(o.attributes, attrLogLevel7Shorten(false))
	}
}

// The logger allows you to define log levels for different modules with
// flexible granularity. Log levels can be set explicitly using this configuration
// options.
//
//	log.WithLogLevelForMod(map[string]slog.Level{
//		"github.com/fogfish/logger": log.INFO,
//		"github.com/you/application": log.DEBUG,
//	})
//
// The logger uses prefix matching to determine the appropriate log level based
// on the source code path:
//
// * Per File: A log level defined for a specific file
// (e.g., github.com/fogfish/logger/logger.go) applies only to that file.
//
// * Per Module: A log level set for a module (e.g., github.com/fogfish/logger)
// applies to all files within that module.
//
// * Per User/Namespace: A log level defined at a higher level
// (e.g., github.com/fogfish) applies to all modules under that namespace.
func WithLogLevelForMod(mods map[string]slog.Level) Option {
	return func(o *opts) {
		o.trie = trie.New()
		for mod, lvl := range mods {
			o.trie.Append(mod, lvl)
		}
	}
}

// The logger allows you to define log levels for different modules with
// flexible granularity. Log levels can be set explicitly using environment
// variables CONFIG_LOG_LEVEL_{NAME}:
//
//	export CONFIG_LOG_LEVEL_DEBUG=github.com/fogfish/logger:github.com/your/app
//	export CONFIG_LOG_LEVEL_INFO=github.com/fogfish
//
// * Per File: A log level defined for a specific file
// (e.g., github.com/fogfish/logger/logger.go) applies only to that file.
//
// * Per Module: A log level set for a module (e.g., github.com/fogfish/logger)
// applies to all files within that module.
//
// * Per User/Namespace: A log level defined at a higher level
// (e.g., github.com/fogfish) applies to all modules under that namespace.

func WithLogLevelForModFromEnv() Option {
	return func(o *opts) {
		root := trie.New()
		for lvl, name := range levelLongName {
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

// Logs file name of the source file only
func WithSourceFileName() Option {
	return func(o *opts) {
		o.addSource = true
		o.attributes = append(o.attributes, attrSourceFileName)
	}
}

// Shorten Source file to few letters only
func WithSourceShorten() Option {
	return func(o *opts) {
		o.addSource = true
		o.attributes = append(o.attributes, attrSourceShorten)
	}
}

// Enable default logging of source file
func WithSource() Option {
	return func(o *opts) {
		o.addSource = true
	}
}

// Disable logging of source file
func WithoutSource() Option {
	return func(o *opts) {
		o.addSource = false
	}
}

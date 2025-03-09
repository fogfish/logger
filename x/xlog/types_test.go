//
// Copyright (C) 2021 - 2024 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/logger
//

package xlog_test

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"regexp"
	"testing"
	"time"

	"github.com/fogfish/logger/x/xlog"
)

func TestSince(t *testing.T) {
	itShouldBe(t, "since", "10[0-9].[0-9]+ms",
		func(log *slog.Logger) {
			v := xlog.SinceNow()
			time.Sleep(100 * time.Millisecond)
			log.Info("", "since", v)
		},
	)
}

func TestPerSecond(t *testing.T) {
	itShouldBe(t, "rate", "1[0-9][0-9][0-9][0-9].[0-9]+",
		func(log *slog.Logger) {
			v := xlog.PerSecondNow()
			v.Acc += 1000
			time.Sleep(90 * time.Millisecond)
			log.Info("", "rate", v)
		},
	)
}

func TestMillisecondOp(t *testing.T) {
	itShouldBe(t, "ms/op", "1[0-9][0-9]",
		func(log *slog.Logger) {
			v := xlog.MillisecondOpNow()
			v.Acc += 1
			time.Sleep(100 * time.Millisecond)
			log.Info("", "ms/op", v)
		},
	)
}

//------------------------------------------------------------------------------

func itShouldBe(t *testing.T, key string, expected string, f func(log *slog.Logger)) {
	t.Helper()

	buf := &bytes.Buffer{}
	log := slog.New(slog.NewJSONHandler(buf, nil))

	f(log)

	var val map[string]any
	if err := json.Unmarshal(buf.Bytes(), &val); err != nil {
		t.Errorf("invalid format, json expected %v", buf.String())
	}

	v, has := val[key]
	if !has {
		t.Errorf("no `%s` found in %v", key, val)
	}

	if matched, err := regexp.MatchString(expected, v.(string)); err != nil || !matched {
		t.Errorf("invalid `%s` value %v, expected %s", key, v, expected)
	}
}

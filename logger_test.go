//
// Copyright (C) 2021 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/logger
//

package logger

import "testing"

func TestShorten(t *testing.T) {
	for in, expect := range map[string]string{
		"afoo/bfoo/cfoo/def": "a.b.c/def",
		"a/b/c/d":            "a.b.c/d",
		"bfoo/cfoo/def":      "b.c/def",
		"cfoo/def":           "c/def",
		"def":                "def",
		"main.main":          "main.main",
		"github.com/fogfish/logger/internal/codec.(*CodecRaw).Process": "g.f.l.i/codec.(*CodecRaw).Process",
	} {
		if shorten(in) != expect {
			t.Errorf("%s not shorten to %s", in, expect)
		}
	}
}

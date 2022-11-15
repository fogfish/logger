//
// Copyright (C) 2021 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/logger
//

package logger

import "strings"

type Message struct {
	key string
	val string
	ext string
}

func (msg Message) toJSON(buf *strings.Builder) {
	if len(msg.key) > 0 {
		buf.WriteString(`"`)
		buf.WriteString(msg.key)
		buf.WriteString(`": `)
	}
	buf.WriteString(msg.ext)
	buf.WriteString(msg.val)
	buf.WriteString(msg.ext)
}

func (msg Message) len() int {
	return len(msg.key) + len(msg.val) + 2*len(msg.ext) + 4
}

type Messages []Message

func (seq Messages) toJSON(buf *strings.Builder) {
	if len(seq) == 0 {
		return
	}

	size := 0
	for i := 0; i < len(seq); i++ {
		size = size + seq[i].len() + 2
	}

	buf.Grow(size)

	if len(seq) > 0 {
		seq[0].toJSON(buf)
		for i := 1; i < len(seq); i++ {
			buf.WriteString(`, `)
			seq[i].toJSON(buf)
		}
	}
}

func (seq Messages) toString(level string) string {
	buf := strings.Builder{}
	buf.WriteString(`{"lvl": "`)
	buf.WriteString(level)
	buf.WriteString(`", `)
	seq.toJSON(&buf)
	buf.WriteString(`}`)
	return buf.String()
}

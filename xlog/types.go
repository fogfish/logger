//
// Copyright (C) 2021 - 2024 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/logger
//

package xlog

import (
	"encoding/json"
	"strconv"
	"time"
)

// Since logs duration passed since given time
type Since time.Time

func (s Since) String() string               { return time.Since(time.Time(s)).String() }
func (s Since) MarshalJSON() ([]byte, error) { return json.Marshal(s.String()) }

// SinceNow return instance of [Since] type initialized with [time.Now]
func SinceNow() Since { return Since(time.Now()) }

//------------------------------------------------------------------------------

// PerSecond logs rates per second since given time
type PerSecond struct {
	T   time.Time
	Acc int
}

func (r PerSecond) String() string {
	v := float64(r.Acc) / time.Since(r.T).Seconds()
	return strconv.FormatFloat(v, 'f', 4, 64)
}

func (r PerSecond) MarshalJSON() ([]byte, error) { return json.Marshal(r.String()) }

// PerSecondNow return instance of [SecondNow] type initialized with [time.Now]
func PerSecondNow() PerSecond { return PerSecond{T: time.Now(), Acc: 0} }

//------------------------------------------------------------------------------

// MillisecondOp logs millisecond required for operation
type MillisecondOp struct {
	T   time.Time
	Acc int
}

func (r MillisecondOp) String() string {
	v := time.Since(r.T).Milliseconds() / int64(r.Acc)
	return strconv.Itoa(int(v))
}

func (r MillisecondOp) MarshalJSON() ([]byte, error) { return json.Marshal(r.String()) }

// MillisecondOpNow return instance of [MillisecondOp] type initialized with [time.Now]
func MillisecondOpNow() MillisecondOp { return MillisecondOp{T: time.Now(), Acc: 0} }

package main

import (
	"log/slog"
	"time"

	log "github.com/fogfish/logger/v3"
	"github.com/fogfish/logger/x/xlog"
)

func do() {
	defer slog.Info("done something", slog.Any("duration", xlog.SinceNow()))

	time.Sleep(100 * time.Millisecond)
}

func rate() {
	ops := xlog.PerSecondNow()
	defer slog.Info("done something", slog.Any("op/sec", ops))

	time.Sleep(100 * time.Millisecond)
	ops.Acc += 100
}

func demand() {
	ops := xlog.MillisecondOpNow()
	defer slog.Info("done something", slog.Any("ms/op", ops))

	time.Sleep(100 * time.Millisecond)
	ops.Acc += 10
}

func main() {
	slog.SetDefault(log.New())
	do()
	rate()
	demand()
}

package main

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

var (
	// This timestamp format differs from RFC3339Nano by using .000 instead
	// of .999999999 which changes the timestamp from 9 variable to 3 fixed
	// decimals (.130 instead of .130987456).
	timestampFormat = log.TimestampFormat(
		func() time.Time { return time.Now().UTC() },
		"2006-01-02T15:04:05.000Z07:00",
	)
)

func info() string {
	return fmt.Sprintf("(os=%s, arch=%s, version=%s)", runtime.GOOS,
		runtime.GOARCH, runtime.Version())
}

func main() {
	var l log.Logger
	// logfmt
	l = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	// json
	// l = log.NewJSONLogger(log.NewSyncWriter(os.Stderr))
	l = log.With(l, "ts", timestampFormat, "caller", log.DefaultCaller)
	level.Info(l).Log("msg", "Start testing")
	level.Info(l).Log("msg", "Info", "build", info())
}

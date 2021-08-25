// Copyright 2021 Kien Nguyen-Tuan <kiennt2609@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
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

func info() []string {
	return []string{"os", runtime.GOOS, "arch", runtime.GOARCH, "version", runtime.Version()}
}

func infoWithMsg(msg, info []string) []interface{} {
	var r []interface{}
	for _, v := range msg {
		r = append(r, v)
	}
	for _, v := range info {
		r = append(r, v)
	}
	return r
}

func main() {
	var l log.Logger
	// logfmt
	l = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	// json
	// l = log.NewJSONLogger(log.NewSyncWriter(os.Stderr))
	l = log.With(l, "ts", timestampFormat, "caller", log.DefaultCaller)
	level.Info(l).Log("msg", "Start testing")
	level.Info(l).Log(infoWithMsg([]string{"msg", "Info"}, info())...)
}

// Copyright 2022 Kien Nguyen-Tuan <kiennt2609@gmail.com>
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
	"fmt"
	"io"
	"os"

	"golang.org/x/term"
)

func main() {
	fd := int(os.Stdin.Fd())
	if !term.IsTerminal(fd) {
		panic("stdin/stdout should be terminal")
	}
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		panic(err)
	}
	defer term.Restore(fd, oldState)
	screen := struct {
		io.Reader
		io.Writer
	}{os.Stdin, os.Stdout}
	t := term.NewTerminal(screen, "")
	t.SetPrompt(string(t.Escape.Red) + ">" + string(t.Escape.Reset))
	rePrefix := string(t.Escape.Cyan) + "You type:" + string(t.Escape.Reset)

	for {
		line, err := t.ReadLine()
		if err != nil && err != io.EOF {
			panic(err)
		}
		if err == io.EOF {
			fmt.Fprintln(t, string(t.Escape.Green)+"Bye bye!")
			break
		}
		if line == "" {
			continue
		}
		fmt.Fprintln(t, rePrefix, line)
	}
}

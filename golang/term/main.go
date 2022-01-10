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

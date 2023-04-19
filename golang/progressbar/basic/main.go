package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/schollz/progressbar/v3"
)

func main() {
	req, _ := http.NewRequest("GET", "https://github.com/ryanoasis/nerd-fonts/releases/download/v2.3.3/VictorMono.zip", nil)
	resp, _ := http.DefaultClient.Do(req)
	defer check(resp.Body.Close)

	f, _ := os.OpenFile("VictorMono.zip", os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()

	bar := progressbar.DefaultBytes(
		resp.ContentLength,
		"downloading",
	)

	io.Copy(io.MultiWriter(f, bar), resp.Body)
}

// check checks the returned error of a function.
func check(f func() error) {
	if err := f(); err != nil {
		fmt.Fprintf(os.Stderr, "received error: %v\n", err)
	}
}

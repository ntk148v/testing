// Stole from here: https://gist.github.com/IceflowRE/e4c2b9163a697105a3e72f35f0cd12a5
package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/schollz/progressbar/v3"
)

type LineWriter struct {
	*MultiProgressBar
	id int
}

func (lw *LineWriter) Write(p []byte) (n int, err error) {
	lw.guard.Lock()
	defer lw.guard.Unlock()
	lw.Move(lw.id, lw.output)
	return lw.output.Write(p)
}

type MultiProgressBar struct {
	output  io.Writer
	curLine int
	Bars    []*progressbar.ProgressBar
	guard   sync.Mutex
}

func NewMultiProgressBar(pBars []*progressbar.ProgressBar, output io.Writer) *MultiProgressBar {
	mpb := &MultiProgressBar{
		curLine: 0,
		Bars:    pBars,
		guard:   sync.Mutex{},
		output:  output,
	}
	for id, pb := range mpb.Bars {
		progressbar.OptionSetWriter(&LineWriter{
			MultiProgressBar: mpb,
			id:               id,
		})(pb)
	}

	return mpb
}

func (mpb *MultiProgressBar) Add(b *progressbar.ProgressBar) {
	mpb.Bars = append(mpb.Bars, b)
	id := len(mpb.Bars) - 1
	progressbar.OptionSetWriter(&LineWriter{
		MultiProgressBar: mpb,
		id:               id,
	})(b)
}

// Move cursor to the beginning of the current progressbar.
func (mpb *MultiProgressBar) Move(id int, writer io.Writer) (int, error) {
	bias := mpb.curLine - id
	mpb.curLine = id
	if bias > 0 {
		// move up
		return fmt.Fprintf(writer, "\r\033[%dA", bias)
	} else if bias < 0 {
		// move down
		return fmt.Fprintf(writer, "\r\033[%dB", -bias)
	}
	return 0, nil
}

// End Move cursor to the end of the Progressbars.
func (mpb *MultiProgressBar) End() {
	mpb.Move(len(mpb.Bars), mpb.output)
}

func main() {
	var (
		urls = []string{
			"https://github.com/ryanoasis/nerd-fonts/releases/download/v2.3.3/Ubuntu.zip",
			"https://github.com/ryanoasis/nerd-fonts/releases/download/v2.3.3/VictorMono.zip",
			"https://github.com/ryanoasis/nerd-fonts/releases/download/v2.3.3/iA-Writer.zip",
		}
		wg sync.WaitGroup
	)

	mpb := NewMultiProgressBar(make([]*progressbar.ProgressBar, 0), os.Stderr)

	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			// download
			resp, _ := http.Get(url)
			defer check(resp.Body.Close)

			zipName := filepath.Join(os.TempDir(),
				strings.TrimPrefix(url, "https://github.com/ryanoasis/nerd-fonts/releases/download/v2.3.3/")) // I know this is hardcode
			f, _ := os.OpenFile(zipName, os.O_CREATE|os.O_WRONLY, 0644)
			defer f.Close()

			bar := progressbar.DefaultBytes(
				resp.ContentLength,
				fmt.Sprintf("downloading %-25s", zipName),
			)
			mpb.Add(bar)
			io.Copy(io.MultiWriter(f, bar), resp.Body)
		}(url)
	}

	wg.Wait()

	fmt.Println("All done")
}

// check checks the returned error of a function.
func check(f func() error) {
	if err := f(); err != nil {
		fmt.Fprintf(os.Stderr, "received error: %v\n", err)
	}
}

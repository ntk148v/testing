package main

import (
	"fmt"
	"github.com/schollz/progressbar/v3"
	"io"
	"os"
	"sync"
	"time"
)

type LineWriter struct {
	*MultiProgressBar
	id int
}

func (lw *LineWriter) Write(p []byte) (n int, err error) {
	lw.guard.Lock()
	defer lw.guard.Unlock()
	lw.move(lw.id, lw.output)
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

// Move cursor to the beginning of the current progressbar.
func (mpb *MultiProgressBar) move(id int, writer io.Writer) (int, error) {
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
	mpb.move(len(mpb.Bars), mpb.output)
}

func main() {
	mpb := NewMultiProgressBar(
		[]*progressbar.ProgressBar{
			progressbar.New(100),
			progressbar.New(100),
			progressbar.New(100),
		},
		os.Stdout,
	)
	mpb.Bars[0].Describe("Bar Zero")
	mpb.Bars[1].Describe("Bar One")
	mpb.Bars[2].Describe("Bar Two")

	length := len(mpb.Bars) - 1
	for val := 0; val < 300; val++ {
		time.Sleep(25 * time.Millisecond)
		barId := val % 3
		mpb.Bars[length-barId].Add(1)
	}
	mpb.End()
	fmt.Printf("TADA")
}

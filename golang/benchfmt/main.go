package main

import (
	"fmt"
	"os/exec"

	"github.com/pterm/pterm"
	"golang.org/x/perf/benchfmt"
)

func main() {
	area, _ := pterm.DefaultArea.WithFullscreen().WithCenter().Start()
	defer area.Stop()

	var bars pterm.Bars

	// Read from file
	// files := benchfmt.Files{Paths: []string{"sample/output.txt"}, AllowStdin: true, AllowLabels: true}
	// for files.Scan() {
	// 	switch rec := files.Result(); rec := rec.(type) {
	// 	case *benchfmt.SyntaxError:
	// 		// Non-fatal result parse error. Warn
	// 		// but keep going.
	// 		fmt.Println(rec)
	// 	case *benchfmt.Result:
	// 		for _, c := range rec.Config {
	// 			fmt.Println(c.Key, string(c.Value))
	// 		}
	// 		fmt.Println(rec.Name, rec.Iters, rec.Values)
	// 	}
	// }
	// if err := files.Err(); err != nil {
	// 	panic(err)
	// }

	// Read directly from stdout
	// cmd := exec.Command("go", "test", "-bench=.", "-benchtime=10s", "-benchmem")
	cmd := exec.Command("go", "test", "-bench=.")
	cmd.Dir = "sample"
	stdout, err := cmd.StdoutPipe()

	if err != nil {
		panic(err)
	}

	if err := cmd.Start(); err != nil {
		panic(err)
	}

	r := benchfmt.NewReader(stdout, "")
	for r.Scan() {
		switch rec := r.Result(); rec := rec.(type) {
		case *benchfmt.SyntaxError:
			// Non-fatal result parse error. Warn
			// but keep going.
			fmt.Println(rec)
		case *benchfmt.Result:
			// for _, c := range rec.Config {
			// 	fmt.Println(c.Key, string(c.Value))
			// }
			fmt.Println(rec.Name, rec.Iters, rec.Values[0].OrigValue, rec.Values[0].Value)
			bars = append(bars, pterm.Bar{
				Label: string(rec.Name),
				Value: int(rec.Values[0].OrigValue), // the first only
			})
		}
	}

	if err := cmd.Wait(); err != nil {
		panic(err)
	}

	barchart := pterm.DefaultBarChart.WithBars(bars).
		WithHorizontal().WithShowValue()
	content, _ := barchart.Srender()
	area.Update(content)
}

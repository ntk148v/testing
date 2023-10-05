package main

import (
	"fmt"
	"log"
	"os"

	"github.com/pterm/pterm"
	"golang.org/x/perf/benchfmt"
)

func main() {
	var bars pterm.Bars

	f, err := os.Open("log.bench")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Create a filter
	// filter, err := benchproc.NewFilter(".name:DisabledWithoutFields")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	r := benchfmt.NewReader(f, "example")
	for r.Scan() {
		var res *benchfmt.Result
		switch rec := r.Result(); rec := rec.(type) {
		case *benchfmt.Result:
			res = rec
		case *benchfmt.SyntaxError:
			// Report a non-fatal parse error.
			log.Print(err)
			continue
		default:
			// Unknown record type. Ignore.
			continue
		}
		// if match, err := filter.Apply(res); !match {
		// 	// Result was fully excluded by the filter.
		// 	if err != nil {
		// 		// Print the reason we rejected this result.
		// 		log.Print(err)
		// 	}
		// 	continue
		// }

		fmt.Println(res.Name, res.Iters, res.Values[0].OrigValue, res.Values[0].Value)
		bars = append(bars, pterm.Bar{
			Label: string(res.Name),
			Value: int(res.Values[0].OrigValue), // the first only
		})
	}
	// Check for I/O errors.
	if err := r.Err(); err != nil {
		log.Fatal(err)
	}

	pterm.DefaultBarChart.WithBars(bars).
		WithHorizontal().WithShowValue().
		Render()
}

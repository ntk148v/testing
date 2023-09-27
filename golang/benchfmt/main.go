package main

import (
	"fmt"

	"golang.org/x/perf/benchfmt"
)

func main() {
	files := benchfmt.Files{Paths: []string{"sample/output.txt"}, AllowStdin: true, AllowLabels: true}
	for files.Scan() {
		switch rec := files.Result(); rec := rec.(type) {
		case *benchfmt.SyntaxError:
			// Non-fatal result parse error. Warn
			// but keep going.
			fmt.Println(rec)
		case *benchfmt.Result:
			for _, c := range rec.Config {
				fmt.Println(c.Key, string(c.Value))
			}
			fmt.Println(rec.Name, rec.Iters, rec.Values)
		}
	}
	if err := files.Err(); err != nil {
		panic(err)
	}
}

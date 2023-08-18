package main

import (
	"fmt"
	"os"

	"golang.org/x/exp/slices"
	"golang.org/x/exp/slog"
)

func main() {
	// new built-in: clear
	// delete all elements from a map or zeros out all
	// elements of a slice, depending on the input
	m := map[string]string{"foo": "bar"}
	clear(m)
	fmt.Println(m)

	s := []string{"foo", "bar"}
	clear(s)
	fmt.Println(s)

	// Structured logging
	// Previously, you can use it through golang.org/x/exp/slog
	// The new log/slog package provides structure logging with levels,
	// emitting key=value pairs that enable fast and accurate machine processing
	// of data with minima allocations.

	// simple
	logger := slog.New(slog.NewTextHandler(os.Stderr))
	logger.Info("this is structured logging", "foo", "bar")
	// time=2023-08-18T11:16:39.149+07:00 level=INFO msg="this is structured logging" foo=bar

	// slices and maps packages
	// you can use them on earlier versions from golang.org/x/exp
	s = []string{"A", "A", "B", "C"}

	// Clone returns a shallow-copy clone of the slice.
	cs := slices.Clone(s)            // => ["A", "A", "B", "C"]
	fmt.Println(slices.Equal(s, cs)) // compare - true

	// Contains returns true if the value exists, or false otherwise
	fmt.Println(slices.Contains(s, "A"))
	fmt.Println(slices.Contains(s, "Z"))

	// Compact removes consecutive runs of duplicates elements
	cs = slices.Compact(s)
	fmt.Println(s, cs)

	// similar with maps.Clone ...
}

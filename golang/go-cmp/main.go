package main

import (
	"fmt"

	"github.com/google/go-cmp/cmp"
)

type Person struct {
	Name    string
	Age     int
	Complex Complex
}

type Complex struct {
	Field1 map[string]int
	Field2 []string
}

func main() {
	m := make(map[string]int)
	m["k1"] = 7
	m["k2"] = 8
	a := Person{
		Name: "Tui",
		Age:  26,
		Complex: Complex{
			Field1: m,
			Field2: []string{"test"},
		},
	}
	b := Person{
		Name: "Nguoikhac",
		Age:  26,
	}
	if diff := cmp.Diff(a, b); diff != "" {
		fmt.Printf("%+v\n", diff)
	}
}

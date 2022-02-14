package main

import (
	"fmt"

	flag "github.com/spf13/pflag"
)

func main() {
	var input string
	flag.StringVarP(&input, "name", "n", "default-value-ne", "Put the name here")
	flag.Parse()
	fmt.Println(input)
}

package main

import (
	"fmt"

	"github.com/atotto/clipboard"
)

func main() {
	// Write into clipboard
	err := clipboard.WriteAll("Here we go")
	if err != nil {
		panic(err)
	}
	str, err := clipboard.ReadAll()
	if err != nil {
		panic(err)
	}
	fmt.Println(str)
}

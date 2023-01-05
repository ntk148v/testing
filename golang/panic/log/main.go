package main

import (
	"errors"
	"fmt"
	"sync"
)

func somethingThatPanic() {
	panic(errors.New("error in somethingThatPanic function")) // a demo-purpose panic
}

func main() {
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer func() {
			wg.Done()
			if v := recover(); v != nil {
				// print the error
				fmt.Println("Capture a panic", v)
				fmt.Println("Avoid crashing the program")
			}
		}()
		somethingThatPanic()
	}()

	wg.Wait()
}

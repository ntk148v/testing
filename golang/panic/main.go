package main

import (
	"errors"
	"sync"
)

func somethingThatPanic() {
	panic(errors.New("error in somethingThatPanic function")) // a demo-purpose panic
}

func main() {
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		somethingThatPanic()
	}()

	wg.Wait()
}

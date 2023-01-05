package main

import (
	"errors"
	"fmt"
	"runtime/debug"
	"sync"
)

type caughtPanicError struct {
	val   interface{}
	stack []byte
}

func (e caughtPanicError) Error() string {
	return fmt.Sprintf("panic: %q\n%s", e.val, string(e.stack))
}

func somethingThatPanic() {
	panic(errors.New("error in somethingThatPanic function")) // a demo-purpose panic
}

func main() {
	var wg sync.WaitGroup
	wg.Add(1)

	done := make(chan error)
	go func() {
		// Turn goroutine's panic to an error and return to the goroutine spawner
		// through channel
		defer func() {
			wg.Done()
			if v := recover(); v != nil {
				done <- caughtPanicError{
					val:   v,
					stack: debug.Stack(),
				}
			} else {
				done <- nil
			}
		}()
		somethingThatPanic()
	}()

	err := <-done
	// Propagate the panic to the goroutine spawner
	if err != nil {
		panic(err)
	}

	wg.Wait()
}

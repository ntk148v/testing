package main

import (
    "errors"

    "github.com/sourcegraph/conc"
)

func somethingThatPanic() {
    panic(errors.New("error in somethingThatPanic function")) // a demo-purpose panic
}

func main() {
    var wg conc.WaitGroup
    wg.Go(somethingThatPanic)

    // Panics with a nice stacktrace
    wg.Wait()
}

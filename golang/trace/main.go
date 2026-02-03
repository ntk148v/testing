package main

import (
	"fmt"
	"os"
	"runtime"
	"runtime/trace"
	"sync"
)

func sum(a, b int) {
	_ = a + b
}

func subtract(a, b int) {
	_ = a - b
}

func multiply(a, b int) {
	_ = a * b
}

func divide(a, b int) {
	_ = a / b
}

// Define functions fn1 through fn5, each performing simple operations // and printing results.
func fn1(wg *sync.WaitGroup) {
	defer wg.Done() // Signal completion of goroutine
	fmt.Println("Executing function 1")
}

func fn2(wg *sync.WaitGroup) {
	defer wg.Done()
	sum(1, 2) // Calls a function to calculate the sum
	fmt.Println("Executing function 2")
}

func fn3(wg *sync.WaitGroup) {
	defer wg.Done()
	subtract(2, 1) // Calls a function to calculate the difference
	fmt.Println("Executing function 3")
}

func fn4(wg *sync.WaitGroup) {
	defer wg.Done()
	multiply(2, 2) // Calls a function to calculate the product
	fmt.Println("Executing function 4")
}

func fn5(wg *sync.WaitGroup) {
	defer wg.Done()
	divide(4, 2) // Calls a function to calculate the division
	fmt.Println("Executing function 5")
}

// Define arithmetic operations: sum, subtract, multiply, divide

func main() {
	var wg sync.WaitGroup
	runtime.GOMAXPROCS(3) // Limit the number of OS threads

	// Create and open a file to store trace data
	f, err := os.Create("trace.out")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// Start tracing
	err = trace.Start(f)
	if err != nil {
		panic(err)
	}
	defer trace.Stop()

	// Adding 5 goroutines to the WaitGroup
	wg.Add(5)

	// Starting each function as a goroutine
	go fn1(&wg)
	go fn2(&wg)
	go fn3(&wg)
	go fn4(&wg)
	go fn5(&wg)

	// Waiting for all goroutines to complete
	wg.Wait()

	// Stop tracing
	fmt.Println("All goroutines completed")
}

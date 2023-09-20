package main

import "fmt"

func main() {
	done := make(chan bool)

	values := []string{"a", "b", "c"}
	for _, v := range values {
		go func() {
			fmt.Println(v)
			done <- true
		}()
	}

	// wait for all goroutines to complete before exiting
	for _ = range values {
		<-done
	}
}

// usually print "c", "c", "c"
// instead of printing "a", "b", "c" in some order
// GOEXPERIMENT=loopvar go run main.go
// print "a", "b", "c" exactly like we expected.

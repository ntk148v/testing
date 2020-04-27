package main

import (
	"fmt"
	"time"
)

func forever() {
	for {
		// Empty, just do nothing
	}
}

func show() {
	for i := 1; i < 9696969; i++ {
		time.Sleep(5 * time.Second)
		fmt.Println(i)
	}
}

func main() {
	go show()
	forever()
	fmt.Println("OK we're done")
}

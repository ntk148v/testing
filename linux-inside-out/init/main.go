package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	fmt.Println("Hello from Go init!")
	fmt.Println("PID:", os.Getpid()) // printing the PID (process ID)

	for i := 0; ; i++ {
		// every two seconds printing the text "tick {tick number}"
		fmt.Println("tick", i)
		time.Sleep(2 * time.Second)
	}
}

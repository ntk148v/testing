package main

import (
	"fmt"
	"sync"
	"time"
)

func forever() {
	m := sync.Mutex{} // Same with sync.RWMutex
	m.Lock()
	m.Lock()
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

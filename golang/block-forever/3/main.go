package main

import (
	"fmt"
	"sync"
	"time"
)

var wg sync.WaitGroup

func forever() {
	wg.Wait()
}

func show() {
	for i := 1; i < 9696969; i++ {
		time.Sleep(5 * time.Second)
		fmt.Println(i)
	}
	wg.Done()
}

func main() {
	wg.Add(1)
	go show()
	forever()
	fmt.Println("OK we're done")
}

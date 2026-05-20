package main

import (
	"fmt"
	"math/rand/v2"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/sync/singleflight"
)

var (
	callCount atomic.Int32
	wg        sync.WaitGroup
)

type Result struct {
	Val    interface{}
	Err    error
	Shared bool
}

// Simulate a function that fetches data from a database
func fetchData() (interface{}, error) {
	callCount.Add(1)
	time.Sleep(100 * time.Millisecond)
	return rand.IntN(100), nil
}

// Wrap the fetchData function with singleflight
func fetchDataWrapperWithDo(g *singleflight.Group, id int) error {
	defer wg.Done()

	time.Sleep(time.Duration(id) * 40 * time.Millisecond)
	v, err, shared := g.Do("key-fetch-data", fetchData)
	if err != nil {
		return err
	}

	fmt.Printf("Goroutine %d: result: %v, shared: %v\n", id, v, shared)
	return nil
}

// Wrap the fetchData function with singleflight using DoChan
func fetchDataWrapperWithDoChan(g *singleflight.Group, id int) error {
	defer wg.Done()

	ch := g.DoChan("key-fetch-data", fetchData)
	select {
	case res := <-ch:
		if res.Err != nil {
			return res.Err
		}
		fmt.Printf("Goroutine %d: result: %v, shared: %v\n", id, res.Val, res.Shared)
	case <-time.After(50 * time.Millisecond):
		return fmt.Errorf("timeout waiting for result")
	}

	return nil
}

func fetchDataWrapperWithForget(g *singleflight.Group, id int, forget bool) error {
	defer wg.Done()

	// Forget the key before fetching
	if forget {
		g.Forget("key-fetch-data")
	}

	v, err, shared := g.Do("key-fetch-data", fetchData)
	if err != nil {
		return err
	}

	fmt.Printf("Goroutine %d: result: %v, shared: %v\n", id, v, shared)
	return nil
}

func main() {
	var g singleflight.Group

	// 5 goroutines to fetch the same data
	const numGoroutines = 5
	wg.Add(numGoroutines)

	for i := range numGoroutines {
		go fetchDataWrapperWithDo(&g, i)
	}

	wg.Add(3)

	// 2 goroutines fetch the data
	go fetchDataWrapperWithForget(&g, 0, false)
	go fetchDataWrapperWithForget(&g, 1, false)

	// Wait a bit and launch 1 more goroutine
	// Ensures goroutines 0, 1, and 2 overlap
	time.Sleep(10 * time.Millisecond)
	go fetchDataWrapperWithForget(&g, 2, true)

	wg.Wait()
	fmt.Printf("Function was called %d times\n", callCount.Load())
}

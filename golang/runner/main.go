package main

import (
	"fmt"
	"sync"
	"time"
)

type scaler struct {
	sync.RWMutex
	id       int
	interval time.Duration
	stopChan chan struct{}
}

func (s *scaler) do() {
	fmt.Printf("Tick tick: %d - %s\n", s.id, time.Now())
}

func (s *scaler) run(wg *sync.WaitGroup) {
	defer wg.Done()
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			s.do()
		case <-s.stopChan:
			time.Sleep(2 * time.Second)
			return
		}
	}
}

func (s *scaler) stop() {
	fmt.Println("stop day")
	close(s.stopChan)
	fmt.Println("stop roi day")
}

func main() {
	ss := make(map[int]*scaler, 0)
	for i := 1; i < 5; i++ {
		s := scaler{
			id:       i,
			interval: time.Duration(i) * time.Second,
			stopChan: make(chan struct{}),
		}
		ss[i] = &s
		fmt.Println(ss)
	}

	var wg sync.WaitGroup

	for _, s := range ss {
		wg.Add(1)
		go s.run(&wg)
	}

	time.Sleep(10 * time.Second)
	ss[1].stop()
	delete(ss, 1)
	wg.Wait()
}

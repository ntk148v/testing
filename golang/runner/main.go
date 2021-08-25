// Copyright 2021 Kien Nguyen-Tuan <kiennt2609@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"sync"
	"time"
)

type scaler struct {
	sync.RWMutex
	id          int
	interval    time.Duration
	stopChan    chan struct{}
	stoppedChan chan struct{}
}

func (s *scaler) do() {
	fmt.Printf("Tick tick: %d - %s\n", s.id, time.Now())
}

func (s *scaler) run(wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(s.stoppedChan)
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()
	for {
		select {
		case <-s.stopChan:
			return
		default:
			select {
			case <-ticker.C:
				s.do()
			case <-s.stopChan:
				return
			}
		}
	}
}

func (s *scaler) stop() {
	fmt.Println("stop day")
	close(s.stopChan)
	<-s.stoppedChan
	fmt.Println("stop roi day")
}

func main() {
	ss := make(map[int]*scaler, 0)
	for i := 1; i < 5; i++ {
		s := scaler{
			id:          i,
			interval:    time.Duration(i) * time.Second,
			stopChan:    make(chan struct{}),
			stoppedChan: make(chan struct{}),
		}
		ss[i] = &s
		fmt.Println(ss)
	}

	var wg sync.WaitGroup

	for _, s := range ss {
		wg.Add(1)
		go s.run(&wg)
	}

	time.Sleep(5 * time.Second)
	ss[1].stop()
	delete(ss, 1)
	wg.Wait()
}

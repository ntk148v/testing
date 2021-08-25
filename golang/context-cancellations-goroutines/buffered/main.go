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
	"log"
	"time"
)

func work() error {
	for i := 0; i < 1000; i++ {
		select {
		case <-time.After(2 * time.Second):
			fmt.Println("Doing some work ", i)
		}
	}
	return nil
}

func main() {
	fmt.Println("Hey, I'm going to do some work")

	ch := make(chan error, 1)
	go func() {
		ch <- work()
	}()

	select {
	case err := <-ch:
		if err != nil {
			log.Fatal("Something went wrong :(", err)
		}
	case <-time.After(4 * time.Second):
		fmt.Println("Life is to short to wait that long")
	}

	fmt.Println("Finished. I'm going home")
}

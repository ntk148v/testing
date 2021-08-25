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
	"strconv"
)

var i int

type Result struct {
	index  int
	result string
}

type Job struct {
	index int
	task  func() (string, error)
	retry int
}

func ConcurrentRetry(tasks []func() (string, error), concurrent int, retry int) <-chan Result {
	jobs := make(chan Job, len(tasks))
	results := make(chan Result, len(tasks))
	counts := make(chan bool, len(tasks))

	producer := func() {
		for i, task := range tasks {
			jobs <- Job{i, task, retry}
		}
	}

	worker := func() {
		for job := range jobs {
			if result, err := job.task(); err != nil {
				job.retry--
				if job.retry > 0 {
					jobs <- job
				} else {
					counts <- true
				}
			} else {
				results <- Result{job.index, result}
				counts <- true
			}
		}
	}

	cleaner := func() {
		for i := 0; i < len(tasks); i++ {
			<-counts
		}
		close(counts)
		close(results)
		close(jobs)
	}

	go producer()
	go cleaner()
	for i := 0; i < concurrent; i++ {
		go worker()
	}

	return results
}

func match() (string, error) {
	if i >= 4 {
		return strconv.Itoa(i), nil
	} else {
		i++
		return strconv.Itoa(i - 1), nil
		//return strconv.Itoa(i - 1), fmt.Errorf("False false")
	}
}

func main() {
	i = 0
	var matchs []func() (string, error)
	matchs = append(matchs, match)
	matchs = append(matchs, match)
	matchs = append(matchs, match)
	matchs = append(matchs, match)
	matchs = append(matchs, match)
	matchs = append(matchs, match)
	results := ConcurrentRetry(matchs, 1, 2)
	for j := 0; j <= len(results); j++ {
		r := <-results
		fmt.Println(r)
	}
}

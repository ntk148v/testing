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

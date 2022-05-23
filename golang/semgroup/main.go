package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/fatih/semgroup"
)

func main() {
	const maxWorkers = 2
	s := semgroup.NewGroup(context.Background(), maxWorkers)

	visitors := []int{5, 2, 10, 8, 9, 3, 1}

	for _, v := range visitors {
		v := v

		s.Go(func() error {
			if v%2 == 0 {
				return errors.New("invalid visitors")
			}
			fmt.Println("Visits: ", v)
			return nil
		})
	}

	// Wait for all visits to complete. Any errors are accumulated.
	if err := s.Wait(); err != nil {
		fmt.Println("Something went wrong")
		fmt.Println(err)
	}
}

// Visits:  9
// Visits:  3
// Visits:  1
// Visits:  5
// 3 error(s) occurred:
// * invalid visitors
// * invalid visitors
// * invalid visitors

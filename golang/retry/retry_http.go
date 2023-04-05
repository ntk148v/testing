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
	"math/rand"
	"net/http"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func retry(attempts int, sleep time.Duration, f func() error) error {
	if err := f(); err != nil {
		if s, ok := err.(stop); ok {
			// Return the original error for later checking
			return s.error
		}

		if attempts--; attempts > 0 {
			// Add some randomness to prevent creating a Thundering Herd
			jitter := time.Duration(rand.Int63n(int64(sleep)))
			sleep = sleep + jitter/2

			time.Sleep(sleep)
			return retry(attempts, 2*sleep, f)
		}
		return err
	}

	return nil
}

type stop struct {
	error
}

// DeleteThing attempts to delete a thing. It will try a maximum of three times.
func DeleteThing(id string) error {
	// Build the request
	req, err := http.NewRequest(
		"DELETE",
		fmt.Sprintf("https://unreliable-api/things/%s", id),
		nil,
	)
	if err != nil {
		return fmt.Errorf("unable to make request: %s", err)
	}

	// Execute the request
	return retry(3, time.Second, func() error {
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			// This error will result in a retry
			return err
		}
		defer resp.Body.Close()

		s := resp.StatusCode
		switch {
		case s >= 500:
			// Retry
			return fmt.Errorf("server error: %v", s)
		case s >= 400:
			// Don't retry, it was client's fault
			return stop{fmt.Errorf("client error: %v", s)}
		default:
			// Happy
			return nil
		}
	})
}

// Copyright 2022 Kien Nguyen-Tuan <kiennt2609@gmail.com>
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
	"io/ioutil"
	"time"

	"github.com/gojek/heimdall/v7/httpclient"
	"github.com/gojek/heimdall/v7/hystrix"
)

func main() {
	// Create a new HTTP client with a default timeout
	timeout := 1000 * time.Millisecond
	client := httpclient.NewClient(httpclient.WithHTTPTimeout(timeout))

	// Use the clients GET method to create and execute the request
	res, err := client.Get("http://google.com", nil)
	if err != nil {
		panic(err)
	}

	// Heimdall returns the standard *http.Response object
	body, err := ioutil.ReadAll(res.Body)
	fmt.Println(string(body))

	// hystrix-like circuit breaker
	// Create a new hystrix-wrapped HTTP client with the command name, along with other required options
	hystrixClient := hystrix.NewClient(
		hystrix.WithHTTPTimeout(10*time.Millisecond),
		hystrix.WithCommandName("google_get_request"),
		hystrix.WithHystrixTimeout(1000*time.Millisecond),
		hystrix.WithMaxConcurrentRequests(30),
		hystrix.WithErrorPercentThreshold(20),
		hystrix.WithStatsDCollector("localhost:8125", "myapp.hystrix"),
	)
	res, err = hystrixClient.Get("http://google.com.wrong", nil)
	if err != nil {
		panic(err)
	}

	body, err = ioutil.ReadAll(res.Body)
	fmt.Println(string(body))
}

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
	"bytes"
	"fmt"
)

type Logger struct{}

func (logger *Logger) Log(message string) {
	fmt.Println(message)
}

type HttpClient struct {
	logger *Logger
}

func (client *HttpClient) Get(url string) string {
	client.logger.Log("Getting " + url)

	// make an HTTP request
	return "my response from " + url
}

func NewHttpClient(logger *Logger) *HttpClient {
	return &HttpClient{logger}
}

type ConcatService struct {
	logger *Logger
	client *HttpClient
}

func (service *ConcatService) GetAll(urls ...string) string {
	service.logger.Log("Running GetAll")

	var result bytes.Buffer

	for _, url := range urls {
		result.WriteString(service.client.Get(url))
	}

	return result.String()
}

func NewConcatService(logger *Logger, client *HttpClient) *ConcatService {
	return &ConcatService{logger, client}
}

func main() {
	service := CreateConcatService()
	result := service.GetAll(
		"http://example.com",
		"http://kiennt.com",
	)
	fmt.Println(result)
}

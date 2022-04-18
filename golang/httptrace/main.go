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

// Source: https://golang.cafe/blog/golang-httptrace-example.html
package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptrace"
)

// Source: https://golang.cafe/blog/golang-httptrace-example.html
func main() {
	req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
	clientTrace := &httptrace.ClientTrace{
		GetConn:      func(hostPort string) { fmt.Println("starting to create conn ", hostPort) },
		DNSStart:     func(info httptrace.DNSStartInfo) { fmt.Println("starting to look up dns", info) },
		DNSDone:      func(info httptrace.DNSDoneInfo) { fmt.Println("done looking up dns", info) },
		ConnectStart: func(network, addr string) { fmt.Println("starting tcp connection", network, addr) },
		ConnectDone:  func(network, addr string, err error) { fmt.Println("tcp connection created", network, addr, err) },
		GotConn:      func(info httptrace.GotConnInfo) { fmt.Println("connection established", info) },
	}
	clientTraceCtx := httptrace.WithClientTrace(req.Context(), clientTrace)
	req = req.WithContext(clientTraceCtx)
	_, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
}

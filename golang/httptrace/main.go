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

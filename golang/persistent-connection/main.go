package main

import (
	"context"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptrace"
)

func main() {
	// client trace to log whether the request's underlying tcp connection was re-used
	clientTrace := &httptrace.ClientTrace{
		GotConn: func(info httptrace.GotConnInfo) { log.Printf("conn was reused: %t", info.Reused) },
	}
	traceCtx := httptrace.WithClientTrace(context.Background(), clientTrace)

	// 1st request
	req, err := http.NewRequestWithContext(traceCtx, http.MethodGet, "http://example.com", nil)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	// The http Client and Transport guarantee that Body is always non-nil, even on
	// responses without a body or responses with a zero-length body. It is the
	// caller's responsibility to close Body. The default HTTP client's Transport
	// may not reuse HTTP/1.x "keep-alive" TCP connections if the Body is not read
	// to completion and closed.
	if _, err := io.Copy(ioutil.Discard, resp.Body); err != nil {
		log.Fatal(err)
	}
	resp.Body.Close()
	// 2nd request
	req, err = http.NewRequestWithContext(traceCtx, http.MethodGet, "http://example.com", nil)
	if err != nil {
		log.Fatal(err)
	}
	_, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
}

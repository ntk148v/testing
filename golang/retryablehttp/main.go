package main

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"net"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

func defaultPooledTransport() *http.Transport {
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		MaxIdleConnsPerHost:   runtime.GOMAXPROCS(0) + 1,
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
	}
	return transport
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func main() {
	url := os.Getenv("URL")
	username := os.Getenv("USERNAME")
	password := os.Getenv("PASSWORD")
	client := retryablehttp.NewClient()
	client.HTTPClient = &http.Client{
		Transport: defaultPooledTransport(),
	}
	req, err := retryablehttp.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", "Basic"+basicAuth(username, password))
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	fmt.Println(resp.Request.URL.String())
}

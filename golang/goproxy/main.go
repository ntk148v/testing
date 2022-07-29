package main

import (
	"log"
	"net/http"

	"github.com/elazarl/goproxy"
)

func main() {
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = true

	// Get proxies from env
	// no fancy handling
	// envProxy := os.Getenv("http_proxy")
	// proxyUrl, _ := url.Parse(envProxy)
	// proxy.Tr = &http.Transport{
	// 	Proxy:           http.ProxyURL(proxyUrl),
	// 	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	// }

	log.Fatal(http.ListenAndServe(":8080", proxy))
}

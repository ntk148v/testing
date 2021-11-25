package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
)

// WriteCounter counts the number of bytes written to it.
type WriteCounter struct {
	Total int64 // Total # of bytes transferred
}

// Write implements the io.Writer interface.
//
// Always completes and never returns an error.
func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += int64(n)
	fmt.Printf("\rDownloading %d MB", wc.Total/1024)
	return n, nil
}

func main() {
	url := flag.String("u", "http://localhost:8100/servefile", "the url to get file")
	count := flag.Int("p", runtime.NumCPU(), "connection count")
	flag.Parse()
	// Use default http client, do not use it in production
	client := http.Client{}
	// Construct a request
	req, err := http.NewRequest("GET", *url, nil)
	if err != nil {
		panic(err)
	}
	// Check range support and get file size
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	log.Println("Check headers...")
	headers := resp.Header
	if accept_ranges, supported := headers["Accept-Ranges"]; !supported {
		log.Fatalln("Doesn't support header `Accept-Ranges`")
	} else if accept_ranges[0] != "bytes" {
		log.Fatalln("Support `Accept-Ranges`, but value is not `bytes`")
	}
	fileSize, _ := strconv.ParseInt(headers["Content-Length"][0], 10, 64)
	partialSize := int(fileSize) / *count
	// Get the 1st part
	req.Header.Add("Range", fmt.Sprintf("bytes=%d-%d", 0, partialSize))
	resp, err = client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	f, err := os.OpenFile("/tmp/1stpart", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	defer f.Close()
	if err != nil {
		panic(err)
	}
	// No handle at all for testing
	written, err := io.Copy(f, io.TeeReader(resp.Body, &WriteCounter{}))
	if err != nil {
		panic(err)
	}
	log.Println("")
	log.Printf("Downloaded %d MB of total %d MB", written/1024, fileSize/1024)
}

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
	fileSizeTmp, _ := strconv.ParseInt(headers["Content-Length"][0], 10, 64)
	fileSize := int(fileSizeTmp)
	partialSize := fileSize / *count
	// Get the 1st part
	req.Header.Add("Range", fmt.Sprintf("bytes=%d-%d", 0, partialSize))
	resp, err = client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	f, err := os.OpenFile("/tmp/1stpart", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// No handle at all for testing
	written, err := io.Copy(f, io.TeeReader(resp.Body, &WriteCounter{}))
	if err != nil {
		panic(err)
	}
	log.Println("")
	log.Printf("Downloaded 1st part %d MB of total %d MB", written/1024, fileSize/1024)

	// // Parallel download
	// var (
	// 	start, end int
	// 	wg         sync.WaitGroup
	// )

	// f, err = os.OpenFile("/tmp/parallel", os.O_CREATE|os.O_WRONLY, 0666)
	// if err != nil {
	// 	panic(err)
	// }
	// defer f.Close()

	// for i := 0; i < *count; i++ {
	// 	if i == *count {
	// 		end = int(fileSize)
	// 	} else {
	// 		end = start + partialSize
	// 	}

	// 	wg.Add(1)
	// 	go func(partIndex, start, end int, wg sync.WaitGroup) {
	// 		defer wg.Done()
	// 		// Perform request
	// 		reqPart, err := http.NewRequest("GET", *url, nil)
	// 		if err != nil {
	// 			log.Printf("Error - part %d: %s\n", partIndex, err)
	// 			return
	// 		}

	// 		req.Header.Add("Range", fmt.Sprintf("bytes=%d-%d", start, end))
	// 		respPart, err := client.Do(reqPart)
	// 		if err != nil {
	// 			log.Printf("Error - part %d: %s\n", partIndex, err)
	// 			return
	// 		}

	// 		defer respPart.Body.Close()
	// 		partSize, _ := strconv.ParseInt(resp.Header["Content-Length"][0], 10, 64)
	// 		log.Printf("Downloading part %d [%d-%d] size %d\n", partIndex, start, end, partSize/1024)

	// 		// make a buffer to keep chunks that are read
	// 		buf := make([]byte, 32*1024)
	// 		for {
	// 			nr, er := respPart.Body.Read(buf)
	// 			if nr > 0 {
	// 				nw, ew := f.WriteAt(buf[0:nr], int64(start))
	// 				if nw < 0 || nr < nw {
	// 					nw = 0
	// 					if ew == nil {
	// 						log.Fatalf("Error - part %d: invalid write\n", partIndex)
	// 						break
	// 					}
	// 				}
	// 				if nr != nw {
	// 					log.Fatalf("Error - part %d, short writting\n", partIndex)
	// 				}

	// 				start += nw
	// 			}

	// 			if er != nil {
	// 				log.Printf("Part %d got err %s", partIndex, er)
	// 				if er == io.EOF || er.Error() == "EOF" {
	// 					// Download successfully
	// 					break
	// 				}
	// 				log.Fatalf("Error - part %d: %s\n", partIndex, err)
	// 			}
	// 		}
	// 	}(i, start, end, wg)
	// 	start = end
	// }

	// wg.Wait()
	// log.Println("Downloaded")
}

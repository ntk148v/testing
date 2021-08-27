package main

import (
	"fmt"
	"os"

	"github.com/hashicorp/go-retryablehttp"
)

func main() {
	url := os.Getenv("URL")
	resp, err := retryablehttp.Get(url)
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
}

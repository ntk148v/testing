package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/avast/retry-go"
)

func main() {
	//url := "http://example.com"
	url := "http://testretrythoima.com"
	var body []byte

	err := retry.Do(
		func() error {
			resp, err := http.Get(url)
			fmt.Println(time.Now())

			if err == nil {
				defer func() {
					if err := resp.Body.Close(); err != nil {
						panic(err)
					}
				}()
				body, err = ioutil.ReadAll(resp.Body)
			}

			return err
		},
		retry.Attempts(10),
		retry.DelayType(retry.BackOffDelay),
	)
	fmt.Println(err)
}

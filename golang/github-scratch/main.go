package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type Label struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Color       string `json:"color"`
}

func main() {
	owner := "vCloud-DFTBA"
	repo := "faythe"
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/labels", owner, repo)

	file, err := ioutil.ReadFile("/tmp/labels.json")
	if err != nil {
		panic(err)
	}
	labels := make([]Label, 10)
	_ = json.Unmarshal([]byte(file), &labels)
	for _, l := range labels {
		jsonValue, _ := json.Marshal(l)
		token := "token " + os.Getenv("TOKEN")
		// Create a new request using http
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonValue))

		// add authorization header to the req
		req.Header.Add("Authorization", token)

		// Send req using http Client
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		body, err := ioutil.ReadAll(resp.Body)
		fmt.Print(string(body))
		resp, err = http.Get(url)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		body, err = ioutil.ReadAll(resp.Body)
		fmt.Print(string(body))
	}
}

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	owner := "ntk148v"
	repo := "testing"
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/labels", owner, repo)
	values := map[string]string{
		"name":        "newbug",
		"description": "Something isn't working",
		"color":       "f29513",
	}
	jsonValue, _ := json.Marshal(values)
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

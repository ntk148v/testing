// Copyright 2021 Kien Nguyen-Tuan <kiennt2609@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

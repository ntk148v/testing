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
	"fmt"
	"io/ioutil"

	"github.com/Netflix/go-env"
	"github.com/andygrunwald/go-jira"
)

type JiraConfig struct {
	URL      string `env:"JIRA_URL"`
	Username string `env:"JIRA_USERNAME"`
	Password string `env:"JIRA_PASSWORD"`
}

func main() {
	var jiraCfg JiraConfig
	fmt.Println("Get variable from environment")
	_, err := env.UnmarshalFromEnviron(&jiraCfg)
	if err != nil {
		panic(err)
	}

	// Oauth config
	tp := jira.BasicAuthTransport{
		Username: jiraCfg.Username,
		Password: jiraCfg.Password,
	}
	fmt.Println("Setup Jira client")
	jiraClient, err := jira.NewClient(tp.Client(), jiraCfg.URL)
	if err != nil {
		panic(err)
	}
	issue, resp, err := jiraClient.Issue.Get("JIRALERT-51060", nil)
	if err != nil {
		fmt.Println(resp.StatusCode)
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		bodyString := string(bodyBytes)
		fmt.Println(bodyString)
		panic(err)
	}
	currentStatus := issue.Fields.Status.Name
	fmt.Printf("Current status: %s\n", currentStatus)
	fmt.Println("Current projects:")
	req, _ := jiraClient.NewRequest("GET", "rest/api/2/project", nil)

	projects := new([]jira.Project)
	_, err = jiraClient.Do(req, projects)
	if err != nil {
		panic(err)
	}

	for _, project := range *projects {
		fmt.Printf("%s: %s\n", project.Key, project.Name)
	}
}

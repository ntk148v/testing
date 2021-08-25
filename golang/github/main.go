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
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func main() {
	token := os.Getenv("TOKEN")
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	// list all repositories for the authenticated user
	// repos, _, err := client.Repositories.List(ctx, "", nil)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(repos)

	b, err := ioutil.ReadFile("/home/kiennt65/Workspace/github.com/ntk148v/til/okd/overview.md")
	if err != nil {
		panic(err)
	}
	content := string(b)
	title := "Just for testing"

	issue, _, err := client.Issues.Create(ctx, "ntk148v", "testing", &github.IssueRequest{
		Title: &title,
		Body:  &content,
	})

	if err != nil {
		panic(err)
	}
	fmt.Println(issue)
}

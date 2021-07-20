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

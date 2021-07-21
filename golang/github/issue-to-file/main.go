package main

import (
	"bytes"
	"context"
	"fmt"
	"os"

	"github.com/google/go-github/github"
	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
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

	// list all issues by repo
	issues, _, err := client.Issues.ListByRepo(ctx, "ntk148v", "testing", &github.IssueListByRepoOptions{})
	if err != nil {
		panic(err)
	}
	markdown := goldmark.New(
		goldmark.WithExtensions(extension.GFM, meta.Meta),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
		),
	)
	for _, issue := range issues {
		var buf bytes.Buffer
		context := parser.NewContext()
		if err := markdown.Convert([]byte(*issue.Body), &buf, parser.WithContext(context)); err != nil {
			panic(err)
		}
		metaData := meta.Get(context)
		if path, ok := metaData["path"]; ok {
			fmt.Println(path)
		}
	}
}

package main

import (
	"bytes"
	"io/ioutil"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

func main() {
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM, extension.Linkify),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
		),
	)
	source, err := ioutil.ReadFile("/home/kiennt65/Workspace/github.com/ntk148v/til/okd/overview.md")
	if err != nil {
		panic(err)
	}
	var buf bytes.Buffer
	if err := md.Convert(source, &buf); err != nil {
		panic(err)
	}
}

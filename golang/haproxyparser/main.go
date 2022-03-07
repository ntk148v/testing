package main

import (
	"fmt"

	parser "github.com/haproxytech/config-parser/v4"
	"github.com/haproxytech/config-parser/v4/options"
)

func main() {
	p, err := parser.New(options.Path("haproxy.cfg"))
	if err != nil {
		panic(err)
	}

	backendSections, err := p.SectionsGet(parser.Backends)
	if err != nil {
		panic(err)
	}

	for _, backendSection := range backendSections {
		data, _ := p.Get(parser.Backends, backendSection, "server")
		fmt.Printf("%+v\n", data)
	}
}

package main

import (
	"fmt"

	"github.com/coreos/go-iptables/iptables"
)

func main() {
	t, err := iptables.NewWithProtocol(iptables.ProtocolIPv4)
	if err != nil {
		panic(err)
	}

	cs, err := t.ListChains("filter")
	if err != nil {
		panic(err)
	}

	for _, c := range cs {
		rs, err := t.List("filter", c)
		if err != nil {
			fmt.Println(err)
		}

		for _, r := range rs {
			fmt.Println(r)
		}
	}
}

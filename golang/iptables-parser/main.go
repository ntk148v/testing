package main

import (
	"fmt"

	ipt "github.com/coreos/go-iptables/iptables"
	iptp "github.com/kilo-io/iptables_parser"
)

func main() {
	t, err := ipt.NewWithProtocol(ipt.ProtocolIPv4)
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
			fmt.Printf("error listing chains: %v\n", err)
			continue
		}

		for _, r := range rs {
			tr, err := iptp.NewFromString(r)
			if err != nil {
				fmt.Printf("error parsing rules: %v\n", err)
				continue
			}

			switch r := tr.(type) {
			case iptp.Rule:
				fmt.Printf("rule parsed: %v\n", r)
			case iptp.Policy:
				fmt.Printf("policy parsed: %v\n", r)
			default:
				fmt.Printf("something else happend: %v\n", r)
			}
		}
	}
}

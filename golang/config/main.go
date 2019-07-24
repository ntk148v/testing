package main

import "fmt"

func main() {
	cfg, err := LoadFile("config.yml")
	if err != nil {
		fmt.Print(err)
	}
	fmt.Println(cfg)
	for _, ops := range cfg.OpenStackConfigs {
		fmt.Println(ops.Name)
	}
}

package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kevinburke/ssh_config"
)

func main() {
	f, _ := os.Open(filepath.Join(os.Getenv("HOME"), ".ssh", "config"))
	cfg, _ := ssh_config.Decode(f)
	for _, host := range cfg.Hosts {
		fmt.Println("patterns:", host.Patterns)
		for _, node := range host.Nodes {
			// Manipulate the nodes as you see fit, or use a type switch to
			// distinguish between Empty, KV, and Include nodes.
			fmt.Println(node.String())
		}
	}

	// Print the config to stdout:
	fmt.Println(cfg.String())

	files := ssh_config.Get("somehost", "IdentityFile")
	fmt.Println(files)
}

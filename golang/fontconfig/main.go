package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

func main() {
	out, err := exec.Command("fc-list", "--format=%{family[0]}\\n").Output()
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(string(out), "\n")
	for _, l := range lines {
		fmt.Println(l)
	}
}

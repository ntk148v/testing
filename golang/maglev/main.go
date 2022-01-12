package main

import (
	"fmt"
	"log"

	"github.com/kkdai/maglev"
)

func main() {
	sizeN := 5
	var lookupSizeM uint64
	lookupSizeM = 13 //(must be prime number)

	var names []string
	for i := 0; i < sizeN; i++ {
		names = append(names, fmt.Sprintf("backend-%d", i))
	}
	//backend-0 ~ backend-4

	mm, err := maglev.NewMaglev(names, lookupSizeM)
	if err != nil {
		log.Fatal("Init maglev", err)
	}
	v, _ := mm.Get("IP1")
	log.Println("node1:", v)
	//node1: backend-2
	v, _ = mm.Get("IP2")
	log.Println("node2:", v)
	//node2: backend-1
	v, _ = mm.Get("IPasdasdwni2")
	log.Println("node3:", v)
	//node3: backend-0

	if err := mm.Remove("backend-0"); err != nil {
		log.Fatal("Remove failed", err)
	}
	v, _ = mm.Get("IPasdasdwni2")
	log.Println("node3-D:", v)
	//node3-D: Change from "backend-0" to "backend-1"
}

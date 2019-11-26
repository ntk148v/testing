package main

import (
	"crypto/md5"
	"fmt"
	"math/rand"
	"time"

	"github.com/ntk148v/hashring"
	"stathat.com/c/consistent"
)

func hash(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func main() {
	servers := []string{"node1", "node2", "node3"}
	ring := hashring.New(servers)
	c := consistent.New()
	c.Add("node1")
	c.Add("node2")
	c.Add("node3")
	var server string
	for i := 0; i < 10; i++ {
		s := randStringRunes(3)
		hs := hash(s)
		fmt.Println("----------------HashRing----------------")
		server, _ = ring.GetNode(s)
		fmt.Printf("%s-%s\n", s, server)
		server, _ = ring.GetNode(hs)
		fmt.Printf("%s-%s\n", hs, server)
		fmt.Println("----------------Consistent----------------")
		server, _ = c.Get(s)
		fmt.Printf("%s-%s\n", s, server)
		server, _ = c.Get(hs)
		fmt.Printf("%s-%s\n", hs, server)
		fmt.Println("--------------------------------------")
	}
}

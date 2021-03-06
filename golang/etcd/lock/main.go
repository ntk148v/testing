package main

import (
	"context"
	"fmt"
	"log"

	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/clientv3/concurrency"
)

func main() {
	cli, err := clientv3.New(clientv3.Config{Endpoints: []string{"127.0.0.1:2379"}})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	// create 2 separate sessions for lock competition
	s1, err := concurrency.NewSession(cli)
	if err != nil {
		log.Fatal(err)
	}
	defer s1.Close()
	m1 := concurrency.NewMutex(s1, "/my-lock")

	s2, err := concurrency.NewSession(cli)
	if err != nil {
		log.Fatal(err)
	}
	defer s2.Close()
	m2 := concurrency.NewMutex(s2, "/my-lock")

	// acquire lock for s1
	if err = m1.Lock(context.TODO()); err != nil {
		log.Fatal(err)
	}
	fmt.Println("acquired lock for s1")

	if err = m2.TryLock(context.TODO()); err == nil {
		log.Fatal("should not acquire lock")
	}
	if err == concurrency.ErrLocked {
		fmt.Println("cannot acquire lock for s2, already locked in another session")
	}
	if err = m1.Unlock(context.TODO()); err != nil {
		log.Fatal(err)
	}
	fmt.Println("released lock for s1")
	if err = m2.TryLock(context.TODO()); err != nil {
		log.Fatal(err)
	}
	fmt.Println("acquired lock for s2")
}

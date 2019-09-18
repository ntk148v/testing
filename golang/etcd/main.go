package main

import (
	"context"
	"fmt"
	"time"

	"go.etcd.io/etcd/clientv3"
)

func main() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints: []string{"10.240.201.235:8379", "10.240.201.236:8379"},
	})

	if err != nil {
		panic(err)
	}

	defer cli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	_, err = cli.Put(ctx, "test1/test2", "dm")
	if err != nil {
		panic(err)
	}
	resp, err := cli.Get(ctx, "test1/test2")
	if err != nil {
		panic(err)
	}

	for _, ev := range resp.Kvs {
		fmt.Printf("%s : %s\n", ev.Key, ev.Value)
	}
	resp, err = cli.Get(ctx, "test1/test2/test")
	if err != nil {
		panic(err)
	}

	for _, ev := range resp.Kvs {
		fmt.Printf("%s : %s\n", ev.Key, ev.Value)
	}
	cancel()
}

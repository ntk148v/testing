package main

import (
	"context"
	"fmt"
	"time"

	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/clientv3/namespace"
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
	defer cancel()

	unprefixedKV := cli.KV
	cli.KV = namespace.NewKV(cli.KV, "prefixne/")
	cli.Watcher = namespace.NewWatcher(cli.Watcher, "prefixne/")
	cli.Lease = namespace.NewLease(cli.Lease, "prefixne/")

	cli.Put(ctx, "foo", "bar")
	resp, _ := unprefixedKV.Get(ctx, "prefixne/foo")
	fmt.Printf("%s\n", resp.Kvs[0].Value)

	unprefixedKV.Put(ctx, "prefixne/foo", "notbar")
	resp, _ = cli.Get(ctx, "foo")
	fmt.Printf("%s\n", resp.Kvs[0].Value)
}

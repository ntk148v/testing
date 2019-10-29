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

	cli.KV = namespace.NewKV(cli.KV, "prefixne/")
	cli.Watcher = namespace.NewWatcher(cli.Watcher, "prefixne/")
	cli.Lease = namespace.NewLease(cli.Lease, "prefixne/")
	// Test 5 seconds
	lresp, err := cli.Grant(ctx, 6)
	if err != nil {
		panic(err)
	}

	watch := cli.Watch(ctx, "foo", clientv3.WithPrefix())

	for i := 0; i < 5; i++ {
		fmt.Println(i)
		_, err = cli.Put(ctx, fmt.Sprintf("foo/%d", i), fmt.Sprintf("bar%d", i), clientv3.WithLease(lresp.ID))
		if err != nil {
			fmt.Println(err)
		}
		time.Sleep(2 * time.Second)
	}

	fmt.Println("Put xong roi day")

	for stay, timeout := true, time.After(10*time.Second); stay; {
		select {
		case wresp := <-watch:
			for _, e := range wresp.Events {
				if e.Type == clientv3.EventTypeDelete {
					fmt.Printf("%+v\n", e)
				}
			}
		case <-timeout:
			stay = false
		}
	}
}

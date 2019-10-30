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
		Endpoints: []string{"127.0.0.1:2379"},
	})

	if err != nil {
		panic(err)
	}

	defer cli.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cli.KV = namespace.NewKV(cli.KV, "prefixne/")
	cli.Watcher = namespace.NewWatcher(cli.Watcher, "prefixne/")
	cli.Lease = namespace.NewLease(cli.Lease, "prefixne/")
	// Test 5 seconds
	lresp, err := cli.Grant(ctx, 5)
	if err != nil {
		panic(err)
	}
	fmt.Println("Grant lease:", lresp.ID)

	_, err = cli.Put(ctx, "foo", "bar", clientv3.WithLease(lresp.ID))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Put foo with lease:", lresp.ID)

	watch := cli.Watch(ctx, "foo", clientv3.WithPrefix())
	done := make(chan bool)
	ticker := time.NewTicker(time.Second)

	go func() {
		for {
			select {
			case <-done:
				// Make sure lease & key be deleted.
				cli.Revoke(ctx, lresp.ID)
				return
			case wresp := <-watch:
				for _, e := range wresp.Events {
					if e.Type == clientv3.EventTypeDelete {
						fmt.Printf("Deleted key: %+v\n", e)
					}
				}
			case <-ticker.C:
				time.Sleep(time.Second)
				resp, err := cli.TimeToLive(ctx, lresp.ID, clientv3.WithAttachedKeys())
				if err != nil {
					panic(err)
				}
				fmt.Printf("Time to live: %d - %s\n", resp.TTL, resp.Keys)
				_, err = cli.KeepAliveOnce(ctx, lresp.ID)
				if err != nil {
					panic(err)
				}
				fmt.Println("Keep alive lease:", lresp.ID)
			}
		}
	}()

	time.Sleep(10 * time.Second)
	ticker.Stop()
	time.Sleep(10 * time.Second)
	done <- true
	fmt.Println("Stopped")
}

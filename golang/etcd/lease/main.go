// Copyright 2021 Kien Nguyen-Tuan <kiennt2609@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"fmt"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/namespace"
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

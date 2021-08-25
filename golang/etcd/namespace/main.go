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

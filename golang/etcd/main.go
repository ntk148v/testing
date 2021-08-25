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
)

type Person struct {
	Name string
	Age  int
}

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

	// p := Person{
	// 	Name: "Kien",
	// 	Age:  26,
	// }

	// value, err := json.Marshal(&p)
	// if err != nil {
	// 	panic(err)
	// }

	// _, err = cli.Put(ctx, "/cloud1/test1", string(value))
	// if err != nil {
	// 	panic(err)
	// }

	// p = Person{
	// 	Name: "Kien Khac",
	// 	Age:  69,
	// }

	// value, err = json.Marshal(&p)
	// if err != nil {
	// 	panic(err)
	// }

	// _, err = cli.Put(ctx, "/cloud1/test2", string(value))
	// if err != nil {
	// 	panic(err)
	// }

	resp, err := cli.Get(ctx, "/clouds", clientv3.WithPrefix())
	if err != nil {
		panic(err)
	}

	for _, ev := range resp.Kvs {
		fmt.Printf("%s : %s\n", ev.Key, ev.Value)
		// p = Person{}
		// err = json.Unmarshal(ev.Value, &p)
		// if err != nil {
		// 	panic(err)
		// }
		// fmt.Printf("%+v", p)
	}

	// resp, err = cli.Get(ctx, "/clouds/", clientv3.WithPrefix())
	// if err != nil {
	// 	panic(err)
	// }

	// if len(resp.Kvs) > 0 {
	// 	fmt.Println("Found")
	// }

	_, _ = cli.Delete(ctx, "/clouds", clientv3.WithPrefix())
	_, _ = cli.Delete(ctx, "/scalers", clientv3.WithPrefix())

}

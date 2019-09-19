package main

import (
	"context"
	"encoding/json"
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

	p := Person{
		Name: "Kien",
		Age:  26,
	}

	value, err := json.Marshal(&p)
	if err != nil {
		panic(err)
	}

	_, err = cli.Put(ctx, "/cloud/test1", string(value))
	if err != nil {
		panic(err)
	}

	p = Person{
		Name: "Kien Khac",
		Age:  69,
	}

	value, err = json.Marshal(&p)
	if err != nil {
		panic(err)
	}

	_, err = cli.Put(ctx, "/cloud/test2", string(value))
	if err != nil {
		panic(err)
	}

	resp, err = cli.Get(ctx, "cloud", clientv3.WithPrefix())
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

	cancel()
}
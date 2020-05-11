package main

import (
	"context"
	"flag"
	"fmt"
	"strings"
	"time"

	"go.etcd.io/etcd/clientv3"
)

func main() {
	etcdHost := flag.String("etcdHost", "127.0.0.1:2379", "etcd host")
	etcdWatchKey := flag.String("etcdWatchKey", "/foo", "etcd key to watch")

	flag.Parse()

	fmt.Println("connecting to etcd - " + *etcdHost)

	etcd, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"http://" + *etcdHost},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("connected to etcd - " + *etcdHost)

	defer etcd.Close()

	watchChan := etcd.Watch(context.Background(), *etcdWatchKey, clientv3.WithPrefix())
	fmt.Println("set WATCH on " + *etcdWatchKey)

	go func() {
		fmt.Println("started goroutine for PUT...")
		var k string
		for i := 0; i <= 20; i++ {
			if i%2 == 0 {
				k = fmt.Sprintf("%s/%d/%s", *etcdWatchKey, i, "subfoo")
			} else {
				k = fmt.Sprintf("%s/%d/%s", *etcdWatchKey, i, "notsubfoo")
			}
			etcd.Put(context.Background(), k, time.Now().String())
			fmt.Println("populated " + k + " with a value..")
			time.Sleep(2 * time.Second)
		}

	}()

	for watchResp := range watchChan {
		for _, event := range watchResp.Events {
			if strings.Contains(string(event.Kv.Key), "/subfoo") {
				fmt.Printf("Event received! %s executed on %q with value %q\n", event.Type, event.Kv.Key, event.Kv.Value)
			}
		}
	}
}

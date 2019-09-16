package main

import (
	"context"
	"fmt"
	"time"

	prometheusclient "github.com/prometheus/client_golang/api"
	prometheus "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

func main() {
	client, err := prometheusclient.NewClient(prometheusclient.Config{
		Address: "http://10.240.226.12:9090",
	})

	if err != nil {
		panic(err)
	}

	api := prometheus.NewAPI(client)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	query := "absent(up) or sum by(instance, job) (up) < 1"
	fmt.Println(query)
	val, warn, err := api.Query(ctx, query, time.Now())
	fmt.Println(warn)
	if err != nil {
		panic(err)
	}
	// fmt.Println(val)
	switch v := val.(type) {
	case model.Vector:
		for _, i := range v {
			fmt.Println(i.Metric["instance"])
		}
	default:
		fmt.Println("Khong phai cai nao")
	}
}

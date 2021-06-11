package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	env "github.com/Netflix/go-env"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

type InfluxDBConfig struct {
	URL      string `env:"INFLUXDB_URL"`
	Username string `env:"INFLUXDB_USERNAME"`
	Password string `env:"INFLUXDB_PASSWORD"`
}

func main() {
	log.Println("Setup InfluxDB client...")
	var influxCfg InfluxDBConfig
	_, err := env.UnmarshalFromEnviron(&influxCfg)
	if err != nil {
		log.Fatal(err)
	}
	if influxCfg.URL == "" {
		log.Fatal(errors.New("Missing InfluxDB configurations"))
	}
	client := influxdb2.NewClient(influxCfg.URL, fmt.Sprintf("%s:%s", influxCfg.Username, influxCfg.Password))
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer func() {
		cancel()
		client.Close()
	}()

	log.Println("Check health")
	check, err := client.Health(ctx)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(*check.Version)

	log.Println("Test query")
	result, err := client.QueryAPI("").QueryRaw(ctx, `buckets()`, nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(result)
}

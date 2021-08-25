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

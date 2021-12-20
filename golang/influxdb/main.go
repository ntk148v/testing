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
	fmt.Println("## Setup InfluxDB client")
	var influxCfg InfluxDBConfig
	_, err := env.UnmarshalFromEnviron(&influxCfg)
	if err != nil {
		panic(err)
	}
	if influxCfg.URL == "" {
		panic(errors.New("missing InfluxDB configurations"))
	}
	client := influxdb2.NewClient(influxCfg.URL, fmt.Sprintf("%s:%s", influxCfg.Username, influxCfg.Password))
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer func() {
		cancel()
		client.Close()
	}()

	check, err := client.Health(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println("## Check health and version", *check.Version)

	// Get query client
	queryAPI := client.QueryAPI("")
	resources := make(map[string][]string)

	// Get customer
	getCustomerQuery := `
	from(bucket: "telegraf/autogen")
	|> range(start: -5m)
	|> filter(fn: (r) => r._measurement == "system")
	|> keyValues(keyColumns: ["customer"])
	|> group()
	|> keep(columns: ["customer"])
	|> distinct(column: "customer")`
	fmt.Printf("## Get customer (query %s)\n", getCustomerQuery)
	customerResult, err := queryAPI.Query(ctx, getCustomerQuery)

	if err == nil {
		if customerResult.Err() != nil {
			fmt.Printf("## Query error: %s\n", customerResult.Err().Error())
		}
		for customerResult.Next() {
			if customerResult.TableChanged() {
				fmt.Printf("table: %s\n", customerResult.TableMetadata().String())
			}
			customer := customerResult.Record().Value().(string)

			// Get host by customer
			getHostQuery := fmt.Sprintf(`
			from(bucket: "telegraf/autogen")
			|> range(start: -5m)
			|> filter(fn: (r) => r._measurement == "system")
			|> filter(fn: (r) => r["customer"] == "%s")
			|> keyValues(keyColumns: ["host"])
			|> group()
			|> keep(columns: ["host"])
			|> distinct(column: "host")
			`, customer)
			fmt.Printf("## Get host by customer %s (query %s)\n", customer, getHostQuery)

			hostResult, err := queryAPI.Query(ctx, getHostQuery)
			if err != nil {
				fmt.Printf("Query error: %s\n", err.Error())
			}
			if hostResult.Err() != nil {
				fmt.Printf("Query error: %s\n", hostResult.Err().Error())
			}
			for hostResult.Next() {
				if hostResult.TableChanged() {
					fmt.Printf("table: %s\n", hostResult.TableMetadata().String())
				}
				resources[customer] = append(resources[customer], hostResult.Record().Value().(string))
			}

		}
	} else {
		fmt.Printf("## Query error: %s\n", err.Error())
	}
	fmt.Println("## Final results:", resources)
}

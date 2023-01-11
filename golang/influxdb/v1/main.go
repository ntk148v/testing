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
	"errors"
	"fmt"
	"time"

	env "github.com/Netflix/go-env"
	_ "github.com/influxdata/influxdb1-client" // this is important because of the bug in go mod
	influxdb1 "github.com/influxdata/influxdb1-client/v2"
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

	client, err := influxdb1.NewHTTPClient(influxdb1.HTTPConfig{
		Addr:     influxCfg.URL,
		Username: influxCfg.Username,
		Password: influxCfg.Password,
	})
	if err != nil {
		panic(err)
	}

	defer client.Close()

	fmt.Println("## Check health with ping")
	pingTime, version, err := client.Ping(10 * time.Second)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Reach influxDB version %s after %s\n", version, pingTime)

	fmt.Println("## Query something")
	cmd := "SHOW SERIES EXACT CARDINALITY"
	fmt.Println("Execute query:", cmd)

	resp, err := client.Query(influxdb1.NewQueryWithRP(cmd, "telegraf", "autogen", "rfc3339"))
	if err != nil {
		panic(err)
	}

	if resp.Error() != nil {
		panic(resp.Error())
	}

	fmt.Println("Result")
	for _, r := range resp.Results {
		fmt.Println(r.StatementId, r.Messages, r.Series, r.Err)
	}
}

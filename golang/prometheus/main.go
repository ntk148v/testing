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
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

// BasicAuthTransport is an http.RoundTripper that authenticates all requests
// using HTTP Basic Authentication with the provided username and password
type BasicAuthTransport struct {
	Username string
	Password string
	// Transport is the underlying HTTP transport to use when making requests.
	// It will default to http.DefaultTransport if nil
	Tranport http.RoundTripper
}

// RoundTrip implements the RoundTripper interface.
func (t *BasicAuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// To set extra headers, we must make a copy of the Request so
	// that we don't modify the Request we were given. This is required by the
	// specification of http.RoundTripper.
	//
	// Since we are going to modify only req.Header here, we only need a deep copy
	// of req.Header.
	clnReq := new(http.Request)
	*clnReq = *req
	clnReq.Header = make(http.Header, len(req.Header))
	for k, s := range req.Header {
		clnReq.Header[k] = append([]string(nil), s...)
	}

	clnReq.SetBasicAuth(t.Username, t.Password)
	return t.transport().RoundTrip(clnReq)
}

func (t *BasicAuthTransport) transport() http.RoundTripper {
	if t.Tranport == nil {
		return http.DefaultTransport
	}
	return t.Tranport
}

func configFromEnv() (*api.Config, error) {
	cfg := &api.Config{}
	url := os.Getenv("PROM_URL")
	url = "http://10.240.202.201:9092"
	if url == "" {
		return nil, errors.New("Prometheus address is missing")
	}
	cfg.Address = url
	user := os.Getenv("PROM_USER")
	pass := os.Getenv("PROM_PASS")
	// Both username and password are non empty string
	if user != "" && pass != "" {
		rt := &BasicAuthTransport{
			Username: user,
			Password: pass,
		}
		cfg.RoundTripper = rt
	}

	return cfg, nil
}

func main() {
	cfg, err := configFromEnv()
	if err != nil {
		panic(err)
	}

	client, err := api.NewClient(*cfg)

	if err != nil {
		panic(err)
	}

	papi := v1.NewAPI(client)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	query := "absent(up) or sum by(instance, job) (up) < 1"
	fmt.Println(query)
	val, warn, err := papi.Query(ctx, query, time.Now())
	fmt.Println(warn)
	if err != nil {
		fmt.Printf("%+v\n", err)
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

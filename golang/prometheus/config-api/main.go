package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"

	amclient "github.com/prometheus/alertmanager/api/v2/client"
	prometheusclient "github.com/prometheus/client_golang/api"
	prometheus "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/prometheus/config"
	"github.com/prometheus/prometheus/discovery"
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

func main() {
	cfg := prometheusclient.Config{
		Address: "http://10.240.201.174:9090",
	}
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
	client, err := prometheusclient.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	api := prometheus.NewAPI(client)
	ams, err := api.AlertManagers(context.Background())
	if err != nil {
		panic(err)
	}
	for _, a := range ams.Active {
		u, err := url.Parse(a.URL)
		if err != nil {
			panic(err)
		}
		fmt.Println(u.Host)
		amCli := amclient.NewHTTPClientWithConfig(nil, &amclient.TransportConfig{
			Host:     u.Host,
			BasePath: amclient.DefaultBasePath,
			Schemes:  amclient.DefaultSchemes,
		})
		// resp, err := amCli.Silence.GetSilences(silence.NewGetSilencesParams())
		resp, err := amCli.General.GetStatus()
		if err != nil {
			panic(err)
		}
		fmt.Println(resp)
	}
	cfgRaw, err := api.Config(context.Background())
	if err != nil {
		panic(err)
	}
	promCfg, err := config.Load(cfgRaw.YAML)
	fmt.Println(promCfg.AlertingConfig.AlertmanagerConfigs[0].ServiceDiscoveryConfigs[0].(discovery.StaticConfig)[0].Targets)
}

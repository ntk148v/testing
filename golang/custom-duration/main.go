package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/prometheus/common/model"
	"gopkg.in/yaml.v2"
)

type TestDefaultDuration struct {
	Duration time.Duration `yaml:"duration",json:"duration"`
}

type TestPrometheusModelDuration struct {
	Duration model.Duration `yaml:"duration",json:"duration"`
}

type TestOpenAPIDuration struct {
	Duration strfmt.Duration `yaml:"duration",json:"duration"`
}

type Duration model.Duration

func (d Duration) String() string {
	return model.Duration(d).String()
}

func (d Duration) MarshalJSON() (interface{}, error) {
	return json.Marshal(model.Duration(d).String())
}

func (d *Duration) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	var dstr string
	if err := json.Unmarshal(data, &dstr); err != nil {
		return err
	}
	tt, err := model.ParseDuration(dstr)
	if err != nil {
		return err
	}
	*d = Duration(tt)
	return nil
}

type TestCustomDuration struct {
	Duration Duration `yaml:"duration",json:"duration"`
}

func main() {
	rt := []byte(`{"duration": "90d"}`)
	var t1 TestDefaultDuration
	err := yaml.Unmarshal(rt, &t1)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(t1.Duration)
	}

	var t2 TestPrometheusModelDuration
	err = yaml.Unmarshal(rt, &t2)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(t2.Duration)
	}

	var t3 TestOpenAPIDuration
	err = json.Unmarshal(rt, &t3)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(t3.Duration)
	}

	var t4 TestCustomDuration
	err = json.Unmarshal(rt, &t4)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(t4.Duration)
	}
}

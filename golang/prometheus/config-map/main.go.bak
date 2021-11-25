// Copyright (c) 2020 Kien Nguyen-Tuan <kiennt2609@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	prommodel "github.com/prometheus/common/model"
	promcfg "github.com/prometheus/prometheus/config"
	promtgr "github.com/prometheus/prometheus/discovery/targetgroup"
	promrule "github.com/prometheus/prometheus/pkg/rulefmt"
	"golang.org/x/sync/errgroup"
	"gopkg.in/yaml.v2"
)

type PromCfgMap struct {
	Main   map[string]*promcfg.Config
	Rule   map[string]*promrule.RuleGroups
	FileSD map[string][]*promtgr.Group
}

// New inits a PromCfgMap instance.
func New(file string) (*PromCfgMap, error) {
	// It may take a while
	mainCfg, err := promcfg.LoadFile(file)
	if err != nil {
		return nil, err
	}
	//parent := filepath.Dir(file)

	var (
		errgr errgroup.Group
		cm    = &PromCfgMap{}
		main  = make(map[string]*promcfg.Config)
	)
	main[file] = mainCfg
	cm.Main = main
	cm.Rule = make(map[string]*promrule.RuleGroups)
	cm.FileSD = make(map[string][]*promtgr.Group)

	// Rule files
	errgr.Go(func() error {
		for _, f := range mainCfg.RuleFiles {
			rgs, errs := promrule.ParseFile(f)
			if errs != nil {
				// Combine errors into one.
				var errStr []string
				for _, err = range errs {
					errStr = append(errStr, err.Error())
				}
				return fmt.Errorf(strings.Join(errStr, "\n"))
			}

			cm.Rule[f] = rgs
		}
		return nil
	})

	// FileSD
	errgr.Go(func() error {
		for _, sc := range mainCfg.ScrapeConfigs {
			for _, fc := range sc.ServiceDiscoveryConfig.FileSDConfigs {
				for _, f := range fc.Files {
					tgrs, err := ReadFileSD(f)
					if err != nil {
						return err
					}
					cm.FileSD[f] = tgrs
				}
			}
		}
		return nil
	})

	return cm, errgr.Wait()
}

// ReadFileSD reads a JSON or YAML list of targets groups from the
// the file, depending on its file extension. It returns full
// configuration target groups.
func ReadFileSD(filename string) ([]*promtgr.Group, error) {
	fd, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	content, err := ioutil.ReadAll(fd)
	if err != nil {
		return nil, err
	}

	var targetGroups []*promtgr.Group

	switch ext := filepath.Ext(filename); strings.ToLower(ext) {
	case ".json":
		if err := json.Unmarshal(content, &targetGroups); err != nil {
			return nil, err
		}
	case ".yml", ".yaml":
		if err := yaml.UnmarshalStrict(content, &targetGroups); err != nil {
			return nil, err
		}
	default:
		panic(errors.Errorf("ReadFileSD: unhandled file extension %q", ext))
	}

	for i, tg := range targetGroups {
		if tg == nil {
			err = errors.New("nil target group item found")
			return nil, err
		}

		tg.Source = fileSource(filename, i)
		//if tg.Labels == nil {
		//	tg.Labels = prommodel.LabelSet{}
		//}
		//tg.Labels[fileSDFilepathLabel] = prommodel.LabelValue(filename)
	}

	return targetGroups, nil
}

const fileSDFilepathLabel = prommodel.MetaLabelPrefix + "filepath"

// fileSource returns a source ID for the i-th target group in the file.
func fileSource(filename string, i int) string {
	return fmt.Sprintf("%s:%d", filename, i)
}

func main() {
	p, err := New("/warehouse/prompoc-test/udcntt-cloud/prometheus_server/prometheus.yml")
	if err != nil {
		panic(err)
	}
	fmt.Print("%+v", p.Rule["/warehouse/prompoc-test/udcntt-cloud/prometheus_server/udcntt_process_alert.rules"].Groups)
}

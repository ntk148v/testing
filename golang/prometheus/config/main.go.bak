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
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/prometheus/config"
	fsnotify "gopkg.in/fsnotify/fsnotify.v1"
)

type PromConfig struct {
	watcher    *fsnotify.Watcher
	lock       sync.RWMutex
	config     *config.Config
	configPath string
	logger     log.Logger
}

func newPromConfig(conf *config.Config, logger log.Logger, confPath string) *PromConfig {
	if logger == nil {
		logger = log.NewNopLogger()
	}
	watcher, _ := fsnotify.NewWatcher()
	return &PromConfig{config: conf, logger: logger, configPath: confPath, watcher: watcher}
}

func (c *PromConfig) listFiles() []string {
	var paths []string
	paths = append(paths, c.configPath)
	paths = append(paths, c.config.RuleFiles...)
	for _, sc := range c.config.ScrapeConfigs {
		for _, filecfg := range sc.ServiceDiscoveryConfig.FileSDConfigs {
			paths = append(paths, filecfg.Files...)
		}
	}
	return paths
}

func (c *PromConfig) watchFiles() {
	paths := c.listFiles()
	dirs := make(map[string]int)
	for _, p := range paths {
		// Grouping
		if idx := strings.LastIndex(p, "/"); idx > -1 {
			p = p[:idx]
		} else {
			cur, err := os.Getwd()
			if err != nil {
				p = "./"
			} else {
				p = cur
			}
		}

		// Ignore the already-watched paths
		if count, ok := dirs[p]; ok {
			count++
			continue
		}
		dirs[p] = 1
		if err := c.watcher.Add(p); err != nil {
			level.Error(c.logger).Log("msg", "Error adding file path to watch", "path", p, "err", err)
		} else {
			level.Info(c.logger).Log("msg", "Add file path to watch", "path", p)
		}
	}
}

func (c *PromConfig) stop() {
	level.Info(c.logger).Log("msg", "Stopping watching...")
	done := make(chan struct{})
	defer close(done)

	// Closing the watcher will deadlock unless all events and errors are drained
	go func() {
		for {
			select {
			case <-c.watcher.Errors:
			case <-c.watcher.Events:
			// Drain all events and errors
			case <-done:
				return
			}
		}
	}()

	if err := c.watcher.Close(); err != nil {
		level.Error(c.logger).Log("msg", "Error closing file watcher", "err", err)
	}
	level.Info(c.logger).Log("msg", "File watcher stopped")
}

func (c *PromConfig) Run(ctx context.Context) {
	defer c.stop()

	ticker := time.NewTicker(time.Second * 2)
	defer ticker.Stop()

	c.watchFiles()

	for {
		select {
		case <-ctx.Done():
			return
		case event := <-c.watcher.Events:
			// fsnotify sometimes sends a bunch of events without name or operation.
			// It's unclear what they are and why they are sent - filter them out.
			if len(event.Name) == 0 {
				break
			}
			// Everything but a chmod requires rereading.
			if event.Op^fsnotify.Chmod == 0 {
				break
			}
			level.Info(c.logger).Log("msg", "Change change", "path", filepath.Clean(event.Name))
		case <-ticker.C:
			level.Info(c.logger).Log("msg", "Tick tick")
		case err := <-c.watcher.Errors:
			if err != nil {
			}
			level.Error(c.logger).Log("msg", "Error watching file", "err", err)
		}
	}

}

func main() {
	confPath := "/tmp/prometheus/config.yml"
	conf, err := config.LoadFile(confPath)
	if err != nil {
		panic(err)
	}

	w := log.NewSyncWriter(os.Stderr)
	logger := log.NewLogfmtLogger(w)
	prom := newPromConfig(conf, log.With(logger, "component", "promconfig"), confPath)
	//paths := prom.listFiles()
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cancel()
		os.Exit(1)
	}()

	prom.Run(ctx)
}

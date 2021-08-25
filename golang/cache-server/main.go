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
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"gopkg.in/alecthomas/kingpin.v2"
)

var globalCache *cache.Cache

const (
	defaultCacheFile            = "/tmp/cached"
	defaultListenAddr           = ":9010"
	defaultCacheExpiration      = "5m"
	defaultCacheClearupInterval = "10m"
	defaultWriteTimeout         = "10s"
	defaultReadTimeout          = "10s"
	// defaultMaxBodySize is the default maximum request body size, in bytes.
	// if the request body is over this size, we will return an HTTP 413 error.
	// 500 MB
	defaultMaxBodySize = 500 * 1024 * 1024
)

func main() {
	// Options
	a := kingpin.New(filepath.Base(os.Args[0]), "A simple HTTP server for caching, powered by gin and go-cache.")
	a.HelpFlag.Short('h')
	var (
		debug                bool
		listenAddr           string
		cacheFile            string
		cacheExpiration      time.Duration
		cacheCleanupInterval time.Duration
		serverReadTimeout    time.Duration
		serverWriteTimeout   time.Duration
	)
	a.Flag("listen-address",
		"Address to listen on by default ':8080'").
		Default(defaultListenAddr).
		StringVar(&listenAddr)
	a.Flag("cache-file",
		"The path to cache file, '/tmp/cache' by default").
		Default(defaultCacheFile).
		StringVar(&cacheFile)
	a.Flag("cache-expiration",
		"The cache expiration time in minutes").
		Default(defaultCacheExpiration).
		DurationVar(&cacheExpiration)
	a.Flag("cache-cleanup-interval",
		"The cache clean-up interval in minutes").
		Default(defaultCacheClearupInterval).
		DurationVar(&cacheCleanupInterval)
	a.Flag("server-readtimeout",
		"The server read timeout").
		Default(defaultReadTimeout).
		DurationVar(&serverReadTimeout)
	a.Flag("server-writetimeout",
		"The server write timeout").
		Default(defaultWriteTimeout).
		DurationVar(&serverWriteTimeout)
	a.Flag("debug-mode", "Enable debug mode").BoolVar(&debug)
	_, err := a.Parse(os.Args[1:])
	if err != nil {
		a.Usage(os.Args[1:])
		log.Fatalln("Error parsing commandline arguments", err)
	}

	// Create a cache with a default expiration time of 5 minutes, and which
	// purges expired items every 10 minutes (default value)
	globalCache = cache.New(cacheExpiration, cacheCleanupInterval)
	log.Println("Load from cached file", cacheFile)
	globalCache.LoadFile(cacheFile)

	// Create context that listens for the interrupt signal from the OS
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer func() {
		log.Println("Save to cached file", cacheFile)
		globalCache.SaveFile(cacheFile)
		stop()
	}()

	router := gin.Default()
	router.Use(gzip.Gzip(gzip.DefaultCompression, gzip.WithExcludedPathsRegexs([]string{"/debug/pprof"})))
	if !debug {
		gin.SetMode(gin.ReleaseMode)
	}

	// Handlers

	// For testing
	router.GET("/", func(c *gin.Context) {
		time.Sleep(10 * time.Second)
		c.String(http.StatusOK, "Welcome HTTP Cache Server")
	})

	// For debugging & profiling
	pprof.Register(router)

	// Set handler
	router.POST("/set/:key", func(c *gin.Context) {
		if c.Request.ContentLength > defaultMaxBodySize {
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{"status": "failed", "message": "Request body too large"})
			return
		}
		reqData, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			log.Println("Error when setting value...", err)
			c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": err.Error()})
			return
		}
		key := c.Param("key")
		// If there is an exist cache, combine then save them both
		if cacheDataInt, found := globalCache.Get(key); found {
			cacheData := cacheDataInt.([]byte)
			cacheData = append(cacheData, reqData...)
			globalCache.SetDefault(key, cacheData)
		} else {
			globalCache.SetDefault(key, reqData)
		}
		c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Set OK"})
	})

	// Get handler
	router.GET("/get/:key", func(c *gin.Context) {
		key := c.Param("key")
		if cacheDataInt, found := globalCache.Get(key); found {
			cacheData := cacheDataInt.([]byte)
			c.String(http.StatusOK, "%s", cacheData)
			globalCache.Delete(key)
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": "Unknow key"})
	})

	// Flush handler
	router.POST("/flush", func(c *gin.Context) {
		globalCache.Flush()
		c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Flush OK"})
	})

	srv := &http.Server{
		Addr:         listenAddr,
		Handler:      router,
		ReadTimeout:  serverReadTimeout,
		WriteTimeout: serverWriteTimeout,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Listen for the interrupt signal.
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	log.Println("Shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}

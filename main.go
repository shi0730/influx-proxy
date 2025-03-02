// Copyright 2021 Shiwen Cheng. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/chengshiwen/influx-proxy/backend"
	"github.com/chengshiwen/influx-proxy/service"
	"github.com/chengshiwen/influx-proxy/util"
)

var (
	ConfigFile string
	Version    bool
	GitCommit  = "not build"
	BuildTime  = "not build"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
	log.SetOutput(os.Stdout)
	flag.StringVar(&ConfigFile, "config", "proxy.json", "proxy config file")
	flag.BoolVar(&Version, "version", false, "proxy version")
	flag.Parse()
}

func main() {
	if Version {
		fmt.Printf("Version:    %s\n", backend.Version)
		fmt.Printf("Git commit: %s\n", GitCommit)
		fmt.Printf("Build time: %s\n", BuildTime)
		fmt.Printf("Go version: %s\n", runtime.Version())
		fmt.Printf("OS/Arch:    %s/%s\n", runtime.GOOS, runtime.GOARCH)
		return
	}

	cfg, err := backend.NewFileConfig(ConfigFile)
	if err != nil {
		fmt.Printf("illegal config file: %s\n", err)
		return
	}
	log.Printf("version: %s, commit: %s, build: %s", backend.Version, GitCommit, BuildTime)
	cfg.PrintSummary()

	err = util.MakeDir(cfg.DataDir)
	if err != nil {
		log.Fatalln("create data dir error")
		return
	}

	mux := http.NewServeMux()
	service.NewHttpService(cfg).Register(mux)

	server := &http.Server{
		Addr:        cfg.ListenAddr,
		Handler:     mux,
		IdleTimeout: time.Duration(cfg.IdleTimeout) * time.Second,
	}
	if cfg.HTTPSEnabled {
		log.Printf("https service start, listen on %s", server.Addr)
		err = server.ListenAndServeTLS(cfg.HTTPSCert, cfg.HTTPSKey)
	} else {
		log.Printf("http service start, listen on %s", server.Addr)
		err = server.ListenAndServe()
	}
	if err != nil {
		log.Print(err)
		return
	}
}

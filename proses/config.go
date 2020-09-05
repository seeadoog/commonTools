package proses

import (
	"fmt"
	"sync"
	"time"
)

type TcpServer struct {
	Listen   string   `json:"listen"`
	Targets  []string `json:"targets"`
	Timeout  int      `json:"timeout"`
	Retry    int      `json:"retry"`
	Debug    bool     `json:"debug"`
	UpDump   string   `json:"up_dump"`
	DownDump string   `json:"down_dump"`
}

type Config struct {
	Tcp map[string]TcpServer `json:"tcp"`
}

func (c *Config) Run() {
	wg := sync.WaitGroup{}
	for k, v := range c.Tcp {
		if v.Retry <= 0 {
			v.Retry = 3
		}
		if v.Timeout == 0 {
			v.Timeout = 5
		}

		if len(v.Targets) == 0 {
			fmt.Println("server config error: Targets cannot be nil ", k)
			continue
		}
		if v.Listen == "" {
			fmt.Println("server config error: listen cannot be nil", k)
		}

		wg.Add(1)

		server := &tcpsrv{
			addrs:   v.Targets,
			idx:     0,
			addr:    v.Listen,
			ctx:     nil,
			timeout: time.Duration(v.Timeout) * time.Second,
			retry:   v.Retry,
			logger:  &loggerImp{debug: v.Debug},
			upstreamDump:v.UpDump,
			downstreamDump:v.DownDump,
		}
		go func(k string) {
			fmt.Println("[run-server]", k, "at", v.Listen, " targets =", v.Targets)
			if err := server.runServer(); err != nil {
				fmt.Println(k + ":run server error:" + err.Error())
			}
			wg.Done()
		}(k)

	}
	wg.Wait()
}

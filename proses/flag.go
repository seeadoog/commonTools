package proses

import (
	"flag"
	"strings"
	"time"
)

func Run(){
	var ls = flag.String("ls",":80","listen addr")
	var upstream = flag.String("addr","10.1.87.70:8888","upstreams")
	var(
		timeout = flag.Duration("t",5,"timeout")
		retry = flag.Int("r",5,"retry")
		debug = flag.Bool("debug",false,"debug")
	)
	flag.Parse()
	addrs:=strings.Split(*upstream,",")

	if len(addrs) == 0{
		panic("addr cannot be empty")
	}
	for _, v := range addrs {
		if v == ""{
			panic("addr item cannot be empty")
		}
	}
		server:=&tcpsrv{
		addrs:  addrs,
		idx:     0,
		addr:    *ls,
		ctx:     nil,
		timeout: *timeout*time.Second,
		retry:   *retry,
		logger:  &loggerImp{debug:*debug},
	}
	panic(server.runServer())
}
package proses

import (
	"flag"
	"github.com/seeadoog/commonTools/ngcfg"
	"io/ioutil"
)

func Run() {
	//var ls = flag.String("ls", ":80", "listen addr")
	//var upstream = flag.String("addr", "10.1.87.70:8888", "upstreams")
	var (
		//timeout = flag.Duration("t", 5, "timeout")
		//retry   = flag.Int("r", 5, "retry")
		//debug   = flag.Bool("debug", false, "debug")
		config   = flag.String("c","prose.cfg","cfgname")

	)
	//flag.Parse()
	//addrs := strings.Split(*upstream, ",")

	//if len(addrs) == 0 {
	//	panic("addr cannot be empty")
	//}
	//for _, v := range addrs {
	//	if v == "" {
	//		panic("addr item cannot be empty")
	//	}
	//}

	data,err:=ioutil.ReadFile(*config)
	if err != nil{
		panic(err)
	}
	cfg:=Config{}
	err = ngcfg.UnmarshalFromBytes(data,&cfg)
	if err != nil{
		panic(err)
	}
	cfg.Run()

	//
	//server := &tcpsrv{
	//	addrs:   addrs,
	//	idx:     0,
	//	addr:    *ls,
	//	ctx:     nil,
	//	timeout: *timeout * time.Second,
	//	retry:   *retry,
	//	logger:  &loggerImp{debug: *debug},
	//}
	//panic(server.runServer())
}

package main

import (
	"flag"
	"net/http"
)

func main(){


	var(
		dir = flag.String("d","./","dir")
		addr = flag.String("ls",":8080","addr")
	)
	flag.Parse()
	panic(http.ListenAndServe(*addr,http.FileServer(http.Dir(*dir))))
}


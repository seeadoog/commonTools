package main

import (
	"fmt"
	"net"
	"time"
)

func main(){
	ls,err:=net.Listen("tcp4",":8080")
	if err != nil{
		panic(err)
	}

	for{
		conn,err:=ls.Accept()
		if err != nil{
			continue
		}
		go serveConn(conn)
	}
}

func serveConn(c net.Conn){
	buf:=make([]byte,4096)
	count:=0
	for{
		n,err:=c.Read(buf)
		if err != nil{
			fmt.Println("read error,",err)
			return
		}
		count+=n
		fmt.Println("recv",n,count)
		time.Sleep(1000*time.Microsecond)
	}
}


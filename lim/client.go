package main

import (
	"fmt"
	"net"
	"time"
)

func main(){
	conn,err:=net.Dial("tcp4","127.0.0.1:8080")
	if err != nil{
		panic(err)
	}
	st:=time.Now()
	data:=make([]byte,1024)
	for i:=0;i<10240;i++{
		conn.Write(data)
	}
	fmt.Println(time.Since(st))
	time.Sleep(2*time.Second)

}
package proses

import (
	"context"
	"errors"
	"net"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

type logger interface {
	Error(s string,args ...interface{})
}

type tcpsrv struct {
	addrs []string
	idx int64
	addr string
	ctx context.Context
	timeout time.Duration
	retry int
	logger
	upstreamDump string
	downstreamDump string
}

func (s *tcpsrv)getUpstream()string{
	return s.addrs[int(atomic.AddInt64(&s.idx,1))%len(s.addrs)]
}

func (s *tcpsrv)dialWithRetry(num int)(net.Conn,error){
	for i:=0 ;i<=num ;i++{
		addr:=s.getUpstream()
		conn,err:=net.DialTimeout("tcp4",addr,s.timeout)
		if err != nil{
			s.logger.Error("dail tcp error:%v",err)
			continue
		}
		return conn,nil
	}
	return nil,errors.New("no connection could be made to target server")
}

func (s *tcpsrv)runServer()error{
	ls,err:=net.Listen("tcp4",s.addr)
	if err != nil{
		return err
	}

	for{
		conn,err:=ls.Accept()
		if err != nil{
			//todo handler accept error
			s.logger.Error("accept error:%v",err)
			continue
		}
		go s.handleConn(conn)

	}
}

func (s *tcpsrv)handleConn(c net.Conn){
	upConn,err:=s.dialWithRetry(s.retry)
	if err != nil{
		c.Write([]byte(err.Error()))
		s.logger.Error("dail error",err)
		c.Close()
		return
	}




	var up,down *os.File
	if s.upstreamDump != ""{
		up,err = os.OpenFile(s.upstreamDump,os.O_WRONLY|os.O_RDONLY|os.O_CREATE|os.O_APPEND,0666)
		if err != nil{
			s.logger.Error("create upstream dump error",err)
		}
	}
	if s.downstreamDump != ""{
		down,err = os.OpenFile(s.downstreamDump,os.O_WRONLY|os.O_RDONLY|os.O_CREATE|os.O_APPEND,0666)
		if err != nil{
			s.logger.Error("create upstream dump error",err)
		}
	}

	go copyBuffer(upConn,c,up)
	copyBuffer(c,upConn,down)
}

var bufferPool = sync.Pool{}

func init(){
	bufferPool.New = func() interface{} {
		buf:=make([]byte,2048)
		return buf
	}
}


func copyBuffer(dst net.Conn,src net.Conn,dump *os.File){
	buf:=bufferPool.Get().([]byte)
	for{
		n,err:=src.Read(buf) // 读的连接出现异常，关闭写
		if err != nil{
			dst.Close()
			src.Close()
			goto end
		}
		if dump != nil{
			dump.Write(buf[:n])
		}
		_,err=dst.Write(buf[:n])
		if err != nil{
			goto end
		}
	}

	end:
		bufferPool.Put(buf)
}
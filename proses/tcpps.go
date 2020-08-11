package proses

import (
	"context"
	"errors"
	"net"
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
	return nil,errors.New("not connection could be made to target server")
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
		return
	}
	go copyBuffer(upConn,c)
	copyBuffer(c,upConn)
}

var bufferPool = sync.Pool{}

func init(){
	bufferPool.New = func() interface{} {
		buf:=make([]byte,2048)
		return buf
	}
}


func copyBuffer(dst net.Conn,src net.Conn){
	buf:=bufferPool.Get().([]byte)
	for{
		n,err:=src.Read(buf) // 读的连接出现异常，关闭写
		if err != nil{
			dst.Close()
			src.Close()
			goto end
		}
		_,err=dst.Write(buf[:n])
		if err != nil{
			goto end
		}
	}

	end:
		bufferPool.Put(buf)
}
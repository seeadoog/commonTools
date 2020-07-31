package localrpc

import (
	"fmt"
	"strconv"
	"sync"
	"sync/atomic"
)

type ApplicationServer interface {
	RegisterHandler(funcName string,handler Handler)
}

type HandleFunc func(funcName string,message Message,callback Callback)(Message,error)

type ApplicationClient interface {
	Call(funcName string,message Message,callback Callback)(Message,error)
}

type localRpcServer struct {
	callBacks sync.Map
	netCaller Client
	handlers sync.Map
	netServer Server
	logger logger
	callId int64
}

func NewRpcServer(netCaller Client,netServer Server,logger logger)*localRpcServer{
	return &localRpcServer{netCaller:netCaller,netServer:netServer,logger:logger}
}

func (c *localRpcServer)Start()error{
	c.netServer.HandleFunc(c.handle)
	return c.netServer.Start()
}

func (c *localRpcServer) RegisterHandler(funcName string, handler HandleFunc) {
	c.handlers.Store(funcName,handler)
}

func (c *localRpcServer)handle(message Message)(Message,error){
	handler,ok:=c.handlers.Load(message.GetFuncName())
	if !ok{
		return nil, fmt.Errorf("invalid func name:",message.GetFuncName())
	}
	hf:=handler.(HandleFunc)

	switch message.GetCallType() {
	case CallTypeSyncCall: // 同步调用， callback 为nil
		return hf(message.GetFuncName(),message,nil)
	case CallTypeDowncall: // 回调函数调用，直接调用回调
		call,ok:=c.callBacks.Load(message.GetCallId())
		if ok && call != nil{
			call.(Callback)(message)
		}else {
			return nil, fmt.Errorf("downcall error,callback is nil")
		}
	case CallTypeAsyncCall: // 异步调用 ，设置好回调函数
		return hf(message.GetFuncName(),message, func(msg Message) {
			msg.SetFuncName(message.GetFuncName())
			msg.SetCallId(message.GetCallId())
			_,err:=c.downCall(msg.GetFuncName(),msg)
			if err != nil && c.logger!= nil{
				c.logger.Errorf(fmt.Sprintf("error call back :%s",err.Error()))
			}
		})
	default:

	}
	return nil,fmt.Errorf("invalid call type:%v",message.GetCallType())
}

func (c *localRpcServer) downCall(funcName string, message Message) (Message, error){
	message.SetCallType(CallTypeDowncall)
	return c.netCaller.Call(funcName,message)
}

func (c *localRpcServer)generateCallId()string{
	return strconv.Itoa(int(atomic.AddInt64(&c.callId,1)))
}

func (c *localRpcServer) AsyncCall(funcName string, message Message, callback Callback) (Message, error) {
	message.SetCallType(CallTypeAsyncCall)
	if callback == nil{
		return nil, fmt.Errorf("async call callback cannot be nil")
	}
	message.SetCallId(c.generateCallId())
	c.callBacks.Store(message.GetCallId(), callback)

	return c.netCaller.Call(funcName,message)
}

func  (c *localRpcServer) Call(funcName string, message Message) (Message, error)  {
	message.SetCallType(CallTypeSyncCall)
	return c.netCaller.Call(funcName,message)
}

/**
server(message){
	a.subscribe(title,func(data){
		message.callback(data)

	})
}

 */

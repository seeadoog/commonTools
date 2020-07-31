package localrpc

import "context"

type Data interface {
	GetBytes()[]byte
	GetArgs()map[string]string
}

type CallType byte

const(
	CallTypeSyncCall = iota
	CallTypeAsyncCall
	CallTypeDowncall
)

type Message interface {
	SetCallId(id string)
	GetCallId()string
	SetFuncName(fname string)
	GetFuncName()string
	GetHeader(key string)string
	SetHeader(key,val string)
	GetCallType()CallType
	SetCallType(ctp CallType)
	GetData(idx int)[]byte
	Bytes()([]byte,error)
}

type Caller interface {
	Call(ctx context.Context,data Message)(Message,error)
	AsyncCall(ctx context.Context,data Message,callback func(message Message))(Message,error)
}
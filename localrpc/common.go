package localrpc

type Request interface {

}

type Response interface {

}

type Handler func(request Message)(Message,error)

type Server interface {
	Start()error
	HandleFunc(handler Handler)
}

type Client interface {
	Call(funcName string,request Message)(Message,error)
}
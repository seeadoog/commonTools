package netproto

type Request interface {

}

type Response interface {

}

type Server interface {
	Call(funcName string,request Request)(Response,error)
}

type Client interface {
	Call(funcName string,request Request)(Response,error)
}
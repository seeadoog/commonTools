package localrpc

type logger interface {
	Debugf(f string,args ...interface{})
	Errorf(f string,args ...interface{})
}

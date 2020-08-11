package proses

import "fmt"

type loggerImp struct {
	debug bool
}

func (l *loggerImp) Error(s string, args ...interface{}) {
	if !l.debug{
		return
	}
	fmt.Printf(s,args...)
}

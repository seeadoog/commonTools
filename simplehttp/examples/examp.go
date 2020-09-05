package main

import (
	"fmt"
	"github.com/seeadoog/commonTools/simplehttp"
)

func main() {

	//var a uint64 = 0
	//for i:=0 ;i< 10;i++{
	//	go func() {
	//		for{
	//			t:=time.Now()
	//			rsp := simplehttp.GET().Timeout(1000 * time.Millisecond).Url("http://172.31.98.182:80").Do()
	//			//rsp := simplehttp.NewRequest().Timeout(1000 * time.Millisecond).GET().Url("http://10.1.87.69:8888").Do()
	//			fmt.Println(rsp.Err(),time.Since(t))
	//			atomic.AddUint64(&a,1)
	//		}
	//	}()
	//}
	//time.Sleep(1*time.Second)
	//
	//fmt.Println(a)
	////fmt.Println(string(rsp.ReadBody()))



}
type item struct {
	w int
	v string
	counter int
}

type weightLoop struct {
	data []item
	i int64
}

func (l *weightLoop)getAddr()string{
	i:=l.i
	c:=l.data[i]

}
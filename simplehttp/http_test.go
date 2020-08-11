package simplehttp

import (
	"fmt"
	"testing"
)

func TestNewRequest(t *testing.T) {
	req:=NewRequest().GET().Url("http://www.baidu.com").Success(400).Do()
	if req.Err() != nil{
		panic(req.Err())
	}


	fmt.Println(string(req.ReadBody()))
	fmt.Println(req.ResponseJson().Get("message"))
}

func TestReq(t *testing.T){
	r:=map[string]interface{}{}
	err:=NewRequest().GET().Url("http://ws-api.xfyun.cn").Success(404).Do().Into(&r)
	if err != nil{
		panic(err)
	}
	fmt.Println(r)

	fmt.Println(string(NewRequest().GET().Url("http://ws-api.xfyun.cn").Do().ReadBody()))
}
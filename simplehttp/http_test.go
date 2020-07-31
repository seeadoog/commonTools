package simplehttp

import (
	"fmt"
	"testing"
)

func TestNewRequest(t *testing.T) {
	req:=NewRequest("GET","http://ws-api.xfyun.cn",nil).Do()
	if req.Errors() != nil{
		panic(req.Errors())
	}
	fmt.Println(string(req.ReadBody()))
	fmt.Println(req.ResponseJson().Get("message"))
}
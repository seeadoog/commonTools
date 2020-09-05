package simplehttp

import (
	"fmt"
	"testing"
	"time"
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

func TestReq2(t *testing.T){
	rsp := NewRequest().Timeout(100 * time.Millisecond).GET().Url("http://10.1.87.66:8099/hbase/ist_res_hot/f:d/4CC5779A_123457891_jiangyin20205_ist_nil_nil_nil/ist").Do()
	fmt.Println(rsp.Err())
	fmt.Println(string(rsp.ReadBody()))
	fmt.Println(rsp.SubHeader("Content-type","charset"))
}
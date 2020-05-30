package main

import (
	"encoding/json"
	"fmt"
	"github.com/seeadoog/commonTools/dnscache"
	"github.com/seeadoog/commonTools/jsonschema"
	"io/ioutil"
	"net/http"
	"time"
)

func main(){
	c:=dnscache.NewDnsCache(10*time.Second)
	req,_:=http.NewRequest("GET","https://ws-api.xfyun.cn/v2/iat",nil)
	resp,err:=c.DoHttpRequest(req)
	if err != nil{
		panic(err)
	}
	b,_:=ioutil.ReadAll(resp.Body)
	fmt.Println(string(b))
	sc:=jsonschema.Schema{}
}

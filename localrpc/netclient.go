package localrpc

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
)

func wrapError(header,msg string,err error)error{
	return fmt.Errorf("%s %s %s",header,msg,err.Error())
}

type httpNetCaller struct {
	client *http.Client
	messageBuilder func([]byte)(Message,error)
	raddr string
}

func (h *httpNetCaller)init()error{
	remote,err:=net.ResolveUnixAddr("unix",h.raddr)
	if err != nil{
		return err
	}

	h.client = &http.Client{
		Transport:&http.Transport{
			Dial: func(network, addr string) (conn net.Conn, e error) {
				local,err:=net.ResolveUnixAddr("unix","/tmp/unix_go_client")
				if err != nil{
					return nil,e
				}
				return net.DialUnix("unix",local,remote)
			},
		},
	}
	return nil
}

func (h *httpNetCaller) Call(funcName string, request Message) (Message, error) {
	bs,err:=request.Bytes()
	header:="call|"
	if err != nil{
		return nil,fmt.Errorf("%s marshal request error:%s",header,err.Error())
	}
	req,err:=http.NewRequest("POST","/",bytes.NewReader(bs))
	if err != nil{
		return nil,fmt.Errorf("%s make request error:%s",header,err.Error())
	}
	resp,err:=h.client.Do(req)
	if err != nil{
		return nil,fmt.Errorf("%s do request error:%s",header,err.Error())
	}
	bodyBytes,err:=ioutil.ReadAll(resp.Body)
	if err != nil{
		return nil,fmt.Errorf("%s read body error:%s",header,err.Error())
	}
	rspMsg,err:=h.messageBuilder(bodyBytes)
	if err != nil{
		return nil, wrapError(header,"unmarshal resp error:",err)
	}
	return rspMsg, nil
}


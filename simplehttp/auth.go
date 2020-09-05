package simplehttp

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/url"
	"time"
)

func AssemblyRequestHeader(requrl,method string,key,secret string, body []byte)map[string]string {
	u,err:=url.Parse(requrl)
	if err != nil{
		return nil
	}
	header:=map[string]string{}
	header["Host"] = u.Host
	currentTime := time.Now().UTC().Format(time.RFC1123)
	header["Date"]=  currentTime
	digest := signBody(body)
	header["Digest"] = "SHA-256="+digest
	sign := generateSignature(u.Host, currentTime,
		method, u.Path, "HTTP/1.1", secret)
	var authHeader = `hmac username="` + key + `", algorithm="hmac-sha256", headers="host date request-line", signature="` + sign + `"`
	header["Authorization"] =  authHeader
	return header
}

func generateSignature(host, date, httpMethod, requestUri, httpProto, secret string) string {
	var signatureStr string
	if len(host) != 0 {
		signatureStr = "host: " + host + "\n"
	}
	signatureStr += "date: " + date + "\n"
	signatureStr += httpMethod + " " + requestUri + " " + httpProto
	return hmacsign(signatureStr, secret)
}



func hmacsign(data, secret string) string {
	if data == ""{
		return "NIL"
	}
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(data))
	encodeData := mac.Sum(nil)
	return base64.StdEncoding.EncodeToString(encodeData)
}

func signBody(data []byte) string {
	sha := sha256.New()
	sha.Write(data)
	encodeData := sha.Sum(nil)
	return base64.StdEncoding.EncodeToString(encodeData)
}

func md5sum(data []byte) string {
	return fmt.Sprintf("%x", md5.Sum(data))
}

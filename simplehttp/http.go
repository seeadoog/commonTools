package simplehttp

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
	"unsafe"
)

type Request struct {
	r                 *http.Request
	errors            []error
	client            *http.Client
	method            string
	url               string
	body              interface{}
	headers           map[string]string
	successStatusCode []int
	timeout time.Duration
}

func (r *Request) GET() *Request {
	r.method = "GET"
	return r
}

func (r *Request) POST() *Request {
	r.method = "POST"
	return r
}

func (r *Request) PUT() *Request {
	r.method = "PUT"
	return r
}

func (r *Request) DELETE() *Request {
	r.method = "DELETE"
	return r
}

func (r *Request) HEAD() *Request {
	r.method = "HEAD"
	return r
}

func (r *Request) OPTIONS() *Request {
	r.method = "OPTIONS"
	return r
}

func (r *Request) PATCH() *Request {
	r.method = "PATCH"
	return r
}


func (r *Request) Url(u string) *Request {
	r.url = u
	return r
}

func (r *Request)Timeout(timeout time.Duration)*Request{
	r.timeout = timeout
	return r
}

func (r *Request) build() *Request {
	body := r.body
	var reader io.Reader
	switch body.(type) {
	case []byte:
		reader = bytes.NewReader(body.([]byte))
	case io.Reader:
		reader = body.(io.Reader)
	case nil:
		reader = nil
	default:
		bf := &bytes.Buffer{}
		e := json.NewEncoder(bf)
		err := e.Encode(body)
		if err != nil {
			r.errors = append(r.errors, err)
			return r
		}
		reader = bf
	}
	var req *http.Request
	var err error
	if r.timeout >0 {
		ctx,_:=context.WithTimeout(context.Background(),r.timeout)
		req ,err = http.NewRequestWithContext(ctx,r.method,r.url,reader)
	}else{
		req, err = http.NewRequest(r.method, r.url, reader)
	}
	if err != nil {
		r.errors = append(r.errors, err)
		return r
	}
	r.r = req

	for k, v := range r.headers {
		req.Header.Set(k, v)
	}
	return r
}

func NewRequest() *Request {
	return &Request{
		headers:           map[string]string{},
		successStatusCode: []int{200, 201, 204},
	}
}

func (r *Request) ApplicationJson() *Request {
	r.headers["Content-type"] = "application/json"
	return r
}

func (r *Request) Header(key, val string) *Request {
	if len(r.errors) > 0 {
		return r
	}
	r.headers[key] = val
	return r
}

func (r *Request) Client(c *http.Client) *Request {
	r.client = c
	return r
}

func (r *Request) Do() (rsp *Response) {
	r.build()
	rsp = &Response{
		errors: r.errors,
	}
	if len(r.errors) > 0 {
		return
	}
	client := r.client
	if client == nil {
		client = http.DefaultClient
	}
	resp, err := client.Do(r.r)
	if err != nil {
		rsp.errors = append(r.errors, err)
		return
	}
	rsp.statusCode = resp.StatusCode
	rsp.statusCode = resp.StatusCode
	rsp.respBody = resp.Body
	rsp.header = resp.Header
	if !in(rsp.statusCode, r.successStatusCode...) {
		rsp.errors = append(r.errors, fmt.Errorf("response status code does not expected:%d,body is:%s", rsp.statusCode, stringOf(rsp.ReadBody())))
	}
	return rsp
}


func (r *Request) Body(b interface{}) *Request {
	r.body = b
	return r
}



func (r *Request) Errors() error {
	if len(r.errors) == 0 {
		return nil
	}
	return mergeError(r.errors)
}

func in(code int, cs ...int) bool {
	for _, v := range cs {
		if code == v {
			return true
		}
	}
	return false
}

func (r *Request) Success(successCodes ...int) *Request {
	r.successStatusCode = append(r.successStatusCode, successCodes...)
	return r
}

func mergeError(errs []error) error {
	bf := bytes.Buffer{}
	for _, v := range errs {
		bf.WriteString(v.Error())
		bf.WriteString("  ")
	}
	return errors.New(stringOf(bf.Bytes()))
}

func stringOf(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

type Response struct {
	errors        []error
	statusCode    int
	respBody      io.ReadCloser
	respBodyBytes []byte
	header http.Header
}

func (r *Response) ResponseJson() *JsonElement {
	data := r.ReadBody()
	if data == nil {
		return nil
	}
	var i map[string]interface{}
	err := json.Unmarshal(data, &i)
	if err != nil {
		return nil
	}
	return &JsonElement{data: i}
}
func (r *Response) Err() error {
	if len(r.errors) == 0 {
		return nil
	}
	return mergeError(r.errors)
}

func (r *Response) ReadBody() []byte {
	if r.respBodyBytes != nil {
		return r.respBodyBytes
	}
	if r.respBody != nil {
		bytes, err := ioutil.ReadAll(r.respBody)
		if err != nil {
			r.errors = append(r.errors, err)
			return nil
		}
		r.respBodyBytes = bytes
		return bytes
	}
	return nil
}

func (r *Response)Headers()http.Header{
	return r.header
}

func (r *Response) StatusCode() int {
	return r.statusCode
}

func (r *Response) Text() string {
	return stringOf(r.ReadBody())
}

func (r *Response) Into(v interface{}) error {
	if len(r.errors) > 0 {
		return mergeError(r.errors)
	}
	b := r.ReadBody()
	err := json.Unmarshal(b, v)
	if err != nil {
		return fmt.Errorf("unmarshal body error:%w; jsonstring:%s", err, *(*string)(unsafe.Pointer(&b)))
	}
	return nil
}



func GET()*Request{
	return NewRequest().GET()
}

func POST()*Request{
	return NewRequest().POST()
}




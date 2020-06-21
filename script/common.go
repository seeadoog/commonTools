package script

import (
	"fmt"
	"strconv"
)

func String(v interface{}) string {
	switch v.(type) {
	case string:
		return v.(string)
	case bool:
		if v.(bool) {
			return "true"
		}
		return "false"
	case float64:
		return strconv.FormatFloat(v.(float64), 'f', -1, 64)
	case nil:
		return ""

	}
	return fmt.Sprintf("%v", v)
}

func Number(v interface{}) float64 {
	switch v.(type) {
	case float64:
		return v.(float64)
	case bool:
		if v.(bool) {
			return 1
		}
		return 0
	case string:
		i, err := strconv.ParseFloat(v.(string), 64)
		if err != nil {
			return i
		}
		if v.(string) == "true" {
			return 1
		}
		return 0
	}
	return 0
}

func Bool(v interface{}) bool {
	switch v.(type) {
	case float64:
		return v.(float64) > 0
	case string:
		return v.(string) == "true"
	case bool:
		return v.(bool)
	}
	return false
}
func Equal(a, b interface{}) bool {
	return String(a) == String(b)
}

type Error struct {
	Message string
}

func Errorf(s string, args ...interface{}) *Error {
	return &Error{
		Message: fmt.Sprintf(s, args...),
	}
}

type Context struct {
	kvs    map[string]interface{}
	parent *Context
}

func (c *Context) Set(key string, val interface{}) bool {
	_, ok := c.kvs[key]
	if ok {
		c.kvs[key] = val
		return true
	}
	if c.parent != nil {
		if c.parent.Set(key, val) {
			return true
		}
	}
	c.kvs[key] = val
	return true
}
func (c *Context) Get(key string) interface{} {
	v, ok := c.kvs[key]
	if ok {
		return v
	}
	if c.parent != nil {
		return c.parent.Get(key)
	}
	return nil
}

func (c *Context) Next() *Context {
	return &Context{
		kvs:    map[string]interface{}{},
		parent: c,
	}
}

type Value interface {
	Get(c *Context) interface{}
}

type Expression interface {
	Exec(c *Context) *Error
}

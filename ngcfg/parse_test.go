package ngcfg

import (
	"fmt"
	"io"
	"testing"
)

/**
aa bb
cc dd
 */
var cfg = []byte(` 
server {  # server config 
	host 127.0.0.1  # server listen host
	port 8080  # server listen port
	tps aaa bbb ccc \
		ddd eee fff \
		ggg hhh eee   # sf
     bingo df


    ans{
		name 123 
 		dfg 123 456 \
			123 123
		du "hw sdf ok thanks"
	}
    handlers{
		auth_by_lua_block "
			if ctx.param.app_id=='4cc5779a'
			then
				ctx.exit(403,'file doesnot exits')
			end
		"
		
		log_by_lua_block "
			ctx.log.info('log .... ',\"hello\")
		"	
		
		content_by_lua "
			
			ctx.writer.write(200,{
				code = 0,
				data = {
					id = ctx.response.consumer_id
				}
			})
		"
		
	}
}


`)
func Test_parse(t *testing.T) {
	e,err:=parse(cfg)
	fmt.Println(e,err)
	ser:=e.Get("server").(*Elem)
	fmt.Println(ser.Get("host"))
	fmt.Println(ser.Get("port"))
	fmt.Println(ser.Get("tps"))
	fmt.Println(ser.GetBool("bingo"))
	ans:=ser.Get("ans").(*Elem)
	fmt.Println(ans.GetNumber("name"))
	fmt.Println(ans.Get("dfg"))
	fmt.Println(ans.GetString("du"))
	handlers:=ser.Get("handlers").(*Elem)
	fmt.Println(handlers.GetString("auth_by_lua_block"))
	fmt.Println(handlers.GetString("log_by_lua_block"))
	fmt.Println(handlers.GetString("content_by_lua"))

}

type HttpScanner struct {
	r io.Reader
	buf []byte
}


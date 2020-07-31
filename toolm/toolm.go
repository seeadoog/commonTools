package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)
var i = 0
func handler(file string,w *bytes.Buffer){
	f,err:=os.Open(file)
	if err != nil{
		panic(err)
	}
	sc:=bufio.NewScanner(f)
	for sc.Scan() && i < *maxLine{
		t:=sc.Text()
		w.WriteString(fmt.Sprintf("%s:*\r\n",t))
		i++
	}
}

var(
	maxLine = flag.Int("n",4000,"maxline")
)

func main(){
	flag.Parse()
	target:=bytes.Buffer{}
	files:=flag.Args()
	flag.Args()
	for _, v := range files {
		fmt.Println("handler=>",v)
		handler(v,&target)
	}

	err:=ioutil.WriteFile("replace_list.txt",target.Bytes(),0666)
	if err != nil{
		panic(err)
	}

}

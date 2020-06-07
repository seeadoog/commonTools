package main

import (
	"bufio"
	"encoding/json"
	"flag"
	jsonscpt "github.com/seeadoog/json_script"
	"io"
	"io/ioutil"
	"os"
)

var(
	input = flag.String("f","","input file,emtpy from stdout")
	script = flag.String("s","","script")
	scriptFile = flag.String("sf","","script")
)

func getInputReader()io.Reader{
	var in io.Reader
	if *input == ""{
		in = os.Stdin
	}
	f,err:=os.Open(*input)
	if err != nil{
		panic(err)
	}
	in = f
	return in
}

func compileScript(str string)jsonscpt.Exp{
	var i interface{}
	if err:=json.Unmarshal([]byte(str),&i);err != nil{
		exp,err:=jsonscpt.ParseExp(str)
		if err != nil{
			panic(err)
		}
		return exp
	}

	exp,err:=jsonscpt.CompileExpFromJsonObject(i)
	if err != nil{
		panic(err)
	}
	return exp
}

func compile()jsonscpt.Exp{
	str:=""
	if *scriptFile != ""{
		f,err:=ioutil.ReadFile(*scriptFile)
		if err != nil{
			panic(err)
		}
		str = string(f)
	}else{
		str = *script
	}
	return compileScript(str)
}

func main(){
	flag.Parse()
	sc:=bufio.NewScanner(getInputReader())
	vm:=jsonscpt.NewVm()
	exp:=compile()
	for sc.Scan(){
		var i interface{}
		err:=json.Unmarshal(sc.Bytes(),&i)
		if err != nil{
			vm.Set("$",sc.Text())
		}else{
			vm.Set("$",i)
		}
		exp.Exec(vm)
	}
}
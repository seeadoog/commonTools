package main

import (
	"bytes"
	"fmt"
	"strconv"
	"time"
)

type printer struct {

}

func main(){
	b:=NewBar(&PercentView{},1000)

	for i:=0;i<1000;i++{
		b.Fresh(i)
		time.Sleep(300*time.Millisecond)
	}
}
// -------------------
// +==>
// ------------------

type View interface {
	Draw(percent int)
}

type Bar struct {
	max int
	view View
}

func NewBar(view View,max int)*Bar{
	return &Bar{
		max:  max,
		view: view,
	}
}

var DefaultView = &bresView{}

func (b *Bar)Fresh(currentNum int)error{
	if currentNum >b.max{
		return fmt.Errorf("fresh view error, currentNum:%d > max:%d",currentNum,b.max)
	}
	percent:= 1000*currentNum/b.max
	b.view.Draw(percent)
	return nil
}

type bresView struct {

}
//
func (b bresView) Draw(percent int) {
	bf:=bytes.Buffer{}
	bf.WriteString("|")
	r:=false
	for i:=0;i<100;i++{
		if i<percent{
			bf.WriteString("=")
		}else{
			if r{
				bf.WriteString(" ")
			}else{
				bf.WriteString(">")
				bf.WriteString(strconv.Itoa(percent))
				bf.WriteString("%")
				r= true
			}

		}
	}
	bf.WriteString("|")

	if percent == 100{
		bf.WriteString("100%")
	}
	fmt.Printf("%s\r",bf.String())
}


type PercentView struct {

}

func (p PercentView) Draw(percent int) {
	fmt.Printf("%.1f %%",float64(percent)/float64(10))
	fmt.Printf("\r")
}

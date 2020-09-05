package main

import (
	"fmt"
	"time"
)

type SpkerInfoTmp struct {
	NVectordim int32
	PfIvector [512]float32
	NSpkId   int32
	NSpkNum  int32
	NBlockNum int32
	Ntb int32
	SpkName   [128]byte
	BUpdate  byte
}

func febo(n int) int {
	f:=make([]int,101)
	f[0]=0
	f[1] =1;
	for i:=2;i<=n;i++{
		f[i] = f[i-1]+f[i-2]
	}
	return f[n]
}

func main(){
	st:=time.Now()
	for i:=0;i<1000000;i++{
		febo(100)
	}
	fmt.Println(time.Since(st))
}
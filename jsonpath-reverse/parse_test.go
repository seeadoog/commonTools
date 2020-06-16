package jsonpath_reverse

import (
	"fmt"
	"testing"
)

func Test_compileTokens(t *testing.T) {
	fmt.Println(compileTokens("a.b.crrr[12]"))
}

func TestReg(t *testing.T){
	str:="[44]"
	fmt.Println(extractTokenRegexp.MatchString(str))
	res:=extractTokenRegexp.FindAllSubmatch([]byte(str),-1)
	for i, re := range res {
		for _, s := range re {
			fmt.Println(i,string(s))
		}
	}

}
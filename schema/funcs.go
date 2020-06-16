package schema

import "strings"

func init(){
	SetFunc("split",funcSplit)
	SetFunc("join",funcJoin)
}

func funcAppend(ctx Context, args ...Value) interface{} {
	bf := strings.Builder{}
	for _, arg := range args {
		v := arg.Get(ctx)
		bf.WriteString(String(v))
	}
	return bf.String()
}

func funcAdd(ctx Context, args ...Value) interface{} {
	var sum float64 = 0
	for _, arg := range args {
		sum += Number(arg.Get(ctx))
	}
	return sum
}

func funcSplit(ctx Context, args ...Value) interface{}{
	if len(args) <2{
		return nil
	}
	str:=String(args[0].Get(ctx))
	sep:=String(args[1].Get(ctx))
	num:=-1
	if len(args)>=3{
		num = int(Number(args[2].Get(ctx)))
	}
	return strings.SplitN(str,sep,num)
}

func funcJoin(ctx Context, args ...Value) interface{}{
	if len(args) <2{
		return ""
	}
	arri,ok:=args[0].Get(ctx).([]string)
	sep:=String(args[1].Get(ctx))
	if ok{
		return strings.Join(arri,sep)
	}
	arr,ok:=args[0].Get(ctx).([]interface{})
	if !ok{
		return ""
	}
	arrs:=make([]string, len(arr))
	for i := range arr {
		arrs[i] = String(arr[i])
	}
	return strings.Join(arrs,sep)
}
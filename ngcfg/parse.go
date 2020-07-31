package ngcfg

import (
	"container/list"
	"fmt"
)

/**
aad{
	sffds sdfds
	dsfs   sdfdsf
}

 */


type Instruct struct {
	Name string
	value []string
}

const(
	statBegin = iota
	statContinue
	statStringBg
	statStringEd
	statAnno
	statSacnKey
	statStartBlock
	statEndBlock

)

const(
	valueLine = iota
	valueObject
)

type scanner struct {
	stack *list.List  // 当前节点保存堆栈
	tk []byte // 当前token
	ltks []string  // 当前行 所有token
	step func(s *scanner,c byte)error // 扫描step
	cvt int  // 当前值类型
	line int  // 扫描行
	rank int  // 扫描列
}
// 重置行号
func (s *scanner)setLine(){
	s.line++
	s.rank=0
}


func parse(b []byte)(*Elem,error){
	sc:=&scanner{
		stack: list.New(),
		tk:    make([]byte,0,5),
		ltks:  make([]string,0,2),
		step:  stepBegin,
		cvt:   valueLine,
		line:1,
		rank:0,
	}
	sc.stack.PushBack(NewElem())
	for _, v := range b {
		sc.rank ++
		if err:=sc.step(sc,v);err !=nil{
			return nil,err
		}
	}
	if sc.stack.Len() != 1{
		return nil,fmt.Errorf("} does not match {")
	}
	return sc.stack.Front().Value.(*Elem),nil
}

//}
func stepObEnd(s *scanner)error{
	if s.stack.Len() <=1{
		return fmt.Errorf("invalid end }:at line:%d rank:%d",s.line,s.rank)
	}
	s.cvt = valueObject
	s.stack.Remove(s.stack.Back())
	s.step = stepEndOb
	return nil
}
//ssss
func stepBegin(s *scanner,c byte)error{
	if isSpace(c){
		return nil
	}
	switch c {
	case '{':
		return fmt.Errorf("invalid begin value:{, line:%d,rank:%d",s.line,s.rank)
	case '#':
		s.step = stepAnno
		return nil
	case '\r','\n':
		if c=='\n'{
			s.setLine()
		}
		return nil
	case '}':
		return stepObEnd(s)
	}
	s.cvt = valueLine
	s.tk = append(s.tk,c)
	s.step = stepContinue
	return nil
}

func stepEndOb(s *scanner,c byte)error{
	if isSpace(c){
		return nil
	}
	switch c {
	case '#':
		s.step = stepAnno
		return nil
	case '\r':
		s.step = stepEscap1
		return nil
	case '\n':
		return stepEscap2(s)
	}
	return fmt.Errorf("invalid c after },line:%d rank:%d",s.line,s.rank)
}

func isSpace(c byte)bool{
	return c == ' ' || c =='\t'
}


func appendLine(s *scanner){
	if len(s.tk) >0{
		s.ltks = append(s.ltks,string(s.tk))
		s.tk = s.tk[:0]
	}
}
//weewr{
func stepContinue(s *scanner,c byte)error{
	switch c {
	case '#':
		appendLine(s)
		s.step = stepAnno
		return nil
	case '\r':
		s.step = stepEscap1
		return nil
	case '\n':
		return stepEscap2(s)
	case '}':
		  return fmt.Errorf("invalid } at start block line:%d,rank:%d",s.line,s.rank)
	case '{':
		s.cvt = valueObject
		e:=NewElem()
		tope:=s.stack.Back().Value.(*Elem)
		appendLine(s)
		if len(s.ltks)!=1{
			return fmt.Errorf("invalid begin value of {,line:%d,rank:%d",s.line,s.rank)
		}
		key:= s.ltks[0]
		if err:=tope.Set(key,e);err != nil{
			return err
		}
		s.stack.PushBack(e)
		s.ltks = []string{}
		s.step = stepStartObject

		return nil
	case '\\':
		s.step = stepEcpSep
		return nil
	case '"':
		s.step = stepInstring
		return nil
	}
	if isSpace(c){
		if len(s.tk) >0{
			s.ltks = append(s.ltks,string(s.tk))
			s.tk = s.tk[:0]
		}
	}else{
		s.tk = append(s.tk,c)
	}
	s.step = stepContinue
	return nil
}
func stepInstring(s *scanner,c byte)error{
	if c == '\n'{
		s.setLine()
	}
	if c == '\\'{
		s.step = stepEcpNext
		return nil
	}

	if c == '"'{
		s.step =stepContinue
		return nil
	}
	s.tk = append(s.tk,c)
	return nil
}

func stepEcpNext(s *scanner,c byte)error{
	s.tk = append(s.tk,c)
	s.step =stepInstring
	return nil
}

func stepEcpSep(s *scanner,c byte)error{
	if isSpace(c){
		return nil
	}
	switch c {
	case '\r','\n':
		if c=='\n'{
			s.setLine()
		}
		return nil
	default:
		if err:=stepContinue(s,c);err !=nil{
			return err
		}
		//s.step = stepContinue
	}
	return nil
}
//{ ....\r\n
func stepStartObject(s *scanner,c byte)error{
	if isSpace(c){
		return nil
	}
	switch c {
	case '#':
		s.step = stepAnno
		return nil
	case '\r':
		s.step = stepEscap1
		return nil
	case '\n':
		return stepEscap2(s)
	}
	return fmt.Errorf("invalid %v after {", rune(c))
}

func stepAnno(s *scanner,c byte)error{
	if c == '\r'{
		s.step = stepEscap1
	}else if c == '\n'{
		return stepEscap2(s)
	}
	return nil
}

//\r
func stepEscap1(s *scanner,c byte)error{
	if c == '\n'{
		return stepEscap2(s)
	}else{
		return fmt.Errorf("invald line sep")
	}
}

func stepEscap2(s *scanner)error{
	s.setLine()
	if s.cvt == valueObject{
		s.step = stepBegin
		s.ltks = []string{}
		return nil
	}else{
		if len(s.tk) >0{
			s.ltks = append(s.ltks,string(s.tk))
			s.tk = s.tk[:0]
		}
		if len(s.ltks) >0{
			if s.stack.Len() ==0{
				return fmt.Errorf("invalid stack")
			}
			tope:=s.stack.Back().Value.(*Elem)
			err:=tope.Set(s.ltks[0],s.ltks[1:])
			if err != nil{
				return err
			}
		}
		//s.ltks = s.ltks[:0]
		s.step = stepBegin
		s.ltks = []string{}
	}
	return nil
}


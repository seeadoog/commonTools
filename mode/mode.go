package mode


type Mode struct {
	modes map[string]int
	current int
}

func (m *Mode)AddMode(name string){
	_,ok:=m.modes[name]
	if ok{
		return
	}
	 m.current <<=1
	 m.modes[name] = m.current
}
// a,b,c,d,e,f

const(
	a = 0x1
	b = 0x2
	c = 0x4
	d = 0x8
	e = 0x10
)


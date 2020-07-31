package ngcfg

type MapElem struct {
	Key  string
	Val  interface{}
	next *MapElem
	pre  *MapElem
	l *LinkedMap
}

func (m *MapElem)Next()*MapElem{
	if m.next == m.l.back{
		return nil
	}
	return m.next
}

func NewLinkedMap()*LinkedMap{
	m:=&LinkedMap{
		data: map[string]*MapElem{},
		front: nil,
		back:  nil,
	}
	m.front = new(MapElem)
	m.back = new(MapElem)
	m.front.next = m.back
	m.back.pre = m.front
	return m
}

type LinkedMap struct {
	data map[string]*MapElem
	front *MapElem
	back *MapElem
}

func (m *LinkedMap)Len()int{
	return len(m.data)
}
//front 1->2->3->back
func (m *LinkedMap)pushBack(e *MapElem){
	pb:=m.back.pre
	pb.next = e
	e.next = m.back
	m.back.pre = e
	e.pre = pb
}


func (m *LinkedMap)Set(key string,val interface{}){
	e,ok:=m.data[key]
	if ok{
		e.Val = val
		return
	}
	e = &MapElem{Val:val,l:m,Key:key}
	m.pushBack(e)
	m.data[key]=e
}

func (m *LinkedMap)Get(key string)(interface{},bool){
	v,ok:=m.data[key]
	if !ok{
		return nil,false
	}
	return v.Val,ok
}

func (m *LinkedMap)MapItem()*MapElem{
	return m.front.next
}
package lrucache

import "sync"

type LRUCache interface {
	Get(key string)interface{}
	Set(key string,v interface{})
}

type node struct {
	pre *node
	next *node
	key string
	val interface{}
}

type lrucacheImp struct {
	head *node
	tail *node
	cap int
	cache map[string]*node
	size int
}

func NewLRUCache(cap int)LRUCache{
	l:= &lrucacheImp{
		head:  nil,
		tail:  nil,
		cap:   cap,
		cache: make(map[string]*node,cap*2),
	}
	l.init()
	return l
}

func (l *lrucacheImp) Get(key string) interface{} {
	return l.get(key)
}

func (l *lrucacheImp) Set(key string, v interface{}) {
	l.set(key,v)
}

func (l *lrucacheImp)init(){
	l.head = &node{}   // 头节点与尾节点各初始化一个哨兵，方便处理节点
	l.tail = &node{}
	l.head.next = l.tail
	l.tail.pre = l.head
}

func (l *lrucacheImp)remove(n *node){
	n.pre.next = n.next
	n.next.pre = n.pre
}


func (l *lrucacheImp)removeTail(){
	l.tail = l.tail.pre
	l.tail.next = nil
}

func (l*lrucacheImp)addToHead(n *node){
	l.head.next.pre = n
	n.next  = l.head.next
	n.pre = l.head
	l.head.next = n
}

func (l *lrucacheImp)moveToHead(n *node){
	l.remove(n)
	l.addToHead(n)
}

func (l *lrucacheImp)get(key string)interface{}{
	n,ok:=l.cache[key]
	if !ok{
		return nil
	}
	l.moveToHead(n)
	return n.val
}

func (l *lrucacheImp)set(key string,val interface{}){
	n,ok:=l.cache[key]
	if ok{ // 当cache 中存在 是，更新 值并移动到head
		n.val = val
		l.moveToHead(n)
		return
	}
	n = &node{key: key,val:val}
	l.size++
	l.cache[key] = n
	l.addToHead(n)
	if l.size > l.cap{
		delete(l.cache,l.tail.pre.key) // 删除尾节点。
		l.removeTail()
		l.size --
	}
}


type concurrencyLRUCache struct {
	cache LRUCache
	lock sync.Mutex
}

func (c *concurrencyLRUCache)Get(key string)interface{}{
	c.lock.Lock()
	v:=c.cache.Get(key)
	c.lock.Unlock()
	return v
}

func (c *concurrencyLRUCache)Set(key string,val interface{}){
	c.lock.Lock()
	c.cache.Set(key,val)
	c.lock.Unlock()
}


func NewConcurrencyLRUCache(cap int)LRUCache{
	return &concurrencyLRUCache{
		cache: NewLRUCache(cap),
	}
}
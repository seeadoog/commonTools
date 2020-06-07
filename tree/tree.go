package main

import "fmt"

func main()  {
	t:=btree{
		less: func(i, j interface{}) bool {
			return i.(int) > j.(int)
		},
	}
	t.add(5)
	t.add(2)
	t.add(6)
	t.add(1)
	t.add(3)
	t.foreach(func(i interface{}) {
		fmt.Println(i)
	})
}



type treeNode struct {
	left *treeNode
	right *treeNode
	parent *treeNode
	data interface{}
}

type less func(i,j interface{})bool

type btree struct {
	root *treeNode
	less less
}

func (t *btree)add(data interface{}){
	if t.root == nil{
		t.root = new(treeNode)
		t.root.data = data
		return
	}

	n:=new(treeNode)
	n.data = data
	t.addNode(t.root,n)
}

func (t *btree)foreach(f func(interface{})){
	t.tral(t.root,f)
}

func (t *btree)tral(r *treeNode,f func(interface{})){
	if r== nil{
		return
	}
	t.tral(r.left,f)
	f(r.data)
	t.tral(r.right,f)
}

func (t *btree)addNode(r,node *treeNode){
	if t.less(r.data,node.data){
		if r.left ==nil{
			r.left = node
			node.parent = r.left
			return
		}
		t.addNode(r.left,node)
		return
	}
	if r.right == nil{
		r.right = node
		node.parent = r.right
		return
	}
	t.addNode(r.right,node)
}


// F-Base
// GPL 3.0 License - bingxio <3106740988@qq.com>
package main

import (
	"fmt"
)

/*
	Tree Structure:

	  	   Tree(T)
		      |
		  	  v					   5 Leaf
	NodeA,  NodeB,  NodeC..		1500 Data
	  |       |		  |
	  v       v		  v
	LA, LB..  L..	  L..
	  \		  |		 /
	  		  v
			Leaf(L)				3 Row
			  |
			  v
			Rows(R)				100 Data

	Table:
		From 	uint64
		At   	uint8
		Rows 	uint64

	1. Rows / 100 =  Row
	2. Row  / 3   =  Leaf
	3. Leaf / 5   =  Node

		T -> Node
*/
type Tree struct {
	Node []Node
}

// Node
type Node struct {
	Leaf []Leaf
}

// Leaf
type Leaf struct {
	Data [][]Row
}

// Stringer : For tree
func (t Tree) Stringer(p int) string {
	x := 0
	y := 0
	for _, v := range t.Node {
		x += len(v.Leaf)
		for _, b := range v.Leaf {
			y += len(b.Data)
		}
	}
	return fmt.Sprintf(
		"<name: '%s' node: %d leaf: %d rows: %d>",
		GlobalEm.tb[p].Name,
		len(t.Node),
		x,
		y,
	)
}

// Back
func (t Tree) BackNode() *Node {
	return &t.Node[len(t.Node)-1]
}

// Back
func (n Node) BackLeaf() *Leaf {
	return &n.Leaf[len(n.Leaf)-1]
}

// Back
func (l Leaf) BackRows() *[]Row {
	return &l.Data[len(l.Data)-1]
}

// Iterator
func (t Tree) Iter(f func(int, Node)) {
	for i, v := range t.Node {
		f(i, v)
	}
}

// Iterator
func (n Node) Iter(f func(int, Leaf)) {
	for i, v := range n.Leaf {
		f(i, v)
	}
}

// Iterator
func (l Leaf) Iter(f func(int, []Row)) {
	for i, v := range l.Data {
		f(i, v)
	}
}

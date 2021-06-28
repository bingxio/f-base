// F-Base
// GPL 3.0 License - bingxio <3106740988@qq.com>
package main

import "fmt"

/*
	Tree Structure:

	  	   Tree(T)
		      |
		  	  v
	NodeA,  NodeB,  NodeC..		5 Leaf
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
	Data []Row
}

// Stringer : For tree
func (t Tree) Stringer() string {
	return fmt.Sprintf("<tree: %d node>", len(t.Node))
}

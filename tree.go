// F-Base
// GPL 3.0 License - bingxio <3106740988@qq.com>
package main

type Tree struct {
}

type Node struct {
	Leaf []Leaf
}

type Leaf struct {
	Data Row
}

// ReadTable : Data reader for table rows in memory
func ReadTable(t *Table) (*Tree, error) {
	return nil, nil
}

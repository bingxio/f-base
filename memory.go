// F-Base
// GPL 3.0 License - bingxio <3106740988@qq.com>
package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

var (
	GlobalMem Memory // Mem
)

/*
	Read Tbs
		|
		v
	Read Tables
		|
		v
  Memory(Buffer)
		|
		v
		G
*/
type Memory struct {
	Tbs  uint    // Count of tbs
	Tb   []Table // Simple table
	Tree []Tree  // Trees of each table rows
}

// Dissemble : Dis
func (m Memory) Dissemble() {
	for _, v := range m.Tree {
		fmt.Println(v.Stringer())
	}
}

// NewMemory : To new memory with em and bytes buffer
func NewMemory(em Em) error {
	GlobalMem.Tbs = uint(len(em.tb))
	GlobalMem.Tb = em.tb

	rows := []Row{}

	// Buf
	for i := 0; i < len(em.tb); i++ {
		t := em.tb[i]
		p := t.From
		// Read rows
		for j := 0; j < int(t.Rows); j++ {
			r, err := ReadRow(int64(p))
			if err != nil {
				return err
			}
			// Push
			rows = append(rows, *r)
			p += RowSize
		}
		// Tree
		l := len(rows)
		if l > 100 {
			// Many Leaf
			fmt.Println(l, float64(l)/3)
		} else {
			// One Leaf, One Node
			GlobalMem.Tree = append(GlobalMem.Tree, Tree{
				Node: []Node{
					{
						Leaf: []Leaf{
							{
								Data: rows,
							},
						},
					},
				},
			})
		}
		// Clear
		rows = rows[:0]
	}
	return nil
}

// ReadRow : Read bytes to row structure
func ReadRow(p int64) (*Row, error) {
	b := make([]byte, RowSize)
	_, err := GlobalEm.file.ReadAt(b, p)
	if err != nil {
		return nil, err
	}
	r := Row{}
	err = gob.NewDecoder(
		bytes.NewBuffer(b)).Decode(&r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

// QuitMemory : Clean cache and save buffer to file
func QuitMemory() error {
	// Close file
	GlobalEm.file.Close()
	return nil
}

// InsertRow : Insert row at specified location
func InsertRow(offset uint64, r Row) error {
	return nil
}

// SelectAll : Query all rows
func SelectAll(rows uint64) ([]Row, error) {
	return nil, nil
}

// SelectOne : Query the row at the specified location
func SelectOne(offset uint64) (*Row, error) {
	return nil, nil
}

// SelectRange : Rows within query scope
func SelectRange(from, to uint64) ([]Row, error) {
	return nil, nil
}

// UpdateRow : Updates the specified row
func UpdateRow(offset uint64, n Row, verify ...Row) error {
	return nil
}

// DeleteOne : Delete the row at the specified location
func DeleteOne(offset uint64, verify ...string) error {
	if len(verify) != 0 {
		fmt.Println("delete verify row data")
	}
	return nil
}

// DeleteAll : Delete all rows of the table
func DeleteAll(from, to uint64) error {
	return nil
}

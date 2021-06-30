// F-Base
// GPL 3.0 License - bingxio <3106740988@qq.com>
package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"strings"
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
	if len(m.Tree) == 0 {
		fmt.Println("<empty tree>")
	} else {
		for k, v := range m.Tree {
			fmt.Println(v.Stringer(k))
		}
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
			// Row
			row := DecimalPlaces(float64(l), 100)
			// Leaf
			leaf := DecimalPlaces(float64(row), 3)
			// Node
			node := DecimalPlaces(float64(leaf), 5)

			// log.Printf("- %d rows: %d leaf: %d node: %d\n", l, row, leaf, node)
			// R
			tr := make([][]Row, row)
			tp := 0
			for p := 100; p <= row*100; p += 100 {
				if p > l {
					tr[tp] = rows[p-100:]
				} else {
					tr[tp] = rows[p-100 : p]
				}
				tp += 1
			}
			// L
			tl := []Leaf{}
			tp = 0
			for p := 0; p < leaf; p++ {
				if (tp + 3) > len(tr) {
					r := [][]Row{}
					for i := tp; i < len(tr); i++ {
						r = append(r, tr[i])
						tp++
					}
					tl = append(tl, Leaf{r})
				} else {
					tl = append(tl, Leaf{
						Data: [][]Row{tr[tp], tr[tp+1], tr[tp+2]}, // 3
					})
					tp += 3
				}
			}
			// N
			tn := []Node{}
			tp = 0
			for p := 0; p < node; p++ {
				if (tp + 5) > len(tl) {
					r := []Leaf{}
					for i := tp; i < len(tl); i++ {
						r = append(r, tl[i])
						tp++
					}
					tn = append(tn, Node{r})
				} else {
					tn = append(tn, Node{
						Leaf: []Leaf{
							tl[tp], tl[tp+1], tl[tp+2], tl[tp+3], tl[tp+4], // 5
						},
					})
					tp += 5
				}
			}
			// T
			GlobalMem.Tree = append(GlobalMem.Tree, Tree{tn})
		} else {
			// One Leaf, One Node
			GlobalMem.Tree = append(GlobalMem.Tree, Tree{
				Node: []Node{
					{
						Leaf: []Leaf{
							{
								Data: [][]Row{
									rows,
								},
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

// DecimalPlaces : Add decimal spaces, 34.24555 -> 35
func DecimalPlaces(raw, to float64) int {
	n := raw / to
	nf := fmt.Sprintf("%f", n)
	r := int(n)
	// Add
	if nf[strings.IndexRune(nf, '.')+1] != '0' {
		r += 1
	}
	return r
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

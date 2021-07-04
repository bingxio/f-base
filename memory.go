// F-Base
// GPL 3.0 License - bingxio <3106740988@qq.com>
package main

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"strings"
	"time"
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
		for k, v := range m.Tree { // Stringer
			fmt.Println(v.Stringer(k))
		}
	}
}

// NewMemory : To new memory with em and bytes buffer
func NewMemory(em Em) error {
	GlobalMem.Tbs = uint(len(em.tb))
	GlobalMem.Tb = em.tb

	var rows []Row

	for i := 0; i < len(em.tb); i++ { // Buf
		t := em.tb[i]
		p := t.From
		for j := 0; j < int(t.Rows); j++ { // Read rows
			r, err := ReadRow(int64(p))
			if err != nil {
				return err
			}
			rows = append(rows, *r) // Push
			p += RowSize
		}
		l := len(rows) // Tree
		if l > 100 {
			row := DecimalPlaces(float64(l), 100)   // Row
			leaf := DecimalPlaces(float64(row), 3)  // Leaf
			node := DecimalPlaces(float64(leaf), 5) // Node

			// log.Printf("- %d rows: %d leaf: %d node: %d\n", l, row, leaf, node)

			tr := make([][]Row, row) // R
			tp := 0
			for p := 100; p <= row*100; p += 100 {
				if p > l {
					tr[tp] = rows[p-100:]
				} else {
					tr[tp] = rows[p-100 : p]
				}
				tp += 1
			}
			var tl []Leaf // L
			tp = 0
			for p := 0; p < leaf; p++ {
				if (tp + 3) > len(tr) {
					var r [][]Row
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
			var tn []Node // N
			tp = 0
			if tl == nil {
				return errors.New("unhandled error")
			}
			for p := 0; p < node; p++ {
				if (tp + 5) > len(tl) {
					var r []Leaf
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
			GlobalMem.Tree = append(GlobalMem.Tree, Tree{tn}) // T
		} else {
			var r [][]Row // One leaf, one node
			l := []Leaf{
				{r},
			}
			n := []Node{
				{l},
			}
			GlobalMem.Tree = append(GlobalMem.Tree, Tree{n})
		}
		rows = nil // Clear
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
	err := GlobalEm.file.Close()
	if err != nil {
		return err
	} // Close file
	return nil
}

// DecimalPlaces : Add decimal spaces, 34.24555 -> 35
func DecimalPlaces(raw, to float64) int {
	n := raw / to
	nf := fmt.Sprintf("%f", n)
	r := int(n)
	if nf[strings.IndexRune(nf, '.')+1] != '0' { // Add
		r += 1
	}
	return r
}

// CountRows : Return count of all rows in tables
func (m Memory) CountRows() uint64 {
	var p uint64
	for _, v := range m.Tb {
		p += v.Rows
	}
	return p
}

// GetBackOffset : Return the back offsets in tables of rows size
func (m Memory) GetBackOffset() uint64 {
	var p = uint64(tbs + (TableSize * len(m.Tb)))
	for _, v := range m.Tb {
		p += v.Rows * RowSize
	}
	return p
}

// NewTable : Add new table and return it
func (m *Memory) NewTable(name string) Table {
	t := Table{
		Name: func() [20]byte {
			x := [20]byte{}
			for i, v := range name {
				if i >= 20 {
					break
				}
				x[i] = byte(v)
			}
			return x
		}(),
		Created: uint32(time.Now().Unix()),
		From:    m.GetBackOffset(),
		Rows:    0,
		At:      uint8(len(m.Tb) + 1),
	}
	m.Tbs += 1
	m.Tb = append(m.Tb, t)
	return t
}

// Insert : SE ? *E
func (m *Memory) Insert(at uint8, fields []string) {
	m.Tb[at].Rows += 1
	if int(at) == len(m.Tree) { // New Tree
		r := []Row{
			{fields},
		}
		l := []Leaf{
			{Data: [][]Row{r}},
		}
		n := []Node{
			{l},
		}
		m.Tree = append(m.Tree, Tree{n})
		return
	}
	r := m.Tree[at].BackNode().BackLeaf().BackRows() // Back
	*r = append(*r, Row{fields})                     // Rows
}

// Selector : GT ? <S> <V>
func (m *Memory) Selector(at, s uint8, v string) {
	e := func(r Row) {
		if r.Len() < s {
			return
		}
		for i, d := range r.Data {
			if i+1 == int(s) && d == v {
				fmt.Println(r.Stringer())
			}
		}
	}
	f := func(i int, r []Row) { // R
		for _, v := range r {
			e(v)
		}
	}
	t := m.Tree[at]              // T
	t.Iter(func(i int, n Node) { // N
		n.Iter(func(i int, l Leaf) { l.Iter(f) }) // L
	})
}

// SelectAll : GE ?
func (m *Memory) SelectAll(at uint8) {
	f := func(i int, r []Row) { // R
		for _, v := range r {
			fmt.Println(v.Stringer()) // Show
		}
	}
	t := m.Tree[at]              // T
	t.Iter(func(i int, n Node) { // N
		n.Iter(func(i int, l Leaf) { l.Iter(f) }) // L
	})
}

// SelectOne : GE ? <F>
func (m *Memory) SelectOne(at uint8, p uint64) {
	f := func(i int, r []Row) {
		for k, _ := range r {
			fmt.Println(k)
		}
	}
	t := m.Tree[at]
	if p > m.Tb[at].Rows {
		return
	}
	n := DecimalPlaces(float64(p), 1500)
	if n > 0 {
		n -= 1
	}
	np := t.Node[n]
	np.Iter(func(i int, l Leaf) { l.Iter(f) })
}

// SelectRange : GE ? <F> -> <T>
func (m *Memory) SelectRange(at uint8, f, t uint64) {}

func (m *Memory) Update(at uint8, n, v string, p uint64) {}

func (m *Memory) Delete(at uint8, p uint64) {}

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
		for k, v := range m.Tree { // Stringer
			fmt.Println(v.Stringer(k))
		}
	}
}

// NewMemory : To new memory with em and bytes buffer
func NewMemory(em Em) error {
	GlobalMem.Tbs = uint(len(em.tb))
	GlobalMem.Tb = em.tb

	rows := []Row{}

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
			tl := []Leaf{} // L
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
			tn := []Node{} // N
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
			GlobalMem.Tree = append(GlobalMem.Tree, Tree{tn}) // T
		} else {
			GlobalMem.Tree = append(GlobalMem.Tree, Tree{ // One leaf, one node
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
	GlobalEm.file.Close() // Close file
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

// Insert : SE ? *E
func (m *Memory) Insert(at uint8, fields []string) {
	r := m.Tree[at].BackNode().BackLeaf().BackRows() // Back
	*r = append(*r, Row{fields})                     // Rows
}

// Selector : GT ? <S> <V>
func (m *Memory) Selector(at uint8, s, v string) error {
	return nil
}

func (m *Memory) SelectAll(at uint8)           {}
func (m *Memory) SelectOne(at uint8, p uint64) {}

func (m *Memory) SelectRange(at uint8, f, t uint64) {}

func (m *Memory) Update(at uint8, n, v string, p uint64) {}

func (m *Memory) Delete(at uint8, p uint64) {}

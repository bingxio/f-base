// F-Base
// GPL 3.0 License - bingxio <3106740988@qq.com>
package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
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
	Tbs uint    // Count of tbs
	Tb  []Table // Simple table
	Tr  []Tree  // Trees of each table rows
}

// Dissemble : Dis
func (m Memory) Dissemble() {
	if len(m.Tr) == 0 {
		fmt.Println("<empty tree>")
	} else {
		for k, v := range m.Tr { // Stringer
			fmt.Println(v.Stringer(k))
		}
	}
}

// NewMemory : To new memory with em and bytes buffer
func NewMemory(em Em) error {
	GlobalMem.Tbs = uint(len(em.Tb))
	GlobalMem.Tb = em.Tb

	var rows []Row

	for i := 0; i < len(em.Tb); i++ { // Buf
		t := em.Tb[i]
		for j := 0; j < int(t.Rows); j++ { // Read rows
			r, err := ReadRow()
			if err != nil {
				return err
			}
			rows = append(rows, *r) // Push
		}
		l := len(rows) // Tr
		if l > 100 {
			row := DecimalPlaces(float64(l), 100)   // R
			leaf := DecimalPlaces(float64(row), 3)  // L
			node := DecimalPlaces(float64(leaf), 5) // N

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

			tl, tp := BuildLeaf(tp, leaf, tr)             // L
			tn := BuildNode(tp, node, tl)                 // N
			GlobalMem.Tr = append(GlobalMem.Tr, Tree{tn}) // T
		} else {
			GlobalMem.Tr = append(GlobalMem.Tr, Tree{ // One leaf, one node
				Node: []Node{
					{
						Leaf: []Leaf{
							{
								Data: [][]Row{rows},
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
func ReadRow() (*Row, error) {
	b := make([]byte, RowSize)
	_, err := GlobalEm.File.Read(b)
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

// BuildLeaf : Build leaf
func BuildLeaf(tp, leaf int, tr [][]Row) ([]Leaf, int) {
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
	return tl, tp
}

// BuildNode : Build node
func BuildNode(tp, node int, tl []Leaf) []Node {
	var tn []Node // N
	tp = 0
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
	return tn
}

// Write : Clean cache and save buffer to File
func (m Memory) Write() error {
	n := GlobalEm.File.Name()
	_, err := GlobalEm.File.Stat()
	if err == nil {
		if err := os.Remove(n); err != nil { // Delete
			return err
		}
		if err := GlobalEm.File.Close(); err != nil {
			return err
		}
	}
	f, err := os.OpenFile(n, DbFlag, 0644) // Create
	if err != nil {
		return err
	}
	b := bytes.NewBuffer([]byte{})
	err = gob.NewEncoder(b).Encode(uint8(m.Tbs))
	if err != nil {
		return err
	}
	_, _ = f.Write(b.Bytes()) // Tbs
	b.Reset()
	for _, v := range m.Tb {
		err = gob.NewEncoder(b).Encode(v)
		if err != nil {
			return err
		}
		_, _ = f.Write(b.Bytes()) // Table
		_, _ = f.Write(make([]byte, TableSize-b.Len()))
		b.Reset()
	}

	// fmt.Println(m.Tbs, m.Tb, m.Tr)

	errC := make(chan string)
	done := make(chan bool)
	go func() {
		rp := func(i int, r []Row) { // R
			for _, v := range r {
				err := gob.NewEncoder(b).Encode(v)
				if err != nil {
					errC <- err.Error()
				} else {
					_, _ = f.Write(b.Bytes())
					_, _ = f.Write(make([]byte, RowSize-b.Len())) // Row
					b.Reset()
				}
			}
		}
		for _, v := range m.Tr { // T
			v.Iter(func(i int, n Node) { // N
				n.Iter(func(i int, l Leaf) { l.Iter(rp) }) // L
			})
		}
		done <- true
	}()
	select {
	case err := <-errC:
		fmt.Println("ERROR: ", err, "!!")
	case <-done:
		if err := f.Close(); err != nil {
			return err
		}
	}
	return nil
}

// DecimalPlaces : Add decimal spaces, 34.24555 -> 35
func DecimalPlaces(raw, to float64) int {
	n := raw / to
	nf := fmt.Sprintf("%f", n)
	r := int(n)
	h := false
	for _, v := range nf {
		if h && (v >= 49 && v <= 57) { // 1 -> 9
			r += 1
			break
		}
		if v == '.' {
			h = true
		}
	}
	return r
}

// NewTable : Add new table and return it
func (m *Memory) NewTable(name string) Table {
	f := func() [20]byte { // String to bytes
		x := [20]byte{}
		for i, v := range name {
			if i >= 20 {
				break
			}
			x[i] = byte(v)
		}
		return x
	}
	t := Table{
		Name:    f(),
		Created: uint32(time.Now().Unix()),
		Rows:    0,
		At:      uint8(len(m.Tb) + 1),
	}
	m.Tbs += 1
	m.Tb = append(m.Tb, t)
	return t
}

// Insert : SE ? *E
func (m *Memory) Insert(at uint8, fields []string) Result {
	m.Tb[at].Rows += 1
	w := Row{fields}
	res := SingleResult{
		Row:    w,
		Offset: m.Tb[at].Rows,
	}
	if int(at) == len(m.Tr) { // New Tr
		r := []Row{
			w,
		}
		l := []Leaf{
			{Data: [][]Row{r}},
		}
		n := []Node{
			{l},
		}
		m.Tr = append(m.Tr, Tree{n})
	} else {
		r := m.Tr[at].BackNode().BackLeaf().BackRows() // Back
		*r = append(*r, w)                             // Rows
	}
	return res
}

// Selector : GT ? <S> <V>
func (m *Memory) Selector(at, s uint8, v string) Result {
	res := MultipleResult{}
	e := func(r Row) {
		if r.Len() < s {
			return
		}
		for i, d := range r.Data {
			if i+1 == int(s) && d == v {
				res.Rows += 1
				res.Data = append(res.Data, r)
				res.Offset = append(res.Offset, res.Rows)
			}
		}
	}
	f := func(i int, r []Row) { // R
		for _, v := range r {
			e(v)
		}
	}
	t := m.Tr[at] // T
	t.Iter(func(i int, n Node) { // N
		n.Iter(func(i int, l Leaf) { l.Iter(f) }) // L
	})
	return res
}

// SelectAll : GE ?
func (m *Memory) SelectAll(at uint8) Result {
	res := MultipleResult{}
	f := func(i int, r []Row) { // R
		for _, v := range r {
			res.Rows += 1
			res.Data = append(res.Data, v)
			res.Offset = append(res.Offset, res.Rows)
		}
	}
	t := m.Tr[at] // T
	t.Iter(func(i int, n Node) { // N
		n.Iter(func(i int, l Leaf) {
			l.Iter(f)
		}) // L
	})
	return res
}

// SelectOne : GE ? <F>
func (m *Memory) SelectOne(at uint8, p uint64) Result {
	t := m.Tr[at]
	if p > m.Tb[at].Rows {
		return nil
	}
	n := DecimalPlaces(float64(p), 1500)
	if n != 0 {
		n -= 1
	}
	d := n * 1500
	row := Row{}
	ok := make(chan bool) // out channel
	f := func(i int, r []Row) {
		for _, v := range r {
			if uint64(d) == p-1 { // found
				row = v
				ok <- true // out
			}
			d += 1
		}
	}
	np := t.Node[n]
	go np.Iter(func(i int, l Leaf) { l.Iter(f) }) // goroutine
	select {
	case <-ok:
		goto end // out
	}
end:
	return SingleResult{Row: row, Offset: p} // found offset
}

// SelectRange : GE ? <F> -> <T>
func (m *Memory) SelectRange(at uint8, fp, tp uint64) Result {
	res := MultipleResult{}
	if fp > m.Tb[at].Rows {
		return nil
	}
	d := 0
	f := func(i int, r []Row) {
		for _, v := range r {
			d += 1
			if uint64(d) >= fp && uint64(d) <= tp {
				res.Data = append(res.Data, v)
				res.Offset = append(res.Offset, uint64(d))
				res.Rows += 1
			}
		}
	}
	t := m.Tr[at] // T
	t.Iter(func(i int, n Node) { // N
		if uint64(d) < tp {
			n.Iter(func(i int, l Leaf) {
				if uint64(d) < tp {
					l.Iter(f)
				}
			}) // L
		}
	})
	return res
}

// Update : UP ? <P> <S> <N> [<V>]
func (m *Memory) Update(at, sp uint8, n, v string, p uint64) Result {
	r := m.SelectOne(at, p)
	if r == nil {
		return nil
	}
	row := r.(SingleResult)
	if v != "" && sp < row.Row.Len() && row.Row.Data[sp-1] != v { // Not verify
		return nil
	}
	if sp > row.Row.Len() { // S
		row.Row.Data = append(row.Row.Data, n)
	} else {
		row.Row.Data[sp-1] = n // Set
	}
	w := DecimalPlaces(float64(p), 1500)
	if w != 0 {
		w -= 1
	}
	d := w * 1500
	f := func(i int, r *[]Row) {
		for j := 0; j < len(*r); j++ {
			if uint64(d) == p-1 {
				(*r)[j] = row.Row // Set
			}
			d += 1
		}
	}
	t := m.Tr[at] // T
	t.Iter(func(i int, n Node) { // N
		n.Iter(func(i int, l Leaf) { l.PointIter(f) })
	}) // L
	return ModifyResult{Rows: 1}
}

// Delete : DE ? <P>
func (m *Memory) Delete(at uint8, p uint64) Result {
	m.Tb[at].Rows -= 1
	w := DecimalPlaces(float64(p), 1500)
	if w != 0 {
		w -= 1
	}
	d := w * 1500
	f := func(l *Leaf, i int, r *[]Row) {
		var dst []Row
		var get = false
		for j := 0; j < len(*r); j++ {
			if uint64(d) == p-1 {
				get = true
			} else {
				dst = append(dst, (*r)[j])
			}
			d += 1
		}
		if get {
			l.Data[i] = dst
		}
	}
	t := m.Tr[at] // T
	t.Iter(func(i int, n Node) { // N
		n.Iter(func(i int, l Leaf) { l.OutPointIter(f) })
	}) // L
	return ModifyResult{Rows: 1}
}

// DeleteAll : DELETE ALL
func (m *Memory) DeleteAll(at uint8) Result {
	var tb []Table
	var tr []Tree
	var rows = m.Tb[at].Rows
	for i := 0; uint(i) < m.Tbs; i++ { // For Mem
		if uint8(i) != at {
			tb = append(tb, m.Tb[i])
			tr = append(tr, m.Tr[i])
		}
	}
	m.Tb = tb
	m.Tr = tr
	m.Tbs -= 1
	tb = nil                                // To nil
	for i := 0; i < len(GlobalEm.Tb); i++ { // For Em
		if uint8(i) != at {
			tb = append(tb, GlobalEm.Tb[i])
		}
	}
	GlobalEm.Tb = tb
	for j := 0; uint(j) < m.Tbs; j++ { // Reset At
		m.Tb[j].At = uint8(j + 1)
		GlobalEm.Tb[j].At = uint8(j + 1)
	}
	return ModifyResult{Rows: rows}
}

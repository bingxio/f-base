// F-Base
// GPL 3.0 License - bingxio <3106740988@qq.com>
package main

import (
	"errors"
	"strconv"
)

// TableSize : Maximum byte length of the structure
const TableSize = 150

// Table : Each data table in the DB
type Table struct {
	Name    [20]byte
	Created uint32
	From    uint64
	Rows    uint64
	At      uint8
}

// Insert : Perform the add operation
func (tb *Table) Insert(fields []string) Result {
	tb.Rows += 1 // For Em table
	return GlobalMem.Insert(tb.At-1, fields)
}

// Selector : Conditional query
func (tb *Table) Selector(s, v string) (Result, error) {
	p, err := strconv.Atoi(s)
	if err != nil {
		return nil, errors.New("need receive integer offset")
	}
	return GlobalMem.Selector(tb.At-1, uint8(p), v), nil
}

// Select : Perform query operation
func (tb *Table) Select(f, t string) (Result, error) {
	rf := 0
	rt := 0
	if f != "" { // F
		i, err := strconv.Atoi(f)
		if err != nil {
			return nil, errors.New("need receive integer offset")
		}
		rf = i
	}
	if t != "" { // T
		i, err := strconv.Atoi(t)
		if err != nil {
			return nil, errors.New("need receive integer offset")
		}
		rt = i
	}
	if rf == 0 && rt == 0 {
		return GlobalMem.SelectAll(tb.At - 1), nil // All
	} else {
		if rf != 0 && rt == 0 {
			return GlobalMem.SelectOne(tb.At-1, uint64(rf)), nil // One
		} else {
			if rt != 0 && rf > rt {
				return nil, errors.New("index limit exceeded")
			}
			return GlobalMem.SelectRange(tb.At-1, uint64(rf), uint64(rt)), nil // Range
		}
	}
}

// Update : Perform update operation
func (tb *Table) Update(p, s, n, v string) (Result, error) {
	rp := -1
	rs := -1
	if p != "" { // P
		i, err := strconv.Atoi(p)
		if err != nil {
			return nil, errors.New("need receive integer offset")
		}
		rp = i
	}
	if s != "" { // S
		i, err := strconv.Atoi(s)
		if err != nil {
			return nil, errors.New("need receive integer offset")
		}
		rs = i
	}
	if rp <= 0 || rs <= 0 {
		return nil, errors.New("index range")
	}
	return GlobalMem.Update(tb.At-1, uint8(rs), n, v, uint64(rp)), nil
}

// Delete : Perform the delete operation
func (tb *Table) Delete(p string) (Result, error) {
	if p == "" {
		return GlobalMem.DeleteAll(tb.At - 1), nil
	}
	rp := -1
	if p != "" {
		i, err := strconv.Atoi(p)
		if err != nil {
			return nil, errors.New("need receive integer offset")
		}
		rp = i
	}
	return GlobalMem.Delete(tb.At-1, uint64(rp)), nil
}

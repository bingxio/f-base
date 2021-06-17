// F-Base
// GPL 3.0 License - bingxio <3106740988@qq.com>
package main

import (
	"errors"
	"log"
	"strconv"
)

// Maximum byte length of the structure
const TableSize = 200

// Table : Each data table in the DB
type Table struct {
	Name    [20]byte
	At      uint8
	Created uint32
	Rows    uint64 // Tbs + Table *? = Offset
}

// Insert : Perform the add operation
func (tb *Table) Insert(fields []string) error {
	tb.Rows++
	return nil
}

// Select : Perform query operation
func (tb *Table) Select(f, t string) ([]Row, error) {
	rf := -1
	rt := -1
	// F
	if f != "" {
		i, err := strconv.Atoi(f)
		if err != nil {
			return nil, errors.New(IntOffset)
		}
		rf = i
	}
	// T
	if t != "" {
		i, err := strconv.Atoi(t)
		if err != nil {
			return nil, errors.New(IntOffset)
		}
		rt = i
	}
	if rt != -1 && rt < rf {
		return nil, errors.New(IndexRange)
	}
	log.Println(rf, rt)
	return nil, nil
}

// Update : Perform update operation
func (tb *Table) Update(p, n, v string) (uint64, error) {
	rp := -1
	// P
	if p != "" {
		i, err := strconv.Atoi(p)
		if err != nil {
			return 0, errors.New(IntOffset)
		}
		rp = i
	}
	log.Println(rp)
	return 0, nil
}

// Delete : Perform the delete operation
func (tb *Table) Delete(p, v string) error {
	// TODO: When p is -1, delete the entire data table
	if p == "-1" {
		log.Println("delete all")
		return nil
	}
	rp := -2
	if p != "" {
		i, err := strconv.Atoi(p)
		if err != nil {
			return errors.New(IntOffset)
		}
		rp = i
	}
	log.Println(rp)
	return nil
}

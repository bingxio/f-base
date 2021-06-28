// F-Base
// GPL 3.0 License - bingxio <3106740988@qq.com>
package main

import (
	"errors"
	"fmt"
	"log"
	"strconv"
)

// Maximum byte length of the structure
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
func (tb *Table) Insert(fields []string) error {
	return nil
}

// Select : Perform query operation
func (tb *Table) Select(f, t string) ([]Row, error) {
	rf := 0
	rt := 0
	// F
	if f != "" {
		i, err := strconv.Atoi(f)
		if err != nil {
			return nil, errors.New("need receive integer offset")
		}
		rf = i
	}
	// T
	if t != "" {
		i, err := strconv.Atoi(t)
		if err != nil {
			return nil, errors.New("need receive integer offset")
		}
		rt = i
	}
	// All
	if rf == 0 && rt == 0 {
		fmt.Println("select all")
	}
	// One
	if rf != 0 && rt == 0 {
		fmt.Println("select one")
	}
	if rt != 0 && rf > rt {
		return nil, errors.New("index range")
	}
	// Range
	if rf != 0 && rt != 0 {
		fmt.Println("select range")
	}
	return nil, nil
}

// Update : Perform update operation
func (tb *Table) Update(p, n, v string) (uint64, error) {
	rp := -1
	// P
	if p != "" {
		i, err := strconv.Atoi(p)
		if err != nil {
			return 0, errors.New("need receive integer offset")
		}
		rp = i
	}
	if rp <= 0 {
		return 0, errors.New("index range")
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
			return errors.New("need receive integer offset")
		}
		rp = i
	}
	log.Println(rp)
	return nil
}

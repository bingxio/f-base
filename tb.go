// F-Base
// GPL 3.0 License - bingxio <3106740988@qq.com>
package main

import (
	"fmt"
)

// Maximum byte length of the structure
const TableSize = 200

// Table : Each data table in the DB
type Table struct {
	Name    [20]byte
	At      uint8
	Len     uint8
	Created uint32
	Rows    uint64 // Tbs + Table *? = Offset
}

// Insert : Perform the add operation
func (tb *Table) Insert(fields []string) error {
	return nil
}

// Select : Perform query operation
func (tb *Table) Select(f, t int64) ([]Row, error) {
	return nil, nil
}

// Update : Perform update operation
func (tb *Table) Update(p, n, v string) (uint64, error) {
	return 0, nil
}

// Delete : Perform the delete operation
func (tb *Table) Delete(p, v string) error {
	// TODO: When p is -1, delete the entire data table
	if p == "-1" {
		fmt.Println("DELETE ALL")
	}
	return nil
}

// Count : Return the count of data in the table
func (tb Table) Count() uint64 {
	return 0
}

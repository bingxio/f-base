// F-Base
// GPL 3.0 License - bingxio <3106740988@qq.com>
package main

// Table : Each data table in the DB
type Table struct {
	name string
	at   int
}

// Insert : Perform the add operation
func (tb *Table) Insert(fields []string) error { return nil }

// Select : Perform query operation
func (tb *Table) Select(f, t int64) ([]Row, error) { return nil, nil }

// Update : Perform update operation
func (tb *Table) Update(p, n, v string) (uint64, error) { return 0, nil }

// Delete : Perform the delete operation
func (tb *Table) Delete(p, v string) error {
	// TODO: When p is -1, delete the entire data table
	return nil
}

// Count : Return the count of data in the table
func (tb Table) Count() uint64 { return 0 }

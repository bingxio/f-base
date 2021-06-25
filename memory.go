// F-Base
// GPL 3.0 License - bingxio <3106740988@qq.com>
package main

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
)

// UpdateTbs : Update the number of tables in memory
func UpdateTbs(n uint8) error {
	b := make([]byte, tbs)
	err := gob.NewEncoder(bytes.NewBuffer(b)).Encode(n)
	if err != nil {
		return err
	}
	c, err := GlobalEm.file.WriteAt(b, 0)
	if err != nil {
		return err
	}
	// Empty
	if c == 0 {
		return errors.New(WriteFile)
	}
	return nil
}

// UpdateTable : Update table in memory
func UpdateTable(at int64, n Table) error {
	offset := tbs + at*TableSize
	fmt.Println(offset)
	b := bytes.NewBuffer([]byte{})
	err := gob.NewEncoder(b).Encode(n)
	if err != nil {
		return err
	}
	b.Write(make([]byte, TableSize-b.Len()))
	c, err := GlobalEm.file.WriteAt(b.Bytes(), offset)
	if err != nil {
		return err
	}
	// Empty
	if c == 0 {
		return errors.New(WriteFile)
	}
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

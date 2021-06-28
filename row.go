// F-Base
// GPL 3.0 License - bingxio <3106740988@qq.com>
package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

// Maximum byte length of the structure
const RowSize = 100

// Row : Row in table
type Row struct {
	Data []string
}

// NewRow : New row byte stream
func NewRow(data []string) ([]byte, error) {
	b := bytes.NewBuffer([]byte{})
	err := gob.NewEncoder(b).Encode(Row{Data: data})
	if err != nil {
		return nil, err
	}
	// 100
	l := RowSize - b.Len()
	b.Write(make([]byte, l))
	return b.Bytes(), nil
}

// Len : Return length of data in row
func (r Row) Len() uint8 {
	return uint8(len(r.Data))
}

// Stringer : stringer
func (r Row) Stringer() string {
	return fmt.Sprintf("%v", r.Data) // %q
}

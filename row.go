// F-Base
// GPL 3.0 License - bingxio <3106740988@qq.com>
package main

// RowSize : Maximum byte length of the structure
const RowSize = 100

// Row : Row in table
type Row struct {
	Data []string
}

// Len : Return length of data in row
func (r Row) Len() uint8 {
	return uint8(len(r.Data))
}

// Stringer : stringer
func (r Row) Stringer() string {
	l := "("
	for k, v := range r.Data {
		l += v
		if k+1 != len(r.Data) {
			l += " "
		}
	}
	l += ")"
	return l
}

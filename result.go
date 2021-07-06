// F-Base
// GPL 3.0 License - bingxio <3106740988@qq.com>
package main

import "fmt"

// Result : Mode of executed expression
type Result interface {
	Stringer() string
}

// MultipleResult : Many data query
type MultipleResult struct {
	Rows   uint64   `json:"rows"`
	Data   []Row    `json:"data"`
	Offset []uint64 `json:"offset"`
}

// SingleResult : Single data
type SingleResult struct {
	Row    Row    `json:"row"`
	Offset uint64 `json:"offset"`
}

// ModifyResult : The rows of limit modified
type ModifyResult struct {
	Rows uint64 `json:"rows"`
}

func (m MultipleResult) Stringer() string {
	return "" // None to display limits
}

func (s SingleResult) Stringer() string {
	return fmt.Sprintf("%d %s", s.Offset, s.Row.Stringer())
}

func (m ModifyResult) Stringer() string {
	return fmt.Sprintf("%d rows modified", m.Rows)
}

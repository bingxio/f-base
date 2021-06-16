// F-Base
// GPL 3.0 License - bingxio <3106740988@qq.com>
package main

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"strings"
	"time"
)

const port = 3742 // PORT of fork mode

// Em : Global database manager
type Em struct {
	db   string
	file *os.File
	tb   []Table
	fork bool
}

// Stringer return the string literal of Em
func (e Em) Stringer(info string) string {
	b := strings.Builder{}
	b.WriteString(fmt.Sprintf("Db: '%s' %sB\n", e.db, info))
	b.WriteString(fmt.Sprintf("Tb: %v", e.TableString()))
	return b.String()
}

// Table : 'table' command to print all of tables in the DB
func (e *Em) Table() {
	for _, v := range e.tb {
		fmt.Printf(
			"<at: %-2d name: '%s' len: %d row: %-5d created: %s>\n",
			v.At,
			v.Name,
			v.Len,
			v.Rows,
			time.Unix(
				int64(v.Created), 0).Format("2006.01.02 15:04"),
		)
	}
}

// TableString : Return a slice of names
func (e Em) TableString() string {
	s := "["
	l := len(e.tb)
	for i := 0; i < l; i++ {
		s += fmt.Sprintf("%d: '%s'",
			i+1, string(e.tb[i].Name[:]))
		if i+1 != l {
			s += " "
		}
	}
	s += "]"
	return s
}

// CountTable : Return the count of tables in the DB
func (e Em) CountTable() int {
	return len(e.tb)
}

// Fork : Fork the program to service into OS
func (e *Em) Fork() {
	fmt.Println(
		"F-Base runs in the background and monitors port:", port)
	os.Exit(0)
}

// LoadDb : Load database
func (e *Em) LoadDb() error {
	info, err := e.file.Stat()
	if err != nil {
		return err
	}
	if info.Size() == 0 {
		// return nil
		e.testData()
	}
	// Read Tbs
	count, err := e.loadTbs()
	if err != nil {
		return err
	}
	if count == 0 {
		return nil
	}
	// Have Tbs
	err = e.loadTables(count)
	if err != nil {
		return err
	}
	return nil
}

// LoadTbs : Returns the counts of table in DB
func (e Em) loadTbs() (uint8, error) {
	e.file.Seek(0, io.SeekStart)
	// Read
	b := make([]byte, 4)
	_, err := e.file.Read(b)
	if err != nil {
		return 0, err
	}
	// Decode
	x := uint8(0)
	err = gob.NewDecoder(
		bytes.NewBuffer(b)).Decode(&x)
	if err != nil {
		return 0, err
	}
	return x, nil
}

// LoadTables : Parse all data tables
func (e *Em) loadTables(tbs uint8) error {
	for tbs > 0 {
		// Read
		b := make([]byte, 200)
		_, err := e.file.Read(b)
		if err != nil {
			return err
		}
		// Decode
		t := Table{}
		err = gob.NewDecoder(
			bytes.NewBuffer(b)).Decode(&t)
		if err != nil {
			return err
		}
		e.tb = append(e.tb, t) // Append
		tbs--
	}
	return nil
}

// NewTable : Store new table
func (e *Em) NewTable(t Table) ([]byte, error) {
	t.At = uint8(len(e.tb) + 1)
	// Write
	b := bytes.NewBuffer([]byte{})
	err := gob.NewEncoder(b).Encode(t)
	if err != nil {
		return nil, err
	}
	if b.Len() >= TableSize {
		return nil, errors.New("size out of specification")
	}
	// b.Len() < 200
	l := TableSize - b.Len()
	_, err = e.file.Write(b.Bytes())
	e.file.Write(make([]byte, l))
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

// ExecuteExpr : Execute command expression
func (e *Em) ExecuteExpr(expr Expr) error {
	if reflect.TypeOf(expr).Name() == "SeExpr" {
		log.Println("will exec se expr")
	}
	fmt.Println(expr.Stringer())
	e.file.Seek(204, io.SeekStart)
	return nil
}

// Export some testing data written to DB file
func (e *Em) testData() {
	// Tbs
	b := bytes.NewBuffer([]byte{})
	err := gob.NewEncoder(b).Encode(uint8(3))
	if err != nil {
		panic(err.Error())
	}
	e.file.Write(b.Bytes())
	fmt.Println(b.Bytes(), b.Len())
	// Table
	b.Reset()
	t := []Table{
		{
			At:      1,
			Name:    [20]byte{117, 115, 101, 114, 115},
			Len:     6,
			Created: uint32(time.Now().Unix()),
			Rows:    89,
		},
		{
			At:      2,
			Name:    [20]byte{116, 111, 100, 111, 115},
			Len:     8,
			Created: uint32(time.Now().Unix()),
			Rows:    144,
		},
		{
			At:      3,
			Name:    [20]byte{112, 97, 121, 101, 100},
			Len:     2,
			Created: uint32(time.Now().Unix()),
			Rows:    99,
		},
	}
	for _, v := range t {
		err = gob.NewEncoder(b).Encode(v)
		if err != nil {
			panic(err.Error())
		}
		e.file.Write(b.Bytes())
		e.file.Write(make([]byte, TableSize-b.Len()))
		log.Println(b.Bytes(), b.Len(), b.Len()+200-b.Len())
		b.Reset()
	}
}

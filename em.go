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
	"strconv"
	"strings"
	"time"
)

const (
	port = 3742 // PORT of fork mode
	tbs  = 4    // Length of bytes for table count
)

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
	b.WriteString("Tb: \n")
	b.WriteString(func() string {
		l := ""
		for i := 0; i < len(e.tb); i++ {
			l += fmt.Sprintf("\t%d: '%s'", e.tb[i].At, e.tb[i].Name[:])
			if i+1 != len(e.tb) {
				l += "\n"
			}
		}
		return l
	}())
	return b.String()
}

// Table : 'table' command to print all of tables in the DB
func (e *Em) Table() {
	for _, v := range e.tb {
		fmt.Printf(
			"<name: '%s' row: %-4d create: %s>\n",
			func(b [20]byte) string {
				s := ""
				for i := 0; i < len(b); i++ {
					if b[i] != 0 {
						s += string(b[i])
					}
				}
				return s
			}(v.Name),
			v.Rows,
			time.Unix(
				int64(v.Created), 0).Format("2006/01/02 15:04"),
		)
	}
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
	NewMemory(*e) // To Memory
	return nil
}

// LoadTbs : Returns the counts of table in DB
func (e Em) loadTbs() (uint8, error) {
	e.file.Seek(0, io.SeekStart)
	// Read
	b := make([]byte, tbs)
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
		b := make([]byte, TableSize)
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

// Exist : Data table exist
func (e Em) Exist(name string) (int, bool) {
	for i, v := range e.tb {
		if func(a, b []byte) bool {
			for i := 0; i < len(b); {
				// [124, 115, 154, 118, 175, 0, 0, 0, 0, 0]
				// [124, 115, 154, 118]
				if a[i] != b[i] {
					return false
				}
				i++
				if i == len(b) && i != len(a) {
					return a[i] == 0
				}
			}
			return true
		}(v.Name[:], []byte(name)) {
			return i, true
		}
	}
	return -1, false
}

// ExecuteExpr : Execute command expression
func (e *Em) ExecuteExpr(expr Expr) error {
	switch reflect.TypeOf(expr).Name() {
	case "SeExpr":
		se := expr.(SeExpr)
		i, exist := e.Exist(se.Table.Literal)
		if !exist {
			log.Println("create new tb", se.Table.Literal)
			return nil
		}
		err := e.tb[i].Insert(func(f []Token) []string {
			e := make([]string, len(f))
			for _, v := range f {
				e = append(e, v.Literal)
			}
			return e
		}(se.F))
		if err != nil {
			return err
		}
	case "UpExpr":
		up := expr.(UpExpr)
		i, exist := e.Exist(up.Table.Literal)
		if !exist {
			return errors.New("table does not exist")
		}
		_, err := e.tb[i].Update(
			up.P.Literal, up.N.Literal, up.V.Literal)
		if err != nil {
			return err
		}
	case "DeExpr":
		de := expr.(DeExpr)
		i, exist := e.Exist(de.Table.Literal)
		if !exist {
			return errors.New("table does not exist")
		}
		err := e.tb[i].Delete(de.P.Literal, de.V.Literal)
		if err != nil {
			return err
		}
	case "GeExpr":
		ge := expr.(GeExpr)
		i, exist := e.Exist(ge.Table.Literal)
		if !exist {
			return errors.New("table does not exist")
		}
		rows, err := e.tb[i].Select(ge.F.Literal, ge.T.Literal)
		if err != nil {
			return err
		}
		for _, v := range rows {
			fmt.Println(v.Stringer())
		}
	}
	return nil
}

// Export some testing data written to DB file
func (e *Em) testData() {
	// Tbs
	b := bytes.NewBuffer([]byte{})
	err := gob.NewEncoder(b).Encode(uint8(2))
	if err != nil {
		panic(err.Error())
	}
	e.file.Write(b.Bytes())
	fmt.Printf("Tbs: %p (%d)\n", b.Bytes(), b.Len())
	// Table
	//
	b.Reset()
	t := []Table{
		{
			Name:    [20]byte{117, 115, 101, 114, 115},
			Created: uint32(time.Now().Unix()),
			From:    304, // 4 + 150*2
			Rows:    522,
			At:      1,
		},
		{
			Name:    [20]byte{116, 111, 100, 111, 115},
			Created: uint32(time.Now().Unix()),
			From:    52504, // 4 + 150*2 + 100*522
			Rows:    143,
			At:      2,
		},
	}
	for _, v := range t {
		err = gob.NewEncoder(b).Encode(v)
		if err != nil {
			panic(err.Error())
		}
		e.file.Write(b.Bytes())
		e.file.Write(make([]byte, TableSize-b.Len()))
		fmt.Printf("Table: %p (%d) >> %d\n",
			b.Bytes(), b.Len(), b.Len()+(TableSize-b.Len()))
		b.Reset()
	}
	r := []Row{
		// {
		// 	Data: []string{"1", "User1", "123456"},
		// },
		// {
		// 	Data: []string{"2", "User2", "123456"},
		// },
		// {
		// 	Data: []string{"3", "User3", "123456"},
		// },
		// {
		// 	Data: []string{"4", "User4", "123456"},
		// },
		// {
		// 	Data: []string{"5", "User5", "123456"},
		// },
		// {
		// 	Data: []string{"1", "Order1", "789101"},
		// },
		// {
		// 	Data: []string{"2", "Order2", "789101"},
		// },
		// {
		// 	Data: []string{"3", "Order3", "789101"},
		// },
		// {
		// 	Data: []string{"4", "Order4", "789101"},
		// },
		// {
		// 	Data: []string{"5", "Order5", "789101"},
		// },
		// {
		// 	Data: []string{"6", "Order6", "789101"},
		// },
		// {
		// 	Data: []string{"7", "Order7", "789101"},
		// },
		// {
		// 	Data: []string{"8", "Order8", "789101"},
		// },
	}
	for i := 0; i < 665; i++ {
		r = append(r, Row{
			Data: []string{strconv.Itoa(i), "test", "what's up"},
		})
	}
	for k, v := range r {
		err = gob.NewEncoder(b).Encode(v)
		if err != nil {
			panic(err.Error())
		}
		e.file.Write(b.Bytes())
		e.file.Write(make([]byte, RowSize-b.Len()))
		fmt.Printf("Row(%d): %p (%d) >> %d\n",
			k+1, b.Bytes(), b.Len(), b.Len()+(RowSize-b.Len()))
		b.Reset()
	}
	e.file.Seek(0, io.SeekStart)
}

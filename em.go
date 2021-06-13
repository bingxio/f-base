// F-Base
// GPL 3.0 License - bingxio <3106740988@qq.com>
package main

import (
  "fmt"
  "os"
  "strings"
  "time"
)

// Em : Global database manager
type Em struct {
  db     string
  file   *os.File
  tb     []Table
  opened time.Time
  fork   bool
}

// Stringer return the string literal of Em
func (e Em) Stringer(info string) string {
  b := strings.Builder{}
  b.WriteString(fmt.Sprintf("Db: %s %s\n", e.db, info))
  b.WriteString(fmt.Sprintf("Tb: %v", e.tb))
  // b.WriteString(fmt.Sprintf("Open: %s",
  //	e.opened.Format("2006.01.02 15:04:05")))
  return b.String()
}

// Table : 'table' command to print all of tables in the DB
func (e *Em) Table() {
  e.tb = append(e.tb, Table{name: "users", at: 1}) // TEST
  e.tb = append(e.tb, Table{name: "todos", at: 2}) // TEST
  for _, v := range e.tb {
    fmt.Printf("(%d '%s')\n", v.at, v.name)
  }
  e.tb = e.tb[0:0]
}

// Fork : Fork the program to service into OS
func (e *Em) Fork() {
  fmt.Println("F-Base runs in the background and monitors port 3742")
  os.Exit(0)
}

// Make a error message
func (e Em) error(err string) {
  fmt.Printf("EM-ERROR: %s", err)
}

// LoadTB : Load database tables in file
func (e *Em) LoadTB() error {
  return nil
}

// ExecuteExpr : Execute command expression
func (e *Em) ExecuteExpr(expr Expr) error {
  fmt.Println(expr.Stringer())
  return nil
}

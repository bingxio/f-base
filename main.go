// F-Base
// GPL 3.0 License - bingxio <3106740988@qq.com>
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"time"
)

// Message : Prompt information of command line
const Message = "F-Base v0.0.1 * MADE AT May 26 2021 18:06:35"

var (
	DbPath   = ""                                      // Database
	DbFlag   = os.O_APPEND | os.O_CREATE | os.O_WRONLY // Open Flags
	ForkMode = false                                   // Fork Mode
	GlobalEm = Em{}                                    // Em
)

/*
  Command   ->   Repl | Fork -> OS <-> Interface-driver
                   |            |
                   v           / - -                   // Front-end
  Stringer  <-    Ast         /    |
                   |         /     |
                   v <-  -  -      v
                  Em(Table, Exec, Fork)                // Back-end
                   |
                   v
   Result   <-    Tb(SE, GE, UP, DE) <- Row <-> {F-Tree}
     |
     v
   Terminal | Driver --> END
*/
func main() {
	flag.Parse()
	arg := flag.Args()

	// Must access the location
	// of the database file
	if len(arg) < 1 {
		_, _ = fmt.Fprintln(os.Stderr, "LOST DATABASE SPECIFIED")
		return
	}
	ForkMode = len(arg) == 2 && arg[1] == "fork"
	DbPath = arg[0]
	f, err := os.OpenFile(DbPath, DbFlag, 0644) // Database File
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
		return
	}
	// Em set
	GlobalEm = Em{
		db:     f.Name(),
		file:   f,
		tb:     nil,
		opened: time.Now(),
		fork:   ForkMode,
	}
	// Load database tables
	err = GlobalEm.LoadTB()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// Fork mode
	if ForkMode {
		GlobalEm.Fork()
	} else {
		repl() // Go to read prompt-line loop
	}
}

// Read Prompt-Line Loop
func repl() {
	fmt.Println(Message)
	// fmt.Printf("DB: '%s' %s\n", DbPath, dbInfo())
	fmt.Println(GlobalEm.Stringer(dbInfo()))

	buf := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")

		l, _, err := buf.ReadLine()
		if err != nil {
			panic(err)
		}
		line := string(l)
		if len(line) == 0 { // Empty line
			continue
		}
		// All black characters line
		j := 0
		for i := 0; i < len(line); i++ {
			if line[i] == ' ' || line[i] == '\t' {
				j++
			}
		}
		if j == len(line) {
			continue
		}
		// 'help', 'exit', 'quit' and other special commands
		if line == "help" {
			usage()
		} else if line == "table" {
			GlobalEm.Table()
		} else if line == "license" {
			license()
		} else if line == "exit" || line == "quit" {
			fmt.Println("bye")
			break
		} else {
			eval(line)
		}
	}
}

// execute each command
func eval(src string) {
	expr := NewLex(src).parse()
	if expr.Kind() == Er { // Error expression
		return // Out
	}
	err := GlobalEm.ExecuteExpr(expr)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

// 'help' command to print usage information
func usage() {
	fmt.Println(
		`
	SE(se) ? *<E>		-> Insert
	GE(ge) ? | <F> | -> <T>	-> Select
	UP(up) ? <P> <N> [<V>]	-> Update
	DE(de) ? <P> [<V>]	-> Delete

	E -> Elements
	F -> F
	T -> T
	P -> Position
	N -> N
	V -> Verify

	table       -> List of tables in the DB
	exit | quit -> Exit the program
	license     -> Show license
	`)
}

// Return the size of the data file
func dbInfo() string {
	stat, err := GlobalEm.file.Stat()
	if err != nil {
		return err.Error()
	}
	return fmt.Sprint(stat.Size()) // Return size of file
}

// Show license
func license() {
	fmt.Println("GPL 3.0 - bingxio(黄菁) <3106740988@qq.com>")
}

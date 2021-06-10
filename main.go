// F-Base
// GPL 3.0 License - bingxio <3106740988@qq.com>
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

// Message : Prompt information of command line
const Message = "F-Base v0.0.1 * MADE AT May 26 2021 18:06:35"

var (
	DbPath = ""                                      // Database
	DbFile *os.File                                  // Db File
	DbFlag = os.O_APPEND | os.O_CREATE | os.O_WRONLY // Open Flags
)

func main() {
	flag.Parse()
	arg := flag.Args()

	// Must access the location
	// of the database file
	if len(arg) < 1 {
		_, _ = fmt.Fprintln(os.Stderr, "LOST DATABASE SPECIFIED")
		return
	}
	DbPath = arg[0]
	f, err := os.OpenFile(DbPath, DbFlag, 0644) // Database File
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
		return
	}
	DbFile = f
	repl() // Go to read prompt-line loop
}

// Read Prompt-Line Loop
func repl() {
	fmt.Println(Message)
	fmt.Printf("DB: '%s' %s\n", DbPath, dbInfo())

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
		j := 0
		for i := 0; i < len(line); i++ {
			if line[i] == ' ' || line[i] == '\t' { // All black characters line
				j++
			}
		}
		if j == len(line) {
			continue
		}
		// 'help', 'exit', 'quit' and other special commands
		if line == "help" {
			usage()
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
	fmt.Println(expr.Stringer())
}

// 'help' command to print usage information
func usage() {
	fmt.Println(
		`Usage: 
	SE ? *<E>
	GE ? | <F> | -> <T>
	UP ? <P> <N> [<V>]
	DE ? <P> [<V>]`)
}

// Return the size of the data file
func dbInfo() string {
	stat, err := DbFile.Stat()
	if err != nil {
		return err.Error()
	}
	return fmt.Sprint(stat.Size())
}

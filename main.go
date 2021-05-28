// F-Base
// GPL 3.0 License - bingxio <3106740988@qq.com>
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

const Message = "F-Base v0.0.1 * MADE AT May 26 2021 18:06:35"

var (
	DbPath = ""     // Database
	DbFile *os.File // Db File
	DbFlag = os.O_APPEND | os.O_CREATE | os.O_WRONLY
)

func main() {
	flag.Parse()
	arg := flag.Args()

	if len(arg) < 1 {
		fmt.Fprintln(os.Stderr, "LOST DATABASE SPECIFIED")
		return
	}
	DbPath = arg[0]
	f, err := os.OpenFile(DbPath, DbFlag, 0644)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}
	DbFile = f
	repl()
}

func repl() {
	fmt.Println(Message)
	fmt.Printf("DB: '%s' %s\n", DbPath, DbInfo())

	buf := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")

		l, _, err := buf.ReadLine()
		if err != nil {
			panic(err)
		}
		line := string(l)
		if len(line) == 0 {
			continue
		}
		j := 0
		for i := 0; i < len(line); i++ {
			if line[i] == ' ' || line[i] == '\t' {
				j++
			}
		}
		if j == len(line) {
			continue
		}
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

func eval(src string) {
	ast := NewLex(src).parse()
	if ast.Kind() == Er {
		return
	}
	fmt.Println(ast.Stringer())
}

func usage() {
	fmt.Println(
		`Usage: 
	SE ? *<E>
	GE ? | <F> | -> <T>
	UP ? <P> <N> [<V>]
	DE ? <P> [<V>]`)
}

func DbInfo() string {
	stat, err := DbFile.Stat()
	if err != nil {
		return err.Error()
	}
	return fmt.Sprint(stat.Size())
}

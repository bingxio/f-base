// F-Base
// GPL 3.0 License - bingxio <3106740988@qq.com>
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

// Message : Prompt information of command line
const Message = "F-Base v0.0.1 * MADE AT May 26 2021 18:06:35"

var (
	DbPath   = ""                      // Database
	DbFlag   = os.O_CREATE | os.O_RDWR // Open Flags
	ForkMode = false                   // Fork Mode
	GlobalEm = Em{}                    // Em
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
		_, _ = fmt.Fprintln(os.Stderr, "lost database specified")
		return
	}
	ForkMode = len(arg) == 2 && arg[1] == "fork"
	DbPath = arg[0]
	f, err := os.OpenFile(DbPath, DbFlag, 0666) // Database File
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
		return
	}
	GlobalEm = Em{ // Em set
		db:   f.Name(),
		file: f,
		tb:   nil,
		fork: ForkMode,
	}
	err = GlobalEm.LoadDb() // Load database
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
		return
	}
	if ForkMode { // Fork mode
		GlobalEm.Fork()
	} else {
		repl() // Go to read prompt-line loop
	}
}

// Read Prompt-Line Loop
func repl() {
	fmt.Println(Message)
	fmt.Println(GlobalEm.Stringer(dbInfo()))

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM) // Control+C or KILL
	go func() {
		for range c {
			fmt.Printf("\nuse 'exit' or 'quit' command\n> ")
		}
	}()

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
		j := 0 // All black character line
		for i := 0; i < len(line); i++ {
			if line[i] == ' ' || line[i] == '\t' {
				j++
			}
		}
		if j == len(line) {
			continue
		}
		// 'help', 'exit',
		// 'quit' and other special commands
		if line == "help" {
			help()
		} else if line == "t" {
			GlobalEm.Table()
		} else if line == "license" {
			license()
		} else if line == "exit" || line == "quit" {
			err := QuitMemory()

			if err != nil {
				_, _ = fmt.Fprintln(os.Stderr, err.Error())
				continue
			}
			fmt.Println("bye")
			break
		} else if line == "p" {
			GlobalMem.Dissemble()
		} else if line == "usage" {
			usage()
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
	fmt.Println("OK")
}

// 'help' command to print help information
func help() {
	fmt.Println(`	SE(se) ? *<E>		-> Insert
	GE(ge) ? | <F> | -> <T>	-> Select
	UP(up) ? <P> <N> [<V>]	-> Update
	DE(de) ? <P>            -> Delete
	GT(gt) ? <S> <V>        -> Selector

	E -> Elements
	F -> From
	T -> To
	P -> Position
	N -> New
	V -> Verify
	S -> Field

	t           -> List of tables in the DB
	exit | quit -> Exit the program
	license     -> Show license
	p           -> Dissemble trees
	usage       -> Show usage`)
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

// Show usage
func usage() {
	fmt.Println(`	SE u 1 bingxio 1234
	
	GE u
	GE u 1
	GE u 1 5
	
	UP u 1 3 new
	UP u 1 3 new old
	
	DE u 1
	
	GT u 3 limit`)
}

// F-Base
// GPL 3.0 License - bingxio <3106740988@qq.com>
package main

import (
    "bufio"
    "fmt"
    "os"
)

const Message = "F-Base v0.0.1 * MADE AT May 26 2021 18:06:35"

func main() {
    fmt.Println(Message)

    buf := bufio.NewReader(os.Stdin)
    for {
        fmt.Print("> ")

        l, _, err := buf.ReadLine()
        if err != nil {
            panic(err)
        }
        line := string(l)
        // Empty line
        if len(line) == 0 {
            continue
        }
        // All empty characters
        j := 0
        for i := 0; i < len(line); i++ {
            if line[i] == ' ' || line[i] == '\t' {
                j++
            }
        }
        if j == len(line) {
            continue
        }
        // Eval
        eval(line)
    }
}

func eval(src string) {
    // fmt.Println("Eval: ", src)
    ast := NewLex(src).parse()
    if ast.Kind() == Er {
        return
    }
    fmt.Println(ast.Stringer())
}

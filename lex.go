// F-Base
// GPL 3.0 License - bingxio <3106740988@qq.com>
package main

import (
    "fmt"
)

type Token struct {
    Literal string
    Kind    int
}

func (t Token) Stringer() string {
    return fmt.Sprintf("<\"%s\" %d>",
        t.Literal, t.Kind)
}

type Lex struct {
    Src string
    Pos int
}

func NewLex(src string) *Lex {
    return &Lex{
        Src: src,
        Pos: 0,
    }
}

func (l Lex) parse() Ast {
    kind := -1
    var expr interface{}

    for l.Pos < len(l.Src) {
        if kind == 4 {
            break
        }
        l.skipWhiteSpace()
        lit := ""
        for {
            if l.end() || l.space() {
                break
            }
            lit += string(l.now())
            l.Pos++
        }
        if lit != "" {
            if kind == -1 {
                switch lit {
                case "se", "SE":
                    kind = 0
                    expr = SeExpr{}
                case "ge", "GE":
                    kind = 1
                    expr = GeExpr{}
                case "up", "UP":
                    kind = 2
                    expr = UpExpr{}
                case "de", "DE":
                    kind = 3
                    expr = DeExpr{}
                default:
                    fmt.Printf("UNKNOWN COMMAND PREFIX: '%s'\n",
                        lit)
                    kind = 4
                    expr = ErExpr{}
                }
            } else {
                fmt.Println("PARSE: ", lit)
            }
            lit = ""
        }
        l.Pos++
    }
    return expr.(Ast)
}

func (l Lex) skipWhiteSpace() {
    for !l.end() && l.space() {
        l.Pos++
    }
}

func (l Lex) space() bool {
    return l.now() == ' ' || l.now() == '\t' || l.now() == '\r'
}

func (l Lex) now() byte {
    return l.Src[l.Pos]
}

func (l Lex) end() bool {
    return l.Pos >= len(l.Src)
}

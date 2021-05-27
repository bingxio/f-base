// F-Base
// GPL 3.0 License - bingxio <3106740988@qq.com>
package main

import (
    "errors"
    "fmt"
)

type Token struct {
    Literal string
}

func (t Token) Stringer() string {
    return fmt.Sprintf("<'%s'>", t.Literal)
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
                ast, err := setParam(&expr, Token{Literal: lit})
                if err != nil {
                    fmt.Println(err.Error())
                    kind = 4
                }
                expr = ast
            }
            lit = ""
        }
        l.Pos++
    }
    if err := verifyExpr(expr.(Ast)); err != nil {
        fmt.Println(err.Error())
        expr = ErExpr{}
    }
    return expr.(Ast)
}

func setParam(expr *interface{}, tok Token) (Ast, error) {
    switch (*expr).(type) {
    case SeExpr:
        seExpr := (*expr).(SeExpr)
        if emptyValue(seExpr.Table) {
            return SeExpr{
                Table: tok,
            }, nil
        } else {
            // FIELDS
            seExpr.Fields = append(seExpr.Fields, tok)
            return seExpr, nil
        }
    case GeExpr:
        geExpr := (*expr).(GeExpr)
        if emptyValue(geExpr.Table) {
            return GeExpr{
                Table: tok,
            }, nil
        } else {
            // FROM
            if emptyValue(geExpr.From) {
                geExpr.From = tok
                return geExpr, nil
            }
            // TO
            if emptyValue(geExpr.To) {
                geExpr.To = tok
                return geExpr, nil
            }
            return ErExpr{}, errors.New("MORE PARAM SPECIFIED")
        }
    case UpExpr:
        upExpr := (*expr).(UpExpr)
        if emptyValue(upExpr.Table) {
            return UpExpr{
                Table: tok,
            }, nil
        } else {
            // POS
            if emptyValue(upExpr.Pos) {
                upExpr.Pos = tok
                return upExpr, nil
            }
            // NEW
            if emptyValue(upExpr.New) {
                upExpr.New = tok
                return upExpr, nil
            }
            // VER
            if emptyValue(upExpr.Ver) {
                upExpr.Ver = tok
                return upExpr, nil
            }
            return ErExpr{}, errors.New("MORE PARAM SPECIFIED")
        }
    case DeExpr:
        deExpr := (*expr).(DeExpr)
        if emptyValue(deExpr.Table) {
            return DeExpr{
                Table: tok,
            }, nil
        } else {
            // POS
            if emptyValue(deExpr.Pos) {
                deExpr.Pos = tok
                return deExpr, nil
            }
            // VER
            if emptyValue(deExpr.Ver) {
                deExpr.Ver = tok
                return deExpr, nil
            }
            return ErExpr{}, errors.New("MORE PARAM SPECIFIED")
        }
    }
    return ErExpr{}, errors.New("PROGRAM ERROR")
}

func emptyValue(t Token) bool {
    return t.Literal == ""
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

func verifyExpr(expr Ast) error {
    switch expr.(type) {
    case SeExpr:
        seExpr := expr.(SeExpr)
        if emptyValue(seExpr.Table) {
            return errors.New("LOST TABLE")
        }
        if len(seExpr.Fields) == 0 {
            return errors.New("LOST FIELDS")
        }
        return nil
    case GeExpr:
        geExpr := expr.(GeExpr)
        if emptyValue(geExpr.Table) {
            return errors.New("LOST TABLE")
        }
        return nil
    case UpExpr:
        upExpr := expr.(UpExpr)
        if emptyValue(upExpr.Table) {
            return errors.New("LOST TABLE")
        }
        if emptyValue(upExpr.Pos) {
            return errors.New("LOST POSITION OF UPDATE LIMIT")
        }
        if emptyValue(upExpr.New) {
            return errors.New("LOST NEW VALUE")
        }
        return nil
    case DeExpr:
        deExpr := expr.(DeExpr)
        if emptyValue(deExpr.Table) {
            return errors.New("LOST TABLE")
        }
        if emptyValue(deExpr.Pos) {
            return errors.New("LOST POSITION OF DELETE LIMIT")
        }
        return nil
    }
    return nil
}

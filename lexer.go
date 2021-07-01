// F-Base
// GPL 3.0 License - bingxio <3106740988@qq.com>
package main

import (
	"errors"
	"fmt"
	"reflect"
)

// Token : Literal quantity of symbol
type Token struct {
	Literal string
}

// Stringer : Return the format output token
func (t Token) Stringer() string {
	return fmt.Sprintf("<'%s'>", t.Literal)
}

// Lex : Lexer
type Lex struct {
	Src string // Source
	Pos int    // Position
}

// NewLex : A new lexical analyzer with source
func NewLex(src string) *Lex {
	return &Lex{
		Src: src,
		Pos: 0,
	}
}

// Analyze command and convert them to expression
// A command line is an expression
// Command does not have a second expression
func (l Lex) parse() Expr {
	kind := -1 // Default, command header not entered
	var expr interface{}

	for l.Pos < len(l.Src) {
		if kind == 5 { // Error expression
			break
		}
		l.skipWhiteSpace()
		lit := ""
		for {
			if l.end() || l.space() {
				break
			}
			lit += string(l.now()) // Token literal
			l.Pos++
		}
		if lit != "" {
			if kind == -1 { // Command header parsing
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
				case "gt", "GT":
					kind = 4
					expr = GtExpr{}
				default:
					fmt.Printf(
						"unknown command prefix '%s'\n",
						lit,
					)
					kind = 5
					expr = ErExpr{}
				}
			} else {
				ast, err := setParam(&expr, Token{Literal: lit}) // Set command params
				if err != nil {
					fmt.Println(err.Error())
					kind = 5
				}
				expr = ast
			}
			lit = ""
		}
		l.Pos++
	}
	if err := verifyExpr(expr.(Expr)); err != nil { // Expr is legal or not
		fmt.Println(err.Error())
		expr = ErExpr{}
	}
	return expr.(Expr)
}

// Set command parameters in turn
// Their heads are all table name
func setParam(expr *interface{}, tok Token) (Expr, error) {
	switch (*expr).(type) {
	case SeExpr:
		seExpr := (*expr).(SeExpr)
		if emptyValue(seExpr.Table) {
			return SeExpr{
				Table: tok,
			}, nil
		} else {
			seExpr.F = append(seExpr.F, tok) // *F
			return seExpr, nil
		}
	case GeExpr:
		geExpr := (*expr).(GeExpr)
		if emptyValue(geExpr.Table) {
			return GeExpr{
				Table: tok,
			}, nil
		} else {
			if emptyValue(geExpr.F) { // F
				geExpr.F = tok
				return geExpr, nil
			}
			if emptyValue(geExpr.T) { // T
				geExpr.T = tok
				return geExpr, nil
			}
			return ErExpr{}, errors.New("more param specified")
		}
	case UpExpr:
		upExpr := (*expr).(UpExpr)
		if emptyValue(upExpr.Table) {
			return UpExpr{
				Table: tok,
			}, nil
		} else {
			if emptyValue(upExpr.P) { // P
				upExpr.P = tok
				return upExpr, nil
			}
			if emptyValue(upExpr.S) { // S
				upExpr.S = tok
				return upExpr, nil
			}
			if emptyValue(upExpr.N) { // N
				upExpr.N = tok
				return upExpr, nil
			}
			if emptyValue(upExpr.V) { // V
				upExpr.V = tok
				return upExpr, nil
			}
			return ErExpr{}, errors.New("more param specified")
		}
	case DeExpr:
		deExpr := (*expr).(DeExpr)
		if emptyValue(deExpr.Table) {
			return DeExpr{
				Table: tok,
			}, nil
		} else {
			if emptyValue(deExpr.P) { // P
				deExpr.P = tok
				return deExpr, nil
			}
			return ErExpr{}, errors.New("more param specified")
		}
	case GtExpr:
		gtExpr := (*expr).(GtExpr)
		if emptyValue(gtExpr.Table) {
			return GtExpr{
				Table: tok,
			}, nil
		} else {
			if emptyValue(gtExpr.S) { // S
				gtExpr.S = tok
				return gtExpr, nil
			}
			if emptyValue(gtExpr.V) { // V
				gtExpr.V = tok
				return gtExpr, nil
			}
			return ErExpr{}, errors.New("more param specified")
		}
	}
	return nil, nil
}

// Default token literal is empty
func emptyValue(t Token) bool {
	return t.Literal == ""
}

// Skip white space
func (l Lex) skipWhiteSpace() {
	for !l.end() && l.space() {
		l.Pos++
	}
}

// Is space of now character
func (l Lex) space() bool {
	return l.now() == ' ' || l.now() == '\t' || l.now() == '\r'
}

// Return character of position in source
func (l Lex) now() byte {
	return l.Src[l.Pos]
}

// FILE end
func (l Lex) end() bool {
	return l.Pos >= len(l.Src)
}

// Determine whether the expression is legal or not
func verifyExpr(expr Expr) error {
	n := reflect.TypeOf(expr).Name()
	switch n {
	case "SeExpr":
		seExpr := expr.(SeExpr)
		if emptyValue(seExpr.Table) {
			return errors.New("lost table")
		}
		if len(seExpr.F) == 0 {
			return errors.New("lost fields")
		}
		return nil
	case "GeExpr":
		geExpr := expr.(GeExpr)
		if emptyValue(geExpr.Table) {
			return errors.New("lost table")
		}
		return nil
	case "UpExpr":
		upExpr := expr.(UpExpr)
		if emptyValue(upExpr.Table) {
			return errors.New("lost table")
		}
		if emptyValue(upExpr.P) {
			return errors.New("lost position of update limit")
		}
		if emptyValue(upExpr.S) {
			return errors.New("lost field of update limit")
		}
		if emptyValue(upExpr.N) {
			return errors.New("lost new value")
		}
		return nil
	case "DeExpr":
		deExpr := expr.(DeExpr)
		if emptyValue(deExpr.Table) {
			return errors.New("lost table")
		}
		if emptyValue(deExpr.P) {
			return errors.New("lost position of delete limit")
		}
		return nil
	case "GtExpr":
		gtExpr := expr.(GtExpr)
		if emptyValue(gtExpr.Table) {
			return errors.New("lost table")
		}
		if emptyValue(gtExpr.S) {
			return errors.New("lost field of selector limit")
		}
		if emptyValue(gtExpr.V) {
			return errors.New("lost value of selector limit")
		}
	}
	return nil
}

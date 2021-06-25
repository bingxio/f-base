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
		if kind == 4 { // Error expression
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
			// Command header parsing
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
					fmt.Printf(
						"UNKNOWN COMMAND PREFIX: '%s'\n",
						lit,
					)
					kind = 4
					expr = ErExpr{}
				}
			} else {
				// Set command parameters
				// In turn
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
	// Determine whether the expression is legal or not
	if err := verifyExpr(expr.(Expr)); err != nil {
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
			// FIELDS
			seExpr.F = append(seExpr.F, tok)
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
			if emptyValue(geExpr.F) {
				geExpr.F = tok
				return geExpr, nil
			}
			// TO
			if emptyValue(geExpr.T) {
				geExpr.T = tok
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
			if emptyValue(upExpr.P) {
				upExpr.P = tok
				return upExpr, nil
			}
			// NEW
			if emptyValue(upExpr.N) {
				upExpr.N = tok
				return upExpr, nil
			}
			// VER
			if emptyValue(upExpr.V) {
				upExpr.V = tok
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
			if emptyValue(deExpr.P) {
				deExpr.P = tok
				return deExpr, nil
			}
			// VER
			if emptyValue(deExpr.V) {
				deExpr.V = tok
				return deExpr, nil
			}
			return ErExpr{}, errors.New("MORE PARAM SPECIFIED")
		}
	}
	return ErExpr{}, errors.New("PROGRAM ERROR")
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
			return errors.New("LOST TABLE")
		}
		if len(seExpr.F) == 0 {
			return errors.New("LOST FIELDS")
		}
		return nil
	case "GeExpr":
		geExpr := expr.(GeExpr)
		if emptyValue(geExpr.Table) {
			return errors.New("LOST TABLE")
		}
		return nil
	case "UpExpr":
		upExpr := expr.(UpExpr)
		if emptyValue(upExpr.Table) {
			return errors.New("LOST TABLE")
		}
		if emptyValue(upExpr.P) {
			return errors.New("LOST POSITION OF UPDATE LIMIT")
		}
		if emptyValue(upExpr.N) {
			return errors.New("LOST NEW VALUE")
		}
		return nil
	case "DeExpr":
		deExpr := expr.(DeExpr)
		if emptyValue(deExpr.Table) {
			return errors.New("LOST TABLE")
		}
		if emptyValue(deExpr.P) {
			return errors.New("LOST POSITION OF DELETE LIMIT")
		}
		return nil
	}
	return nil
}

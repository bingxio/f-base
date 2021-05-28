// F-Base
// GPL 3.0 License - bingxio <3106740988@qq.com>
package main

import "fmt"

const (
	Se = iota // INSERT
	Ge        // SELECT
	Up        // UPDATE
	De        // DELETE
	Er        // ERROR
)

// SE ? *<E>
// GE ? | <F> | -> <T>
// UP ? <P> <N> [<V>]
// DE ? <P> [<V>]
type Ast interface {
	Stringer() string
	Kind() int
}

// INSERT
type SeExpr struct {
	Table  Token
	Fields []Token
}

func (s SeExpr) Stringer() string {
	ele := func(l []Token) string {
		lit := "["
		for i := 0; i < len(l); i++ {
			lit += l[i].Literal

			if i+1 != len(l) {
				lit += ", "
			}
		}
		lit += "]"
		return lit
	}
	return fmt.Sprintf("SeExpr \n\tTable=%s \n\tFields=%s",
		s.Table.Stringer(), ele(s.Fields))
}

func (s SeExpr) Kind() int { return Se }

// SELECT
type GeExpr struct {
	Table Token
	From  Token
	To    Token
}

func (g GeExpr) Stringer() string {
	return fmt.Sprintf("GeExpr \n\tTable=%s \n\tFrom=%s \n\tTo=%s",
		g.Table.Stringer(), g.From.Stringer(), g.To.Stringer())
}

func (g GeExpr) Kind() int { return Ge }

// UPDATE
type UpExpr struct {
	Table Token
	Pos   Token
	New   Token
	Ver   Token
}

func (u UpExpr) Stringer() string {
	return fmt.Sprintf(
		"UpExpr \n\tTable=%s \n\tPos=%s \n\tNew=%s \n\tVer=%s",
		u.Table.Stringer(),
		u.Pos.Stringer(),
		u.New.Stringer(),
		u.Ver.Stringer())
}

func (u UpExpr) Kind() int { return Up }

// DELETE
type DeExpr struct {
	Table Token
	Pos   Token
	Ver   Token
}

func (d DeExpr) Stringer() string {
	return fmt.Sprintf("DeExpr \n\tTable=%s \n\tPos=%s \n\tVer=%s",
		d.Table.Stringer(), d.Pos.Stringer(), d.Ver.Stringer())
}

func (d DeExpr) Kind() int { return De }

// ERROR
type ErExpr struct{}

func (e ErExpr) Stringer() string { return "ErExpr" }
func (e ErExpr) Kind() int        { return Er }

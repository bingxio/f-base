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

// Expr
// SE ? *<E>
// GE ? | <F> | -> <T>
// UP ? <P> <N> [<V>]
// DE ? <P> [<V>]
type Expr interface {
  Stringer() string
  Kind() int
}

// SeExpr INSERT
type SeExpr struct {
  Table Token
  F     []Token
}

func (s SeExpr) Stringer() string {
  return fmt.Sprintf("SeExpr \n\tTable=%s \n\t*F=%s",
    s.Table.Stringer(), func(l []Token) string {
      lit := "["
      for i := 0; i < len(l); i++ {
        lit += l[i].Literal

        if i+1 != len(l) {
          lit += ", "
        }
      }
      lit += "]"
      return lit
    }(s.F),
  )
}

func (s SeExpr) Kind() int { return Se }

// GeExpr SELECT
type GeExpr struct {
  Table Token
  F     Token
  T     Token
}

func (g GeExpr) Stringer() string {
  return fmt.Sprintf("GeExpr \n\tTable=%s \n\tF=%s \n\tT=%s",
    g.Table.Stringer(), g.F.Stringer(), g.T.Stringer())
}

func (g GeExpr) Kind() int { return Ge }

// UpExpr UPDATE
type UpExpr struct {
  Table Token
  P     Token
  N     Token
  V     Token
}

func (u UpExpr) Stringer() string {
  return fmt.Sprintf(
    "UpExpr \n\tTable=%s \n\tP=%s \n\tN=%s \n\tV=%s",
    u.Table.Stringer(),
    u.P.Stringer(),
    u.N.Stringer(),
    u.V.Stringer())
}

func (u UpExpr) Kind() int { return Up }

// DeExpr DELETE
type DeExpr struct {
  Table Token
  P     Token
  V     Token
}

func (d DeExpr) Stringer() string {
  return fmt.Sprintf("DeExpr \n\tTable=%s \n\tP=%s \n\tV=%s",
    d.Table.Stringer(), d.P.Stringer(), d.V.Stringer())
}

func (d DeExpr) Kind() int { return De }

// ErExpr ERROR
type ErExpr struct{}

func (e ErExpr) Stringer() string { return "ErExpr" }
func (e ErExpr) Kind() int        { return Er }

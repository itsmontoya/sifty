package matcher

import (
	"github.com/itsmontoya/sifty/docview"
	"github.com/itsmontoya/sifty/query"
)

func makeNotNode(in *query.Clause) (out notNode) {
	out.child = toNode(*in)
	return out
}

type notNode struct {
	child node
}

func (n notNode) eval(doc docview.DocView) (ok bool, err error) {
	ok, err = n.child.eval(doc)
	return !ok, err
}

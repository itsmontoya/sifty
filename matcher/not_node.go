package matcher

import "github.com/itsmontoya/sifty/query"

func makeNotNode(in *query.Clause) (out notNode, err error) {
	out.child, err = toNode(*in)
	return out, err
}

type notNode struct {
	child node
}

func (n notNode) eval(doc DocView) (ok bool, err error) {
	ok, err = n.child.eval(doc)
	return !ok, err
}

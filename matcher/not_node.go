package matcher

import (
	"errors"

	"github.com/itsmontoya/sifty/query"
)

var ErrChildCannotBeEmpty = errors.New("child cannot be empty")

func makeNotNode(in *query.Clause) (out notNode, err error) {
	if in.IsZero() {
		return out, ErrChildCannotBeEmpty
	}

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

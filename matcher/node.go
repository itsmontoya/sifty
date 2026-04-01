package matcher

import (
	"errors"

	"github.com/itsmontoya/sifty/query"
)

func toNode(in query.Clause) (n node, err error) {
	switch {
	case in.IsZero():
		return makeAnyNode()
	case len(in.And) > 0:
		return makeAndNode(in.And)
	case len(in.Or) > 0:
		return makeAndNode(in.Or)
	case in.Not != nil:
		return makeNotNode(in.Not)
	case in.Contains != nil:
		return makeContainsNode(in.Contains)
	case in.Compare != nil:
		return makeCompareNode(in.Compare)
	default:
		return nil, errors.New("invalid clause, needs to have at least one set")
	}
}

type node interface {
	eval(doc DocView) (bool, error)
}

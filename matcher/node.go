package matcher

import (
	"github.com/itsmontoya/sifty/docview"
	"github.com/itsmontoya/sifty/query"
)

func toNode(in query.Clause) (n node) {
	switch {
	case len(in.And) > 0:
		return makeAndNode(in.And)
	case len(in.Or) > 0:
		return makeOrNode(in.Or)
	case in.Not != nil:
		return makeNotNode(in.Not)
	case in.Contains != nil:
		return makeContainsNode(in.Contains)
	case in.Compare != nil:
		return makeCompareNode(in.Compare)
	default:
		panic("matcher.toNode: invalid clause after validation")
	}
}

type node interface {
	eval(doc docview.DocView) (bool, error)
}

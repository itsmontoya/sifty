package matcher

import "github.com/itsmontoya/sifty/query"

func makeCompareNode(in *query.CompareExpr) (out compareNode) {
	out.field = in.Field
	out.gt = in.Gt
	out.gte = in.Gte
	out.lt = in.Lt
	out.lte = in.Lte
	return out
}

type compareNode struct {
	field string

	eq  any
	gt  any
	gte any
	lt  any
	lte any
}

func (n compareNode) eval(doc DocView) (ok bool, err error) {
	var val any
	if val, ok, err = doc.Get(n.field); !ok || err != nil {
		return ok, err
	}

	switch {
	case n.eq != nil:
		return isEqualTo(val, n.eq), nil
	case n.gt != nil:
		return isGreaterThan(val, n.gt), nil
	case n.gte != nil:
		return isGreaterThanOrEqualTo(val, n.gt), nil
	case n.lt != nil:
		return isLessThan(val, n.gt), nil
	case n.lte != nil:
		return isLessThanOrEqualTo(val, n.gt), nil
	default:
		return false, nil
	}
}

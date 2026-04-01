package matcher

import (
	"strings"

	"github.com/itsmontoya/sifty/query"
)

func makeContainsNode(in *query.ContainsExpr) (out containsNode, err error) {
	out.field = in.Field
	out.value = in.Value
	return out, nil
}

type containsNode struct {
	field string
	value string
}

func (n containsNode) eval(doc DocView) (ok bool, err error) {
	var val any
	if val, ok, err = doc.Get(n.field); !ok || err != nil {
		return false, err
	}

	var str string
	if str, ok = val.(string); !ok {
		return false, nil
	}

	return strings.Contains(str, n.value), nil
}

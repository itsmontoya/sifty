package matcher

import "github.com/itsmontoya/sifty/query"

func makeOrNode(in []query.Clause) (out orNode, err error) {
	for _, c := range in {
		var n node
		if n, err = toNode(c); err != nil {
			return out, err
		}

		out.children = append(out.children, n)
	}

	return
}

type orNode struct {
	children []node
}

func (n orNode) eval(doc DocView) (ok bool, err error) {
	for _, child := range n.children {
		ok, err = child.eval(doc)
		switch {
		case err != nil:
			return false, err
		case ok:
			return true, nil
		}
	}

	return false, nil
}

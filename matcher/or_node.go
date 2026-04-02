package matcher

import "github.com/itsmontoya/sifty/query"

func makeOrNode(in []query.Clause) (out orNode) {
	for _, c := range in {
		out.children = append(out.children, toNode(c))
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

package matcher

import "github.com/itsmontoya/sifty/query"

func makeAndNode(in []query.Clause) (out andNode, err error) {
	for _, c := range in {
		var n node
		if n, err = toNode(c); err != nil {
			return out, err
		}

		out.children = append(out.children, n)
	}

	return
}

type andNode struct {
	children []node
}

func (n andNode) eval(doc DocView) (ok bool, err error) {
	for _, child := range n.children {
		if ok, err = child.eval(doc); !ok || err != nil {
			return false, err
		}
	}

	return true, nil
}

package matcher

import "github.com/itsmontoya/sifty/query"

func makeAndNode(in []query.Clause) (out andNode) {
	for _, c := range in {
		out.children = append(out.children, toNode(c))
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

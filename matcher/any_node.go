package matcher

import "github.com/itsmontoya/sifty/docview"

func makeAnyNode() (out anyNode) {
	return
}

type anyNode struct{}

func (n anyNode) eval(doc docview.DocView) (ok bool, err error) {
	return true, nil
}

package matcher

func makeAnyNode() (out anyNode) {
	return
}

type anyNode struct{}

func (n anyNode) eval(doc DocView) (ok bool, err error) {
	return true, nil
}

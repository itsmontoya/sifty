package matcher

type testDocView struct {
	values map[string]any
	errs   map[string]error
}

func (d testDocView) Get(path string) (out any, ok bool, err error) {
	if err = d.errs[path]; err != nil {
		return nil, false, err
	}

	if d.values == nil {
		return nil, false, nil
	}

	out, ok = d.values[path]
	return out, ok, nil
}

type testNode struct {
	ok  bool
	err error
}

func (n testNode) eval(doc DocView) (ok bool, err error) {
	return n.ok, n.err
}

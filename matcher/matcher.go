package matcher

import (
	"github.com/itsmontoya/sifty/query"
)

func Compile(q query.Query) (out *Matcher, err error) {
	if err = q.Validate(); err != nil {
		return
	}

	var m Matcher
	if m.root, err = toNode(q.Filter); err != nil {
		return nil, err
	}

	return &m, nil
}

type Matcher struct {
	root node
}

func (m *Matcher) IsMatch(in DocView) (ok bool, err error) {
	return m.root.eval(in)
}

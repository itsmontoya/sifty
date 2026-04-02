package matcher

import (
	"fmt"

	"github.com/itsmontoya/sifty/query"
)

func Compile(q query.Query) (out *Matcher, err error) {
	if err = q.Validate(); err != nil {
		fmt.Println("Oh it happens here")
		return
	}

	var m Matcher
	if q.Filter.IsZero() {
		m.root = makeAnyNode()
		return &m, err
	}

	m.root = toNode(q.Filter)
	return &m, nil
}

type Matcher struct {
	root node
}

func (m *Matcher) IsMatch(in DocView) (ok bool, err error) {
	return m.root.eval(in)
}

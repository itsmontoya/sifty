package matcher

import (
	"fmt"

	"github.com/itsmontoya/sifty/query"
)

func Compile(q query.Query) (out *Matcher, err error) {
	if err = q.Validate(); err != nil {
		err = fmt.Errorf("cannot compile, invalid query: %w", err)
		return nil, err
	}

	var m Matcher
	if q.Filter.IsZero() {
		m.root = makeAnyNode()
		return &m, nil
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

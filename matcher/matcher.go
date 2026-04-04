package matcher

import (
	"fmt"

	"github.com/itsmontoya/sifty/query"
)

// Compile validates a query and compiles its filter into an executable Matcher.
// A zero-value filter compiles to a matcher that matches all documents.
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

// Matcher evaluates compiled query filters against a DocView input.
type Matcher struct {
	root node
}

// IsMatch evaluates the compiled filter against in.
// It returns any error emitted by the underlying DocView.
func (m *Matcher) IsMatch(in DocView) (ok bool, err error) {
	return m.root.eval(in)
}

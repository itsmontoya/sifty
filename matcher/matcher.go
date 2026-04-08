package matcher

import (
	"fmt"
	"time"

	"github.com/itsmontoya/sifty/docview"
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
	m.tr = q.TimeRange
	if q.Filter.IsZero() {
		m.root = makeAnyNode()
		return &m, nil
	}

	m.root = toNode(q.Filter)
	return &m, nil
}

// Matcher evaluates compiled query filters against a DocView input.
type Matcher struct {
	tr *query.TimeRange

	root node
}

// IsMatch evaluates the compiled filter against in.
// It returns any error emitted by the underlying DocView.
func (m *Matcher) IsMatch(ts time.Time, in docview.DocView) (ok bool, err error) {
	if !m.isInRange(ts) {
		return false, nil
	}

	return m.root.eval(in)
}

func (m *Matcher) RangeBounds(ts time.Time) (compare int) {
	switch {
	case m.tr == nil:
		return 0
	case m.tr.From != nil && m.tr.From.After(ts):
		return -1
	case m.tr.To != nil && m.tr.To.Before(ts):
		return 1
	default:
		return 0
	}
}

func (m *Matcher) isInRange(ts time.Time) (ok bool) {
	return m.RangeBounds(ts) == 0
}

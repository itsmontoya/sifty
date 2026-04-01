package query

import "fmt"

type Clause struct {
	And []Clause `json:"and,omitempty"`
	Or  []Clause `json:"or,omitempty"`
	Not *Clause  `json:"not,omitempty"`

	Contains *ContainsExpr `json:"contains,omitempty"`
	Compare  *CompareExpr  `json:"compare,omitempty"`
}

func (c Clause) Validate() error {
	var set int
	if len(c.And) > 0 {
		set++
	}

	if len(c.Or) > 0 {
		set++
	}

	if c.Not != nil {
		set++
	}

	if c.Contains != nil {
		set++
	}

	if c.Compare != nil {
		set++
	}

	if set != 1 {
		return fmt.Errorf("clause must define exactly one operator")
	}

	if len(c.And) > 0 {
		for i, sub := range c.And {
			if err := sub.Validate(); err != nil {
				return fmt.Errorf("invalid and clause at index %d: %w", i, err)
			}
		}
	}

	if len(c.Or) > 0 {
		for i, sub := range c.Or {
			if err := sub.Validate(); err != nil {
				return fmt.Errorf("invalid or clause at index %d: %w", i, err)
			}
		}
	}

	if c.Not != nil {
		if err := c.Not.Validate(); err != nil {
			return fmt.Errorf("invalid not clause: %w", err)
		}
	}

	if c.Contains != nil {
		if err := c.Contains.Validate(); err != nil {
			return err
		}
	}

	if c.Compare != nil {
		if err := c.Compare.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func (c Clause) IsZero() bool {
	switch {
	case len(c.And) > 0:
		return false
	case len(c.Or) > 0:
		return false
	case c.Not != nil && !c.Not.IsZero():
		return false
	case c.Contains != nil:
		return false
	case c.Compare != nil:
		return false
	default:
		return true
	}
}

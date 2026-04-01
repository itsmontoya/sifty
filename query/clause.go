package query

import "fmt"

type Clause struct {
	And []Clause `json:"and,omitempty"`
	Or  []Clause `json:"or,omitempty"`
	Not *Clause  `json:"not,omitempty"`

	Term     *TermExpr     `json:"term,omitempty"`
	Contains *ContainsExpr `json:"contains,omitempty"`
	Range    *RangeExpr    `json:"range,omitempty"`
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

	if c.Term != nil {
		set++
	}

	if c.Contains != nil {
		set++
	}

	if c.Range != nil {
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

	if c.Term != nil {
		if err := c.Term.Validate(); err != nil {
			return err
		}
	}

	if c.Contains != nil {
		if err := c.Contains.Validate(); err != nil {
			return err
		}
	}

	if c.Range != nil {
		if err := c.Range.Validate(); err != nil {
			return err
		}
	}

	return nil
}

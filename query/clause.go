package query

import (
	"errors"
	"fmt"
)

type Clause struct {
	And []Clause `json:"and,omitempty"`
	Or  []Clause `json:"or,omitempty"`
	Not *Clause  `json:"not,omitempty"`

	Contains *ContainsExpr `json:"contains,omitempty"`
	Compare  *CompareExpr  `json:"compare,omitempty"`
}

func (c Clause) Validate() (err error) {
	if c.countOperators() != 1 {
		return fmt.Errorf("clause must define exactly one operator")
	}

	var errs []error
	errs = append(errs, c.validateAnd())
	errs = append(errs, c.validateOr())
	errs = append(errs, c.validateNot())
	errs = append(errs, c.validateContains())
	errs = append(errs, c.validateCompare())
	return errors.Join(errs...)
}

func (c Clause) IsZero() bool {
	switch {
	case len(c.And) > 0:
		return false
	case len(c.Or) > 0:
		return false
	case c.Not != nil:
		return false
	case c.Contains != nil:
		return false
	case c.Compare != nil:
		return false
	default:
		return true
	}
}

func (c Clause) countOperators() (n int) {
	if len(c.And) > 0 {
		n++
	}

	if len(c.Or) > 0 {
		n++
	}

	if c.Not != nil {
		n++
	}

	if c.Contains != nil {
		n++
	}

	if c.Compare != nil {
		n++
	}

	return n
}

func (c Clause) validateAnd() (err error) {
	if len(c.And) == 0 {
		return nil
	}

	var errs []error
	for i, sub := range c.And {
		if err = sub.Validate(); err != nil {
			errs = append(errs, fmt.Errorf("invalid AND clause at index %d: %w", i, err))
		}
	}

	return errors.Join(errs...)
}

func (c Clause) validateOr() (err error) {
	if len(c.Or) == 0 {
		return nil
	}

	var errs []error
	for i, sub := range c.Or {
		if err = sub.Validate(); err != nil {
			errs = append(errs, fmt.Errorf("invalid OR clause at index %d: %w", i, err))
		}
	}

	return errors.Join(errs...)
}

func (c Clause) validateNot() (err error) {
	if c.Not == nil {
		return nil
	}

	if err = c.Not.Validate(); err != nil {
		return fmt.Errorf("invalid NOT clause: %w", err)
	}

	return nil
}

func (c Clause) validateContains() (err error) {
	if c.Contains == nil {
		return nil
	}

	if err = c.Contains.Validate(); err != nil {
		return fmt.Errorf("invalid CONTAINS expression: %w", err)
	}

	return nil
}

func (c Clause) validateCompare() (err error) {
	if c.Compare == nil {
		return nil
	}

	if err = c.Compare.Validate(); err != nil {
		return fmt.Errorf("invalid COMPARE clause: %w", err)
	}

	return nil
}

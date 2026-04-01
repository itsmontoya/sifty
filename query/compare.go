package query

import "fmt"

type CompareExpr struct {
	Field string `json:"field"`

	Eq  any `json:"eq,omitempty"`
	Gt  any `json:"gt,omitempty"`
	Gte any `json:"gte,omitempty"`
	Lt  any `json:"lt,omitempty"`
	Lte any `json:"lte,omitempty"`
}

func (c CompareExpr) Validate() error {
	if c.Field == "" {
		return fmt.Errorf("compare.field is required")
	}

	switch {
	case c.Eq != nil:
	case c.Gt != nil:
	case c.Gte != nil:
	case c.Lt != nil:
	case c.Lte != nil:
	default:
		return fmt.Errorf("compare requires at least one bound")
	}

	return nil
}

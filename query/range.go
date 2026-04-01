package query

import "fmt"

type RangeExpr struct {
	Field string `json:"field"`
	Gt    any    `json:"gt,omitempty"`
	Gte   any    `json:"gte,omitempty"`
	Lt    any    `json:"lt,omitempty"`
	Lte   any    `json:"lte,omitempty"`
}

func (r RangeExpr) Validate() error {
	if r.Field == "" {
		return fmt.Errorf("range.field is required")
	}

	if r.Gt == nil && r.Gte == nil && r.Lt == nil && r.Lte == nil {
		return fmt.Errorf("range requires at least one bound")
	}

	return nil
}

package query

import "fmt"

type TermExpr struct {
	Field string `json:"field"`
	Value any    `json:"value"`
}

func (t TermExpr) Validate() error {
	if t.Field == "" {
		return fmt.Errorf("term.field is required")
	}

	if t.Value == nil {
		return fmt.Errorf("term.value is required")
	}

	return nil
}

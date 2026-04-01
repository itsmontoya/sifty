package query

import "fmt"

type ContainsExpr struct {
	Field string `json:"field"`
	Value string `json:"value"`
}

func (c ContainsExpr) Validate() error {
	if c.Field == "" {
		return fmt.Errorf("contains.field is required")
	}

	if c.Value == "" {
		return fmt.Errorf("contains.value is required")
	}

	return nil
}

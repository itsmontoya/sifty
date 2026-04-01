package query

import "testing"

func TestClauseValidateInvalidMultipleOperators(t *testing.T) {
	c := Clause{
		Term: &TermExpr{
			Field: "status",
			Value: "active",
		},
		Contains: &ContainsExpr{
			Field: "title",
			Value: "foo",
		},
	}

	if err := c.Validate(); err == nil {
		t.Fatal("expected validation error")
	}
}

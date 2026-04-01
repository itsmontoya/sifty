package query

import "testing"

func TestSortFieldValidateInvalidDirection(t *testing.T) {
	s := SortField{
		Field:     "created_at",
		Direction: "sideways",
	}

	if err := s.Validate(); err == nil {
		t.Fatal("expected validation error")
	}
}

package query

import "fmt"

type SortField struct {
	Field     string        `json:"field"`
	Direction SortDirection `json:"direction,omitempty"`
}

func (s SortField) Validate() error {
	if s.Field == "" {
		return fmt.Errorf("sort.field is required")
	}

	switch s.Direction {
	case "", SortDirectionAsc, SortDirectionDesc:
		return nil
	default:
		return fmt.Errorf("sort.direction must be %q or %q", SortDirectionAsc, SortDirectionDesc)
	}
}

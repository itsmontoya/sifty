package query

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
)

type Query struct {
	Filter    Clause      `json:"filter,omitempty"`
	TimeRange *TimeRange  `json:"timeRange,omitempty"`
	Sort      []SortField `json:"sort,omitempty"`

	Limit  *int `json:"limit,omitempty"`
	Offset int  `json:"offset,omitempty"`
}

func (q *Query) Validate() error {
	if q.Limit != nil && *q.Limit < 0 {
		return fmt.Errorf("limit must be >= 0")
	}

	if q.Offset < 0 {
		return fmt.Errorf("offset must be >= 0")
	}

	if !q.Filter.IsZero() {
		if err := q.Filter.Validate(); err != nil {
			return fmt.Errorf("invalid filter: %w", err)
		}
	}

	for i, sort := range q.Sort {
		if err := sort.Validate(); err != nil {
			return fmt.Errorf("invalid sort at index %d: %w", i, err)
		}
	}

	return nil
}

func (q *Query) Hash() (out string, err error) {
	h := sha256.New()
	v := *q
	v.Limit = nil
	v.Offset = 0

	if err = json.NewEncoder(h).Encode(v); err != nil {
		return "", err
	}

	return string(h.Sum(nil)), nil
}

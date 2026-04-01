package query

import (
	"encoding/json"
	"testing"
)

func TestQueryValidate(t *testing.T) {
	limit := 25
	q := Query{
		Filter: &Clause{
			And: []Clause{
				{
					Term: &TermExpr{
						Field: "status",
						Value: "active",
					},
				},
				{
					Range: &RangeExpr{
						Field: "score",
						Gte:   10,
						Lt:    100,
					},
				},
			},
		},
		Sort: []SortField{
			{Field: "created_at", Direction: SortDirectionDesc},
		},
		Limit:  &limit,
		Offset: 10,
	}

	if err := q.Validate(); err != nil {
		t.Fatalf("query should be valid: %v", err)
	}
}

func TestQueryJSONUnmarshal(t *testing.T) {
	in := []byte(`{
		"filter": {
			"or": [
				{"term": {"field": "status", "value": "active"}},
				{"contains": {"field": "title", "value": "golang"}}
			]
		},
		"sort": [
			{"field": "created_at", "direction": "desc"}
		],
		"limit": 10,
		"offset": 0
	}`)

	var q Query
	if err := json.Unmarshal(in, &q); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if err := q.Validate(); err != nil {
		t.Fatalf("query should be valid after unmarshal: %v", err)
	}
}

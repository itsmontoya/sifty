package query

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestQueryValidate(t *testing.T) {
	var (
		validLimit   = 25
		invalidLimit = -1
	)

	tt := []struct {
		name      string
		in        Query
		errSubstr string
	}{
		{
			name: "valid",
			in: Query{
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
				Limit:  &validLimit,
				Offset: 10,
			},
		},
		{
			name: "invalid limit",
			in: Query{
				Limit: &invalidLimit,
			},
			errSubstr: "limit must be >= 0",
		},
		{
			name: "invalid offset",
			in: Query{
				Offset: -1,
			},
			errSubstr: "offset must be >= 0",
		},
		{
			name: "invalid filter",
			in: Query{
				Filter: &Clause{},
			},
			errSubstr: "invalid filter:",
		},
		{
			name: "invalid sort",
			in: Query{
				Sort: []SortField{
					{Field: "created_at", Direction: "sideways"},
				},
			},
			errSubstr: "invalid sort at index 0:",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var err error
			err = tc.in.Validate()
			if tc.errSubstr == "" && err != nil {
				t.Fatalf("unexpected validation error: %v", err)
			}

			if tc.errSubstr != "" && err == nil {
				t.Fatal("expected validation error")
			}

			if tc.errSubstr != "" && !strings.Contains(err.Error(), tc.errSubstr) {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestQueryJSONUnmarshal(t *testing.T) {
	tt := []struct {
		name string
		in   []byte
	}{
		{
			name: "valid query json",
			in: []byte(`{
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
	}`),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var (
				q   Query
				err error
			)

			err = json.Unmarshal(tc.in, &q)
			if err != nil {
				t.Fatalf("unmarshal failed: %v", err)
			}

			err = q.Validate()
			if err != nil {
				t.Fatalf("query should be valid after unmarshal: %v", err)
			}
		})
	}
}

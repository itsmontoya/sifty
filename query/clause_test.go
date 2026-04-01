package query

import (
	"strings"
	"testing"
)

func TestClauseValidate(t *testing.T) {
	tt := []struct {
		name      string
		in        Clause
		errSubstr string
	}{
		{
			name: "invalid multiple operators",
			in: Clause{
				Term: &TermExpr{
					Field: "status",
					Value: "active",
				},
				Contains: &ContainsExpr{
					Field: "title",
					Value: "foo",
				},
			},
			errSubstr: "clause must define exactly one operator",
		},
		{
			name: "invalid no operator",
			in:   Clause{},
			errSubstr: "clause must define exactly one operator",
		},
		{
			name: "invalid and subclause",
			in: Clause{
				And: []Clause{
					{Term: &TermExpr{Field: "status", Value: "active"}},
					{},
				},
			},
			errSubstr: "invalid and clause at index 1:",
		},
		{
			name: "invalid or subclause",
			in: Clause{
				Or: []Clause{
					{Term: &TermExpr{Field: "status", Value: "active"}},
					{},
				},
			},
			errSubstr: "invalid or clause at index 1:",
		},
		{
			name: "invalid not subclause",
			in: Clause{
				Not: &Clause{},
			},
			errSubstr: "invalid not clause:",
		},
		{
			name: "valid not",
			in: Clause{
				Not: &Clause{
					Term: &TermExpr{Field: "status", Value: "active"},
				},
			},
		},
		{
			name: "invalid term expr",
			in: Clause{
				Term: &TermExpr{},
			},
			errSubstr: "term.field is required",
		},
		{
			name: "invalid contains expr",
			in: Clause{
				Contains: &ContainsExpr{},
			},
			errSubstr: "contains.field is required",
		},
		{
			name: "invalid range expr",
			in: Clause{
				Range: &RangeExpr{},
			},
			errSubstr: "range.field is required",
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

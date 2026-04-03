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
				Compare: &CompareExpr{
					Field: "status",
					Eq:    "active",
				},
				Contains: &ContainsExpr{
					Field: "title",
					Value: "foo",
				},
			},
			errSubstr: "clause must define exactly one operator",
		},
		{
			name:      "invalid no operator",
			in:        Clause{},
			errSubstr: "clause must define exactly one operator",
		},
		{
			name: "invalid and subclause",
			in: Clause{
				And: []Clause{
					{Compare: &CompareExpr{Field: "status", Eq: "active"}},
					{},
				},
			},
			errSubstr: "invalid AND clause at index 1:",
		},
		{
			name: "invalid or subclause",
			in: Clause{
				Or: []Clause{
					{Compare: &CompareExpr{Field: "status", Eq: "active"}},
					{},
				},
			},
			errSubstr: "invalid OR clause at index 1:",
		},
		{
			name: "invalid not subclause",
			in: Clause{
				Not: &Clause{},
			},
			errSubstr: "invalid NOT clause:",
		},
		{
			name: "valid not",
			in: Clause{
				Not: &Clause{
					Compare: &CompareExpr{Field: "status", Eq: "active"},
				},
			},
		},
		{
			name: "invalid contains expr",
			in: Clause{
				Contains: &ContainsExpr{},
			},
			errSubstr: "contains.field is required",
		},
		{
			name: "invalid compare expr",
			in: Clause{
				Compare: &CompareExpr{},
			},
			errSubstr: "compare.field is required",
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

func TestClause_IsZero(t *testing.T) {
	tests := []struct {
		name string
		c    Clause
		want bool
	}{
		{
			name: "empty clause is zero",
			c:    Clause{},
			want: true,
		},
		{
			name: "and set is non-zero",
			c: Clause{
				And: []Clause{
					{
						Compare: &CompareExpr{
							Field: "status",
							Eq:    "active",
						},
					},
				},
			},
			want: false,
		},
		{
			name: "or set is non-zero",
			c: Clause{
				Or: []Clause{
					{
						Compare: &CompareExpr{
							Field: "status",
							Eq:    "active",
						},
					},
				},
			},
			want: false,
		},
		{
			name: "contains set is non-zero",
			c: Clause{
				Contains: &ContainsExpr{
					Field: "title",
					Value: "go",
				},
			},
			want: false,
		},
		{
			name: "compare set is non-zero",
			c: Clause{
				Compare: &CompareExpr{
					Field: "status",
					Eq:    "active",
				},
			},
			want: false,
		},
		{
			name: "not with non-zero child is non-zero",
			c: Clause{
				Not: &Clause{
					Compare: &CompareExpr{
						Field: "status",
						Eq:    "active",
					},
				},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.c.IsZero()
			if got != tt.want {
				t.Errorf("IsZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

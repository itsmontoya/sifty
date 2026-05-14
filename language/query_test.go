package language

import (
	"reflect"
	"testing"
	"time"

	"github.com/itsmontoya/sifty/query"
)

func TestToQuery_OperatorMapping(t *testing.T) {
	var now time.Time
	now = time.Date(2026, time.May, 13, 9, 30, 0, 0, time.UTC)

	tests := []struct {
		name string
		ast  AST
		want query.Query
	}{
		{
			name: "eq maps to compare eq",
			ast: AST{
				Filter: &ConditionExpr{
					Field: "status",
					Op:    OpEq,
					Value: "active",
				},
			},
			want: query.Query{
				Filter: query.Clause{
					Compare: &query.CompareExpr{
						Field: "status",
						Eq:    "active",
					},
				},
			},
		},
		{
			name: "neq maps to not compare eq",
			ast: AST{
				Filter: &ConditionExpr{
					Field: "status",
					Op:    OpNeq,
					Value: "active",
				},
			},
			want: query.Query{
				Filter: query.Clause{
					Not: &query.Clause{
						Compare: &query.CompareExpr{
							Field: "status",
							Eq:    "active",
						},
					},
				},
			},
		},
		{
			name: "gt maps to compare gt",
			ast: AST{
				Filter: &ConditionExpr{
					Field: "score",
					Op:    OpGt,
					Value: 10,
				},
			},
			want: query.Query{
				Filter: query.Clause{
					Compare: &query.CompareExpr{
						Field: "score",
						Gt:    10,
					},
				},
			},
		},
		{
			name: "gte maps to compare gte",
			ast: AST{
				Filter: &ConditionExpr{
					Field: "score",
					Op:    OpGte,
					Value: 10,
				},
			},
			want: query.Query{
				Filter: query.Clause{
					Compare: &query.CompareExpr{
						Field: "score",
						Gte:   10,
					},
				},
			},
		},
		{
			name: "lt maps to compare lt",
			ast: AST{
				Filter: &ConditionExpr{
					Field: "score",
					Op:    OpLt,
					Value: 10,
				},
			},
			want: query.Query{
				Filter: query.Clause{
					Compare: &query.CompareExpr{
						Field: "score",
						Lt:    10,
					},
				},
			},
		},
		{
			name: "lte maps to compare lte",
			ast: AST{
				Filter: &ConditionExpr{
					Field: "score",
					Op:    OpLte,
					Value: 10,
				},
			},
			want: query.Query{
				Filter: query.Clause{
					Compare: &query.CompareExpr{
						Field: "score",
						Lte:   10,
					},
				},
			},
		},
		{
			name: "contains maps to contains clause",
			ast: AST{
				Filter: &ConditionExpr{
					Field: "customer",
					Op:    OpContains,
					Value: "acme",
				},
			},
			want: query.Query{
				Filter: query.Clause{
					Contains: &query.ContainsExpr{
						Field: "customer",
						Value: "acme",
					},
				},
			},
		},
		{
			name: "not contains maps to not contains clause",
			ast: AST{
				Filter: &ConditionExpr{
					Field: "customer",
					Op:    OpNotContains,
					Value: "acme",
				},
			},
			want: query.Query{
				Filter: query.Clause{
					Not: &query.Clause{
						Contains: &query.ContainsExpr{
							Field: "customer",
							Value: "acme",
						},
					},
				},
			},
		},
		{
			name: "is true maps to compare eq true",
			ast: AST{
				Filter: &ConditionExpr{
					Field: "enabled",
					Op:    OpIsTrue,
				},
			},
			want: query.Query{
				Filter: query.Clause{
					Compare: &query.CompareExpr{
						Field: "enabled",
						Eq:    true,
					},
				},
			},
		},
		{
			name: "is false maps to compare eq false",
			ast: AST{
				Filter: &ConditionExpr{
					Field: "enabled",
					Op:    OpIsFalse,
				},
			},
			want: query.Query{
				Filter: query.Clause{
					Compare: &query.CompareExpr{
						Field: "enabled",
						Eq:    false,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				got query.Query
				err error
			)

			got, err = ToQuery(tt.ast, now)
			if err != nil {
				t.Fatalf("ToQuery() failed: %v", err)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("ToQuery() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestToQuery_BooleanNestingAndNegation(t *testing.T) {
	var now time.Time
	now = time.Date(2026, time.May, 13, 9, 30, 0, 0, time.UTC)

	var (
		got query.Query
		err error
	)

	got, err = ToQuery(AST{
		Filter: &AndExpr{
			Left: &OrExpr{
				Left: &ConditionExpr{
					Field: "a",
					Op:    OpEq,
					Value: "x",
				},
				Right: &NotExpr{
					Inner: &ConditionExpr{
						Field: "b",
						Op:    OpEq,
						Value: "y",
					},
				},
			},
			Right: &ConditionExpr{
				Field: "c",
				Op:    OpEq,
				Value: "z",
			},
		},
	}, now)
	if err != nil {
		t.Fatalf("ToQuery() failed: %v", err)
	}

	var want query.Query
	want = query.Query{
		Filter: query.Clause{
			And: []query.Clause{
				{
					Or: []query.Clause{
						{
							Compare: &query.CompareExpr{
								Field: "a",
								Eq:    "x",
							},
						},
						{
							Not: &query.Clause{
								Compare: &query.CompareExpr{
									Field: "b",
									Eq:    "y",
								},
							},
						},
					},
				},
				{
					Compare: &query.CompareExpr{
						Field: "c",
						Eq:    "z",
					},
				},
			},
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("ToQuery() = %#v, want %#v", got, want)
	}
}

func TestToQuery_SortLimitSkipPropagation(t *testing.T) {
	var now time.Time
	now = time.Date(2026, time.May, 13, 9, 30, 0, 0, time.UTC)

	var (
		limit = 50
		got   query.Query
		err   error
	)

	got, err = ToQuery(AST{
		Filter: &ConditionExpr{
			Field: "status",
			Op:    OpEq,
			Value: "active",
		},
		Sort: []SortExpr{
			{
				Field: "created",
				Desc:  false,
			},
			{
				Field: "score",
				Desc:  true,
			},
		},
		Limit: &limit,
		Skip:  100,
	}, now)
	if err != nil {
		t.Fatalf("ToQuery() failed: %v", err)
	}

	var want query.Query
	want = query.Query{
		Filter: query.Clause{
			Compare: &query.CompareExpr{
				Field: "status",
				Eq:    "active",
			},
		},
		Sort: []query.SortField{
			{
				Field:     "created",
				Direction: query.SortDirectionAsc,
			},
			{
				Field:     "score",
				Direction: query.SortDirectionDesc,
			},
		},
		Limit:  &limit,
		Offset: 100,
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("ToQuery() = %#v, want %#v", got, want)
	}
}

func TestToQuery_TemporalMappings(t *testing.T) {
	var now time.Time
	now = time.Date(2026, time.May, 13, 15, 4, 5, 0, time.UTC)

	dayStart := time.Date(2026, time.May, 13, 0, 0, 0, 0, time.UTC)
	dayEnd := time.Date(2026, time.May, 13, 23, 59, 59, 0, time.UTC)
	weekStart := time.Date(2026, time.May, 11, 0, 0, 0, 0, time.UTC)
	monthStart := time.Date(2026, time.May, 1, 0, 0, 0, 0, time.UTC)
	last7 := now.Add(-7 * 24 * time.Hour)
	before := time.Date(2026, time.May, 1, 12, 0, 0, 0, time.UTC)
	after := time.Date(2026, time.April, 15, 18, 30, 0, 0, time.UTC)

	tests := []struct {
		name string
		ast  AST
		want *query.TimeRange
	}{
		{
			name: "before maps to TimeRange.To",
			ast: AST{
				Filter: &ConditionExpr{
					Field: "created_at",
					Op:    OpBefore,
					Value: before,
				},
			},
			want: &query.TimeRange{
				To: &before,
			},
		},
		{
			name: "after maps to TimeRange.From",
			ast: AST{
				Filter: &ConditionExpr{
					Field: "created_at",
					Op:    OpAfter,
					Value: after,
				},
			},
			want: &query.TimeRange{
				From: &after,
			},
		},
		{
			name: "today maps to day range",
			ast: AST{
				Filter: &ConditionExpr{
					Field: "created_at",
					Op:    OpToday,
				},
			},
			want: &query.TimeRange{
				From: &dayStart,
				To:   &dayEnd,
			},
		},
		{
			name: "this week maps to week start through now",
			ast: AST{
				Filter: &ConditionExpr{
					Field: "created_at",
					Op:    OpThisWeek,
				},
			},
			want: &query.TimeRange{
				From: &weekStart,
				To:   &now,
			},
		},
		{
			name: "this month maps to month start through now",
			ast: AST{
				Filter: &ConditionExpr{
					Field: "created_at",
					Op:    OpThisMonth,
				},
			},
			want: &query.TimeRange{
				From: &monthStart,
				To:   &now,
			},
		},
		{
			name: "in last N days maps to now minus N days through now",
			ast: AST{
				Filter: &ConditionExpr{
					Field: "created_at",
					Op:    OpInLast,
					Value: 7,
				},
			},
			want: &query.TimeRange{
				From: &last7,
				To:   &now,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				got query.Query
				err error
			)

			got, err = ToQuery(tt.ast, now)
			if err != nil {
				t.Fatalf("ToQuery() failed: %v", err)
			}

			if !reflect.DeepEqual(got.TimeRange, tt.want) {
				t.Fatalf("ToQuery().TimeRange = %#v, want %#v", got.TimeRange, tt.want)
			}
		})
	}
}

func TestToQuery_TemporalConflictsAndInvalidUsage(t *testing.T) {
	var now time.Time
	now = time.Date(2026, time.May, 13, 15, 4, 5, 0, time.UTC)

	tests := []struct {
		name string
		ast  AST
	}{
		{
			name: "reject temporal operator on non-canonical temporal field",
			ast: AST{
				Filter: &ConditionExpr{
					Field: "updated_at",
					Op:    OpToday,
				},
			},
		},
		{
			name: "reject contradictory temporal range",
			ast: AST{
				Filter: &AndExpr{
					Left: &ConditionExpr{
						Field: "created_at",
						Op:    OpAfter,
						Value: time.Date(2026, time.May, 10, 0, 0, 0, 0, time.UTC),
					},
					Right: &ConditionExpr{
						Field: "created_at",
						Op:    OpBefore,
						Value: time.Date(2026, time.May, 1, 0, 0, 0, 0, time.UTC),
					},
				},
			},
		},
		{
			name: "reject temporal constraints across multiple fields",
			ast: AST{
				Filter: &AndExpr{
					Left: &ConditionExpr{
						Field: "created_at",
						Op:    OpAfter,
						Value: time.Date(2026, time.May, 1, 0, 0, 0, 0, time.UTC),
					},
					Right: &ConditionExpr{
						Field: "updated_at",
						Op:    OpBefore,
						Value: time.Date(2026, time.May, 10, 0, 0, 0, 0, time.UTC),
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			_, err = ToQuery(tt.ast, now)
			if err == nil {
				t.Fatal("ToQuery() succeeded unexpectedly")
			}
		})
	}
}

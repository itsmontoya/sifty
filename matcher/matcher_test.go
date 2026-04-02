package matcher

import (
	"errors"
	"strings"
	"testing"

	"github.com/itsmontoya/sifty/query"
)

func TestCompile(t *testing.T) {
	tt := []struct {
		name      string
		in        query.Query
		errSubstr string
	}{
		{
			name: "valid query with empty filter compiles",
			in: query.Query{
				Filter: query.Clause{},
			},
		},
		{
			name: "invalid query returns validation error",
			in: query.Query{
				Filter: query.Clause{Not: &query.Clause{}},
			},
			errSubstr: "invalid filter:",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var (
				m   *Matcher
				err error
			)

			m, err = Compile(tc.in)
			if tc.errSubstr == "" && err != nil {
				t.Fatalf("unexpected compile error: %v", err)
			}

			if tc.errSubstr != "" && err == nil {
				t.Fatal("expected compile error")
			}

			if tc.errSubstr != "" && !strings.Contains(err.Error(), tc.errSubstr) {
				t.Fatalf("unexpected error: %v", err)
			}

			if tc.errSubstr == "" && m == nil {
				t.Fatal("expected matcher")
			}
		})
	}
}

func TestMatcherIsMatch(t *testing.T) {
	var (
		errGet = errors.New("read failed")
		errExp = errors.New("explode")
	)

	tt := []struct {
		name    string
		query   query.Query
		doc     testDocView
		wantOK  bool
		wantErr error
	}{
		{
			name: "empty filter matches any doc",
			query: query.Query{
				Filter: query.Clause{},
			},
			doc:    testDocView{},
			wantOK: true,
		},
		{
			name: "contains matches",
			query: query.Query{
				Filter: query.Clause{
					Contains: &query.ContainsExpr{
						Field: "title",
						Value: "go",
					},
				},
			},
			doc:    testDocView{values: map[string]any{"title": "golang"}},
			wantOK: true,
		},
		{
			name: "contains missing field is no match",
			query: query.Query{
				Filter: query.Clause{
					Contains: &query.ContainsExpr{
						Field: "title",
						Value: "go",
					},
				},
			},
			doc:    testDocView{values: map[string]any{}},
			wantOK: false,
		},
		{
			name: "contains doc read error bubbles up",
			query: query.Query{
				Filter: query.Clause{
					Contains: &query.ContainsExpr{
						Field: "title",
						Value: "go",
					},
				},
			},
			doc:     testDocView{errs: map[string]error{"title": errGet}},
			wantErr: errGet,
		},
		{
			name: "and short-circuits false",
			query: query.Query{
				Filter: query.Clause{
					And: []query.Clause{
						{Contains: &query.ContainsExpr{Field: "title", Value: "go"}},
						{Contains: &query.ContainsExpr{Field: "title", Value: "rust"}},
					},
				},
			},
			doc:    testDocView{values: map[string]any{"title": "golang"}},
			wantOK: false,
		},
		{
			name: "or returns true when one child matches",
			query: query.Query{
				Filter: query.Clause{
					Or: []query.Clause{
						{Contains: &query.ContainsExpr{Field: "title", Value: "rust"}},
						{Contains: &query.ContainsExpr{Field: "title", Value: "go"}},
					},
				},
			},
			doc:    testDocView{values: map[string]any{"title": "golang"}},
			wantOK: true,
		},
		{
			name: "not negates child",
			query: query.Query{
				Filter: query.Clause{
					Not: &query.Clause{
						Contains: &query.ContainsExpr{Field: "title", Value: "rust"},
					},
				},
			},
			doc:    testDocView{values: map[string]any{"title": "golang"}},
			wantOK: true,
		},
		{
			name: "or child error bubbles up",
			query: query.Query{
				Filter: query.Clause{
					Or: []query.Clause{
						{Contains: &query.ContainsExpr{Field: "title", Value: "go"}},
						{Contains: &query.ContainsExpr{Field: "boom", Value: "x"}},
					},
				},
			},
			doc:     testDocView{values: map[string]any{"title": "none"}, errs: map[string]error{"boom": errExp}},
			wantErr: errExp,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var (
				m   *Matcher
				err error
				ok  bool
			)

			m, err = Compile(tc.query)
			if err != nil {
				t.Fatalf("compile failed: %v", err)
			}

			ok, err = m.IsMatch(tc.doc)
			if tc.wantErr == nil && err != nil {
				t.Fatalf("unexpected match error: %v", err)
			}

			if tc.wantErr != nil && !errors.Is(err, tc.wantErr) {
				t.Fatalf("expected error %v, got %v", tc.wantErr, err)
			}

			if ok != tc.wantOK {
				t.Fatalf("IsMatch() = %v, want %v", ok, tc.wantOK)
			}
		})
	}
}

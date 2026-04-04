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

func TestCompileClauseTypes(t *testing.T) {
	tt := []struct {
		name string
		in   query.Query
		want any
	}{
		{
			name: "and clause compiles to andNode root",
			in: query.Query{
				Filter: query.Clause{
					And: []query.Clause{
						{Contains: &query.ContainsExpr{Field: "title", Value: "go"}},
					},
				},
			},
			want: andNode{},
		},
		{
			name: "or clause compiles to orNode root",
			in: query.Query{
				Filter: query.Clause{
					Or: []query.Clause{
						{Contains: &query.ContainsExpr{Field: "title", Value: "go"}},
					},
				},
			},
			want: orNode{},
		},
		{
			name: "not clause compiles to notNode root",
			in: query.Query{
				Filter: query.Clause{
					Not: &query.Clause{
						Contains: &query.ContainsExpr{Field: "title", Value: "go"},
					},
				},
			},
			want: notNode{},
		},
		{
			name: "contains clause compiles to containsNode root",
			in: query.Query{
				Filter: query.Clause{
					Contains: &query.ContainsExpr{Field: "title", Value: "go"},
				},
			},
			want: containsNode{},
		},
		{
			name: "compare clause compiles to compareNode root",
			in: query.Query{
				Filter: query.Clause{
					Compare: &query.CompareExpr{Field: "score", Gte: 10},
				},
			},
			want: compareNode{},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var (
				m   *Matcher
				err error
			)

			m, err = Compile(tc.in)
			if err != nil {
				t.Fatalf("unexpected compile error: %v", err)
			}

			switch tc.want.(type) {
			case andNode:
				if _, ok := m.root.(andNode); !ok {
					t.Fatalf("root type = %T, want andNode", m.root)
				}
			case orNode:
				if _, ok := m.root.(orNode); !ok {
					t.Fatalf("root type = %T, want orNode", m.root)
				}
			case notNode:
				if _, ok := m.root.(notNode); !ok {
					t.Fatalf("root type = %T, want notNode", m.root)
				}
			case containsNode:
				if _, ok := m.root.(containsNode); !ok {
					t.Fatalf("root type = %T, want containsNode", m.root)
				}
			case compareNode:
				if _, ok := m.root.(compareNode); !ok {
					t.Fatalf("root type = %T, want compareNode", m.root)
				}
			default:
				t.Fatalf("unexpected wanted type %T", tc.want)
			}
		})
	}
}

func TestCompileNestedBooleanCombinations(t *testing.T) {
	var (
		m   *Matcher
		err error
	)

	m, err = Compile(query.Query{
		Filter: query.Clause{
			And: []query.Clause{
				{
					Or: []query.Clause{
						{Contains: &query.ContainsExpr{Field: "title", Value: "go"}},
						{Compare: &query.CompareExpr{Field: "score", Gte: 10}},
					},
				},
				{
					Not: &query.Clause{
						Contains: &query.ContainsExpr{Field: "category", Value: "draft"},
					},
				},
			},
		},
	})
	if err != nil {
		t.Fatalf("unexpected compile error: %v", err)
	}

	andRoot, ok := m.root.(andNode)
	if !ok {
		t.Fatalf("root type = %T, want andNode", m.root)
	}

	if len(andRoot.children) != 2 {
		t.Fatalf("and children = %d, want 2", len(andRoot.children))
	}

	orChild, ok := andRoot.children[0].(orNode)
	if !ok {
		t.Fatalf("first child type = %T, want orNode", andRoot.children[0])
	}

	if len(orChild.children) != 2 {
		t.Fatalf("or children = %d, want 2", len(orChild.children))
	}

	if _, ok = orChild.children[0].(containsNode); !ok {
		t.Fatalf("or child[0] type = %T, want containsNode", orChild.children[0])
	}

	if _, ok = orChild.children[1].(compareNode); !ok {
		t.Fatalf("or child[1] type = %T, want compareNode", orChild.children[1])
	}

	notChild, ok := andRoot.children[1].(notNode)
	if !ok {
		t.Fatalf("second child type = %T, want notNode", andRoot.children[1])
	}

	if _, ok = notChild.child.(containsNode); !ok {
		t.Fatalf("not child type = %T, want containsNode", notChild.child)
	}
}

func TestCompileInvalidQueryErrorWrapping(t *testing.T) {
	var (
		m   *Matcher
		err error
	)

	m, err = Compile(query.Query{
		Filter: query.Clause{
			Compare: &query.CompareExpr{
				Field: "score",
			},
		},
	})

	if m != nil {
		t.Fatalf("matcher = %v, want nil", m)
	}

	if err == nil {
		t.Fatal("expected compile error")
	}

	if !strings.Contains(err.Error(), "cannot compile, invalid query:") {
		t.Fatalf("error missing compile wrapper: %v", err)
	}

	if !strings.Contains(err.Error(), "invalid filter:") {
		t.Fatalf("error missing filter wrapper: %v", err)
	}

	if !strings.Contains(err.Error(), "compare requires at least one bound") {
		t.Fatalf("error missing compare validation cause: %v", err)
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

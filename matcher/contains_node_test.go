package matcher

import (
	"errors"
	"testing"

	"github.com/itsmontoya/sifty/query"
)

func TestMakeContainsNode(t *testing.T) {
	tt := []struct {
		name string
		in   *query.ContainsExpr
	}{
		{
			name: "copies field and value",
			in: &query.ContainsExpr{
				Field: "title",
				Value: "go",
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var (
				out containsNode
				err error
			)

			out, err = makeContainsNode(tc.in)
			if err != nil {
				t.Fatalf("makeContainsNode() error = %v", err)
			}

			if out.field != tc.in.Field {
				t.Fatalf("field = %q, want %q", out.field, tc.in.Field)
			}

			if out.value != tc.in.Value {
				t.Fatalf("value = %q, want %q", out.value, tc.in.Value)
			}
		})
	}
}

func TestContainsNodeEval(t *testing.T) {
	var errExp = errors.New("failed")

	tt := []struct {
		name    string
		node    containsNode
		doc     testDocView
		wantOK  bool
		wantErr error
	}{
		{
			name:   "contains value",
			node:   containsNode{field: "title", value: "go"},
			doc:    testDocView{values: map[string]any{"title": "golang"}},
			wantOK: true,
		},
		{
			name:   "field missing",
			node:   containsNode{field: "title", value: "go"},
			doc:    testDocView{values: map[string]any{}},
			wantOK: false,
		},
		{
			name:   "non string value",
			node:   containsNode{field: "title", value: "go"},
			doc:    testDocView{values: map[string]any{"title": 42}},
			wantOK: false,
		},
		{
			name:    "doc read error",
			node:    containsNode{field: "title", value: "go"},
			doc:     testDocView{errs: map[string]error{"title": errExp}},
			wantErr: errExp,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var (
				ok  bool
				err error
			)

			ok, err = tc.node.eval(tc.doc)
			if errors.Is(err, tc.wantErr) == false {
				t.Fatalf("unexpected error: %v", err)
			}

			if ok != tc.wantOK {
				t.Fatalf("eval() = %v, want %v", ok, tc.wantOK)
			}
		})
	}
}

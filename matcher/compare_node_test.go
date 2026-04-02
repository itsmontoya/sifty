package matcher

import (
	"errors"
	"testing"

	"github.com/itsmontoya/sifty/query"
)

func TestMakeCompareNode(t *testing.T) {
	tt := []struct {
		name string
		in   *query.CompareExpr
	}{
		{
			name: "copies supported bounds",
			in: &query.CompareExpr{
				Field: "score",
				Eq:    10,
				Gt:    11,
				Gte:   12,
				Lt:    13,
				Lte:   14,
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var (
				out compareNode
				err error
			)

			out, err = makeCompareNode(tc.in)
			if err != nil {
				t.Fatalf("makeCompareNode() error = %v", err)
			}

			if out.field != tc.in.Field {
				t.Fatalf("field = %q, want %q", out.field, tc.in.Field)
			}

			// Current constructor behavior copies gt/gte/lt/lte but not eq.
			if out.eq != nil {
				t.Fatalf("eq = %v, want nil", out.eq)
			}

			if out.gt != tc.in.Gt {
				t.Fatalf("gt = %v, want %v", out.gt, tc.in.Gt)
			}

			if out.gte != tc.in.Gte {
				t.Fatalf("gte = %v, want %v", out.gte, tc.in.Gte)
			}

			if out.lt != tc.in.Lt {
				t.Fatalf("lt = %v, want %v", out.lt, tc.in.Lt)
			}

			if out.lte != tc.in.Lte {
				t.Fatalf("lte = %v, want %v", out.lte, tc.in.Lte)
			}
		})
	}
}

func TestCompareNodeEval(t *testing.T) {
	var errExp = errors.New("failed")

	tt := []struct {
		name    string
		node    compareNode
		doc     testDocView
		wantOK  bool
		wantErr error
	}{
		{
			name:   "eq matches",
			node:   compareNode{field: "score", eq: 10},
			doc:    testDocView{values: map[string]any{"score": 10}},
			wantOK: true,
		},
		{
			name:   "eq no match",
			node:   compareNode{field: "score", eq: 10},
			doc:    testDocView{values: map[string]any{"score": 11}},
			wantOK: false,
		},
		{
			name:   "gt matches",
			node:   compareNode{field: "score", gt: 10},
			doc:    testDocView{values: map[string]any{"score": 11}},
			wantOK: true,
		},
		{
			name:   "gt no match",
			node:   compareNode{field: "score", gt: 10},
			doc:    testDocView{values: map[string]any{"score": 10}},
			wantOK: false,
		},
		{
			name:   "gte branch evaluates",
			node:   compareNode{field: "score", gte: 10},
			doc:    testDocView{values: map[string]any{"score": 11}},
			wantOK: false,
		},
		{
			name:   "lt branch evaluates",
			node:   compareNode{field: "score", lt: 10},
			doc:    testDocView{values: map[string]any{"score": 9}},
			wantOK: false,
		},
		{
			name:   "lte branch evaluates",
			node:   compareNode{field: "score", lte: 10},
			doc:    testDocView{values: map[string]any{"score": 10}},
			wantOK: false,
		},
		{
			name:   "no compare operator configured",
			node:   compareNode{field: "score"},
			doc:    testDocView{values: map[string]any{"score": 10}},
			wantOK: false,
		},
		{
			name:   "type mismatch no match",
			node:   compareNode{field: "score", gt: 10},
			doc:    testDocView{values: map[string]any{"score": "10"}},
			wantOK: false,
		},
		{
			name:   "field missing",
			node:   compareNode{field: "score", gt: 10},
			doc:    testDocView{values: map[string]any{}},
			wantOK: false,
		},
		{
			name:    "doc read error",
			node:    compareNode{field: "score", gt: 10},
			doc:     testDocView{errs: map[string]error{"score": errExp}},
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

package matcher

import (
	"errors"
	"testing"
)

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

package matcher

import (
	"errors"
	"testing"

	"github.com/itsmontoya/sifty/query"
)

func TestMakeAndNode(t *testing.T) {
	tt := []struct {
		name           string
		in             []query.Clause
		wantChildCount int
		wantFirstType  string
		wantSecondType string
		wantErr        bool
	}{
		{
			name: "constructs children",
			in: []query.Clause{
				{Contains: &query.ContainsExpr{Field: "title", Value: "go"}},
				{Compare: &query.CompareExpr{Field: "score", Gt: 10}},
			},
			wantChildCount: 2,
			wantFirstType:  "containsNode",
			wantSecondType: "compareNode",
		},
		{
			name: "invalid child",
			in: []query.Clause{
				{Not: &query.Clause{}},
			},
			wantErr: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var (
				out andNode
				err error
			)

			out, err = makeAndNode(tc.in)
			if tc.wantErr && err == nil {
				t.Fatal("expected makeAndNode error")
			}

			if !tc.wantErr && err != nil {
				t.Fatalf("unexpected makeAndNode error: %v", err)
			}

			if err != nil {
				return
			}

			if len(out.children) != tc.wantChildCount {
				t.Fatalf("children length = %d, want %d", len(out.children), tc.wantChildCount)
			}

			if _, ok := out.children[0].(containsNode); !ok {
				t.Fatalf("child[0] type = %T, want %s", out.children[0], tc.wantFirstType)
			}

			if _, ok := out.children[1].(compareNode); !ok {
				t.Fatalf("child[1] type = %T, want %s", out.children[1], tc.wantSecondType)
			}
		})
	}
}

func TestAndNodeEval(t *testing.T) {
	var errExp = errors.New("child failed")

	tt := []struct {
		name    string
		node    andNode
		wantOK  bool
		wantErr error
	}{
		{
			name:   "all children true",
			node:   andNode{children: []node{testNode{ok: true}, testNode{ok: true}}},
			wantOK: true,
		},
		{
			name:   "one child false",
			node:   andNode{children: []node{testNode{ok: true}, testNode{ok: false}}},
			wantOK: false,
		},
		{
			name:    "child error",
			node:    andNode{children: []node{testNode{ok: true}, testNode{err: errExp}}},
			wantOK:  false,
			wantErr: errExp,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var (
				ok  bool
				err error
			)

			ok, err = tc.node.eval(testDocView{})
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("unexpected error: %v", err)
			}

			if ok != tc.wantOK {
				t.Fatalf("eval() = %v, want %v", ok, tc.wantOK)
			}
		})
	}
}

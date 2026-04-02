package matcher

import (
	"errors"
	"testing"

	"github.com/itsmontoya/sifty/query"
)

func TestMakeOrNode(t *testing.T) {
	tt := []struct {
		name           string
		in             []query.Clause
		wantErr        bool
		wantChildCount int
		wantFirstType  string
		wantSecondType string
	}{
		{
			name: "constructs children",
			in: []query.Clause{
				{Contains: &query.ContainsExpr{Field: "title", Value: "go"}},
				{Not: &query.Clause{Contains: &query.ContainsExpr{Field: "title", Value: "rust"}}},
			},
			wantErr:        false,
			wantChildCount: 2,
			wantFirstType:  "containsNode",
			wantSecondType: "notNode",
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
				out orNode
				err error
			)

			out, err = makeOrNode(tc.in)
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

			switch tc.wantFirstType {
			case "containsNode":
				if _, ok := out.children[0].(containsNode); !ok {
					t.Fatalf("child[0] type = %T, want %s", out.children[0], tc.wantFirstType)
				}
			case "anyNode":
				if _, ok := out.children[0].(anyNode); !ok {
					t.Fatalf("child[0] type = %T, want %s", out.children[0], tc.wantFirstType)
				}
			}

			switch tc.wantSecondType {
			case "containsNode":
				if _, ok := out.children[1].(containsNode); !ok {
					t.Fatalf("child[1] type = %T, want %s", out.children[1], tc.wantSecondType)
				}
			case "notNode":
				if _, ok := out.children[1].(notNode); !ok {
					t.Fatalf("child[1] type = %T, want %s", out.children[1], tc.wantSecondType)
				}
			}
		})
	}
}

func TestOrNodeEval(t *testing.T) {
	var errExp = errors.New("child failed")

	tt := []struct {
		name    string
		node    orNode
		wantOK  bool
		wantErr error
	}{
		{
			name:   "one child true",
			node:   orNode{children: []node{testNode{ok: false}, testNode{ok: true}}},
			wantOK: true,
		},
		{
			name:   "all children false",
			node:   orNode{children: []node{testNode{ok: false}, testNode{ok: false}}},
			wantOK: false,
		},
		{
			name:    "child error",
			node:    orNode{children: []node{testNode{ok: false}, testNode{err: errExp}}},
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

package matcher

import (
	"errors"
	"testing"

	"github.com/itsmontoya/sifty/query"
)

func TestMakeNotNode(t *testing.T) {
	tt := []struct {
		name string
		in   *query.Clause
	}{
		{
			name: "constructs child",
			in: &query.Clause{
				Contains: &query.ContainsExpr{
					Field: "title",
					Value: "go",
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var (
				out notNode
				err error
			)

			out, err = makeNotNode(tc.in)
			if err != nil {
				t.Fatalf("makeNotNode() error = %v", err)
			}

			if _, ok := out.child.(containsNode); !ok {
				t.Fatalf("child type = %T, want containsNode", out.child)
			}
		})
	}
}

func TestNotNodeEval(t *testing.T) {
	var errExp = errors.New("child failed")

	tt := []struct {
		name    string
		node    notNode
		wantOK  bool
		wantErr error
	}{
		{
			name:   "negates true",
			node:   notNode{child: testNode{ok: true}},
			wantOK: false,
		},
		{
			name:   "negates false",
			node:   notNode{child: testNode{ok: false}},
			wantOK: true,
		},
		{
			name:    "passes through errors",
			node:    notNode{child: testNode{err: errExp}},
			wantOK:  true,
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
			if errors.Is(err, tc.wantErr) == false {
				t.Fatalf("unexpected error: %v", err)
			}

			if ok != tc.wantOK {
				t.Fatalf("eval() = %v, want %v", ok, tc.wantOK)
			}
		})
	}
}

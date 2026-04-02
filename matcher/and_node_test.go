package matcher

import (
	"errors"
	"testing"
)

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

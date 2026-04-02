package matcher

import (
	"errors"
	"testing"
)

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

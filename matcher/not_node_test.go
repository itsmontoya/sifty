package matcher

import (
	"errors"
	"testing"
)

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

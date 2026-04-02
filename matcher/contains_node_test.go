package matcher

import (
	"errors"
	"testing"
)

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

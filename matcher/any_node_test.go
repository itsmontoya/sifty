package matcher

import "testing"

func TestMakeAnyNode(t *testing.T) {
	tt := []struct {
		name string
	}{
		{name: "constructs any node"},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var (
				out anyNode
				err error
			)

			out, err = makeAnyNode()
			if err != nil {
				t.Fatalf("makeAnyNode() error = %v", err)
			}

			_ = out
		})
	}
}

func TestAnyNodeEval(t *testing.T) {
	var (
		n   anyNode
		ok  bool
		err error
	)

	ok, err = n.eval(testDocView{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !ok {
		t.Fatal("anyNode should always match")
	}
}

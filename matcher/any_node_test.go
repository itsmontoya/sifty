package matcher

import "testing"

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

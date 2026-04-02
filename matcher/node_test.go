package matcher

import (
	"testing"

	"github.com/itsmontoya/sifty/query"
)

func TestToNode(t *testing.T) {
	tt := []struct {
		name string
		in   query.Clause
		want any
	}{
		{
			name: "zero clause returns anyNode",
			in:   query.Clause{},
			want: anyNode{},
		},
		{
			name: "and clause returns andNode",
			in: query.Clause{And: []query.Clause{
				{Contains: &query.ContainsExpr{Field: "title", Value: "go"}},
			}},
			want: andNode{},
		},
		{
			name: "or clause returns orNode",
			in: query.Clause{Or: []query.Clause{
				{Contains: &query.ContainsExpr{Field: "title", Value: "go"}},
			}},
			want: orNode{},
		},
		{
			name: "not clause returns notNode",
			in:   query.Clause{Not: &query.Clause{Contains: &query.ContainsExpr{Field: "title", Value: "go"}}},
			want: notNode{},
		},
		{
			name: "contains clause returns containsNode",
			in:   query.Clause{Contains: &query.ContainsExpr{Field: "title", Value: "go"}},
			want: containsNode{},
		},
		{
			name: "compare clause returns compareNode",
			in:   query.Clause{Compare: &query.CompareExpr{Field: "score", Gt: 10}},
			want: compareNode{},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var (
				got node
				err error
			)

			got, err = toNode(tc.in)
			if err != nil {
				t.Fatalf("toNode() error = %v", err)
			}

			switch tc.want.(type) {
			case anyNode:
				if _, ok := got.(anyNode); !ok {
					t.Fatalf("toNode() type = %T, want anyNode", got)
				}
			case andNode:
				if _, ok := got.(andNode); !ok {
					t.Fatalf("toNode() type = %T, want andNode", got)
				}
			case orNode:
				if _, ok := got.(orNode); !ok {
					t.Fatalf("toNode() type = %T, want orNode", got)
				}
			case notNode:
				if _, ok := got.(notNode); !ok {
					t.Fatalf("toNode() type = %T, want notNode", got)
				}
			case containsNode:
				if _, ok := got.(containsNode); !ok {
					t.Fatalf("toNode() type = %T, want containsNode", got)
				}
			case compareNode:
				if _, ok := got.(compareNode); !ok {
					t.Fatalf("toNode() type = %T, want compareNode", got)
				}
			default:
				t.Fatalf("unexpected wanted type %T", tc.want)
			}
		})
	}
}

package matcher

import "testing"

func TestCompare(t *testing.T) {
	tt := []struct {
		name string
		a    any
		b    any
		want comparison
	}{
		{name: "string less than", a: "abc", b: "abd", want: comparisonLessThan},
		{name: "string equal", a: "abc", b: "abc", want: comparisonEqualTo},
		{name: "string greater than", a: "abd", b: "abc", want: comparisonGreaterThan},
		{name: "int less than", a: 1, b: 2, want: comparisonLessThan},
		{name: "int equal", a: 2, b: 2, want: comparisonEqualTo},
		{name: "int greater than", a: 3, b: 2, want: comparisonGreaterThan},
		{name: "int8 compare", a: int8(4), b: int8(1), want: comparisonGreaterThan},
		{name: "int16 compare", a: int16(1), b: int16(2), want: comparisonLessThan},
		{name: "int32 compare", a: int32(7), b: int32(7), want: comparisonEqualTo},
		{name: "int64 compare", a: int64(9), b: int64(3), want: comparisonGreaterThan},
		{name: "uint compare", a: uint(1), b: uint(2), want: comparisonLessThan},
		{name: "uint8 compare", a: uint8(2), b: uint8(2), want: comparisonEqualTo},
		{name: "uint16 compare", a: uint16(5), b: uint16(1), want: comparisonGreaterThan},
		{name: "uint32 compare", a: uint32(4), b: uint32(9), want: comparisonLessThan},
		{name: "uint64 compare", a: uint64(10), b: uint64(10), want: comparisonEqualTo},
		{name: "float32 compare", a: float32(2.5), b: float32(2.5), want: comparisonEqualTo},
		{name: "float64 compare", a: float64(9.1), b: float64(8.2), want: comparisonGreaterThan},
		{name: "uintptr compare", a: uintptr(3), b: uintptr(7), want: comparisonLessThan},
		{name: "type mismatch", a: 1, b: int64(1), want: comparisonNoMatch},
		{name: "unsupported lhs type", a: true, b: true, want: comparisonNoMatch},
		{name: "nil lhs", a: nil, b: 1, want: comparisonNoMatch},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got := compare(tc.a, tc.b)
			if got != tc.want {
				t.Fatalf("compare(%v, %v) = %v, want %v", tc.a, tc.b, got, tc.want)
			}
		})
	}
}

func TestCompareAny(t *testing.T) {
	tt := []struct {
		name string
		a    int
		b    any
		want comparison
	}{
		{name: "less than", a: 1, b: 2, want: comparisonLessThan},
		{name: "equal", a: 2, b: 2, want: comparisonEqualTo},
		{name: "greater than", a: 3, b: 2, want: comparisonGreaterThan},
		{name: "type mismatch", a: 3, b: int64(3), want: comparisonNoMatch},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got := compareany(tc.a, tc.b)
			if got != tc.want {
				t.Fatalf("compareany(%v, %v) = %v, want %v", tc.a, tc.b, got, tc.want)
			}
		})
	}
}

func TestCompareHelpers(t *testing.T) {
	tt := []struct {
		name string
		a    any
		b    any
		eq   bool
		gt   bool
		gte  bool
		lt   bool
		lte  bool
	}{
		{
			name: "equal",
			a:    5,
			b:    5,
			eq:   true,
			gte:  true,
			lte:  true,
		},
		{
			name: "greater than",
			a:    6,
			b:    5,
			gt:   true,
			gte:  true,
		},
		{
			name: "less than",
			a:    4,
			b:    5,
			lt:   true,
			lte:  true,
		},
		{
			name: "no match type mismatch",
			a:    5,
			b:    int64(5),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if got := isEqualTo(tc.a, tc.b); got != tc.eq {
				t.Fatalf("isEqualTo(%v, %v) = %v, want %v", tc.a, tc.b, got, tc.eq)
			}

			if got := isGreaterThan(tc.a, tc.b); got != tc.gt {
				t.Fatalf("isGreaterThan(%v, %v) = %v, want %v", tc.a, tc.b, got, tc.gt)
			}

			if got := isGreaterThanOrEqualTo(tc.a, tc.b); got != tc.gte {
				t.Fatalf("isGreaterThanOrEqualTo(%v, %v) = %v, want %v", tc.a, tc.b, got, tc.gte)
			}

			if got := isLessThan(tc.a, tc.b); got != tc.lt {
				t.Fatalf("isLessThan(%v, %v) = %v, want %v", tc.a, tc.b, got, tc.lt)
			}

			if got := isLessThanOrEqualTo(tc.a, tc.b); got != tc.lte {
				t.Fatalf("isLessThanOrEqualTo(%v, %v) = %v, want %v", tc.a, tc.b, got, tc.lte)
			}
		})
	}
}

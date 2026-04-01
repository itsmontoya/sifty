package matcher

import "cmp"

func compare(a, b any) (c comparison) {
	switch aVal := a.(type) {
	case string:
		return compareany(aVal, b)

	case int:
		return compareany(aVal, b)
	case int8:
		return compareany(aVal, b)
	case int16:
		return compareany(aVal, b)
	case int32:
		return compareany(aVal, b)
	case int64:
		return compareany(aVal, b)

	case uint:
		return compareany(aVal, b)
	case uint8:
		return compareany(aVal, b)
	case uint16:
		return compareany(aVal, b)
	case uint32:
		return compareany(aVal, b)
	case uint64:
		return compareany(aVal, b)
	case uintptr:
		return compareany(aVal, b)

	case float32:
		return compareany(aVal, b)
	case float64:
		return compareany(aVal, b)

	default:
		return comparisonNoMatch
	}
}

func compareany[T cmp.Ordered](a T, b any) (c comparison) {
	bVal, ok := b.(T)
	if !ok {
		return comparisonNoMatch
	}

	switch {
	case a < bVal:
		return comparisonLessThan
	case a > bVal:
		return comparisonGreaterThan
	default:
		return comparisonEqualTo
	}
}

func isGreaterThan(a, b any) (ok bool) {
	c := compare(a, b)
	return c == comparisonGreaterThan
}

func isGreaterThanOrEqualTo(a, b any) (ok bool) {
	c := compare(a, b)
	return c == comparisonGreaterThan || c == comparisonEqualTo
}

func isLessThan(a, b any) (ok bool) {
	c := compare(a, b)
	return c == comparisonLessThan
}

func isLessThanOrEqualTo(a, b any) (ok bool) {
	c := compare(a, b)
	return c == comparisonLessThan || c == comparisonEqualTo
}

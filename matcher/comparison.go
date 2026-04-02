package matcher

const (
	comparisonNoMatch comparison = iota
	comparisonLessThan
	comparisonEqualTo
	comparisonGreaterThan
)

type comparison uint8

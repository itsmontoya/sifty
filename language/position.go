package language

// Position identifies where a token starts in the source input.
// Offset is a zero-based rune offset (not a byte offset).
// Line and Column are one-based.
type Position struct {
	Offset int
	Line   int
	Column int
}

package language

import "testing"

func TestScannerPeekAtEndReturnsZeroRune(t *testing.T) {
	var (
		s    scanner
		char rune
	)

	s = makeScanner("a")
	s.advance()

	char = s.peek()
	if char != 0 {
		t.Fatalf("peek() = %q, want %q", char, rune(0))
	}
}

func TestScannerAdvanceNewlineUpdatesLineAndColumn(t *testing.T) {
	var (
		s scanner
	)

	s = makeScanner("\n")
	s.advance()

	if s.line != 2 {
		t.Fatalf("line = %d, want %d", s.line, 2)
	}

	if s.column != 1 {
		t.Fatalf("column = %d, want %d", s.column, 1)
	}
}

func TestScannerAdvanceTabIncrementsColumnByOne(t *testing.T) {
	var (
		s scanner
	)

	s = makeScanner("\t")
	s.advance()

	if s.line != 1 {
		t.Fatalf("line = %d, want %d", s.line, 1)
	}

	if s.column != 2 {
		t.Fatalf("column = %d, want %d", s.column, 2)
	}
}

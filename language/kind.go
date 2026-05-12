package language

import (
	"strings"
)

const (
	// KindUnknown indicates an unclassified token kind.
	KindUnknown Kind = iota
	// KindEOF marks end-of-input.
	KindEOF

	// punctuation
	// KindLParen is "(".
	KindLParen
	// KindRParen is ")".
	KindRParen
	// KindComma is ",".
	KindComma

	// literals/words
	// KindNumber is a numeric literal.
	KindNumber
	// KindString is a quoted string literal.
	KindString
	// KindIdentifier is a non-keyword word token.
	KindIdentifier

	// keywords
	// KindAnd is keyword "and".
	KindAnd
	// KindOr is keyword "or".
	KindOr
	// KindNot is keyword "not".
	KindNot
	// KindIs is keyword "is".
	KindIs
	// KindContains is keyword "contains".
	KindContains
	// KindContain is keyword "contain".
	KindContain
	// KindDoes is keyword "does".
	KindDoes
	// KindGreater is keyword "greater".
	KindGreater
	// KindLess is keyword "less".
	KindLess
	// KindThan is keyword "than".
	KindThan
	// KindAt is keyword "at".
	KindAt
	// KindLeast is keyword "least".
	KindLeast
	// KindMost is keyword "most".
	KindMost
	// KindBefore is keyword "before".
	KindBefore
	// KindAfter is keyword "after".
	KindAfter
	// KindIn is keyword "in".
	KindIn
	// KindThe is keyword "the".
	KindThe
	// KindLast is keyword "last".
	KindLast
	// KindToday is keyword "today".
	KindToday
	// KindThis is keyword "this".
	KindThis
	// KindWeek is keyword "week".
	KindWeek
	// KindMonth is keyword "month".
	KindMonth
	// KindSorted is keyword "sorted".
	KindSorted
	// KindBy is keyword "by".
	KindBy
	// KindAscending is keyword "ascending".
	KindAscending
	// KindDescending is keyword "descending".
	KindDescending
	// KindLimit is keyword "limit".
	KindLimit
	// KindSkip is keyword "skip".
	KindSkip
	// KindTrue is keyword "true".
	KindTrue
	// KindFalse is keyword "false".
	KindFalse
	// KindDays is keyword "days".
	KindDays
)

// Kind describes a token's lexical category.
type Kind uint8

func getKind(in string) (out Kind) {
	in = strings.ToLower(in)
	switch in {
	case "contains":
		return KindContains
	case "contain":
		return KindContain
	case "and":
		return KindAnd
	case "or":
		return KindOr
	case "not":
		return KindNot
	case "is":
		return KindIs

	case "true":
		return KindTrue
	case "false":
		return KindFalse

	case "greater":
		return KindGreater
	case "less":
		return KindLess
	case "than":
		return KindThan

	case "at":
		return KindAt
	case "least":
		return KindLeast
	case "most":
		return KindMost
	case "does":
		return KindDoes
	case "before":
		return KindBefore
	case "after":
		return KindAfter

	case "in":
		return KindIn
	case "the":
		return KindThe
	case "last":
		return KindLast
	case "this":
		return KindThis

	case "today":
		return KindToday
	case "week":
		return KindWeek
	case "month":
		return KindMonth
	case "days":
		return KindDays

	case "limit":
		return KindLimit
	case "skip":
		return KindSkip

	case "sorted":
		return KindSorted
	case "by":
		return KindBy

	case "ascending":
		return KindAscending
	case "descending":
		return KindDescending

	default:
		return KindIdentifier
	}
}

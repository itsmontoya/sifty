package language

import "fmt"

func makeLexError(message, offender string, position Position) (e lexError) {
	e.Message = message
	e.Offender = offender
	e.Position = position
	return e
}

// lexerror represents a tokenizer/lexer error with source context.
type lexError struct {
	Message  string
	Offender string // optional: rune or token text; empty when unavailable

	Position Position
}

func (e lexError) Error() string {
	if e.Offender != "" {
		return e.withOffenderMessage()
	}

	return e.withoutOffenderMessage()
}

func (e lexError) withOffenderMessage() string {
	return fmt.Sprintf("%s at line %d, column %d (offset %d): %q",
		e.Message,
		e.Position.Line,
		e.Position.Column,
		e.Position.Offset,
		e.Offender,
	)
}

func (e lexError) withoutOffenderMessage() string {
	return fmt.Sprintf("%s at line %d, column %d (offset %d)",
		e.Message,
		e.Position.Line,
		e.Position.Column,
		e.Position.Offset,
	)
}

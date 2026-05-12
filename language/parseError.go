package language

import "fmt"

func makeParseError(message string, at Token) (p parseError) {
	p.Message = message
	p.Position = at.Position
	p.Lexeme = at.Lexeme
	p.Kind = at.Kind
	return p
}

type parseError struct {
	Message  string
	Position Position
	Lexeme   string
	Kind     Kind
}

func (e parseError) Error() string {
	if e.Lexeme != "" {
		return e.withLexeme()
	}

	return e.withoutLexeme()
}

func (e parseError) withLexeme() string {
	return fmt.Sprintf("%s at line %d, column %d (offset %d): %q",
		e.Message,
		e.Position.Line,
		e.Position.Column,
		e.Position.Offset,
		e.Lexeme,
	)
}

func (e parseError) withoutLexeme() string {
	return fmt.Sprintf("%s at line %d, column %d (offset %d)",
		e.Message,
		e.Position.Line,
		e.Position.Column,
		e.Position.Offset,
	)
}

package language

var emptyToken = Token{
	Kind: KindEOF,
}

// ToTokens tokenizes input into lexer tokens with positional metadata.
func ToTokens(input string) (out []Token, err error) {
	var (
		token   Token
		emitted bool
	)

	s := makeScanner(input)
	out = make([]Token, 0, 32)

	for !s.isAtEnd() {
		token, emitted, err = s.scanNextToken()
		switch {
		case err != nil:
			// Return encountered error
			return nil, err
		case !emitted:
			// Continue in loop without action
		default:
			// Append token to output
			out = append(out, token)
		}
	}

	out = append(out, makeToken(KindEOF, "", s.Position()))
	return out, nil
}

func makeToken(kind Kind, lexeme string, position Position) (t Token) {
	t.Kind = kind
	t.Lexeme = lexeme
	t.Position = position
	return t
}

// Token is a single lexical unit emitted by ToTokens.
type Token struct {
	Kind     Kind
	Lexeme   string
	Position Position
}

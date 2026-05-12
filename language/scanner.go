package language

func makeScanner(input string) (s scanner) {
	s.src = []rune(input)
	s.line = 1
	s.column = 1
	return s
}

type scanner struct {
	src    []rune
	i      int // current rune index / offset
	line   int // 1-based
	column int // 1-based
}

func (s *scanner) Position() (p Position) {
	p.Offset = s.i
	p.Line = s.line
	p.Column = s.column
	return p
}

func (s *scanner) isAtEnd() bool {
	return s.i >= len(s.src)
}

func (s *scanner) peek() (char rune) {
	if s.isAtEnd() {
		return 0
	}

	return s.src[s.i]
}

func (s *scanner) advance() rune {
	ch := s.src[s.i]
	s.i++

	if ch == '\n' {
		s.line++
		s.column = 1
		return ch
	}

	s.column++
	return ch
}

func (s *scanner) scanQuotedString() (string, error) {
	start := s.i // opening quote index
	s.advance()  // consume opening quote

	for !s.isAtEnd() {
		if s.peek() == '"' {
			end := s.i  // index of closing quote
			s.advance() // consume closing quote
			return string(s.src[start+1 : end]), nil
		}

		s.advance()
	}

	p := s.Position()
	return "", makeLexError("unterminated string", "", p)
}

func (s *scanner) scanNumber() string {
	start := s.i
	for !s.isAtEnd() && isDigit(s.peek()) {
		s.advance()
	}

	return string(s.src[start:s.i])
}

func (s *scanner) scanWord() string {
	start := s.i
	for !s.isAtEnd() && isWordPart(s.peek()) {
		s.advance()
	}

	return string(s.src[start:s.i])
}

func (s *scanner) scanNextToken() (token Token, emitted bool, err error) {
	var (
		start Position
		ch    rune
		lex   string
	)

	start = s.Position()
	ch = s.peek()

	switch {
	case isWhitespace(ch):
		s.advance()
		return token, false, nil

	case ch == '(':
		s.advance()
		return makeToken(KindLParen, "(", start), true, nil

	case ch == ')':
		s.advance()
		return makeToken(KindRParen, ")", start), true, nil

	case ch == ',':
		s.advance()
		return makeToken(KindComma, ",", start), true, nil

	case ch == '"':
		lex, err = s.scanQuotedString()
		if err != nil {
			return token, false, err
		}
		return makeToken(KindString, lex, start), true, nil

	case isDigit(ch):
		lex = s.scanNumber()
		return makeToken(KindNumber, lex, start), true, nil

	case isWordStart(ch):
		lex = s.scanWord()
		return makeToken(getKind(lex), lex, start), true, nil

	default:
		err = makeLexError("invalid character", string(ch), start)
		return token, false, err
	}
}

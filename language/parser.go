package language

import (
	"fmt"
	"slices"
	"strconv"
)

func makeParser(tokens []Token) (p parser) {
	if len(tokens) == 0 {
		tokens = append(tokens, emptyToken)
	}

	p.tokens = tokens
	p.current = 0
	return p
}

type parser struct {
	tokens  []Token
	current int
}

func (p *parser) isAtEnd() bool {
	return p.peek().Kind == KindEOF
}

func (p *parser) peek() (t Token) {
	if p.current < 0 || p.current >= len(p.tokens) {
		return Token{Kind: KindEOF}
	}

	return p.tokens[p.current]
}

func (p *parser) previous() (t Token) {
	if p.current-1 < 0 || p.current-1 >= len(p.tokens) {
		return Token{Kind: KindEOF}
	}

	return p.tokens[p.current-1]
}

func (p *parser) advance() (t Token) {
	if !p.isAtEnd() {
		p.current++
	}

	return p.previous()
}

func (p *parser) check(kind Kind) bool {
	if p.isAtEnd() {
		return kind == KindEOF
	}

	return p.peek().Kind == kind
}

func (p *parser) match(kinds ...Kind) (match bool) {
	return slices.ContainsFunc(kinds, func(kind Kind) bool {
		if p.check(kind) {
			p.advance()
			return true
		}

		return false
	})
}

func (p *parser) expect(kind Kind, message string) (t Token, err error) {
	if p.check(kind) {
		return p.advance(), nil
	}

	return t, makeParseError(message, p.peek())
}

// ----- Filter parsing (recursive descent) -----

// parseFilterExpr parses filter expressions using precedence layers:
// low -> high: OR, AND, unary NOT, primary.
func (p *parser) parseFilterExpr() (out Expr, err error) {
	return p.parseOr()
}

// parseOr parses left-associative OR chains.
func (p *parser) parseOr() (out Expr, err error) {
	var left, right Expr
	if left, err = p.parseAnd(); err != nil {
		return nil, err
	}

	for p.match(KindOr) {
		if right, err = p.parseAnd(); err != nil {
			return nil, err
		}

		left = OrExpr{Left: left, Right: right}
	}

	return left, nil
}

// parseAnd parses left-associative AND chains.
func (p *parser) parseAnd() (out Expr, err error) {
	var left, right Expr
	if left, err = p.parseNot(); err != nil {
		return nil, err
	}

	for p.match(KindAnd) {
		if right, err = p.parseNot(); err != nil {
			return nil, err
		}

		left = AndExpr{Left: left, Right: right}
	}

	return left, nil
}

// parseNot parses unary NOT recursively (right-associative).
func (p *parser) parseNot() (out Expr, err error) {
	if !p.match(KindNot) {
		return p.parsePrimary()
	}

	var inner Expr
	if inner, err = p.parseNot(); err != nil {
		return nil, err
	}

	return NotExpr{Inner: inner}, nil
}

func (p *parser) parsePrimary() (out Expr, err error) {
	if !p.match(KindLParen) {
		return p.parseCondition()
	}

	var inner Expr
	if inner, err = p.parseFilterExpr(); err != nil {
		return nil, err
	}

	if _, err = p.expect(KindRParen, "expected ')' after grouped expression"); err != nil {
		return nil, err
	}

	return inner, nil
}

func (p *parser) parseCondition() (out ConditionExpr, err error) {
	var fieldTok Token
	if fieldTok, err = p.expect(KindIdentifier, "expected field name"); err != nil {
		return out, err
	}

	switch {
	case p.match(KindIs):
		out, err = p.parseIsCondition()

	case p.match(KindContains), p.match(KindContain):
		out, err = p.parseContainsCondition()

	case p.match(KindDoes):
		out, err = p.parseContainsCondition()

	case p.match(KindBefore):
		out, err = p.parseTemporalCondition()

	case p.match(KindAfter):
		out, err = p.parseTemporalCondition()

	case p.match(KindToday):
		out, err = p.parseTemporalCondition()

	case p.match(KindThis):
		out, err = p.parseTemporalCondition()

	case p.match(KindIn):
		out, err = p.parseTemporalCondition()
	default:
		return out, makeParseError("expected condition operator after field", p.peek())
	}

	out.Field = fieldTok.Lexeme
	return out, nil
}

func (p *parser) parseIsCondition() (out ConditionExpr, err error) {
	// is not <value>
	if p.match(KindNot) {
		var v any
		if v, err = p.parseConditionValue(); err != nil {
			return out, err
		}

		out.Op = OpNeq
		out.Value = v
		return out, nil
	}

	// is true / is false / is <value>
	if p.match(KindTrue) {
		out.Op = OpIsTrue
		out.Value = nil
		return out, nil
	}

	if p.match(KindFalse) {
		out.Op = OpIsFalse
		out.Value = nil
		return out, nil
	}

	var v any
	if v, err = p.parseConditionValue(); err != nil {
		return out, err
	}

	out.Op = OpEq
	out.Value = v
	return out, nil
}

func (p *parser) parseContainsCondition() (out ConditionExpr, err error) {
	if p.previous().Kind == KindContains {
		var v any
		if v, err = p.parseConditionValue(); err != nil {
			return out, err
		}

		out.Op = OpContains
		out.Value = v
		return out, nil
	}

	// assumes caller matched KindDoes before calling
	if _, err = p.expect(KindNot, "expected 'not' after 'does'"); err != nil {
		return out, err
	}

	if _, err = p.expect(KindContain, "expected 'contain' after 'does not'"); err != nil {
		return out, err
	}

	var v any
	if v, err = p.parseConditionValue(); err != nil {
		return out, err
	}

	out.Op = OpNotContains
	out.Value = v
	return out, nil
}

func (p *parser) parseTemporalCondition() (out ConditionExpr, err error) {
	switch p.previous().Kind {
	case KindBefore:
		var v any
		if v, err = p.parseConditionValue(); err != nil {
			return out, err
		}

		out.Op = OpBefore
		out.Value = v
		return out, nil

	case KindAfter:
		var v any
		if v, err = p.parseConditionValue(); err != nil {
			return out, err
		}

		out.Op = OpAfter
		out.Value = v
		return out, nil

	case KindToday:
		out.Op = OpToday
		out.Value = nil
		return out, nil

	case KindThis:
		if p.match(KindWeek) {
			out.Op = OpThisWeek
			out.Value = nil
			return out, nil
		}

		if p.match(KindMonth) {
			out.Op = OpThisMonth
			out.Value = nil
			return out, nil
		}

		return out, makeParseError("expected 'week' or 'month' after 'this'", p.peek())

	case KindIn:
		if _, err = p.expect(KindThe, "expected 'the' after 'in'"); err != nil {
			return out, err
		}

		if _, err = p.expect(KindLast, "expected 'last' after 'in the'"); err != nil {
			return out, err
		}

		nTok, nErr := p.expect(KindNumber, "expected number after 'in the last'")
		if nErr != nil {
			return out, nErr
		}

		if _, err = p.expect(KindDays, "expected 'days' after duration number"); err != nil {
			return out, err
		}

		n, convErr := parseIntToken(nTok)
		if convErr != nil {
			return out, makeParseError(convErr.Error(), nTok)
		}

		out.Op = OpInLast
		out.Value = n
		return out, nil
	}

	return out, makeParseError("expected temporal condition keyword", p.previous())
}

func (p *parser) parseConditionValue() (out any, err error) {
	t := p.peek()
	switch t.Kind {
	case KindString, KindIdentifier:
		p.advance()
		return t.Lexeme, nil
	case KindNumber:
		p.advance()
		return parseIntToken(t)
	case KindTrue:
		p.advance()
		return true, nil
	case KindFalse:
		p.advance()
		return false, nil
	default:
		return nil, makeParseError("expected condition value", t)
	}
}

// ----- Sort / limit / skip -----

func (p *parser) parseSortClause() (out []SortExpr, err error) {
	// TODO: parse "sorted by ..."
	return out, nil
}

func (p *parser) parseLimitClause() (out *int, err error) {
	// TODO: parse "limit <number>"
	return out, nil
}

func (p *parser) parseSkipClause() (out *int, err error) {
	// TODO: parse "skip <number>"
	return out, nil
}

func parseIntToken(t Token) (out int, err error) {
	out, err = strconv.Atoi(t.Lexeme)
	if err != nil {
		return out, fmt.Errorf("invalid number %q: %w", t.Lexeme, err)
	}

	return out, nil
}

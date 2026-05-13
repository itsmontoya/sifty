package language

func ToAST(ts []Token) (out AST, err error) {
	p := makeParser(ts)
	if p.isAtEnd() {
		return out, nil
	}

	if out.Filter, err = p.parseFilterExpr(); err != nil {
		return out, err
	}

	if out.Sort, err = p.parseSortClause(); err != nil {
		return out, err
	}

	if out.Limit, err = p.parseLimitClause(); err != nil {
		return out, err
	}

	if out.Skip, err = p.parseSkipClause(); err != nil {
		return out, err
	}

	if !p.isAtEnd() {
		return out, makeParseError("unexpected trailing tokens", p.peek())
	}

	return out, nil
}

// AST is the parsed query representation.
type AST struct {
	Filter Expr
	Sort   []SortExpr
	Limit  *int
	Skip   *int
}

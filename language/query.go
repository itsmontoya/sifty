package language

import (
	"time"

	"github.com/itsmontoya/sifty/query"
)

func ToQuery(ast AST, now time.Time) (out query.Query, err error) {
	out.Filter = compileExpr(ast.Filter)
	out.Limit = ast.Limit
	out.Offset = ast.Skip
	return out, err
}

func compileExpr(e Expr) (out query.Clause) {
	switch t := e.(type) {
	case *AndExpr:
		out.And = append(out.And, compileExpr(t.Left), compileExpr(t.Right))
	case *OrExpr:
		out.Or = append(out.Or, compileExpr(t.Left), compileExpr(t.Right))
	case *NotExpr:
		c := compileExpr(t.Inner)
		out.Not = &c
	case *ConditionExpr:
		out = compileConditionExpr(t)
	}

	return out
}

func compileConditionExpr(t *ConditionExpr) (out query.Clause) {
	switch t.Op {
	case OpContains:
		out.Contains = &query.ContainsExpr{
			Field: t.Field,
			//	Value: t.Value,
		}

		return out
	case OpNotContains:
		c := compileExpr(t)
		out.Not = &c
		return out
	default:
		out.Compare = compileCompareExpr(t)
		out.Compare.Field = t.Field
		return
	}
}

func compileCompareExpr(t *ConditionExpr) (out *query.CompareExpr) {
	var c query.CompareExpr
	c.Field = t.Field

	switch t.Op {
	case OpGt:
		c.Gt = t.Value
	case OpGte:
		c.Gte = t.Value
	case OpLt:
		c.Lt = t.Value
	case OpLte:
		c.Lte = t.Value
	case OpEq:
		c.Eq = t.Value
	}

	return &c
}

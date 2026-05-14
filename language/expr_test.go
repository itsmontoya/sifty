package language

import "testing"

func TestExprIsExpr(t *testing.T) {
	var (
		_ Expr = expr{}
		_ Expr = AndExpr{}
		_ Expr = OrExpr{}
		_ Expr = NotExpr{}
		_ Expr = ConditionExpr{}
		_ Expr = SortExpr{}
	)

	expr{}.isExpr()
}

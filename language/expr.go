package language

var (
	_ Expr = &AndExpr{}
	_ Expr = &OrExpr{}
	_ Expr = &NotExpr{}
	_ Expr = &ConditionExpr{}
	_ Expr = &SortExpr{}
)

// Expr is the marker interface for all filter expression nodes.
type Expr interface {
	isExpr() bool
}

type expr struct{}

func (e expr) isExpr() bool {
	return true
}

// AndExpr represents a logical AND.
type AndExpr struct {
	expr

	Left  Expr
	Right Expr
}

// OrExpr represents a logical OR.
type OrExpr struct {
	expr

	Left  Expr
	Right Expr
}

// NotExpr represents a logical NOT.
type NotExpr struct {
	expr

	Inner Expr
}

// ConditionExpr is a field/operator/value predicate.
type ConditionExpr struct {
	expr

	Field string
	Op    ConditionOp
	Value any
}

// SortExpr defines one sort clause.
type SortExpr struct {
	expr

	Field string
	Desc  bool
}

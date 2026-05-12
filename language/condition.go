package language

const (
	// OpEq is equality.
	OpEq          ConditionOp = "eq"
	// OpNeq is inequality.
	OpNeq         ConditionOp = "neq"
	// OpContains checks substring/set containment.
	OpContains    ConditionOp = "contains"
	// OpNotContains checks non-containment.
	OpNotContains ConditionOp = "not_contains"
	// OpGt is greater-than.
	OpGt          ConditionOp = "gt"
	// OpGte is greater-than-or-equal.
	OpGte         ConditionOp = "gte"
	// OpLt is less-than.
	OpLt          ConditionOp = "lt"
	// OpLte is less-than-or-equal.
	OpLte         ConditionOp = "lte"
	// OpBefore is before (time/date).
	OpBefore      ConditionOp = "before"
	// OpAfter is after (time/date).
	OpAfter       ConditionOp = "after"
	// OpInLast is relative "in last <n>".
	OpInLast      ConditionOp = "in_last"
	// OpToday is "today".
	OpToday       ConditionOp = "today"
	// OpThisWeek is "this week".
	OpThisWeek    ConditionOp = "this_week"
	// OpThisMonth is "this month".
	OpThisMonth   ConditionOp = "this_month"
	// OpIsTrue checks boolean true.
	OpIsTrue      ConditionOp = "is_true"
	// OpIsFalse checks boolean false.
	OpIsFalse     ConditionOp = "is_false"
)

// ConditionOp identifies the normalized condition operator in the AST.
type ConditionOp string

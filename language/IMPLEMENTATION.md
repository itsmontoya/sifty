# Sifty Friendly Query Language: Implementation Foundation

## Purpose

This document is the implementation foundation for the first version of Sifty's friendly query language.

The primary goal of this language is approachability. It should feel readable and usable for project managers, business owners, support staff, and engineers alike. The first version should favor clarity and correctness over cleverness or shorthand.

For v1, we will require explicit parentheses for grouped boolean expressions. This keeps parsing deterministic, reduces ambiguity, and gives us a solid base to extend later.

---

## Goals

* Build a friendly, readable filter language
* Compile the language into Sifty query objects
* Keep the parser deterministic and easy to reason about
* Favor explicitness over hidden precedence rules
* Make future ergonomic improvements possible without rewriting the core

---

## Non-Goals for v1

The following are intentionally out of scope for the first implementation:

* Implicit operator precedence between `and` and `or`
* No-paren mixing of grouped boolean logic
* Fuzzy matching
* Wildcards
* Regex
* Free-text search without a field
* Natural-language sentence interpretation beyond the defined grammar
* Array-specific operators
* Advanced date math beyond a small, explicit set

---

## Core Design Principles

### 1. Friendly over technical

The language should read like instructions, not code.

Good:

```text
status is active
customer name contains acme
score is at least 70
```

Not for v1:

```text
status:"active"
customer.name ~ "acme"
score >= 70
```

### 2. Explicit grouping

When users mix boolean expressions, grouping must be explicit with parentheses.

Valid:

```text
(status is active or status is pending)
and priority is high
```

Invalid:

```text
status is active or status is pending and priority is high
```

### 3. Internal simplicity

The parser should produce a small intermediate AST, which is then compiled into Sifty query objects.

Do not parse directly into the final Sifty query type.

### 4. Future-friendly

The v1 design should make it easy to later add:

* implicit precedence
* optional parens in common cases
* sugar like `status is any of active, pending`
* more relaxed multi-line input

---

## V1 User-Facing Query Examples

### Simple

```text
status is active
priority is high
customer name contains acme
score is at least 70
not archived
```

### Grouped boolean logic

```text
(status is active or status is pending)
and priority is high
and not archived
```

### More complex

```text
(status is active or status is pending)
and (priority is high or priority is urgent)
and customer name contains acme
and score is at least 70
and not archived
and created in the last 30 days
sorted by created descending, score descending
limit 50
skip 100
```

### Time-oriented

```text
created after January 1, 2026
created before April 1, 2026
created in the last 30 days
updated today
updated this week
```

---

## V1 Syntax Rules

### Boolean rules

* `and` joins two expressions
* `or` joins two expressions
* `not` negates the next condition or grouped expression
* when `and` and `or` are mixed, parentheses are required to show grouping
* parentheses may also be used around a single expression for clarity

### Result controls

The filter expression may be followed by optional result controls:

* `sorted by <field> [ascending|descending] [, <field> [ascending|descending]]*`
* `limit <number>`
* `skip <number>`

These should be parsed separately from the filter expression and compiled into top-level query settings.

---

## Recommended V1 Grammar

This is the conceptual grammar, not necessarily the final parser implementation syntax.

```text
query         = filterExpr [sortClause] [limitClause] [skipClause]

filterExpr    = unaryExpr { booleanOp unaryExpr }
unaryExpr     = ["not"] primaryExpr
primaryExpr   = condition | "(" filterExpr ")"
booleanOp     = "and" | "or"

condition     = field equalityCondition
              | field containsCondition
              | field comparisonCondition
              | field dateCondition
              | bareBooleanCondition

equalityCondition   = "is" value
                    | "is not" value

containsCondition   = "contains" value
                    | "does not contain" value

comparisonCondition = "is greater than" value
                    | "is less than" value
                    | "is at least" value
                    | "is at most" value

dateCondition       = "before" dateValue
                    | "after" dateValue
                    | "in the last" durationValue
                    | "today"
                    | "this week"
                    | "this month"

bareBooleanCondition = "not" field
                     | field

sortClause    = "sorted by" sortField { "," sortField }
sortField     = field ["ascending" | "descending"]
limitClause   = "limit" number
skipClause    = "skip" number
```

---

## Important Clarification About Boolean Parsing

There are two possible approaches:

### Option A: fully require parens whenever boolean operators are chained

Example:

```text
((status is active or status is pending) and priority is high)
```

### Option B: allow flat chaining, but require parens whenever grouping intent matters

Example:

```text
(status is active or status is pending)
and priority is high
```

Recommended approach for v1: **Option B**.

This keeps the language readable while still requiring explicit grouping whenever `or` is involved in a meaningful way.

Validation rule:

* reject ambiguous mixes of `and` and `or` at the same level without parentheses

So this is valid:

```text
status is active and priority is high and not archived
```

This is valid:

```text
status is active or status is pending
```

This is valid:

```text
(status is active or status is pending) and priority is high
```

This is invalid:

```text
status is active or status is pending and priority is high
```

---

## Field Aliasing

Users should not need to know internal field paths.

Introduce a field alias map that resolves friendly field names into Sifty field paths.

Example:

```go
map[string]string{
    "status":        "value.status",
    "priority":      "value.priority",
    "customer name": "value.customer.name",
    "score":         "value.score",
    "created":       "@timestamp",
    "updated":       "@timestamp",
    "archived":      "value.archived",
}
```

### Notes

* Field aliases should be normalized before lookup
* Matching should likely be case-insensitive
* Multiple spaces should collapse into one
* It is acceptable for v1 to require exact alias phrases after normalization

Suggested normalization:

* trim leading and trailing space
* lowercase
* collapse repeated spaces

---

## Tokens

A simple tokenizer should be used rather than trying to parse directly from raw strings.

### Suggested token categories

* left paren
* right paren
* comma
* number
* string
* identifier
* keyword
* end of input

### Keyword examples

* `and`
* `or`
* `not`
* `is`
* `contains`
* `does`
* `greater`
* `less`
* `than`
* `at`
* `least`
* `most`
* `before`
* `after`
* `in`
* `the`
* `last`
* `today`
* `this`
* `week`
* `month`
* `sorted`
* `by`
* `ascending`
* `descending`
* `limit`
* `skip`

### String literals

Support quoted strings for safety and clarity:

```text
customer name contains "acme corp"
status is "in progress"
```

For v1, you may also allow unquoted single-token values:

```text
status is active
priority is high
```

### Recommendation

Allow both:

* quoted string values for multi-word values
* unquoted single-token values for convenience

---

## Parsing Strategy

Use a hand-written recursive descent parser.

This will likely be the clearest and easiest-to-maintain approach for the initial implementation.

### Suggested phases

1. tokenize input
2. parse filter expression into AST
3. parse optional sort, limit, and skip clauses
4. validate AST and clause structure
5. compile AST into Sifty query objects

---

## AST Design

Use a small intermediate AST.

```go
type Expr interface {
    isExpr()
}

type AndExpr struct {
    Left  Expr
    Right Expr
}

type OrExpr struct {
    Left  Expr
    Right Expr
}

type NotExpr struct {
    Inner Expr
}

type ConditionExpr struct {
    Field string
    Op    ConditionOp
    Value any
}

type ConditionOp string

const (
    OpEq             ConditionOp = "eq"
    OpNeq            ConditionOp = "neq"
    OpContains       ConditionOp = "contains"
    OpNotContains    ConditionOp = "not_contains"
    OpGt             ConditionOp = "gt"
    OpGte            ConditionOp = "gte"
    OpLt             ConditionOp = "lt"
    OpLte            ConditionOp = "lte"
    OpBefore         ConditionOp = "before"
    OpAfter          ConditionOp = "after"
    OpInLast         ConditionOp = "in_last"
    OpToday          ConditionOp = "today"
    OpThisWeek       ConditionOp = "this_week"
    OpThisMonth      ConditionOp = "this_month"
    OpIsTrue         ConditionOp = "is_true"
    OpIsFalse        ConditionOp = "is_false"
)

type SortExpr struct {
    Field string
    Desc  bool
}

type ASTQuery struct {
    Filter Expr
    Sort   []SortExpr
    Limit  *int
    Skip   *int
}
```

### Notes

* `Field` should store the resolved internal field path after alias resolution, or the friendly field name before resolution depending on where you want validation to happen
* I recommend keeping the friendly field name in the AST first, then resolving aliases during compilation or a dedicated validation pass

---

## Condition Forms

### Equality

```text
status is active
status is not active
```

### String contains

```text
customer.name contains acme
notes contains "follow up"
notes does not contain urgent
```

### Numeric comparisons

```text
score is greater than 10
score is less than 100
score is at least 70
score is at most 100
```

### Boolean shortcuts

```text
archived
not archived
archived is true
archived is false
```

Recommendation:

* `archived` alone should compile to `archived is true`
* `not archived` should compile to `archived is false` or `not (archived is true)`

For clarity and internal consistency, I recommend compiling `not archived` as `not (archived is true)`.

### Time conditions

```text
created before January 1, 2026
created after April 1, 2026
created in the last 30 days
created today
created this week
created this month
```

---

## Time Handling

Sifty has a top-level time range concept. The friendly query language should compile time-based conditions into that top-level structure where possible.

### Field mapping

For v1, treat friendly time fields like `created` and `updated` as aliases to a reserved timestamp concept.

Example:

* `created` -> `@timestamp`

### Recommended time compilation rules

* `created after X` -> `TimeRange.From = X`
* `created before X` -> `TimeRange.To = X`
* `created in the last 30 days` -> `TimeRange.From = now - 30 days`
* `created today` -> start of current day through end of current day
* `created this week` -> start of current week through now or end of week, depending on desired semantics
* `created this month` -> start of current month through now or end of month

### Important note

If multiple time conditions are present, they should be merged carefully and validated for conflicts.

Example:

```text
created after January 1, 2026 and created before April 1, 2026
```

This should produce a bounded time range.

---

## Compilation to Sifty Query Objects

The AST should compile into Sifty's internal query structure.

### Mapping guidelines

#### `AndExpr`

Compiles to a clause with `And` children.

#### `OrExpr`

Compiles to a clause with `Or` children.

#### `NotExpr`

Compiles to a clause with `Not`.

#### `ConditionExpr`

Compiles based on operator:

* `OpContains` -> `ContainsExpr`
* `OpEq` -> `CompareExpr` with `Eq`
* `OpGt` -> `CompareExpr` with `Gt`
* `OpGte` -> `CompareExpr` with `Gte`
* `OpLt` -> `CompareExpr` with `Lt`
* `OpLte` -> `CompareExpr` with `Lte`
* time operators -> top-level `TimeRange`

### Important implementation note

If Sifty's compare implementation currently handles only one comparator per compare node, then range expressions must compile into separate compare clauses joined by `and`.

Example:

```text
score is at least 70 and score is less than 100
```

Must become:

* compare `score >= 70`
* compare `score < 100`
* wrapped in `and`

Do not combine both bounds into a single compare node unless the matcher supports that correctly.

---

## Validation Rules

Validation should happen after parsing and before compilation, or during compilation if you prefer a combined pass.

### Validation checklist

* field alias exists
* operator is valid for the field type if field typing is known
* date values parse successfully
* duration values parse successfully
* limit is non-negative
* skip is non-negative
* sort fields exist
* sort directions are valid
* boolean groupings are explicit when required

### Ambiguity rejection

Reject ambiguous mixed boolean expressions without parentheses.

Suggested message:

> Ambiguous query: when mixing 'and' and 'or', please use parentheses to group conditions.

### Unknown field rejection

Suggested message:

> Unknown field: 'customer'. Did you mean 'customer name'?

This kind of error will matter a lot for UX.

---

## Error Handling Philosophy

Errors should be written for non-technical users.

Bad:

> unexpected token near IDENT at position 17

Good:

> I could not understand this part of the query: `status is active or status is pending and priority is high`
>
> When mixing `and` and `or`, please use parentheses to show what should be grouped.

### Error principles

* show the problematic part
* explain why it failed
* suggest a corrected form when possible
* avoid parser jargon unless running in debug mode

---

## Suggested Implementation Order

### Phase 1: foundation

* tokenizer
* parser for simple conditions
* AST types
* compiler for equality, contains, numeric comparisons

### Phase 2: boolean logic

* `and`
* `or`
* `not`
* parentheses
* ambiguity detection

### Phase 3: result controls

* `sorted by`
* `limit`
* `skip`

### Phase 4: time support

* `before`
* `after`
* `in the last`
* `today`
* `this week`
* `this month`

### Phase 5: validation and errors

* better user-facing messages
* field suggestions
* type validation

### Phase 6: ergonomics after stabilization

* optional no-paren precedence rules
* sugar like `is any of`
* looser multi-line support

---

## Example End-to-End

### Input

```text
(status is active or status is pending)
and (priority is high or priority is urgent)
and customer name contains "acme"
and score is at least 70
and not archived
and created in the last 30 days
sorted by created descending, score descending
limit 50
skip 100
```

### Parsed AST sketch

```text
And(
  Or(
    Eq(status, active),
    Eq(status, pending),
  ),
  And(
    Or(
      Eq(priority, high),
      Eq(priority, urgent),
    ),
    And(
      Contains(customer name, "acme"),
      And(
        Gte(score, 70),
        And(
          Not(True(archived)),
          InLast(created, 30 days),
        ),
      ),
    ),
  ),
)
```

### Compiled query intent

* grouped `or` status clause
* grouped `or` priority clause
* `contains` on customer name
* `gte` compare on score
* `not archived`
* top-level `TimeRange` for the last 30 days
* top-level sort with two fields
* limit 50
* skip 100

---

## Open Decisions

These do not need to be solved before implementation starts, but they should be tracked.

### 1. Bare field semantics

Should `archived` mean `archived is true`?

Recommended answer: yes.

### 2. `not archived`

Should this compile to `archived is false` or `not (archived is true)`?

Recommended answer: compile to `not (archived is true)` unless your field typing guarantees a strict boolean.

### 3. Multi-word unquoted values

Should `status is in progress` require quotes?

Recommended answer for v1: yes, require quotes for multi-word values.

Example:

```text
status is "in progress"
```

### 4. Field suggestions

Should the parser suggest close field aliases?

Recommended answer: yes, after the first working implementation.

### 5. Time zone behavior

How should `today`, `this week`, and `this month` be interpreted?

Recommended answer: use the application's configured local time zone consistently.

---

## Testing Recommendations

### Unit tests

Write tests for:

* tokenization
* field parsing
* condition parsing
* boolean grouping
* parentheses handling
* sort parsing
* limit parsing
* skip parsing
* time parsing
* AST compilation
* validation failures
* error message quality

### Table-driven tests

Use table-driven tests heavily.

Examples:

* input string
* expected AST
* expected compilation result
* expected error string

### Golden cases

Create a small set of representative full-query examples as long-term regression tests.

---

## Final Recommendation

Keep the first implementation boring in the best possible way.

A strict, explicit, friendly parser that works predictably is much more valuable than a clever parser that guesses.

The right v1 foundation is:

* friendly phrasing
* explicit grouping
* small AST
* clear validation
* clean compiler boundary

Once that is stable, you can safely make the language feel more magical later.

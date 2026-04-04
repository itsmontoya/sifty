// Package matcher compiles query filters into executable matchers.
//
// Behavior specification:
//
// Nil/zero filter behavior:
//
// Compile treats a zero-value query filter (query.Clause.IsZero() == true) as
// a match-all filter. The returned matcher always reports a match unless
// evaluation returns an error from the document view.
//
// Null and missing fields:
//
// A field is considered "missing" when DocView.Get(path) returns ok=false and
// err=nil.
//
// A field may be present with a null value (Go nil) when DocView.Get(path)
// returns ok=true, err=nil, value=nil. This is not treated as missing.
//
// Missing-field semantics:
//
// A field is considered "missing" when DocView.Get(path) returns ok=false and
// err=nil.
//
// Operator behavior for missing fields:
//   - contains: returns false, nil.
//   - compare.eq: returns false, nil.
//   - compare.gt: returns false, nil.
//   - compare.gte: returns false, nil.
//   - compare.lt: returns false, nil.
//   - compare.lte: returns false, nil.
//   - and: the missing-field child is a non-match; evaluation short-circuits on
//     first false.
//   - or: the missing-field child is a non-match; evaluation continues until a
//     true child is found, otherwise false.
//   - not: negates the child result, so a missing-field child non-match becomes
//     true.
//
// Type-mismatch semantics:
//
// A type mismatch is evaluated as a non-match (false, nil), not an error.
//
// Operator behavior for type mismatches:
//   - contains: the document value must be string; non-string values are type
//     mismatches and return false, nil.
//   - compare.eq / compare.gt / compare.gte / compare.lt / compare.lte:
//     comparison requires exact runtime type equality between document value and
//     comparison operand (for example int vs int64 is a mismatch). Mismatches
//     return false, nil.
//   - string and numeric values are never coerced. For example "10" vs 10 and
//     float64(10) vs int(10) are mismatches and return false, nil.
//   - compare with unsupported document-value type (for example bool, map,
//     slice, nil): treated as non-match and returns false, nil.
//
// Errors versus non-match:
//
// Evaluation errors only come from DocView.Get returning err != nil. Those
// errors are propagated directly by IsMatch/eval and are never converted into
// false.
//
// Non-matches are represented as (false, nil), including missing fields, type
// mismatches, unsupported compare types, and clause/operator checks that do not
// satisfy the filter.
package matcher

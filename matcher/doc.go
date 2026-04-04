// Package matcher compiles query filters into executable matchers.
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
//   - compare with unsupported document-value type (for example bool, map,
//     slice, nil): treated as non-match and returns false, nil.
package matcher

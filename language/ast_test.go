package language

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestToAST_FilterExpressions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantTree string
		wantErr  bool
	}{
		{
			name:     "single equality condition",
			input:    "status is active",
			wantTree: "cond(status,eq,active)",
		},
		{
			name:     "single inequality condition using is not",
			input:    "status is not active",
			wantTree: "cond(status,neq,active)",
		},
		{
			name:     "contains condition",
			input:    `customer contains "acme corp"`,
			wantTree: `cond(customer,contains,"acme corp")`,
		},
		{
			name:     "not contains condition",
			input:    `customer does not contain "acme corp"`,
			wantTree: `cond(customer,not_contains,"acme corp")`,
		},
		{
			name:     "precedence not over and over or",
			input:    "a is x or b is y and not c is z",
			wantTree: "or(cond(a,eq,x),and(cond(b,eq,y),not(cond(c,eq,z))))",
		},
		{
			name:     "grouping with parentheses",
			input:    "(a is x or b is y) and c is z",
			wantTree: "and(or(cond(a,eq,x),cond(b,eq,y)),cond(c,eq,z))",
		},
		{
			name:     "boolean true condition normalizes op",
			input:    "enabled is true",
			wantTree: "cond(enabled,is_true,nil)",
		},
		{
			name:     "boolean false condition normalizes op",
			input:    "enabled is false",
			wantTree: "cond(enabled,is_false,nil)",
		},
		{
			name:     "temporal today condition",
			input:    "created today",
			wantTree: "cond(created,today,nil)",
		},
		{
			name:     "temporal this week condition",
			input:    "created this week",
			wantTree: "cond(created,this_week,nil)",
		},
		{
			name:     "temporal this month condition",
			input:    "created this month",
			wantTree: "cond(created,this_month,nil)",
		},
		{
			name:     "temporal in the last days condition",
			input:    "created in the last 7 days",
			wantTree: "cond(created,in_last,7)",
		},
		{
			name:     "comparison greater than condition",
			input:    "score greater than 10",
			wantTree: "cond(score,gt,10)",
		},
		{
			name:     "comparison less than condition",
			input:    "score less than 10",
			wantTree: "cond(score,lt,10)",
		},
		{
			name:     "comparison at least condition",
			input:    "score at least 10",
			wantTree: "cond(score,gte,10)",
		},
		{
			name:     "comparison at most condition",
			input:    "score at most 10",
			wantTree: "cond(score,lte,10)",
		},
		{
			name:    "unmatched left parenthesis returns parse error",
			input:   "(status is active",
			wantErr: true,
		},
		{
			name:    "trailing operator returns parse error",
			input:   "status is active and",
			wantErr: true,
		},
		{
			name:    "greater missing than returns parse error",
			input:   "score greater 10",
			wantErr: true,
		},
		{
			name:    "less missing than returns parse error",
			input:   "score less 10",
			wantErr: true,
		},
		{
			name:    "at missing least or most returns parse error",
			input:   "score at 10",
			wantErr: true,
		},
		{
			name:    "at least missing value returns parse error",
			input:   "score at least",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens, err := ToTokens(tt.input)
			if err != nil {
				t.Fatalf("ToTokens() failed: %v", err)
			}

			got, err := ToAST(tokens)
			if err != nil {
				if !tt.wantErr {
					t.Fatalf("ToAST() failed: %v", err)
				}
				return
			}

			if tt.wantErr {
				t.Fatal("ToAST() succeeded unexpectedly")
			}

			gotTree := exprTreeString(got.Filter)
			if gotTree != tt.wantTree {
				t.Fatalf("ToAST() filter tree = %s, want %s", gotTree, tt.wantTree)
			}
		})
	}
}

func TestToAST_EmptyInputReturnsZeroValue(t *testing.T) {
	got, err := ToAST([]Token{})
	if err != nil {
		t.Fatalf("ToAST() failed: %v", err)
	}

	if !reflect.DeepEqual(got, AST{}) {
		t.Fatalf("ToAST() = %#v, want zero value %#v", got, AST{})
	}
}

func TestToAST_LimitAndSkipClauses(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantLimit *int
		wantSkip  *int
		wantErr   bool
	}{
		{
			name:      "limit only",
			input:     "status is active limit 50",
			wantLimit: intPtr(50),
			wantSkip:  nil,
		},
		{
			name:      "skip only",
			input:     "status is active skip 100",
			wantLimit: nil,
			wantSkip:  intPtr(100),
		},
		{
			name:      "limit and skip",
			input:     "status is active limit 50 skip 100",
			wantLimit: intPtr(50),
			wantSkip:  intPtr(100),
		},
		{
			name:    "limit missing number returns parse error",
			input:   "status is active limit",
			wantErr: true,
		},
		{
			name:    "skip missing number returns parse error",
			input:   "status is active skip",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens, err := ToTokens(tt.input)
			if err != nil {
				t.Fatalf("ToTokens() failed: %v", err)
			}

			got, err := ToAST(tokens)
			if err != nil {
				if !tt.wantErr {
					t.Fatalf("ToAST() failed: %v", err)
				}

				return
			}

			if tt.wantErr {
				t.Fatal("ToAST() succeeded unexpectedly")
			}

			if !reflect.DeepEqual(got.Limit, tt.wantLimit) {
				t.Fatalf("ToAST().Limit = %#v, want %#v", got.Limit, tt.wantLimit)
			}

			if !reflect.DeepEqual(got.Skip, tt.wantSkip) {
				t.Fatalf("ToAST().Skip = %#v, want %#v", got.Skip, tt.wantSkip)
			}
		})
	}
}

func TestToAST_SortClauses(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantSort []SortExpr
		wantErr  bool
	}{
		{
			name:  "single sort field defaults to ascending",
			input: "status is active sorted by created",
			wantSort: []SortExpr{
				{Field: "created", Desc: false},
			},
		},
		{
			name:  "single sort field descending",
			input: "status is active sorted by created descending",
			wantSort: []SortExpr{
				{Field: "created", Desc: true},
			},
		},
		{
			name:  "multiple sort fields mixed direction",
			input: "status is active sorted by created, score descending",
			wantSort: []SortExpr{
				{Field: "created", Desc: false},
				{Field: "score", Desc: true},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens, err := ToTokens(tt.input)
			if err != nil {
				t.Fatalf("ToTokens() failed: %v", err)
			}

			got, err := ToAST(tokens)
			if err != nil {
				if !tt.wantErr {
					t.Fatalf("ToAST() failed: %v", err)
				}

				return
			}

			if tt.wantErr {
				t.Fatal("ToAST() succeeded unexpectedly")
			}

			if !reflect.DeepEqual(got.Sort, tt.wantSort) {
				t.Fatalf("ToAST().Sort = %#v, want %#v", got.Sort, tt.wantSort)
			}
		})
	}
}

func TestToAST_CompleteQueryCombinations(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantTree  string
		wantSort  []SortExpr
		wantLimit *int
		wantSkip  *int
	}{
		{
			name:      "simple filter sort limit skip",
			input:     "status is active sorted by created limit 50 skip 100",
			wantTree:  "cond(status,eq,active)",
			wantSort:  []SortExpr{{Field: "created", Desc: false}},
			wantLimit: intPtr(50),
			wantSkip:  intPtr(100),
		},
		{
			name:      "grouped filter with multiple sorts and paging",
			input:     "(status is active or status is pending) sorted by created descending, score limit 25 skip 10",
			wantTree:  "or(cond(status,eq,active),cond(status,eq,pending))",
			wantSort:  []SortExpr{{Field: "created", Desc: true}, {Field: "score", Desc: false}},
			wantLimit: intPtr(25),
			wantSkip:  intPtr(10),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens, err := ToTokens(tt.input)
			if err != nil {
				t.Fatalf("ToTokens() failed: %v", err)
			}

			got, err := ToAST(tokens)
			if err != nil {
				t.Fatalf("ToAST() failed: %v", err)
			}

			if exprTreeString(got.Filter) != tt.wantTree {
				t.Fatalf("ToAST().Filter = %s, want %s", exprTreeString(got.Filter), tt.wantTree)
			}

			if !reflect.DeepEqual(got.Sort, tt.wantSort) {
				t.Fatalf("ToAST().Sort = %#v, want %#v", got.Sort, tt.wantSort)
			}

			if !reflect.DeepEqual(got.Limit, tt.wantLimit) {
				t.Fatalf("ToAST().Limit = %#v, want %#v", got.Limit, tt.wantLimit)
			}

			if !reflect.DeepEqual(got.Skip, tt.wantSkip) {
				t.Fatalf("ToAST().Skip = %#v, want %#v", got.Skip, tt.wantSkip)
			}
		})
	}
}

func TestToAST_SyntaxErrorsWithPosition(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantMessage string
		wantKind    Kind
		wantLexeme  string
		wantPos     Position
	}{
		{
			name:        "sorted missing by",
			input:       "status is active sorted created",
			wantMessage: "expected 'by' after 'sorted'",
			wantKind:    KindIdentifier,
			wantLexeme:  "created",
			wantPos:     Position{Offset: 24, Line: 1, Column: 25},
		},
		{
			name:        "sort trailing comma",
			input:       "status is active sorted by created,",
			wantMessage: "expected sort field after ','",
			wantKind:    KindEOF,
			wantLexeme:  "",
			wantPos:     Position{Offset: 35, Line: 1, Column: 36},
		},
		{
			name:        "limit missing number",
			input:       "status is active limit",
			wantMessage: "expected number after 'limit'",
			wantKind:    KindEOF,
			wantLexeme:  "",
			wantPos:     Position{Offset: 22, Line: 1, Column: 23},
		},
		{
			name:        "skip missing number",
			input:       "status is active skip",
			wantMessage: "expected number after 'skip'",
			wantKind:    KindEOF,
			wantLexeme:  "",
			wantPos:     Position{Offset: 21, Line: 1, Column: 22},
		},
		{
			name:        "unexpected trailing token",
			input:       "status is active nonsense",
			wantMessage: "unexpected trailing tokens",
			wantKind:    KindIdentifier,
			wantLexeme:  "nonsense",
			wantPos:     Position{Offset: 17, Line: 1, Column: 18},
		},
		{
			name:        "limit numeric overflow returns parse error",
			input:       "status is active limit 999999999999999999999",
			wantMessage: `invalid number "999999999999999999999": strconv.Atoi: parsing "999999999999999999999": value out of range`,
			wantKind:    KindNumber,
			wantLexeme:  "999999999999999999999",
			wantPos:     Position{Offset: 23, Line: 1, Column: 24},
		},
		{
			name:        "skip numeric overflow returns parse error",
			input:       "status is active skip 999999999999999999999",
			wantMessage: `invalid number "999999999999999999999": strconv.Atoi: parsing "999999999999999999999": value out of range`,
			wantKind:    KindNumber,
			wantLexeme:  "999999999999999999999",
			wantPos:     Position{Offset: 22, Line: 1, Column: 23},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens, err := ToTokens(tt.input)
			if err != nil {
				t.Fatalf("ToTokens() failed: %v", err)
			}

			_, err = ToAST(tokens)
			if err == nil {
				t.Fatal("ToAST() succeeded unexpectedly")
			}

			parseErr, ok := err.(parseError)
			if !ok {
				t.Fatalf("error type = %T, want %T", err, parseError{})
			}

			if parseErr.Message != tt.wantMessage {
				t.Fatalf("Message = %q, want %q", parseErr.Message, tt.wantMessage)
			}

			if parseErr.Kind != tt.wantKind {
				t.Fatalf("Kind = %v, want %v", parseErr.Kind, tt.wantKind)
			}

			if parseErr.Lexeme != tt.wantLexeme {
				t.Fatalf("Lexeme = %q, want %q", parseErr.Lexeme, tt.wantLexeme)
			}

			if parseErr.Position != tt.wantPos {
				t.Fatalf("Position = %+v, want %+v", parseErr.Position, tt.wantPos)
			}
		})
	}
}

func exprTreeString(expr any) string {
	if expr == nil {
		return "<nil>"
	}

	value := reflect.ValueOf(expr)
	typ := value.Type()
	if typ.Kind() == reflect.Pointer {
		if value.IsNil() {
			return "<nil>"
		}
		value = value.Elem()
		typ = value.Type()
	}

	switch typ.Name() {
	case "AndExpr":
		return fmt.Sprintf(
			"and(%s,%s)",
			exprTreeString(value.FieldByName("Left").Interface()),
			exprTreeString(value.FieldByName("Right").Interface()),
		)
	case "OrExpr":
		return fmt.Sprintf(
			"or(%s,%s)",
			exprTreeString(value.FieldByName("Left").Interface()),
			exprTreeString(value.FieldByName("Right").Interface()),
		)
	case "NotExpr":
		return fmt.Sprintf("not(%s)", exprTreeString(value.FieldByName("Inner").Interface()))
	case "ConditionExpr":
		return fmt.Sprintf(
			"cond(%s,%s,%s)",
			value.FieldByName("Field").Interface(),
			value.FieldByName("Op").Interface(),
			conditionValueString(value.FieldByName("Value").Interface()),
		)
	default:
		return fmt.Sprintf("<unknown:%T>", expr)
	}
}

func conditionValueString(value any) string {
	switch v := value.(type) {
	case nil:
		return "nil"
	case string:
		if strings.Contains(v, " ") {
			return fmt.Sprintf("%q", v)
		}
		return v
	default:
		return fmt.Sprintf("%v", v)
	}
}

func intPtr(i int) *int {
	return &i
}

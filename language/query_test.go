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
			name:    "unmatched left parenthesis returns parse error",
			input:   "(status is active",
			wantErr: true,
		},
		{
			name:    "trailing operator returns parse error",
			input:   "status is active and",
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

	if !reflect.DeepEqual(got, ASTQuery{}) {
		t.Fatalf("ToAST() = %#v, want zero value %#v", got, ASTQuery{})
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

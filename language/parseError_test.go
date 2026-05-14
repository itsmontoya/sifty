package language

import "testing"

func TestMakeParseError(t *testing.T) {
	var (
		at  Token
		got parseError
	)

	at.Kind = KindIdentifier
	at.Lexeme = "status"
	at.Position = Position{
		Offset: 4,
		Line:   1,
		Column: 5,
	}

	got = makeParseError("expected condition operator after field", at)

	if got.Message != "expected condition operator after field" {
		t.Fatalf("Message = %q, want %q", got.Message, "expected condition operator after field")
	}

	if got.Position != at.Position {
		t.Fatalf("Position = %+v, want %+v", got.Position, at.Position)
	}

	if got.Lexeme != at.Lexeme {
		t.Fatalf("Lexeme = %q, want %q", got.Lexeme, at.Lexeme)
	}

	if got.Kind != at.Kind {
		t.Fatalf("Kind = %v, want %v", got.Kind, at.Kind)
	}
}

func TestParseErrorError(t *testing.T) {
	tests := []struct {
		name string
		err  parseError
		want string
	}{
		{
			name: "with lexeme",
			err: parseError{
				Message: "expected field name",
				Position: Position{
					Offset: 2,
					Line:   1,
					Column: 3,
				},
				Lexeme: "(",
				Kind:   KindLParen,
			},
			want: `expected field name at line 1, column 3 (offset 2): "("`,
		},
		{
			name: "without lexeme",
			err: parseError{
				Message: "unexpected trailing tokens",
				Position: Position{
					Offset: 8,
					Line:   1,
					Column: 9,
				},
			},
			want: "unexpected trailing tokens at line 1, column 9 (offset 8)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.err.Error()
			if got != tt.want {
				t.Fatalf("Error() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestToASTParseError(t *testing.T) {
	var (
		tokens []Token
		err    error
		got    parseError
		ok     bool
	)

	tokens, err = ToTokens("status")
	if err != nil {
		t.Fatalf("ToTokens() failed: %v", err)
	}

	_, err = ToAST(tokens)
	if err == nil {
		t.Fatal("ToAST() succeeded unexpectedly")
	}

	got, ok = err.(parseError)
	if !ok {
		t.Fatalf("error type = %T, want %T", err, parseError{})
	}

	if got.Message != "expected condition operator after field" {
		t.Fatalf("Message = %q, want %q", got.Message, "expected condition operator after field")
	}

	if got.Position != (Position{Offset: 6, Line: 1, Column: 7}) {
		t.Fatalf("Position = %+v, want %+v", got.Position, Position{Offset: 6, Line: 1, Column: 7})
	}

	if got.Lexeme != "" {
		t.Fatalf("Lexeme = %q, want empty string", got.Lexeme)
	}

	if got.Kind != KindEOF {
		t.Fatalf("Kind = %v, want %v", got.Kind, KindEOF)
	}
}

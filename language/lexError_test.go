package language

import "testing"

func TestMakeLexError(t *testing.T) {
	var (
		position Position
		got      lexError
	)

	position.Offset = 9
	position.Line = 2
	position.Column = 4

	got = makeLexError("invalid character", "@", position)

	if got.Message != "invalid character" {
		t.Fatalf("Message = %q, want %q", got.Message, "invalid character")
	}

	if got.Offender != "@" {
		t.Fatalf("Offender = %q, want %q", got.Offender, "@")
	}

	if got.Position != position {
		t.Fatalf("Position = %+v, want %+v", got.Position, position)
	}
}

func TestLexErrorError(t *testing.T) {
	tests := []struct {
		name string
		err  lexError
		want string
	}{
		{
			name: "with offender",
			err: lexError{
				Message:  "invalid character",
				Offender: "@",
				Position: Position{
					Offset: 7,
					Line:   1,
					Column: 8,
				},
			},
			want: "invalid character at line 1, column 8 (offset 7): \"@\"",
		},
		{
			name: "without offender",
			err: lexError{
				Message: "unterminated string",
				Position: Position{
					Offset: 17,
					Line:   1,
					Column: 18,
				},
			},
			want: "unterminated string at line 1, column 18 (offset 17)",
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

func TestToTokensLexError(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		wantMessage  string
		wantOffender string
		wantPos      Position
	}{
		{
			name:         "invalid character includes offender and position",
			input:        "status @ active",
			wantMessage:  "invalid character",
			wantOffender: "@",
			wantPos: Position{
				Offset: 7,
				Line:   1,
				Column: 8,
			},
		},
		{
			name:         "unterminated string has empty offender",
			input:        `status is "active`,
			wantMessage:  "unterminated string",
			wantOffender: "",
			wantPos: Position{
				Offset: 17,
				Line:   1,
				Column: 18,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				err error
				got lexError
				ok  bool
			)

			_, err = ToTokens(tt.input)
			if err == nil {
				t.Fatal("ToTokens() succeeded unexpectedly")
			}

			got, ok = err.(lexError)
			if !ok {
				t.Fatalf("error type = %T, want %T", err, lexError{})
			}

			if got.Message != tt.wantMessage {
				t.Fatalf("Message = %q, want %q", got.Message, tt.wantMessage)
			}

			if got.Offender != tt.wantOffender {
				t.Fatalf("Offender = %q, want %q", got.Offender, tt.wantOffender)
			}

			if got.Position != tt.wantPos {
				t.Fatalf("Position = %+v, want %+v", got.Position, tt.wantPos)
			}
		})
	}
}

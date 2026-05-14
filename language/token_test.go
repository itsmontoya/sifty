package language_test

import (
	"reflect"
	"testing"

	"github.com/itsmontoya/sifty/language"
)

func TestToTokens(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []language.Token
		wantErr bool
	}{
		{
			name:  "empty input emits eof",
			input: "",
			want: []language.Token{
				{
					Kind:   language.KindEOF,
					Lexeme: "",
					Position: language.Position{
						Offset: 0,
						Line:   1,
						Column: 1,
					},
				},
			},
		},
		{
			name:  "whitespace only emits eof at final position",
			input: " \t\n",
			want: []language.Token{
				{
					Kind:   language.KindEOF,
					Lexeme: "",
					Position: language.Position{
						Offset: 3,
						Line:   2,
						Column: 1,
					},
				},
			},
		},
		{
			name:  "newline position tracking single newline",
			input: "a\nb",
			want: []language.Token{
				{
					Kind:   language.KindIdentifier,
					Lexeme: "a",
					Position: language.Position{
						Offset: 0,
						Line:   1,
						Column: 1,
					},
				},
				{
					Kind:   language.KindIdentifier,
					Lexeme: "b",
					Position: language.Position{
						Offset: 2,
						Line:   2,
						Column: 1,
					},
				},
				{
					Kind:   language.KindEOF,
					Lexeme: "",
					Position: language.Position{
						Offset: 3,
						Line:   2,
						Column: 2,
					},
				},
			},
		},
		{
			name:  "newline position tracking consecutive newlines and space",
			input: "a\n\n b",
			want: []language.Token{
				{
					Kind:   language.KindIdentifier,
					Lexeme: "a",
					Position: language.Position{
						Offset: 0,
						Line:   1,
						Column: 1,
					},
				},
				{
					Kind:   language.KindIdentifier,
					Lexeme: "b",
					Position: language.Position{
						Offset: 4,
						Line:   3,
						Column: 2,
					},
				},
				{
					Kind:   language.KindEOF,
					Lexeme: "",
					Position: language.Position{
						Offset: 5,
						Line:   3,
						Column: 3,
					},
				},
			},
		},
		{
			name:  "mixed whitespace tracking tab newline prefix",
			input: "\t\na",
			want: []language.Token{
				{
					Kind:   language.KindIdentifier,
					Lexeme: "a",
					Position: language.Position{
						Offset: 2,
						Line:   2,
						Column: 1,
					},
				},
				{
					Kind:   language.KindEOF,
					Lexeme: "",
					Position: language.Position{
						Offset: 3,
						Line:   2,
						Column: 2,
					},
				},
			},
		},
		{
			name:  "mixed whitespace tracking space tab newline tab",
			input: " \t\n\tb",
			want: []language.Token{
				{
					Kind:   language.KindIdentifier,
					Lexeme: "b",
					Position: language.Position{
						Offset: 4,
						Line:   2,
						Column: 2,
					},
				},
				{
					Kind:   language.KindEOF,
					Lexeme: "",
					Position: language.Position{
						Offset: 5,
						Line:   2,
						Column: 3,
					},
				},
			},
		},
		{
			name:  "punctuation tokens",
			input: "(,)",
			want: []language.Token{
				{
					Kind:   language.KindLParen,
					Lexeme: "(",
					Position: language.Position{
						Offset: 0,
						Line:   1,
						Column: 1,
					},
				},
				{
					Kind:   language.KindComma,
					Lexeme: ",",
					Position: language.Position{
						Offset: 1,
						Line:   1,
						Column: 2,
					},
				},
				{
					Kind:   language.KindRParen,
					Lexeme: ")",
					Position: language.Position{
						Offset: 2,
						Line:   1,
						Column: 3,
					},
				},
				{
					Kind:   language.KindEOF,
					Lexeme: "",
					Position: language.Position{
						Offset: 3,
						Line:   1,
						Column: 4,
					},
				},
			},
		},
		{
			name:  "keywords and identifiers",
			input: "status is active and priority is high",
			want: []language.Token{
				{
					Kind:   language.KindIdentifier,
					Lexeme: "status",
					Position: language.Position{
						Offset: 0,
						Line:   1,
						Column: 1,
					},
				},
				{
					Kind:   language.KindIs,
					Lexeme: "is",
					Position: language.Position{
						Offset: 7,
						Line:   1,
						Column: 8,
					},
				},
				{
					Kind:   language.KindIdentifier,
					Lexeme: "active",
					Position: language.Position{
						Offset: 10,
						Line:   1,
						Column: 11,
					},
				},
				{
					Kind:   language.KindAnd,
					Lexeme: "and",
					Position: language.Position{
						Offset: 17,
						Line:   1,
						Column: 18,
					},
				},
				{
					Kind:   language.KindIdentifier,
					Lexeme: "priority",
					Position: language.Position{
						Offset: 21,
						Line:   1,
						Column: 22,
					},
				},
				{
					Kind:   language.KindIs,
					Lexeme: "is",
					Position: language.Position{
						Offset: 30,
						Line:   1,
						Column: 31,
					},
				},
				{
					Kind:   language.KindIdentifier,
					Lexeme: "high",
					Position: language.Position{
						Offset: 33,
						Line:   1,
						Column: 34,
					},
				},
				{
					Kind:   language.KindEOF,
					Lexeme: "",
					Position: language.Position{
						Offset: 37,
						Line:   1,
						Column: 38,
					},
				},
			},
		},
		{
			name:  "keywords are case insensitive",
			input: "STATUS IS active",
			want: []language.Token{
				{
					Kind:   language.KindIdentifier,
					Lexeme: "STATUS",
					Position: language.Position{
						Offset: 0,
						Line:   1,
						Column: 1,
					},
				},
				{
					Kind:   language.KindIs,
					Lexeme: "IS",
					Position: language.Position{
						Offset: 7,
						Line:   1,
						Column: 8,
					},
				},
				{
					Kind:   language.KindIdentifier,
					Lexeme: "active",
					Position: language.Position{
						Offset: 10,
						Line:   1,
						Column: 11,
					},
				},
				{
					Kind:   language.KindEOF,
					Lexeme: "",
					Position: language.Position{
						Offset: 16,
						Line:   1,
						Column: 17,
					},
				},
			},
		},
		{
			name:  "offsets are rune based for unicode input",
			input: `"é" is`,
			want: []language.Token{
				{
					Kind:   language.KindString,
					Lexeme: "é",
					Position: language.Position{
						Offset: 0,
						Line:   1,
						Column: 1,
					},
				},
				{
					Kind:   language.KindIs,
					Lexeme: "is",
					Position: language.Position{
						Offset: 4,
						Line:   1,
						Column: 5,
					},
				},
				{
					Kind:   language.KindEOF,
					Lexeme: "",
					Position: language.Position{
						Offset: 6,
						Line:   1,
						Column: 7,
					},
				},
			},
		},
		{
			name:  "number token",
			input: "limit 50 skip 100",
			want: []language.Token{
				{
					Kind:   language.KindLimit,
					Lexeme: "limit",
					Position: language.Position{
						Offset: 0,
						Line:   1,
						Column: 1,
					},
				},
				{
					Kind:   language.KindNumber,
					Lexeme: "50",
					Position: language.Position{
						Offset: 6,
						Line:   1,
						Column: 7,
					},
				},
				{
					Kind:   language.KindSkip,
					Lexeme: "skip",
					Position: language.Position{
						Offset: 9,
						Line:   1,
						Column: 10,
					},
				},
				{
					Kind:   language.KindNumber,
					Lexeme: "100",
					Position: language.Position{
						Offset: 14,
						Line:   1,
						Column: 15,
					},
				},
				{
					Kind:   language.KindEOF,
					Lexeme: "",
					Position: language.Position{
						Offset: 17,
						Line:   1,
						Column: 18,
					},
				},
			},
		},
		{
			name:  "quoted string token",
			input: `customer contains "acme corp"`,
			want: []language.Token{
				{
					Kind:   language.KindIdentifier,
					Lexeme: "customer",
					Position: language.Position{
						Offset: 0,
						Line:   1,
						Column: 1,
					},
				},
				{
					Kind:   language.KindContains,
					Lexeme: "contains",
					Position: language.Position{
						Offset: 9,
						Line:   1,
						Column: 10,
					},
				},
				{
					Kind:   language.KindString,
					Lexeme: "acme corp",
					Position: language.Position{
						Offset: 18,
						Line:   1,
						Column: 19,
					},
				},
				{
					Kind:   language.KindEOF,
					Lexeme: "",
					Position: language.Position{
						Offset: 29,
						Line:   1,
						Column: 30,
					},
				},
			},
		},
		{
			name:  "boolean and logical keywords",
			input: "true or not false",
			want: []language.Token{
				{
					Kind:   language.KindTrue,
					Lexeme: "true",
					Position: language.Position{
						Offset: 0,
						Line:   1,
						Column: 1,
					},
				},
				{
					Kind:   language.KindOr,
					Lexeme: "or",
					Position: language.Position{
						Offset: 5,
						Line:   1,
						Column: 6,
					},
				},
				{
					Kind:   language.KindNot,
					Lexeme: "not",
					Position: language.Position{
						Offset: 8,
						Line:   1,
						Column: 9,
					},
				},
				{
					Kind:   language.KindFalse,
					Lexeme: "false",
					Position: language.Position{
						Offset: 12,
						Line:   1,
						Column: 13,
					},
				},
				{
					Kind:   language.KindEOF,
					Lexeme: "",
					Position: language.Position{
						Offset: 17,
						Line:   1,
						Column: 18,
					},
				},
			},
		},
		{
			name:  "boolean and logical keywords are case insensitive",
			input: "TRUE Or nOt FALSE",
			want: []language.Token{
				{
					Kind:   language.KindTrue,
					Lexeme: "TRUE",
					Position: language.Position{
						Offset: 0,
						Line:   1,
						Column: 1,
					},
				},
				{
					Kind:   language.KindOr,
					Lexeme: "Or",
					Position: language.Position{
						Offset: 5,
						Line:   1,
						Column: 6,
					},
				},
				{
					Kind:   language.KindNot,
					Lexeme: "nOt",
					Position: language.Position{
						Offset: 8,
						Line:   1,
						Column: 9,
					},
				},
				{
					Kind:   language.KindFalse,
					Lexeme: "FALSE",
					Position: language.Position{
						Offset: 12,
						Line:   1,
						Column: 13,
					},
				},
				{
					Kind:   language.KindEOF,
					Lexeme: "",
					Position: language.Position{
						Offset: 17,
						Line:   1,
						Column: 18,
					},
				},
			},
		},
		{
			name:    "invalid character returns error",
			input:   "status @ active",
			wantErr: true,
		},
		{
			name:    "unterminated string returns error",
			input:   `status is "active`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := language.ToTokens(tt.input)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("ToTokens() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("ToTokens() succeeded unexpectedly")
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToTokens() = %v, want %v", got, tt.want)
			}
		})
	}
}

package language

import "testing"

func TestIsNumber(t *testing.T) {
	tests := []struct {
		name  string
		input []rune
		want  bool
	}{
		{
			name:  "empty input",
			input: []rune{},
			want:  true,
		},
		{
			name:  "single digit",
			input: []rune{'7'},
			want:  true,
		},
		{
			name:  "multiple digits",
			input: []rune{'1', '2', '3', '4', '5'},
			want:  true,
		},
		{
			name:  "letters are not numbers",
			input: []rune{'a'},
			want:  false,
		},
		{
			name:  "mixed alphanumeric",
			input: []rune{'1', 'a'},
			want:  false,
		},
		{
			name:  "whitespace is not a number",
			input: []rune{' '},
			want:  false,
		},
		{
			name:  "symbol is not a number",
			input: []rune{'-'},
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isNumber(tt.input...)
			if got != tt.want {
				t.Fatalf("isNumber(%q) = %t, want %t", string(tt.input), got, tt.want)
			}
		})
	}
}

func TestIsWordPart(t *testing.T) {
	tests := []struct {
		name  string
		input rune
		want  bool
	}{
		{
			name:  "lowercase letter",
			input: 'a',
			want:  true,
		},
		{
			name:  "uppercase letter",
			input: 'Z',
			want:  true,
		},
		{
			name:  "digit",
			input: '9',
			want:  true,
		},
		{
			name:  "dot",
			input: '.',
			want:  true,
		},
		{
			name:  "left bracket",
			input: '[',
			want:  true,
		},
		{
			name:  "right bracket",
			input: ']',
			want:  true,
		},
		{
			name:  "dash is not allowed",
			input: '-',
			want:  false,
		},
		{
			name:  "underscore is not allowed",
			input: '_',
			want:  false,
		},
		{
			name:  "space is not allowed",
			input: ' ',
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isWordPart(tt.input)
			if got != tt.want {
				t.Fatalf("isWordPart(%q) = %t, want %t", tt.input, got, tt.want)
			}
		})
	}
}

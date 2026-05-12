package language

import "testing"

func TestGetKind_MissingKeywordMappings(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want Kind
	}{
		{
			name: "contains keyword contain",
			in:   "contain",
			want: KindContain,
		},
		{
			name: "contains keyword contains",
			in:   "contains",
			want: KindContains,
		},
		{
			name: "comparison keyword greater",
			in:   "greater",
			want: KindGreater,
		},
		{
			name: "comparison keyword less",
			in:   "less",
			want: KindLess,
		},
		{
			name: "comparison keyword than",
			in:   "than",
			want: KindThan,
		},
		{
			name: "comparison keyword at",
			in:   "at",
			want: KindAt,
		},
		{
			name: "comparison keyword least",
			in:   "least",
			want: KindLeast,
		},
		{
			name: "comparison keyword most",
			in:   "most",
			want: KindMost,
		},
		{
			name: "condition keyword does",
			in:   "does",
			want: KindDoes,
		},
		{
			name: "condition keyword before",
			in:   "before",
			want: KindBefore,
		},
		{
			name: "condition keyword after",
			in:   "after",
			want: KindAfter,
		},
		{
			name: "condition keyword in",
			in:   "in",
			want: KindIn,
		},
		{
			name: "condition keyword the",
			in:   "the",
			want: KindThe,
		},
		{
			name: "condition keyword last",
			in:   "last",
			want: KindLast,
		},
		{
			name: "condition keyword today",
			in:   "today",
			want: KindToday,
		},
		{
			name: "condition keyword this",
			in:   "this",
			want: KindThis,
		},
		{
			name: "condition keyword week",
			in:   "week",
			want: KindWeek,
		},
		{
			name: "condition keyword month",
			in:   "month",
			want: KindMonth,
		},
		{
			name: "condition keyword days",
			in:   "days",
			want: KindDays,
		},
		{
			name: "sort keyword sorted",
			in:   "sorted",
			want: KindSorted,
		},
		{
			name: "sort keyword by",
			in:   "by",
			want: KindBy,
		},
		{
			name: "sort keyword ascending",
			in:   "ascending",
			want: KindAscending,
		},
		{
			name: "sort keyword descending",
			in:   "descending",
			want: KindDescending,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getKind(tt.in)
			if got != tt.want {
				t.Fatalf("getKind(%q) = %v, want %v", tt.in, got, tt.want)
			}
		})
	}
}

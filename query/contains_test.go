package query

import "testing"

func TestContainsExprValidate(t *testing.T) {
	tt := []struct {
		name    string
		in      ContainsExpr
		wantErr bool
	}{
		{
			name: "valid",
			in: ContainsExpr{
				Field: "title",
				Value: "golang",
			},
		},
		{
			name: "missing field",
			in: ContainsExpr{
				Value: "golang",
			},
			wantErr: true,
		},
		{
			name: "missing value",
			in: ContainsExpr{
				Field: "title",
			},
			wantErr: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.in.Validate()
			if tc.wantErr && err == nil {
				t.Fatal("expected validation error")
			}

			if !tc.wantErr && err != nil {
				t.Fatalf("unexpected validation error: %v", err)
			}
		})
	}
}

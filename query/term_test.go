package query

import "testing"

func TestTermExprValidate(t *testing.T) {
	tt := []struct {
		name    string
		in      TermExpr
		wantErr bool
	}{
		{
			name: "valid",
			in: TermExpr{
				Field: "status",
				Value: "active",
			},
		},
		{
			name: "missing field",
			in: TermExpr{
				Value: "active",
			},
			wantErr: true,
		},
		{
			name: "missing value",
			in: TermExpr{
				Field: "status",
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

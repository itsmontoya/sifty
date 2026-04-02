package query

import "testing"

func TestCompareExprValidate(t *testing.T) {
	tt := []struct {
		name    string
		in      CompareExpr
		wantErr bool
	}{
		{
			name: "valid with gte and lt",
			in: CompareExpr{
				Field: "score",
				Gte:   10,
				Lt:    100,
			},
		},
		{
			name: "missing field",
			in: CompareExpr{
				Gte: 10,
			},
			wantErr: true,
		},
		{
			name: "missing bounds",
			in: CompareExpr{
				Field: "score",
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

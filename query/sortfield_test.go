package query

import "testing"

func TestSortFieldValidate(t *testing.T) {
	tt := []struct {
		name    string
		in      SortField
		wantErr bool
	}{
		{
			name: "invalid direction",
			in: SortField{
				Field:     "created_at",
				Direction: "sideways",
			},
			wantErr: true,
		},
		{
			name: "missing field",
			in: SortField{
				Direction: SortDirectionAsc,
			},
			wantErr: true,
		},
		{
			name: "valid empty direction",
			in: SortField{
				Field: "created_at",
			},
		},
		{
			name: "valid asc",
			in: SortField{
				Field:     "created_at",
				Direction: SortDirectionAsc,
			},
		},
		{
			name: "valid desc",
			in: SortField{
				Field:     "created_at",
				Direction: SortDirectionDesc,
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var err error
			err = tc.in.Validate()
			if tc.wantErr && err == nil {
				t.Fatal("expected validation error")
			}

			if !tc.wantErr && err != nil {
				t.Fatalf("unexpected validation error: %v", err)
			}
		})
	}
}

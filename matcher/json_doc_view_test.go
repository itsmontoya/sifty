package matcher

import (
	"errors"
	"testing"
)

func TestNewJSONDocView(t *testing.T) {
	tt := []struct {
		name    string
		in      []byte
		wantErr bool
	}{
		{
			name:    "invalid json",
			in:      []byte(`{"title":`),
			wantErr: true,
		},
		{
			name: "valid json object",
			in:   []byte(`{"title":"golang","meta":{"lang":"go"}}`),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var (
				out *JSONDocView
				err error
			)

			out, err = NewJSONDocView(tc.in)
			if tc.wantErr && err == nil {
				t.Fatal("expected error")
			}

			if !tc.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !tc.wantErr && out == nil {
				t.Fatal("expected JSONDocView")
			}
		})
	}
}

func TestJSONDocViewGet(t *testing.T) {
	var (
		view *JSONDocView
		err  error
	)

	view, err = NewJSONDocView([]byte(`{
		"title":"golang",
		"meta":{"lang":"go","version":1},
		"obj":{"inner":{"name":"sifty"}}
	}`))
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	tt := []struct {
		name    string
		path    string
		want    any
		wantOK  bool
		wantErr error
	}{
		{
			name:   "get root scalar",
			path:   "title",
			want:   "golang",
			wantOK: true,
		},
		{
			name:   "get nested scalar",
			path:   "meta.lang",
			want:   "go",
			wantOK: true,
		},
		{
			name:   "get nested number",
			path:   "meta.version",
			want:   float64(1),
			wantOK: true,
		},
		{
			name:   "get nested object",
			path:   "obj.inner",
			wantOK: true,
		},
		{
			name:   "missing root key",
			path:   "missing",
			wantOK: false,
		},
		{
			name:   "missing nested key",
			path:   "meta.missing",
			wantOK: false,
		},
		{
			name:   "non-object encountered before last path segment",
			path:   "title.value",
			wantOK: false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var (
				got   any
				gotOK bool
				err   error
			)

			got, gotOK, err = view.Get(tc.path)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("unexpected error: %v", err)
			}

			if gotOK != tc.wantOK {
				t.Fatalf("ok = %v, want %v", gotOK, tc.wantOK)
			}

			if !tc.wantOK {
				return
			}

			if tc.want != nil && got != tc.want {
				t.Fatalf("value = %v, want %v", got, tc.want)
			}

			if tc.want == nil && got == nil {
				t.Fatal("expected non-nil value")
			}
		})
	}
}

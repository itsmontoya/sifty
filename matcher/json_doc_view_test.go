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

func BenchmarkNewJSONDocView(b *testing.B) {
	b.Run("flat_object", func(b *testing.B) {
		var in = []byte(`{"title":"golang matcher benchmark","score":12}`)

		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			var (
				out *JSONDocView
				err error
			)

			out, err = NewJSONDocView(in)
			if err != nil {
				b.Fatalf("NewJSONDocView() error = %v", err)
			}

			if out == nil {
				b.Fatal("NewJSONDocView() returned nil view")
			}
		}
	})
}

func BenchmarkJSONDocViewGet(b *testing.B) {
	var (
		view *JSONDocView
		err  error
	)

	view, err = NewJSONDocView([]byte(`{
		"a":"v1",
		"lvl1":{"b":"v2"},
		"lvl1b":{"lvl2":{"c":"v3"}},
		"lvl1c":{"lvl2":{"lvl3":{"d":"v4"}}}
	}`))
	if err != nil {
		b.Fatalf("NewJSONDocView() error = %v", err)
	}

	tt := []struct {
		name string
		path string
	}{
		{
			name: "depth_1",
			path: "a",
		},
		{
			name: "depth_2",
			path: "lvl1.b",
		},
		{
			name: "depth_3",
			path: "lvl1b.lvl2.c",
		},
		{
			name: "depth_4",
			path: "lvl1c.lvl2.lvl3.d",
		},
	}

	for _, tc := range tt {
		b.Run(tc.name, func(b *testing.B) {
			var (
				val any
				ok  bool
			)

			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				val, ok, err = view.Get(tc.path)
				if err != nil {
					b.Fatalf("Get(%q) error = %v", tc.path, err)
				}

				if !ok {
					b.Fatalf("Get(%q) ok = false, want true", tc.path)
				}

				if val == nil {
					b.Fatalf("Get(%q) value = nil, want non-nil", tc.path)
				}
			}
		})
	}
}

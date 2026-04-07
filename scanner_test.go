package sifty

import (
	"strings"
	"testing"

	"github.com/itsmontoya/sifty/matcher"
	"github.com/itsmontoya/sifty/query"
)

func TestScannerProcessRowInvalidJSON(t *testing.T) {
	t.Parallel()

	m, err := matcher.Compile(query.Query{})
	if err != nil {
		t.Fatalf("unexpected compile error: %v", err)
	}

	s := makeScanner(m, 10)
	err = s.processRow([]byte("not-json"))
	if err == nil {
		t.Fatal("expected invalid json error")
	}

	if !strings.Contains(err.Error(), "error unmarshaling bytes as a JSON object") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestScannerAppendLimitNegativeOne(t *testing.T) {
	t.Parallel()

	s := scanner{
		limit: -1,
	}

	if !s.isAtLimit() {
		t.Fatal("expected scanner to be at limit when limit is -1")
	}

	err := s.append([]byte(`{"value":{"foo":1}}`))
	if err != errBreak {
		t.Fatalf("append error = %v, want %v", err, errBreak)
	}

	if got, want := len(s.matches), 1; got != want {
		t.Fatalf("match count = %d, want %d", got, want)
	}
}

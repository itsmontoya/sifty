package sifty

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/itsmontoya/sifty/matcher"
	"github.com/itsmontoya/sifty/query"
)

func TestScannerProcessRowInvalidJSON(t *testing.T) {
	t.Parallel()

	m, err := matcher.Compile(query.Query{})
	if err != nil {
		t.Fatalf("unexpected compile error: %v", err)
	}

	s := scanner{
		m: m,
	}

	err = s.processRow(rawRow{
		Timestamp: time.Now(),
		Value:     json.RawMessage("not-json"),
	})

	if err == nil {
		t.Fatal("expected invalid json error")
	}

	if !strings.Contains(err.Error(), "error unmarshaling bytes as a JSON object") {
		t.Fatalf("unexpected error: %v", err)
	}
}

package sifty_test

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/itsmontoya/sifty"
	"github.com/itsmontoya/sifty/query"
)

type testEntry struct {
	Foo int    `json:"foo"`
	Bar int    `json:"bar"`
	Tag string `json:"tag"`
}

func TestNewWithInvalidPath(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	filePath := filepath.Join(dir, "db-file")
	if err := os.WriteFile(filePath, []byte("x"), 0600); err != nil {
		t.Fatalf("unexpected write error: %v", err)
	}

	s, err := sifty.New(filePath, 10)
	if s != nil {
		t.Fatalf("sifty instance = %#v, expected nil", s)
	}

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestRowJSONMarshaling(t *testing.T) {
	t.Parallel()

	in := sifty.Row{
		Timestamp: time.Date(2026, time.April, 4, 12, 0, 0, 0, time.UTC),
		Value: testEntry{
			Foo: 1,
			Bar: 2,
			Tag: "alpha",
		},
	}

	bs, err := json.Marshal(in)
	if err != nil {
		t.Fatalf("unexpected marshal error: %v", err)
	}

	var out map[string]json.RawMessage
	if err = json.Unmarshal(bs, &out); err != nil {
		t.Fatalf("unexpected unmarshal error: %v", err)
	}

	if _, ok := out["timestamp"]; !ok {
		t.Fatal("missing timestamp field")
	}

	if _, ok := out["value"]; !ok {
		t.Fatal("missing value field")
	}
}

func TestSiftyAppendAndScan(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	s, err := sifty.New(dir, 10)
	if err != nil {
		t.Fatalf("unexpected new error: %v", err)
	}

	entries := []testEntry{
		{Foo: 1, Bar: 10, Tag: "alpha"},
		{Foo: 2, Bar: 20, Tag: "beta"},
		{Foo: 3, Bar: 30, Tag: "alpha"},
	}

	for _, entry := range entries {
		if err = s.Append(entry); err != nil {
			t.Fatalf("unexpected append error: %v", err)
		}
	}

	matches, err := s.Scan(query.Query{Filter: query.Clause{Contains: &query.ContainsExpr{Field: "tag", Value: "alpha"}}}, 10)
	if err != nil {
		t.Fatalf("unexpected scan error: %v", err)
	}

	if got, want := len(matches), 2; got != want {
		t.Fatalf("match count = %d, want %d", got, want)
	}

	rows := decodeScannedEntries(t, matches)
	if got, want := rows[0].Foo, 1; got != want {
		t.Fatalf("rows[0].foo = %d, want %d", got, want)
	}

	if got, want := rows[1].Foo, 3; got != want {
		t.Fatalf("rows[1].foo = %d, want %d", got, want)
	}
}

func TestSiftyScanLimit(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	s, err := sifty.New(dir, 10)
	if err != nil {
		t.Fatalf("unexpected new error: %v", err)
	}

	for i := 0; i < 3; i++ {
		if err = s.Append(testEntry{Foo: i, Tag: "x"}); err != nil {
			t.Fatalf("unexpected append error: %v", err)
		}
	}

	matches, err := s.Scan(query.Query{}, 2)
	if err == nil {
		t.Fatal("expected scan limit error")
	}

	if got, want := err.Error(), "break"; got != want {
		t.Fatalf("scan error = %q, want %q", got, want)
	}

	if got, want := len(matches), 0; got != want {
		t.Fatalf("match count = %d, want %d", got, want)
	}
}

func TestSiftyScanInvalidQuery(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	s, err := sifty.New(dir, 10)
	if err != nil {
		t.Fatalf("unexpected new error: %v", err)
	}

	_, err = s.Scan(query.Query{Filter: query.Clause{Compare: &query.CompareExpr{Field: "foo"}}}, 10)
	if err == nil {
		t.Fatal("expected query validation error")
	}

	if !strings.Contains(err.Error(), "cannot compile, invalid query") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSiftyLoadsExistingFile(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()

	first, err := sifty.New(dir, 10)
	if err != nil {
		t.Fatalf("unexpected new error: %v", err)
	}

	seed := []testEntry{
		{Foo: 10, Bar: 100, Tag: "seed"},
		{Foo: 11, Bar: 110, Tag: "seed"},
	}

	for _, entry := range seed {
		if err = first.Append(entry); err != nil {
			t.Fatalf("unexpected append error on first instance: %v", err)
		}
	}

	second, err := sifty.New(dir, 10)
	if err != nil {
		t.Fatalf("unexpected new error on reopen: %v", err)
	}

	matches, err := second.Scan(
		query.Query{
			Filter: query.Clause{
				Contains: &query.ContainsExpr{
					Field: "tag",
					Value: "seed",
				},
			},
		},
		10,
	)
	if err != nil {
		t.Fatalf("unexpected scan error on reopened instance: %v", err)
	}

	if got, want := len(matches), len(seed); got != want {
		t.Fatalf("match count on reopened instance = %d, want %d", got, want)
	}

	rows := decodeScannedEntries(t, matches)
	if got, want := rows[0].Foo, seed[0].Foo; got != want {
		t.Fatalf("rows[0].foo = %d, want %d", got, want)
	}

	if got, want := rows[1].Foo, seed[1].Foo; got != want {
		t.Fatalf("rows[1].foo = %d, want %d", got, want)
	}
}

func TestSiftyScanTimeRangeIncludesOlderSegments(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	s, err := sifty.New(dir, 5)
	if err != nil {
		t.Fatalf("unexpected new error: %v", err)
	}

	for i := 0; i < 15; i++ {
		if err = s.Append(testEntry{Foo: i, Tag: "time-range"}); err != nil {
			t.Fatalf("unexpected append error: %v", err)
		}
	}

	now := time.Now()
	from := now.Add(-5 * time.Minute)
	to := now.Add(5 * time.Minute)

	nonEmptyLogFiles, err := nonEmptySegmentNames(dir)
	if err != nil {
		t.Fatalf("unexpected non-empty segment read error: %v", err)
	}

	if got, want := len(nonEmptyLogFiles), 3; got != want {
		t.Fatalf("non-empty segment count = %d, want %d", got, want)
	}

	olderSegmentName := nonEmptyLogFiles[0]
	newerSegmentName := nonEmptyLogFiles[len(nonEmptyLogFiles)-1]

	olderSegmentToIgnoreName := fmt.Sprintf("%s.log", from.Add(-2*time.Hour).Format(time.RFC3339Nano))
	if err = os.Rename(filepath.Join(dir, olderSegmentName), filepath.Join(dir, olderSegmentToIgnoreName)); err != nil {
		t.Fatalf("unexpected older segment rename error: %v", err)
	}

	futureSegmentName := fmt.Sprintf("%s.log", to.Add(2*time.Hour).Format(time.RFC3339Nano))
	if err = os.Rename(filepath.Join(dir, newerSegmentName), filepath.Join(dir, futureSegmentName)); err != nil {
		t.Fatalf("unexpected newer segment rename error: %v", err)
	}

	reopened, err := sifty.New(dir, 5)
	if err != nil {
		t.Fatalf("unexpected reopen error: %v", err)
	}

	matches, err := reopened.Scan(
		query.Query{
			Filter: query.Clause{
				Contains: &query.ContainsExpr{
					Field: "tag",
					Value: "time-range",
				},
			},
			TimeRange: &query.TimeRange{
				From: &from,
				To:   &to,
			},
		},
		20,
	)
	if err != nil {
		t.Fatalf("unexpected scan error: %v", err)
	}

	if got, want := len(matches), 5; got != want {
		t.Fatalf("match count = %d, want %d", got, want)
	}

	rows := decodeScannedEntries(t, matches)
	sort.Slice(rows, func(i, j int) bool {
		return rows[i].Foo < rows[j].Foo
	})

	for i := 0; i < 5; i++ {
		if got, want := rows[i].Foo, i+5; got != want {
			t.Fatalf("rows[%d].foo = %d, want %d", i, got, want)
		}
	}
}

func nonEmptySegmentNames(dir string) (out []string, err error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".log" {
			continue
		}

		var info os.FileInfo
		if info, err = entry.Info(); err != nil {
			return nil, err
		}

		if info.Size() > 0 {
			out = append(out, entry.Name())
		}
	}

	sort.Strings(out)
	return out, nil
}

func decodeScannedEntries(t *testing.T, matches []any) (rows []testEntry) {
	t.Helper()

	rows = make([]testEntry, 0, len(matches))
	for i, match := range matches {
		raw, ok := match.(json.RawMessage)
		if !ok {
			t.Fatalf("matches[%d] type = %T, want json.RawMessage", i, match)
		}

		var entry testEntry
		if err := json.Unmarshal(raw, &entry); err != nil {
			t.Fatalf("unexpected value unmarshal error at index %d: %v", i, err)
		}

		rows = append(rows, entry)
	}

	return rows
}

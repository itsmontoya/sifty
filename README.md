# sifty

`sifty` is a lightweight append-and-scan library for JSON-like documents.
It stores rows on disk, adds a timestamp per row, and evaluates filters using
the built-in `query` + `matcher` packages.

## Install

```bash
go get github.com/itsmontoya/sifty
```

## What It Does

- Appends values as rows with:
  - `timestamp`
  - `value`
- Persists rows to local segment files via `iodb`
- Scans rows with compiled query filters

## Core API

- `New(path string, segmentSize int) (*Sifty, error)`
- `(*Sifty).Append(in any) error`
- `(*Sifty).Scan(q query.Query, limit int) ([]any, error)`

## Segment Size

`segmentSize` is always finite.
There is no "no limit" mode or sentinel value for unlimited segment size.
Use a positive `segmentSize`.

## Quick Example

```go
package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/itsmontoya/sifty"
	"github.com/itsmontoya/sifty/query"
)

type Entry struct {
	Title string `json:"title"`
	Score int    `json:"score"`
}

func main() {
	store, err := sifty.New("./data", 1000)
	if err != nil {
		log.Fatal(err)
	}

	if err = store.Append(Entry{Title: "go patterns", Score: 42}); err != nil {
		log.Fatal(err)
	}

	results, err := store.Scan(
		query.Query{
			Filter: query.Clause{
				Contains: &query.ContainsExpr{
					Field: "value.title",
					Value: "go",
				},
			},
		},
		50,
	)
	if err != nil {
		log.Fatal(err)
	}

	for _, raw := range results {
		bs, ok := raw.(json.RawMessage)
		if !ok {
			continue
		}

		fmt.Println(string(bs))
	}
}
```

## Query Notes

- Filter paths target stored row fields, so application data is typically under
  `value.*` (for example `value.title`).
- Query validation errors are returned by `Scan`.

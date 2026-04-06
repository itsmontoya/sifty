package sifty

import (
	"encoding/json"
	"io"
	"time"
)

func makeRow(v any) (r Row) {
	r.Timestamp = time.Now()
	r.Value = v
	return r
}

type Row struct {
	Timestamp time.Time `json:"timestamp"`
	Value     any       `json:"value"`
}

func (r *Row) append(w io.Writer) error {
	return json.NewEncoder(w).Encode(r)
}

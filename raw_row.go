package sifty

import (
	"encoding/json"
	"time"
)

type rawRow struct {
	Timestamp time.Time       `json:"timestamp"`
	Value     json.RawMessage `json:"value"`
}

package sifty

import "time"

type Row struct {
	Timestamp time.Time `json:"timestamp"`
	Value     any       `json:"value"`
}

package sifty

import "time"

type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Value     any       `json:"value"`
}

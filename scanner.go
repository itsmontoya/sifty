package sifty

import (
	"encoding/json"
	"io"

	"github.com/itsmontoya/sifty/matcher"
)

func makeScanner(m *matcher.Matcher, limit int) (s scanner) {
	s.m = m
	s.limit = limit
	return s
}

type scanner struct {
	m     *matcher.Matcher
	limit int

	matches []any
}

func (s *scanner) process(r io.Reader) (err error) {
	return iterateRows(r, s.processRow)
}

func (s *scanner) processRow(raw json.RawMessage) (err error) {
	var view matcher.DocView
	if view, err = matcher.NewJSONDocView(raw); err != nil {
		return err
	}

	var ok bool
	if ok, err = s.m.IsMatch(view); err != nil {
		return err
	}

	if !ok {
		return nil
	}

	return s.append(raw)
}

func (s *scanner) append(value any) (err error) {
	s.matches = append(s.matches, value)
	if s.isAtLimit() {
		return errBreak
	}

	return nil
}

func (s *scanner) isAtLimit() (ok bool) {
	switch s.limit {
	case -1:
		return true
	default:
		return len(s.matches) >= s.limit
	}
}

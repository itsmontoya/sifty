package sifty

import (
	"io"

	"github.com/itsmontoya/sifty/docview"
	"github.com/itsmontoya/sifty/docview/jsondoc"
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

func (s *scanner) processRow(row rawRow) (err error) {
	var view docview.DocView
	if view, err = jsondoc.NewJSONDoc(row.Value); err != nil {
		return err
	}

	var ok bool
	if ok, err = s.m.IsMatch(row.Timestamp, view); !ok || err != nil {
		return err
	}

	return s.append(row.Value)
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

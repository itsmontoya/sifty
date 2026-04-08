package sifty

import (
	"io"

	"github.com/itsmontoya/iodb"
	"github.com/itsmontoya/sifty/docview"
	"github.com/itsmontoya/sifty/docview/jsondoc"
	"github.com/itsmontoya/sifty/matcher"
)

func makeScanner(m *matcher.Matcher, f *iodb.File, matches chan result, limit int) (s scanner) {
	s.m = m
	s.f = f
	s.ch = matches
	s.limit = limit
	return s
}

type scanner struct {
	m *matcher.Matcher
	f *iodb.File

	limit int

	ch     chan result
	result result
}

func (s *scanner) process() {
	s.result.err = s.f.Read(s.processReader)
	s.ch <- s.result
}

func (s *scanner) processReader(r io.Reader) error {
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
	if s.isAtLimit() {
		return errBreak
	}

	s.result.matches = append(s.result.matches, value)
	return nil
}

func (s *scanner) isAtLimit() (ok bool) {
	switch s.limit {
	case -1:
		return true
	default:
		return len(s.result.matches) >= s.limit
	}
}

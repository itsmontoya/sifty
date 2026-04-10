package sifty

import (
	"encoding/json"
	"io"

	"github.com/itsmontoya/iodb"
	"github.com/itsmontoya/sifty/docview"
	"github.com/itsmontoya/sifty/docview/jsondoc"
)

func makeScanner(r *request, f *iodb.File) (s scanner, err error) {
	s.request = r
	s.in = f
	if s.out, err = r.bkt.Create(f.Key()); err != nil {
		return s, err
	}

	return s, nil
}

type scanner struct {
	*request
	result result

	in  *iodb.File
	out *iodb.File
}

func (s *scanner) process() {
	var errs []error
	errs = append(errs, s.in.Read(s.processReader))

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

func (s *scanner) append(value json.RawMessage) (err error) {
	err = s.out.Append(func(w io.Writer) (err error) {
		return json.NewEncoder(w).Encode(value)
	})

	return err
}

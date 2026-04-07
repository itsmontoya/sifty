package sifty

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/itsmontoya/iodb"
	"github.com/itsmontoya/sifty/matcher"
	"github.com/itsmontoya/sifty/query"
)

var errBreak = errors.New("break")

// New creates or opens a sifty store at path.
//
// segmentSize controls when a new segment file is rotated in.
// There is no "no limit" option for segment size, and no sentinel value
// enables unlimited growth of a single segment.
// Pass a positive segmentSize to use finite segment rotation.
func New(path string, segmentSize int) (out *Sifty, err error) {
	var s Sifty
	if s.db, err = iodb.New(path); err != nil {
		return nil, err
	}

	s.maxCount = segmentSize

	if err = s.init(); err != nil {
		return nil, err
	}

	return &s, nil
}

type Sifty struct {
	wmux sync.Mutex

	db *iodb.DB

	f *iodb.File

	count    int
	maxCount int
}

func (s *Sifty) Append(in any) (err error) {
	s.wmux.Lock()
	defer s.wmux.Unlock()
	if err = s.append(in); err != nil {
		return err
	}

	s.count++

	return s.rotate()
}

func (s *Sifty) Scan(q query.Query, limit int) (matches []any, err error) {
	var m *matcher.Matcher
	if m, err = matcher.Compile(q); err != nil {
		return nil, err
	}

	scn := makeScanner(m, limit)
	if err = s.f.Read(scn.process); err == errBreak {
		err = nil
	}

	return scn.matches, err
}

func (s *Sifty) append(in any) (err error) {
	row := makeRow(in)
	return s.f.Append(row.append)
}

func (s *Sifty) rotate() (err error) {
	if s.f != nil && s.count < s.maxCount {
		return nil
	}

	s.count = 0
	key := fmt.Sprintf("%s.log", time.Now().Format(time.RFC3339))
	s.f, err = s.db.Create(key)
	return err
}

func (s *Sifty) init() (err error) {
	if err = s.loadLatest(); err != nil {
		return err
	}

	if err = s.setCountFromFile(); err != nil {
		return err
	}

	return s.rotate()
}

func (s *Sifty) loadLatest() (err error) {
	err = s.db.Cursor(func(c *iodb.Cursor) (err error) {
		s.f, _ = c.Last()
		return nil
	})

	return err
}

func (s *Sifty) setCountFromFile() (err error) {
	if s.f == nil {
		return nil
	}

	return s.f.Read(s.setCountFromRows)
}

func (s *Sifty) setCountFromRows(r io.Reader) error {
	return iterateRows(r, s.setCountFromRow)
}

func (s *Sifty) setCountFromRow(bs json.RawMessage) (err error) {
	s.count++
	return nil
}

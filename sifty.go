package sifty

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
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
	// Timestamp for when the current file was created
	createdFileAt time.Time

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

func (s *Sifty) Scan(q query.Query) (matches []any, err error) {
	var m *matcher.Matcher
	if m, err = matcher.Compile(q); err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	ch := make(chan result, 4)
	var errs []error
	err = s.iterateFilesInReverse(func(f *iodb.File) (err error) {
		var ts time.Time
		if ts, err = keyToTimestamp(f.Key()); err != nil {
			return err
		}

		switch m.RangeBounds(ts) {
		case 0:
		case 1:
			return nil
		case -1:
			return errBreak
		}

		scn := makeScanner(m, f, ch)
		wg.Go(scn.process)
		return nil
	})

	switch err {
	case nil:
	case errBreak:
	default:
		errs = append(errs, err)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for result := range ch {
		if result.err != nil {
			errs = append(errs, result.err)
			continue
		}

		matches = append(matches, result.matches...)
	}

	return matches, errors.Join(errs...)
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
	createdAt := time.Now()
	key := fmt.Sprintf("%s.log", createdAt.Format(time.RFC3339Nano))
	s.f, err = s.db.Create(key)
	s.createdFileAt = createdAt
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

func (s *Sifty) iterateFilesInReverse(fn func(*iodb.File) error) (err error) {
	err = s.db.Cursor(func(c *iodb.Cursor) (err error) {
		for f, ok := c.Last(); ok && err == nil; f, ok = c.Prev() {
			err = fn(f)
		}

		return err
	})

	return err
}

func keyToTimestamp(key string) (out time.Time, err error) {
	stripped := strings.Replace(key, ".log", "", 1)
	if out, err = time.Parse(time.RFC3339Nano, stripped); err != nil {
		err = fmt.Errorf(`error parsing key of "%s": %w`, key, err)
		return out, err
	}

	return out, nil
}

package sifty

import (
	"sync"
	"time"

	"github.com/itsmontoya/iodb"
	"github.com/itsmontoya/sifty/matcher"
)

type request struct {
	bkt *iodb.Bucket
	m   *matcher.Matcher
	wg  sync.WaitGroup
}

func (r *request) process(f *iodb.File) (err error) {
	var ts time.Time
	if ts, err = keyToTimestamp(f.Key()); err != nil {
		return err
	}

	switch r.m.RangeBounds(ts) {
	case 0:
	case 1:
		return nil
	case -1:
		return errBreak
	}

	var scn scanner
	if scn, err = makeScanner(r, f); err != nil {
		return err
	}

	r.wg.Go(scn.process)
	return nil

}

package sifty

import (
	"encoding/json"
	"errors"
	"io"
)

func iterateRows[T any](r io.Reader, fn func(T) error) (err error) {
	dec := json.NewDecoder(r)
	for err == nil {
		var t T
		if err = dec.Decode(&t); err != nil {
			break
		}

		err = fn(t)
	}

	if errors.Is(err, io.EOF) {
		return nil
	}

	return err
}

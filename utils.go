package sifty

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"
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

func keyToTimestamp(key string) (out time.Time, err error) {
	stripped := strings.Replace(key, ".log", "", 1)
	if out, err = time.Parse(time.RFC3339Nano, stripped); err != nil {
		err = fmt.Errorf(`error parsing key of "%s": %w`, key, err)
		return out, err
	}

	return out, nil
}

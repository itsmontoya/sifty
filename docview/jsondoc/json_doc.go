package jsondoc

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/itsmontoya/sifty/docview"
)

var _ docview.DocView = &JSONDoc{}

// NewJSONDoc constructs a DocView backed by a JSON object payload.
// The input must decode into a JSON object at the top level.
func NewJSONDoc(bs []byte) (out *JSONDoc, err error) {
	var j JSONDoc
	if err = json.Unmarshal(bs, &j.m); err != nil {
		return nil, fmt.Errorf("error unmarshaling bytes as a JSON object: %w", err)
	}

	return &j, nil
}

// JSONDoc resolves dot-delimited field paths from a decoded JSON object.
type JSONDoc struct {
	m map[string]any
}

// Get returns the value at path when all path segments resolve through objects.
// If any segment is missing, or a non-object value is encountered before the
// final segment, it returns (nil, false, nil).
func (j *JSONDoc) Get(path string) (v any, ok bool, err error) {
	splitPath := strings.Split(path, ".")
	m := j.m
	for i, part := range splitPath {
		if v, ok, err = getPartValue(m, part); err != nil {
			return nil, false, err
		}

		if i == len(splitPath)-1 {
			break
		}

		if m, ok = v.(map[string]any); !ok {
			return nil, false, nil
		}

	}

	return v, ok, nil
}

func getPartValue(m map[string]any, part string) (v any, ok bool, err error) {
	var index int
	if part, index, err = parseKey(part); err != nil {
		return nil, false, err
	}

	if v, ok = m[part]; !ok {
		return nil, false, nil
	}

	if index == -1 {
		return v, true, nil
	}

	if v, err = getIndexedValue(v, index); err != nil {
		return nil, false, err
	}

	return v, true, nil
}

func getIndexedValue(v any, index int) (out any, err error) {
	var (
		a  []any
		ok bool
	)
	if a, ok = v.([]any); !ok {
		return nil, fmt.Errorf("calling index on non-array type: %T", a)
	}

	return a[index], nil
}

func parseKey(key string) (out string, i int, err error) {
	start := strings.Index(key, "[")
	if start == -1 {
		return key, -1, nil
	}

	end := strings.Index(key, "]")
	iStr := key[start+1 : end]
	if i, err = strconv.Atoi(iStr); err != nil {
		return
	}

	out = key[:start]
	return
}

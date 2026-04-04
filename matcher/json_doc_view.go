package matcher

import (
	"encoding/json"
	"fmt"
	"strings"
)

var _ DocView = &JSONDocView{}

// NewJSONDocView constructs a DocView backed by a JSON object payload.
// The input must decode into a JSON object at the top level.
func NewJSONDocView(bs []byte) (out *JSONDocView, err error) {
	var j JSONDocView
	if err = json.Unmarshal(bs, &j.m); err != nil {
		return nil, fmt.Errorf("error unmarshaling bytes as a JSON object: %w", err)
	}

	return &j, nil
}

// JSONDocView resolves dot-delimited field paths from a decoded JSON object.
type JSONDocView struct {
	m map[string]any
}

// Get returns the value at path when all path segments resolve through objects.
// If any segment is missing, or a non-object value is encountered before the
// final segment, it returns (nil, false, nil).
func (j *JSONDocView) Get(path string) (v any, ok bool, err error) {
	splitPath := strings.Split(path, ".")
	m := j.m
	for i, part := range splitPath {
		v, ok = m[part]
		if !ok {
			return nil, false, nil
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

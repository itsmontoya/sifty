package matcher

import (
	"encoding/json"
	"strings"
)

var _ DocView = &JSONDocView{}

func NewJSONDocView(bs []byte) (out *JSONDocView, err error) {
	var j JSONDocView
	if err = json.Unmarshal(bs, &j.m); err != nil {
		return nil, err
	}

	return &j, nil
}

type JSONDocView struct {
	m map[string]any
}

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

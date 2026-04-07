package jsondoc

import (
	"fmt"
	"strconv"
	"strings"
)

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

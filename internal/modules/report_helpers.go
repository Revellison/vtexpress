package modules

import (
	"fmt"
)

func nestedMap(root map[string]any, path ...string) (map[string]any, error) {
	current := root
	for _, key := range path {
		nextRaw, ok := current[key]
		if !ok {
			return nil, fmt.Errorf("missing key: %s", key)
		}
		next, ok := nextRaw.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("key %s is not an object", key)
		}
		current = next
	}
	return current, nil
}

func pickString(values map[string]any, key string, fallback string) string {
	if values == nil {
		return fallback
	}
	raw, ok := values[key]
	if !ok || raw == nil {
		return fallback
	}
	str, ok := raw.(string)
	if !ok {
		return fallback
	}
	if str == "" {
		return fallback
	}
	return str
}

func pickInt(values map[string]any, key string) int {
	if values == nil {
		return 0
	}
	raw, ok := values[key]
	if !ok || raw == nil {
		return 0
	}
	switch v := raw.(type) {
	case int:
		return v
	case int32:
		return int(v)
	case int64:
		return int(v)
	case float64:
		return int(v)
	default:
		return 0
	}
}

func pickInt64(values map[string]any, key string) int64 {
	if values == nil {
		return 0
	}
	raw, ok := values[key]
	if !ok || raw == nil {
		return 0
	}
	switch v := raw.(type) {
	case int:
		return int64(v)
	case int32:
		return int64(v)
	case int64:
		return v
	case float64:
		return int64(v)
	default:
		return 0
	}
}

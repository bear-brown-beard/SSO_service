package utils

import (
	"encoding/json"
	"strconv"
)

func MapString(data map[string]any, key string) string {
	_, ok := data[key]
	if !ok {
		return ""
	}
	switch v := data[key].(type) {
	case string:
		return v
	case json.Number:
		return v.String()
	case int:
		return strconv.Itoa(v)
	}

	return ""
}

func MapFloat(data map[string]any, key string) float64 {
	if val, ok := data[key].(float64); ok {
		return val
	}
	return 0
}

func MapInt(data map[string]any, key string) int {
	if val, ok := data[key].(float64); ok {
		return int(val)
	}
	return 0
}
func MapJsonNumber(data map[string]any, key string) int {
	_, ok := data[key]
	if !ok {
		return 0
	}

	switch v := data[key].(type) {
	case json.Number:
		valInt, err := v.Int64()
		if err != nil {
			return 0
		}
		return int(valInt)

	case int:
		return v
	case float64:
		return int(v)
	case string:
		valInt, err := strconv.Atoi(v)
		if err != nil {
			return 0
		}
		return valInt
	}
	return 0
}

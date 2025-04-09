package model

import (
	"encoding/json"
	"fmt"
	"github.com/charmbracelet/log"
	"strings"
)

type SqlRow struct {
	EntityName string
	Data       map[string]interface{}
}

func GetRequiredValue[T any](row SqlRow, path string, valuePtr *T) bool {
	return GetRowValue[T](row, path, valuePtr, true)
}

func GetOptionalValue[T any](row SqlRow, path string, valuePtr *T) bool {
	return GetRowValue[T](row, path, valuePtr, false)
}

func GetRowValue[T any](row SqlRow, path string, valuePtr *T, required bool) bool {
	var err error
	var result T
	if result, err = getPathValue[T](row.Data, path); err != nil {
		if required {
			log.Infof("%s -> %s: %s", row.EntityName, path, err.Error())
			if log.GetLevel() >= log.DebugLevel {
				dataStr, _ := json.Marshal(row.Data)
				log.Debugf("data: %s", dataStr)
			}
			return false
		}
	}
	*valuePtr = result
	return true

}

const PathSeparator = "."

func getPathValue[T any](data map[string]interface{}, path string) (T, error) {
	keys := strings.Split(path, PathSeparator)
	var zero T // Default zero value for the type T
	var value interface{} = data
	var ok bool

	for _, key := range keys {
		switch casted := value.(type) {
		case map[string]interface{}:
			if value, ok = casted[key]; !ok {
				return zero, fmt.Errorf("missing value for %s in path %s", key, path)
			}
		case []interface{}:
			if index, ok := parseIndex(key); ok && index >= 0 && index < len(casted) {
				value = casted[index]
			} else {
				return zero, fmt.Errorf("invalid index %s in path %s", key, path)
			}
		default:
			return zero, fmt.Errorf("type mismatch for %s in path %s, expected %T, got %T", key, path, zero, value)
		}
	}

	if finalValue, ok := value.(T); ok {
		return finalValue, nil
	}
	return zero, fmt.Errorf("type mismatch for at %s, expected %T, got %T", path, zero, value)
}

func parseIndex(key string) (int, bool) {
	var index int
	_, err := fmt.Sscanf(key, "[%d]", &index)
	return index, err == nil
}

package adapter

import (
	"fmt"
	"github.com/ctrlplanedev/cli/internal/api"
	"strings"
)

type SteampipeAdapter struct {
	Table     string
	Translate func(data *map[string]interface{}) (api.AgentResource, bool)
}

func getPathValue[T any](data *map[string]interface{}, path string) (T, error) {
	keys := strings.Split(path, ".")
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

package adapter

import (
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/ctrlplanedev/cli/internal/api"
	"strings"
)

type SteampipeAdapter struct {
	Table     string
	Translate func(data *map[string]interface{}) (api.AgentResource, bool)
}

func getPathValue[T any](data *map[string]interface{}, path string) (T, bool) {
	keys := strings.Split(path, ".")
	return getValue[T](data, keys)
}

func getValue[T any](data interface{}, keys []string) (T, bool) {
	var zero T // Default zero value for the type T
	var exists bool
	var value interface{}

	if len(keys) == 0 {
		if value, ok := data.(T); ok {
			return value, true
		}
		return zero, false
	}

	currentKey := keys[0]
	switch casted := data.(type) {
	case map[string]interface{}:
		value, exists = casted[currentKey]
	case []interface{}:
		if index, ok := parseIndex(currentKey); ok && index >= 0 && index < len(casted) {
			value = casted[index]
		} else {
			log.Warn("could not find index in array")
			return zero, false
		}
	}
	if !exists {
		return zero, false
	}

	if len(keys) == 1 {
		if typedValue, ok := value.(T); ok {
			return typedValue, true
		}
		return zero, false
	}

	return getValue[T](value, keys[1:])
}

func parseIndex(key string) (int, bool) {
	var index int
	_, err := fmt.Sscanf(key, "[%d]", &index)
	return index, err == nil
}

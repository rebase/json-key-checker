package jsonutil

import (
	"encoding/json"
	"fmt"
	"os"
)

func collectKeysRecursiveWithPath(data map[string]interface{}, prefix string, keys map[string]struct{}) {
	for key, value := range data {
		fullKey := key
		if prefix != "" {
			fullKey = prefix + "." + key
		}

		keys[fullKey] = struct{}{}

		if nestedMap, ok := value.(map[string]interface{}); ok {
			collectKeysRecursiveWithPath(nestedMap, fullKey, keys)
		}
	}
}

func LoadJsonKeysWithPath(filePath string) (map[string]struct{}, error) {
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("Error reading file '%s': %w", filePath, err)
	}

	var data map[string]interface{}
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return nil, fmt.Errorf("Error parsing JSON '%s': %w", filePath, err)
	}

	keys := make(map[string]struct{})
	collectKeysRecursiveWithPath(data, "", keys)

	return keys, nil
}

func WriteJsonFile(filePath string, data map[string]interface{}) error {
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("Error marshaling JSON '%s': %w", filePath, err)
	}

	err = os.WriteFile(filePath, bytes, 0644)
	if err != nil {
		return fmt.Errorf("Error writing file '%s': %w", filePath, err)
	}
	return nil
}

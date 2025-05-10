package keyadder

import (
	"strings"
)

func AddMissingKeySplitByFirstDot(data map[string]interface{}, fullKey string) (int, error) {
	addedCount := 0
	currentMap := data

	firstDotIndex := strings.Index(fullKey, ".")

	if firstDotIndex == -1 {
		finalKeyName := fullKey
		if _, exists := currentMap[finalKeyName]; !exists {
			currentMap[finalKeyName] = ""
			addedCount++
		}
		return addedCount, nil
	}

	currentSegment := fullKey[:firstDotIndex]
	finalKeyName := fullKey[firstDotIndex+1:]

	value, exists := currentMap[currentSegment]
	if exists {
		if nestedMap, ok := value.(map[string]interface{}); ok {
			currentMap = nestedMap
		} else {
			newMap := make(map[string]interface{})
			currentMap[currentSegment] = newMap
			currentMap = newMap
		}
	} else {
		newMap := make(map[string]interface{})
		currentMap[currentSegment] = newMap
		currentMap = newMap
	}

	if _, exists := currentMap[finalKeyName]; !exists {
		currentMap[finalKeyName] = ""
		addedCount++
	}

	return addedCount, nil
}

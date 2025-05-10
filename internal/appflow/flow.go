package appflow

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/rebase/json-key-checker/internal/jsonutil"
	"github.com/rebase/json-key-checker/internal/keyadder"
	"github.com/rebase/json-key-checker/internal/output"
)

func CompareJsonFiles(filePaths []string) (map[string][]string, error) {
	allKeys := make(map[string]struct{})
	fileKeysMap := make(map[string]map[string]struct{})
	missingKeysMap := make(map[string][]string)
	anyMissing := false

	fmt.Println(output.InfoColor("Analyzing files..."))
	for _, path := range filePaths {

		keys, err := jsonutil.LoadJsonKeysWithPath(path)
		if err != nil {
			return nil, fmt.Errorf(output.ErrorColor(output.Bold("Error: Failed to process '%s' - %v"))+"\n", path, err)
		}
		fileKeysMap[path] = keys
		for key := range keys {
			allKeys[key] = struct{}{}
		}
		fmt.Printf(output.SuccessColor("- '%s' processed (%d keys)")+"\n", path, len(keys))
	}

	fmt.Println(output.HeaderColor("\n--- Missing Keys Comparison ---"))
	for path, currentKeys := range fileKeysMap {
		var missing []string
		for key := range allKeys {
			if _, exists := currentKeys[key]; !exists {
				missing = append(missing, key)
			}
		}
		if len(missing) > 0 {
			missingKeysMap[path] = missing
			anyMissing = true
			fmt.Printf(output.Bold("File '%s' has missing keys (%d total):")+"\n", path, len(missing))
			sort.Strings(missing)
			for _, mKey := range missing {
				fmt.Printf("  %s\n", output.MissingKeyColor("- "+mKey))
			}
		} else {
			fmt.Printf(output.SuccessColor("File '%s' has no missing keys.")+"\n", path)
		}
	}

	if !anyMissing {
		fmt.Println(output.SuccessColor("\nNo missing keys found in any file. Exiting."))
		return nil, fmt.Errorf("no missing keys found")
	}

	return missingKeysMap, nil
}

func ProcessFileWithNestedKeys(path string, missingKeys []string) error {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("Error reading file '%s': %v", path, err)
	}

	var data map[string]interface{}
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return fmt.Errorf("Error parsing JSON '%s': %v", path, err)
	}

	processedCount := 0

	for _, missingKey := range missingKeys {

		_, err := keyadder.AddMissingKeySplitByFirstDot(data, missingKey)
		if err != nil {
			fmt.Printf("  "+output.WarningColor("Warning: Error adding key '%s' - %v")+"\n", missingKey, err)
		}
		processedCount++
	}

	err = jsonutil.WriteJsonFile(path, data)
	if err != nil {
		return fmt.Errorf("Error writing file '%s': %v", path, err)
	}

	fmt.Printf("  "+output.SuccessColor("Success: Added %d keys to '%s'.")+"\n", processedCount, path)
	return nil
}

func PromptAndProcessFiles(missingKeysMap map[string][]string) error {
	fmt.Print(output.Bold("\nAdd missing keys to each file? (y/N): "))
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))

	if input != "y" {
		fmt.Println(output.WarningColor("Key addition cancelled."))
		return fmt.Errorf("user cancelled")
	}

	fmt.Println(output.InfoColor("\nStarting missing key addition..."))
	for path, missingKeys := range missingKeysMap {
		fmt.Printf(output.InfoColor(output.Bold("- Modifying file '%s'...\n")), path)
		err := ProcessFileWithNestedKeys(path, missingKeys)
		if err != nil {
			fmt.Printf("  "+output.ErrorColor("Error: %v")+"\n", err)
		}
	}

	return nil
}

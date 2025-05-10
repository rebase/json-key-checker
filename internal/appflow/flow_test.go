package appflow

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"

	"github.com/rebase/json-key-checker/internal/output"
)

const file1Content = `{
  "app": {
    "name": "json-checker",
    "version": "1.0"
  },
  "settings": {
    "timeout_seconds": 30,
    "retry_count": 5
  },
  "features": {
    "featureA": {
      "enabled": true,
      "details": {
        "level": "high",
        "threshold": 0.8
      }
    },
    "featureB": {
      "enabled": false
    }
  },
  "user_defaults": {
    "theme": "dark",
    "language": "en"
  }
}`

const file2Content = `{
  "app": {
    "name": "json-checker"
  },
  "features": {
    "featureA": {
      "details": {
        "level": "high"
      }
    }
  },
  "user_defaults": {
    "theme": "light"
  }
}`

const file3Content = `{
  "app": {
    "name": "json-checker",
    "version": "1.0"
  },
  "settings": {
    "timeout_seconds": 60
  },
  "user_defaults": {
    "theme": "dark",
    "language": "ko"
  }
}`

func TestMissingKeysIdentification(t *testing.T) {
	tempDir := t.TempDir()

	file1Path := filepath.Join(tempDir, "file1.json")
	file2Path := filepath.Join(tempDir, "file2.json")
	file3Path := filepath.Join(tempDir, "file3.json")

	if err := os.WriteFile(file1Path, []byte(file1Content), 0644); err != nil {
		t.Fatalf("Failed to write %s: %v", file1Path, err)
	}
	if err := os.WriteFile(file2Path, []byte(file2Content), 0644); err != nil {
		t.Fatalf("Failed to write %s: %v", file2Path, err)
	}
	if err := os.WriteFile(file3Path, []byte(file3Content), 0644); err != nil {
		t.Fatalf("Failed to write %s: %v", file3Path, err)
	}

	missingKeysMap, err := CompareJsonFiles([]string{file1Path, file2Path, file3Path})
	if err != nil {
		if err.Error() != "no missing keys found" {
			t.Fatalf("CompareJsonFiles failed: %v", err)
		}
	}

	if len(missingKeysMap[file1Path]) != 0 {
		t.Errorf("file1.json: Expected 0 missing keys, but found %d: %v", len(missingKeysMap[file1Path]), missingKeysMap[file1Path])
	}

	expectedMissingFile2 := []string{
		"app.version",
		"features.featureA.details.threshold",
		"features.featureA.enabled",
		"features.featureB",
		"features.featureB.enabled",
		"settings",
		"settings.retry_count",
		"settings.timeout_seconds",
		"user_defaults.language",
	}
	actualMissingFile2 := missingKeysMap[file2Path]
	sort.Strings(expectedMissingFile2)
	sort.Strings(actualMissingFile2)
	if !reflect.DeepEqual(actualMissingFile2, expectedMissingFile2) {
		t.Errorf("file2.json: Missing keys mismatch.\nExpected: %v\nActual:   %v", expectedMissingFile2, actualMissingFile2)
	} else {
		t.Logf("file2.json: Missing keys match expected: %v", expectedMissingFile2)
	}

	expectedMissingFile3 := []string{
		"features",
		"features.featureA",
		"features.featureA.details",
		"features.featureA.details.level",
		"features.featureA.details.threshold",
		"features.featureA.enabled",
		"features.featureB",
		"features.featureB.enabled",
		"settings.retry_count",
	}
	actualMissingFile3 := missingKeysMap[file3Path]
	sort.Strings(expectedMissingFile3)
	sort.Strings(actualMissingFile3)
	if !reflect.DeepEqual(actualMissingFile3, expectedMissingFile3) {
		t.Errorf("file3.json: Missing keys mismatch.\nExpected: %v\nActual:   %v", expectedMissingFile3, actualMissingFile3)
	} else {
		t.Logf("file3.json: Missing keys match expected: %v", expectedMissingFile3)
	}
}

func TestKeyAdditionEffect(t *testing.T) {
	tempDir := t.TempDir()

	file2Path := filepath.Join(tempDir, "file2_to_process.json")
	file3Path := filepath.Join(tempDir, "file3_to_process.json")

	initialFile2Content := file2Content
	missingKeysForFile2 := []string{
		"app.version",
		"features.featureA.details.threshold",
		"features.featureA.enabled",
		"features.featureB",
		"features.featureB.enabled",
		"settings",
		"settings.retry_count",
		"settings.timeout_seconds",
		"user_defaults.language",
	}
	expectedFile2Content := `{
  "app": {
    "name": "json-checker",
    "version": ""
  },
  "settings": {
    "timeout_seconds": "",
    "retry_count": ""
  },
  "features": {
    "featureA": {
      "details": {
        "level": "high"
      }
    },
    "featureA.details.threshold": "",
    "featureA.enabled": "",
    "featureB": "",
    "featureB.enabled": ""
  },
  "user_defaults": {
    "theme": "light",
    "language": ""
  }
}`

	initialFile3Content := file3Content
	missingKeysForFile3 := []string{
		"features",
		"features.featureA",
		"features.featureA.details",
		"features.featureA.details.level",
		"features.featureA.details.threshold",
		"features.featureA.enabled",
		"features.featureB",
		"features.featureB.enabled",
		"settings.retry_count",
	}
	expectedFile3Content := `{
  "app": {
    "name": "json-checker",
    "version": "1.0"
  },
  "settings": {
    "timeout_seconds": 60,
    "retry_count": ""
  },
  "features": {
    "featureA": "",
    "featureA.details": "",
    "featureA.details.level": "",
    "featureA.details.threshold": "",
    "featureA.enabled": "",
    "featureB": "",
    "featureB.enabled": ""
  },
  "user_defaults": {
    "language": "ko",
    "theme": "dark"
  }
}`

	if err := os.WriteFile(file2Path, []byte(initialFile2Content), 0644); err != nil {
		t.Fatalf("Failed to write %s: %v", file2Path, err)
	}

	if err := ProcessFileWithNestedKeys(file2Path, missingKeysForFile2); err != nil {
		t.Fatalf("ProcessFileWithNestedKeys failed for %s: %v", file2Path, err)
	}

	actualModifiedFile2Bytes, err := os.ReadFile(file2Path)
	if err != nil {
		t.Fatalf("Failed to read modified file %s: %v", file2Path, err)
	}

	var actualMap2, expectedMap2 map[string]interface{}
	if err := json.Unmarshal(actualModifiedFile2Bytes, &actualMap2); err != nil {
		t.Fatalf("Failed to parse modified JSON from %s: %v", file2Path, err)
	}
	if err := json.Unmarshal([]byte(expectedFile2Content), &expectedMap2); err != nil {
		t.Fatalf("Failed to parse expected JSON for file2: %v", err)
	}

	if !reflect.DeepEqual(actualMap2, expectedMap2) {
		actualJSON, _ := json.MarshalIndent(actualMap2, "", "  ")
		expectedJSON, _ := json.MarshalIndent(expectedMap2, "", "  ")
		t.Errorf("file2.json: Modified file content mismatch.\nExpected:\n%s\nActual:\n%s", string(expectedJSON), string(actualJSON))
	} else {
		t.Log(output.SuccessColor("file2.json: Modified file content matches expected."))
	}

	if err := os.WriteFile(file3Path, []byte(initialFile3Content), 0644); err != nil {
		t.Fatalf("Failed to write %s: %v", file3Path, err)
	}

	if err := ProcessFileWithNestedKeys(file3Path, missingKeysForFile3); err != nil {
		t.Fatalf("ProcessFileWithNestedKeys failed for %s: %v", file3Path, err)
	}

	actualModifiedFile3Bytes, err := os.ReadFile(file3Path)
	if err != nil {
		t.Fatalf("Failed to read modified file %s: %v", file3Path, err)
	}

	var actualMap3, expectedMap3 map[string]interface{}
	if err := json.Unmarshal(actualModifiedFile3Bytes, &actualMap3); err != nil {
		t.Fatalf("Failed to parse modified JSON from %s: %v", file3Path, err)
	}
	if err := json.Unmarshal([]byte(expectedFile3Content), &expectedMap3); err != nil {
		t.Fatalf("Failed to parse expected JSON for file3: %v", err)
	}

	if !reflect.DeepEqual(actualMap3, expectedMap3) {
		actualJSON, _ := json.MarshalIndent(actualMap3, "", "  ")
		expectedJSON, _ := json.MarshalIndent(expectedMap3, "", "  ")
		t.Errorf("file3.json: Modified file content mismatch.\nExpected:\n%s\nActual:\n%s", string(expectedJSON), string(actualJSON))
	} else {
		t.Log(output.SuccessColor("file3.json: Modified file content matches expected."))
	}
}

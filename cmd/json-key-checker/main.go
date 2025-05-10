package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/rebase/json-key-checker/internal/appflow"
	"github.com/rebase/json-key-checker/internal/output"
)

func main() {
	fmt.Println(output.InfoColor(output.Bold("Starting JSON Key Checker!")))

	flag.Parse()
	filePaths := flag.Args()

	if len(filePaths) < 1 {
		fmt.Println(output.ErrorColor(output.Bold("Error: Please provide at least one JSON file path.")))
		fmt.Println("Usage: " + output.Bold("json-key-checker <file1.json> <file2.json> ..."))
		os.Exit(1)
	}

	missingKeysMap, err := appflow.CompareJsonFiles(filePaths)
	if err != nil {
		if err.Error() == "no missing keys found" {
			return
		}
		fmt.Printf(output.ErrorColor(output.Bold("Error: Unexpected issue during file comparison: %v"))+"\n", err)
		os.Exit(1)
	}

	err = appflow.PromptAndProcessFiles(missingKeysMap)
	if err != nil {
		if err.Error() == "user cancelled" {
			return
		}
		fmt.Printf(output.ErrorColor(output.Bold("Error: Unexpected issue during file processing: %v"))+"\n", err)
		os.Exit(1)
	}

	fmt.Println(output.SuccessColor(output.Bold("\nAll tasks completed.")))
}

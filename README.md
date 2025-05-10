# JSON Key Checker

A command-line tool to compare JSON files and add missing keys based on a union of all keys found.

## Description

This tool reads multiple JSON files, identifies all unique keys across all files (including nested keys using dot notation), and then checks each file for missing keys compared to the complete set. Optionally, it can add the missing keys to the files with a default empty string value.

## Installation

### Via Homebrew

```bash
brew install rebase/rebase/json-key-checker
```

## Usage

```bash
json-key-checker file1.json file2.json file3.json ...
```

## Use Case: i18n Localization

This tool is especially useful for managing i18n (internationalization) JSON translation files.

For example, when you have:

- `en.json` (English)
- `ko.json` (Korean)
- `ja.json` (Japanese)

You can run this tool to find keys that exist in one language file but are missing in others, and automatically fill them with empty string placeholders.

**Given the following translation files:**

```json
// en.json

{
  "animal": {
    "cat": "Cat",
    "dog": "Dog"
  },
  "human": {
    "age": "Age",
    "name": "Name"
  }
}
```

```json
// ko.json

{
  "animal": {
    "cat": "고양이",
    "dog": "개"
  },
  "human": {
    "name": "이름"
  }
}
```

```json
// ja.json

{
  "animal": {
    "dog": "犬"
  }
}
```

**Run the command:**

```bash
json-key-checker en.json ko.json ja.json
```

**Terminal output:**

```bash
Starting JSON Key Checker!
Analyzing files...
- 'en.json' processed (6 keys)
- 'ko.json' processed (5 keys)
- 'ja.json' processed (2 keys)

--- Missing Keys Comparison ---
File 'ja.json' has missing keys (4 total):
  - animal.cat
  - human
  - human.age
  - human.name
File 'en.json' has no missing keys.
File 'ko.json' has missing keys (1 total):
  - human.age

Add missing keys to each file? (y/N): y

Starting missing key addition...
- Modifying file 'ja.json'...
  Success: Added 4 keys to 'ja.json'.
- Modifying file 'ko.json'...
  Success: Added 1 keys to 'ko.json'.

All tasks completed.
```

**Updated files:**

```json
// ko.json

{
  "animal": {
    "cat": "고양이",
    "dog": "개"
  },
  "human": {
    "age": "",
    "name": "이름"
  }
}
```

```json
// ja.json

{
  "animal": {
    "cat": "",
    "dog": "犬"
  },
  "human": {
    "age": "",
    "name": ""
  }
}
```

Now all translation files contain a consistent set of keys — all that’s left is to fill in the missing values in the "" placeholders.

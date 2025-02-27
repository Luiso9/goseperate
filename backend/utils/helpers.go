package utils

import (
  "os"
  "fmt"

  "path/filepath"
)

func RemoveTempFiles(paths ...string) {
	for _, path := range paths {
		os.RemoveAll(path)
	}
}

func FindExistingFile(dir, filename string) (string, error) {
  pattern := fmt.Sprintf("%s/%s.*", dir, filename)
  matches, err := filepath.Glob(pattern)
  if err != nil || len(matches) == 0 {
    return "", fmt.Errorf("file not found")
  }
  return matches[0], nil
}

package pkg

import (
	"os"
	"path/filepath"
)

// WriteFile writes data to the specified path, creating parent directories if necessary.
func WriteFile(path string, data []byte) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func DeleteDir(path string) error {
	return os.RemoveAll(path)
}

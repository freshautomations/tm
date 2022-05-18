package utils

import (
	"fmt"
	"path/filepath"
)

// GetSlashPath replaces "/" with OS file separator character
func GetSlashPath(format string, args ...any) string {
	return filepath.FromSlash(fmt.Sprintf(format, args...))
}

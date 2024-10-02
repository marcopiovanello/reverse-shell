package paths

import "strings"

func WindowsBase(path string) string {
	if strings.ContainsRune(path, '\\') {
		parts := strings.Split(path, "\\")
		return parts[len(parts)-1]
	}
	return path
}

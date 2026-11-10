package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func FolderExistsInPwd(folderName string) bool {
	info, err := os.Stat(folderName)
	return err == nil && info.IsDir()
}

func FileExists(filepath string) bool {
	info, err := os.Stat(filepath)
	return err == nil && info.Mode().IsRegular()
}

func EnsureUniqueFilename(fullPath string) string {
	if !FileExists(fullPath) {
		return fullPath
	}

	dir := filepath.Dir(fullPath)
	ext := filepath.Ext(fullPath)
	baseName := strings.TrimSuffix(filepath.Base(fullPath), ext)

	counter := 1
	for {
		newBaseName := fmt.Sprintf("%s-%d", baseName, counter)
		newPath := filepath.Join(dir, newBaseName+ext)

		if !FileExists(newPath) {
			return newPath
		}
		counter++
	}
}

var nonHandleChars = regexp.MustCompile(`[^a-z0-9_]+`)

// Handlize converts a display name into a safe file handle: lowercase letters,
// digits and underscores, with any other runs of characters collapsed to "-".
// "My Hero / Banner!" becomes "my-hero-banner".
func Handlize(name string) string {
	handle := nonHandleChars.ReplaceAllString(strings.ToLower(name), "-")
	return strings.Trim(handle, "-")
}

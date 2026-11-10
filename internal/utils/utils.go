package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func FolderExistsInPwd(folderName string) bool {
	info, err := os.Stat(folderName)
	if os.IsNotExist(err) {
		return false
	}
	if err != nil {
		log.Fatalf("folderExistsInPwd: %s", err)
	}
	return info.IsDir()
}

func FileExists(filepath string) bool {
	info, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		return false
	}
	if err != nil {
		log.Fatalf("fileExists: %s", err)
	}
	return info.Mode().IsRegular()
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

func Handlize(name string) string {
	cName := strings.ToLower(strings.ReplaceAll(name, " ", "-"))
	return cName
}

// Package theme writes generated files into a Shopify theme directory.
package theme

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/guitarkeegan/shop-cli/internal/shopify"
)

var (
	ErrMissingName  = errors.New("a name is required (use -n)")
	ErrInvalidName  = errors.New("name must contain at least one letter or number")
	ErrNotThemeRoot = errors.New("run this from the root of a Shopify theme")
)

// WriteSection writes sections/<handle>.liquid and returns its path.
func WriteSection(name, content string) (string, error) {
	return write("sections", name, content)
}

// WriteBlock writes blocks/<handle>.liquid and returns its path.
func WriteBlock(name, content string) (string, error) {
	return write("blocks", name, content)
}

// WriteSnippet writes snippets/<handle>.liquid and returns its path.
func WriteSnippet(name, content string) (string, error) {
	return write("snippets", name, content)
}

// write saves content to dir/<handle>.liquid in the working directory. It
// never overwrites: an existing file gets a numeric suffix instead.
func write(dir, name, content string) (string, error) {
	if strings.TrimSpace(name) == "" {
		return "", ErrMissingName
	}
	handle := shopify.Handle(name)
	if handle == "" {
		return "", fmt.Errorf("%w: %q", ErrInvalidName, name)
	}
	if info, err := os.Stat(dir); err != nil || !info.IsDir() {
		return "", fmt.Errorf("no %q directory found; %w", dir, ErrNotThemeRoot)
	}

	if content != "" && !strings.HasSuffix(content, "\n") {
		content += "\n"
	}

	path := uniquePath(filepath.Join(dir, handle+".liquid"))
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return "", fmt.Errorf("write %s: %w", path, err)
	}
	return path, nil
}

// uniquePath returns path, or path with "-1", "-2", ... before the extension
// if a file already exists there.
func uniquePath(path string) string {
	if !exists(path) {
		return path
	}
	ext := filepath.Ext(path)
	base := strings.TrimSuffix(path, ext)
	for i := 1; ; i++ {
		candidate := fmt.Sprintf("%s-%d%s", base, i, ext)
		if !exists(candidate) {
			return candidate
		}
	}
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

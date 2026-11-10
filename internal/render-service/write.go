package renderservice

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/guitarkeegan/shop-cli/internal/utils"
)

// writeLiquid writes content to dir/<handle>.liquid, adding a numeric suffix
// instead of overwriting an existing file. It returns the path written.
func writeLiquid(dir, name, content string) (string, error) {
	handle := utils.Handlize(name)
	if handle == "" {
		return "", fmt.Errorf("%w: %q", ErrInvalidName, name)
	}

	if !strings.HasSuffix(content, "\n") && content != "" {
		content += "\n"
	}

	path := utils.EnsureUniqueFilename(filepath.Join(dir, handle+".liquid"))
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return "", fmt.Errorf("write %s: %w", path, err)
	}
	return path, nil
}

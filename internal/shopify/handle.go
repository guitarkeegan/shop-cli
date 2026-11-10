package shopify

import (
	"regexp"
	"strings"
)

var nonHandleChars = regexp.MustCompile(`[^a-z0-9_]+`)

// Handle converts a display name into a Shopify-style handle: lowercase
// letters, digits and underscores, with any other runs of characters
// collapsed to "-". "My Hero / Banner!" becomes "my-hero-banner".
func Handle(name string) string {
	handle := nonHandleChars.ReplaceAllString(strings.ToLower(name), "-")
	return strings.Trim(handle, "-")
}

package shopify

import (
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

var (
	nonHandleChars = regexp.MustCompile(`[^a-z0-9_]+`)
	stripAccents   = transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
)

// Handle converts a display name into a Shopify-style handle: accents are
// removed, then anything other than lowercase letters, digits and
// underscores collapses to "-". "Café / Banner!" becomes "cafe-banner".
func Handle(name string) string {
	if folded, _, err := transform.String(stripAccents, name); err == nil {
		name = folded
	}
	handle := nonHandleChars.ReplaceAllString(strings.ToLower(name), "-")
	return strings.Trim(handle, "-")
}

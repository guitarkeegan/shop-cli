package shopify

import "testing"

func TestHandle(t *testing.T) {
	tests := map[string]string{
		"Hero":              "hero",
		"Hero Banner":       "hero-banner",
		"My Hero / Banner!": "my-hero-banner",
		"  padded  ":        "padded",
		"snake_case_name":   "snake_case_name",
		"Ünïcode Café":      "unicode-cafe",
		"日本":                "",
		"a/b":               "a-b",
		"../../etc":         "etc",
		"!!!":               "",
		"":                  "",
	}
	for in, want := range tests {
		if got := Handle(in); got != want {
			t.Errorf("Handle(%q) = %q, want %q", in, got, want)
		}
	}
}

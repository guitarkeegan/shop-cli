package theme

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

// inTheme runs the test from a temporary directory containing dirs.
func inTheme(t *testing.T, dirs ...string) {
	t.Helper()
	root := t.TempDir()
	for _, d := range dirs {
		if err := os.Mkdir(filepath.Join(root, d), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	t.Chdir(root)
}

func read(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

func TestWriteSection(t *testing.T) {
	inTheme(t, "sections")

	path, err := WriteSection("Hero Banner", "content")
	if err != nil {
		t.Fatal(err)
	}
	if want := filepath.Join("sections", "hero-banner.liquid"); path != want {
		t.Errorf("path = %q, want %q", path, want)
	}
	if got := read(t, path); got != "content\n" {
		t.Errorf("content = %q, want a trailing newline added", got)
	}
}

func TestWriteNeverOverwrites(t *testing.T) {
	inTheme(t, "snippets")

	var paths []string
	for _, content := range []string{"first\n", "second\n", "third\n"} {
		path, err := WriteSnippet("card", content)
		if err != nil {
			t.Fatal(err)
		}
		paths = append(paths, path)
	}

	want := []string{"card.liquid", "card-1.liquid", "card-2.liquid"}
	for i, p := range paths {
		if filepath.Base(p) != want[i] {
			t.Errorf("write %d went to %q, want %q", i, p, want[i])
		}
	}
	if got := read(t, paths[0]); got != "first\n" {
		t.Errorf("first file was overwritten: %q", got)
	}
}

func TestWriteErrors(t *testing.T) {
	tests := []struct {
		name    string
		dirs    []string
		write   func(name, content string) (string, error)
		file    string
		wantErr error
	}{
		{"missing name", []string{"blocks"}, WriteBlock, "", ErrMissingName},
		{"blank name", []string{"blocks"}, WriteBlock, "   ", ErrMissingName},
		{"no usable characters", []string{"blocks"}, WriteBlock, "!!!", ErrInvalidName},
		{"not a theme root", nil, WriteSection, "Hero", ErrNotThemeRoot},
		{"wrong folder", []string{"sections"}, WriteBlock, "Promo", ErrNotThemeRoot},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inTheme(t, tt.dirs...)
			if _, err := tt.write(tt.file, "x"); !errors.Is(err, tt.wantErr) {
				t.Errorf("err = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestEmptyContentStaysEmpty(t *testing.T) {
	inTheme(t, "snippets")
	path, err := WriteSnippet("blank", "")
	if err != nil {
		t.Fatal(err)
	}
	if got := read(t, path); got != "" {
		t.Errorf("content = %q, want empty", got)
	}
}

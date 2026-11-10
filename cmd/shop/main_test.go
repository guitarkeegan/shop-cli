package main

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/guitarkeegan/shop-cli/internal/render"
	"github.com/guitarkeegan/shop-cli/internal/theme"
)

// newTheme runs the test from an empty Shopify theme directory.
func newTheme(t *testing.T) {
	t.Helper()
	root := t.TempDir()
	for _, d := range []string{"sections", "blocks", "snippets"} {
		if err := os.Mkdir(filepath.Join(root, d), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	t.Chdir(root)
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

func TestCreateSection(t *testing.T) {
	newTheme(t)

	err := createSection(&sectionCreateOptions{
		Name:     "Hero Banner",
		Settings: "text:title:Title",
		Blocks:   []string{"Slide,image_picker:image:Image"},
	}, nil)(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	got := readFile(t, filepath.Join("sections", "hero-banner.liquid"))
	for _, want := range []string{
		"assign title = section.settings.title",
		"{% when 'slide' %}",
		"assign image = block.settings.image",
		`"id": "color_scheme"`,
		`"id": "padding_top"`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("section is missing %q:\n%s", want, got)
		}
	}
}

func TestCreateBlock(t *testing.T) {
	newTheme(t)

	if err := createBlock(&blockCreateOptions{Name: "Promo", Settings: "url:link:Link"})(context.Background()); err != nil {
		t.Fatal(err)
	}
	if got := readFile(t, filepath.Join("blocks", "promo.liquid")); !strings.Contains(got, "assign link = block.settings.link") {
		t.Errorf("unexpected block:\n%s", got)
	}
}

func TestCreateSnippet(t *testing.T) {
	newTheme(t)

	err := createSnippet(&snippetCreateOptions{
		Name:        "Product Card",
		Params:      "product:object,show_price?:boolean",
		Description: "Renders a product card",
	})(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	got := readFile(t, filepath.Join("snippets", "product-card.liquid"))
	for _, want := range []string{"@param {object} product", "@param {boolean} [show_price]", "Renders a product card"} {
		if !strings.Contains(got, want) {
			t.Errorf("snippet is missing %q:\n%s", want, got)
		}
	}
}

func TestCommandErrors(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name      string
		run       func(context.Context) error
		inTheme   bool
		wantErr   error
		wantUsage bool
	}{
		{"section without name", createSection(&sectionCreateOptions{}, nil), true, theme.ErrMissingName, true},
		{"bad setting", createSection(&sectionCreateOptions{Name: "Hero", Settings: "txt:a:A"}, nil), true, render.ErrInvalidSetting, true},
		{"bad block", createSection(&sectionCreateOptions{Name: "Hero", Blocks: []string{"!!!"}}, nil), true, render.ErrInvalidBlock, true},
		{"outside a theme", createBlock(&blockCreateOptions{Name: "Promo"}), false, theme.ErrNotThemeRoot, false},
		{"parent command only", requireSubcommand, true, ErrMissingSubcommand, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.inTheme {
				newTheme(t)
			} else {
				t.Chdir(t.TempDir())
			}
			err := tt.run(ctx)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if got := isUsageError(err); got != tt.wantUsage {
				t.Errorf("isUsageError = %v, want %v", got, tt.wantUsage)
			}
		})
	}
}

func TestExamplesConfigFromEnv(t *testing.T) {
	for _, env := range []string{envGithubToken, envGithubOwner, envGithubRepo, envExamplesPath} {
		t.Setenv(env, "")
	}
	if _, err := examplesConfigFromEnv(); !errors.Is(err, ErrMissingEnvVars) {
		t.Fatalf("err = %v, want %v", err, ErrMissingEnvVars)
	}

	t.Setenv(envGithubToken, "token")
	t.Setenv(envGithubOwner, "owner")
	t.Setenv(envGithubRepo, "repo")
	t.Setenv(envExamplesPath, "examples")
	cfg, err := examplesConfigFromEnv()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Owner != "owner" || cfg.Repo != "repo" || cfg.FolderPath != "examples" {
		t.Errorf("cfg = %+v", cfg)
	}
}

package render

import (
	"flag"
	"os"
	"path/filepath"
	"testing"

	"github.com/guitarkeegan/shop-cli/internal/shopify"
)

var update = flag.Bool("update", false, "rewrite golden files in testdata/")

// golden compares got with testdata/<name>.liquid, or rewrites the file
// when run with -update.
func golden(t *testing.T, name, got string) {
	t.Helper()
	path := filepath.Join("testdata", name+".liquid")
	if *update {
		if err := os.WriteFile(path, []byte(got), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	want, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read golden file (run with -update to create it): %v", err)
	}
	if got != string(want) {
		t.Errorf("output does not match %s (run with -update to accept):\n--- got ---\n%s\n--- want ---\n%s", path, got, want)
	}
}

func mustParseSettings(t *testing.T, list string) []shopify.Setting {
	t.Helper()
	settings, err := ParseSettings(list)
	if err != nil {
		t.Fatal(err)
	}
	return settings
}

func TestSection(t *testing.T) {
	blocks, err := ParseBlocks([]string{
		"Slide,image_picker:image:Image,text:caption:Caption",
		"Card,text:title:Title",
	})
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name string
		opts []SectionOption
	}{
		{"section_minimal", []SectionOption{WithSectionName("Plain")}},
		{"section_settings", []SectionOption{
			WithSectionName("Hero Banner"),
			WithSectionSettings(mustParseSettings(t, "text:title:Title,richtext:copy:Copy,range:opacity:Opacity")),
		}},
		{"section_blocks", []SectionOption{
			WithSectionName("Slideshow"),
			WithSectionSettings(mustParseSettings(t, "text:heading:Heading")),
			WithSectionBlocks(blocks),
		}},
		{"section_defaults", []SectionOption{
			WithSectionName("Feature"),
			WithSectionSettings(mustParseSettings(t, "header:x:Content,text:title:Title")),
			WithDefaults(),
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Section(tt.opts...)
			if err != nil {
				t.Fatal(err)
			}
			golden(t, tt.name, got)
		})
	}
}

func TestBlock(t *testing.T) {
	got, err := Block(
		WithBlockName("Promo"),
		WithBlockSettings(mustParseSettings(t, "text:heading:Heading,url:link:Link")),
	)
	if err != nil {
		t.Fatal(err)
	}
	golden(t, "block", got)
}

func TestSnippet(t *testing.T) {
	tests := []struct {
		name string
		opts []SnippetOption
	}{
		{"snippet_empty", nil},
		{"snippet_doc", []SnippetOption{
			WithSnippetDescription("Renders a product card"),
			WithSnippetParams(ParseSnippetParams("product:object,show_price?:boolean,title")),
			WithSnippetExample("{% render 'card', product: product %}"),
		}},
		{"snippet_body", []SnippetOption{
			WithSnippetParams(ParseSnippetParams("title")),
			WithSnippetBody("<h2>{{ title }}</h2>\n"),
		}},
		{"snippet_param_descriptions", []SnippetOption{
			WithSnippetParams([]shopify.Param{
				{Name: "title", Type: shopify.StringParam, Required: true, Description: "Heading text"},
				{Name: "count", Type: shopify.NumberParam, Description: "How many to show"},
			}),
			WithSnippetExample("   "),
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			golden(t, tt.name, Snippet(tt.opts...))
		})
	}
}

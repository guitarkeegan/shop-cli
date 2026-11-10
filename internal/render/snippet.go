package render

import (
	"strings"

	"github.com/guitarkeegan/shop-cli/internal/shopify"
)

type snippet struct {
	doc  shopify.Snippet
	body string
}

// SnippetOption configures a rendered snippet.
type SnippetOption func(*snippet)

func WithSnippetDescription(description string) SnippetOption {
	return func(s *snippet) { s.doc.Description = description }
}

func WithSnippetParams(params []shopify.Param) SnippetOption {
	return func(s *snippet) { s.doc.Params = params }
}

// WithSnippetExample adds an @example. Blank examples are ignored.
func WithSnippetExample(example string) SnippetOption {
	return func(s *snippet) {
		if strings.TrimSpace(example) != "" {
			s.doc.Examples = append(s.doc.Examples, example)
		}
	}
}

// WithSnippetBody sets the Liquid that follows the doc block, such as the
// contents of a template.
func WithSnippetBody(body string) SnippetOption {
	return func(s *snippet) { s.body = body }
}

// Snippet renders a snippet file: a {% doc %} block, if there is anything to
// document, followed by the body.
func Snippet(opts ...SnippetOption) string {
	var s snippet
	for _, opt := range opts {
		opt(&s)
	}
	return docBlock(s.doc) + s.body
}

// docBlock renders a LiquidDoc {% doc %} block. It returns "" when there is
// nothing to document.
func docBlock(sn shopify.Snippet) string {
	if sn.DocText == "" && sn.Description == "" && len(sn.Params) == 0 && len(sn.Examples) == 0 {
		return ""
	}

	var doc strings.Builder
	doc.WriteString("{% doc %}\n")

	if sn.DocText != "" {
		doc.WriteString(sn.DocText + "\n\n")
	}

	for _, p := range sn.Params {
		name := p.Name
		if !p.Required {
			name = "[" + name + "]"
		}
		doc.WriteString("@param {" + p.Type.String() + "} " + name)
		if p.Description != "" {
			doc.WriteString(" - " + p.Description)
		}
		doc.WriteByte('\n')
	}

	for _, example := range sn.Examples {
		doc.WriteString("\n@example\n" + example + "\n")
	}

	if sn.Description != "" {
		doc.WriteString("\n@description - " + sn.Description + "\n")
	}

	doc.WriteString("{% enddoc %}\n")
	return doc.String()
}

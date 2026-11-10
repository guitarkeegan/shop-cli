package renderservice

import (
	"context"
	"fmt"
	"os"
	"strings"

	shop2 "github.com/guitarkeegan/shop-cli/internal/shopify-v2"
	"github.com/guitarkeegan/shop-cli/internal/tui"
	"github.com/guitarkeegan/shop-cli/internal/utils"
)

var dbg = func() func(format string, as ...any) {
	if os.Getenv("DEBUG") == "" {
		return func(format string, as ...any) {}
	}
	return func(format string, as ...any) {
		fmt.Printf(format+"\n", as...)
	}
}()

type SnippetService struct {
	FileName    string
	UseTemplate bool
	Snippet     shop2.Snippet
}

// String renders the snippet's {% doc %} block. It returns an empty string
// when there is nothing to document.
func (s *SnippetService) String() string {
	sn := s.Snippet
	if sn.DocText == "" && sn.Description == "" && len(sn.Params) == 0 && len(sn.Examples) == 0 {
		return ""
	}

	var doc strings.Builder
	doc.WriteString("{% doc %}\n")

	if sn.DocText != "" {
		doc.WriteString(sn.DocText)
		doc.WriteString("\n\n")
	}

	for _, param := range sn.Params {
		doc.WriteString("@param {")
		doc.WriteString(param.Type.String())
		doc.WriteString("} ")
		if param.Required {
			doc.WriteString(param.Name)
		} else {
			doc.WriteString("[" + param.Name + "]")
		}
		if param.Description != "" {
			doc.WriteString(" - ")
			doc.WriteString(param.Description)
		}
		doc.WriteByte('\n')
	}

	for _, example := range sn.Examples {
		doc.WriteString("\n@example\n")
		doc.WriteString(example)
		doc.WriteByte('\n')
	}

	if sn.Description != "" {
		doc.WriteString("\n@description - ")
		doc.WriteString(sn.Description)
		doc.WriteByte('\n')
	}

	doc.WriteString("{% enddoc %}\n")
	return doc.String()
}

func NewSnippetService() *SnippetService {
	return &SnippetService{}
}

type SnippetOption func(*SnippetService)

func WithSnippetDocText(docText string) SnippetOption {
	return func(s *SnippetService) {
		s.Snippet.DocText = docText
	}
}

func WithSnippetDescription(description string) SnippetOption {
	return func(s *SnippetService) {
		s.Snippet.Description = description
	}
}

func WithSnippetParams(p []shop2.Param) SnippetOption {
	return func(s *SnippetService) {
		s.Snippet.Params = p
	}
}

func WithSnippetAddParam(p shop2.Param) SnippetOption {
	return func(s *SnippetService) {
		s.Snippet.Params = append(s.Snippet.Params, p)
	}
}

func WithSnippetExamples(examples []string) SnippetOption {
	return func(s *SnippetService) {
		s.Snippet.Examples = examples
	}
}

func WithSnippetTemplate(useTemplate bool) SnippetOption {
	return func(s *SnippetService) {
		s.UseTemplate = useTemplate
	}
}

// WithSnippetAddExample adds an @example. Blank examples are ignored.
func WithSnippetAddExample(example string) SnippetOption {
	return func(s *SnippetService) {
		if strings.TrimSpace(example) == "" {
			return
		}
		s.Snippet.Examples = append(s.Snippet.Examples, example)
	}
}

// this one needs to be required
func WithSnippetFileName(fileName string) SnippetOption {
	return func(s *SnippetService) {
		s.FileName = fileName
	}
}

func parseRequired(arg string) (string, bool) {
	switch true {
	case arg == "":
		return arg, true
	case arg[len(arg)-1] == '?':
		return arg[:len(arg)-1], false
	default:
		return arg, true
	}
}

// ExtractShopifyParams parses CLI params like "title:string,image?:object,copy".
// A trailing "?" on the name marks the param optional; the type defaults to string.
func (s *SnippetService) ExtractShopifyParams(cliArgs string) []shop2.Param {
	var params []shop2.Param

	for param := range strings.SplitSeq(cliArgs, ",") {
		trimmed := strings.TrimSpace(param)
		if trimmed == "" {
			continue
		}
		variable, typ, _ := strings.Cut(trimmed, ":")
		name, required := parseRequired(strings.TrimSpace(variable))
		params = append(params, shop2.Param{
			Type:     shop2.ParamTypeFromString(strings.TrimSpace(typ)),
			Name:     name,
			Required: required,
		})
	}
	return params
}

// use this to pass in the command line args
func WithSnippetParamsAsString(params string) SnippetOption {
	return func(s *SnippetService) {
		s.Snippet.Params = s.ExtractShopifyParams(params)
	}
}

// CreateSnippet renders a snippet and writes it to snippets/<handle>.liquid.
// With a template, the picked file's contents follow the doc block.
// It returns the path of the created file.
func CreateSnippet(ctx context.Context, opts ...SnippetOption) (string, error) {
	snip := NewSnippetService()
	for _, opt := range opts {
		opt(snip)
	}

	if snip.FileName == "" {
		return "", ErrMissingName
	}
	if !utils.FolderExistsInPwd("snippets") {
		return "", ErrMissingSnippetsDir
	}

	content := snip.String()

	if snip.UseTemplate {
		tempPath, err := tui.Run()
		if err != nil {
			return "", err
		}
		if tempPath != "" {
			tFile, err := os.ReadFile(tempPath)
			if err != nil {
				return "", fmt.Errorf("read template %s: %w", tempPath, err)
			}
			content += string(tFile) + "\n"
		}
	}

	return writeLiquid("snippets", snip.FileName, content)
}

// CreateSnippetPreview renders a snippet without writing it. Use this within a TUI.
func CreateSnippetPreview(ctx context.Context, opts ...SnippetOption) (string, error) {
	snip := NewSnippetService()
	for _, opt := range opts {
		opt(snip)
	}

	if snip.FileName == "" {
		return "", ErrMissingName
	}

	return snip.String(), nil
}

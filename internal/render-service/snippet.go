package renderservice

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/guitarkeegan/shop-cli/internal/colors"
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

var (
	missingSnippetDirErr = errors.New("No 'snippets' dir found. Are you at the top level of the project?")
)

type SnippetService struct {
	FileName    string
	UseTemplate bool
	Snippet     shop2.Snippet
}

// Return the string representation fo the SnippetService
// not handling errors currently
func (s *SnippetService) String() string {
	var doc strings.Builder

	doc.WriteString("{% doc %}\n")
	doc.WriteString(s.Snippet.DocText)
	doc.WriteString("\n\n")

	for _, param := range s.Snippet.Params {
		doc.WriteString("@param {")
		doc.WriteString(shop2.ParamTypeToString[int(param.Type)])
		doc.WriteString("} ")

		if param.Required {
			doc.WriteString(param.Name)
			doc.WriteString(" - ")
			doc.WriteString(param.Description)
		} else {
			doc.WriteByte('[')
			doc.WriteString(param.Name)
			doc.WriteByte(']')
			doc.WriteString(" - ")
			doc.WriteString(param.Description)
		}
		doc.WriteByte('\n')
	}

	for _, example := range s.Snippet.Examples {
		doc.WriteString("@example\n")
		doc.WriteString(example)
		doc.WriteString("\n\n")
	}

	if s.Snippet.Description != "" {
		doc.WriteString("@description - ")
		doc.WriteString(s.Snippet.Description)
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

func WithSnippetAddExample(example string) SnippetOption {
	return func(s *SnippetService) {
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

// Use this to transform the string passed via the cli version
func (s *SnippetService) ExtractShopifyParams(cliArgs string) []shop2.Param {
	paramSeq := strings.SplitSeq(cliArgs, ",")
	// wild guess for capacity
	params := make([]shop2.Param, 0, 50)

	for param := range paramSeq {
		trimmed := strings.TrimSpace(param)
		variable, typ, hasColon := strings.Cut(trimmed, ":")
		trimmedVariable, required := parseRequired(variable)
		if hasColon {
			params = append(params, shop2.Param{
				Type:     shop2.ParamTypeFromString(typ),
				Name:     trimmedVariable,
				Required: required,
			})
		} else {
			trimmedVariable, required = parseRequired(param)
			params = append(params, shop2.Param{
				Type:        shop2.ParamTypeFromString(""),
				Name:        trimmedVariable,
				Required:    required,
				Description: "",
			})
		}
	}
	return params

}

// use this to pass in the command line args
func WithSnippetParamsAsString(params string) SnippetOption {
	return func(s *SnippetService) {
		s.Snippet.Params = s.ExtractShopifyParams(params)
	}
}

var CreateSnippet = func(ctx context.Context, opts ...SnippetOption) func(ctx context.Context) error {
	return func(ctx context.Context) error {

		if !utils.FolderExistsInPwd("snippets") {
			return missingSnippetDirErr
		}

		snip := NewSnippetService()
		for _, opt := range opts {
			opt(snip)
		}

		if snip.FileName == "" {
			return missingNameErr
		}

		path := filepath.Join("snippets", utils.Handlize(snip.FileName)+".liquid")
		if utils.FileExists(path) {
			path = utils.EnsureUniqueFilename(path)
		}

		var liquidTemplate string
		if snip.UseTemplate {
			// returns abs path
			tempPath, err := tui.Run()
			if err != nil {
				return err
			}
			fmt.Println(colors.BoldGreen(tempPath))
			if tempPath != "" {
				tFile, err := os.ReadFile(tempPath)
				if err != nil {
					panic("could not read from template!")
				}
				liquidTemplate = string(tFile)
			}
		}

		if snip.Snippet.Params == nil {
			if len(liquidTemplate) > 0 {
				if err := os.WriteFile(path, []byte(liquidTemplate), 0o644); err != nil {
					return fmt.Errorf("handleCreateSnippet: write file: %w", err)
				}

				fmt.Println(colors.Success("new snippet created!"))
				return nil
			}
			if err := os.WriteFile(path, nil, 0o644); err != nil {
				return fmt.Errorf("handleCreateSnippet: write file: %w", err)
			}
			fmt.Println(colors.Success("New Snippet Created!"))
			fmt.Println(colors.Success(snip.FileName))
			return nil
		}

		if err := os.WriteFile(path, []byte(snip.String()+liquidTemplate+"\n"), 0o644); err != nil {
			return fmt.Errorf("handleCreateSnippet: write file with params: %w", err)
		}

		fmt.Println(colors.Success("New Snippet Created!"))
		fmt.Println(colors.Success(snip.FileName))
		return nil
	}
}

var CreateSnippetPreview = func(ctx context.Context, opts ...SnippetOption) (string, error) {

	snip := NewSnippetService()
	for _, opt := range opts {
		opt(snip)
	}

	if snip.FileName == "" {
		return "", missingNameErr
	}

	// var liquidTemplate string
	// if snip.UseTemplate {
	// 	// returns abs path
	// 	tempPath, err := tui.Run()
	// 	if err != nil {
	// 		return "", err
	// 	}
	// 	fmt.Println(colors.BoldGreen(tempPath))
	// 	if tempPath != "" {
	// 		tFile, err := os.ReadFile(tempPath)
	// 		if err != nil {
	// 			return "", err
	// 		}
	// 		liquidTemplate = string(tFile)
	// 	}
	// }

	return snip.String(), nil
}

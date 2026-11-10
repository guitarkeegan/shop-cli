package fulltui

import (
	"errors"
	"regexp"

	"github.com/charmbracelet/huh"
)

const (
	typeKey        = "typeKey"
	textKey        = "textKey"
	optionalKey    = "optionalKey"
	descriptionKey = "descriptionKey"
	exampleKey     = "exampleKey"
	nameKey        = "nameKey"
	actionKey      = "actionKey"
	responseKey    = "responseKey"
	inputTypeKey   = "inputTypeKey"
	idKey          = "idKey"
	labelKey       = "labelKey"
	continueKey    = "continueKey"
)

var (
	errSnakeCase = errors.New("invalid liquid format, must be snake_case")
	errEmpty     = errors.New("field cannot be empty")
)

var snakeCase = regexp.MustCompile("^[a-z_][a-z0-9_]*$")

func snakeCaseValidator(s string) error {
	if !snakeCase.MatchString(s) {
		return errSnakeCase
	}
	return nil
}

func notEmptyValidator(s string) error {
	if s == "" {
		return errEmpty
	}
	return nil
}

// User is asked to select the file type they wish to work with, give it a name
// and select an action.
func newFileSelector() *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Key(typeKey).
				Title("Which file-type are we working with?").
				Options(
					huh.NewOption("Section", "section"),
					huh.NewOption("Snippet", "snippet"),
					huh.NewOption("Block", "block"),
				),
			huh.NewInput().
				Key(nameKey).
				Title("Filename?").
				Validate(notEmptyValidator),
		),
		huh.NewGroup(
			huh.NewSelect[string]().
				Key(actionKey).
				Title("Action?").
				Options(
					huh.NewOption("Create", "create"),
				),
		),
	).WithTheme(ShopTheme())
}

// newSnippetDocForm asks for a snippet doc param. The first form (more is
// false) also asks whether to add params at all, and for the description.
func newSnippetDocForm(more bool) *huh.Form {
	wantsParams := more
	wantsMoreParams := more

	return huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Key(responseKey).
				Title("Add doc params?").
				Description("We can define params here that the caller is able to pass in to the snippet.").
				Value(&wantsParams),
		).WithHideFunc(func() bool {
			return wantsParams
		}),
		huh.NewGroup(
			huh.NewText().
				Key(descriptionKey).
				Title("Description").
				Description("Optionally, provide a description for the snippet."),
			huh.NewText().
				Key(exampleKey).
				Title("Example").
				Description("Optionally, provide an example usage."),
		).WithHideFunc(func() bool {
			return !wantsParams || wantsMoreParams
		}),
		huh.NewGroup(
			huh.NewInput().
				Key(nameKey).
				Title("Name?").
				Description("Name of the variable.").
				Validate(snakeCaseValidator),
			huh.NewSelect[string]().
				Key(typeKey).
				Title("Select a param type.").
				Options(shopifyParamOptions...),
			huh.NewText().
				Key(textKey).
				Title("Description").
				Description("Optionally, include a description for the variable."),
			huh.NewConfirm().
				Key(optionalKey).
				Title("Required?"),
		).WithHideFunc(func() bool {
			return !wantsParams
		}),
	).WithTheme(ShopTheme())
}

func newContinueSnippetParams() *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Key(continueKey).
				Title("Add more params?"),
		),
	).WithTheme(ShopTheme())
}

// newSettingsForm asks for a setting's type, id and label. The first form
// (more is false) also asks whether to include settings at all.
func newSettingsForm(more bool) *huh.Form {
	wantsSettings := more
	return huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Key(responseKey).
				Title("Include settings?").
				Affirmative("Yes").
				Negative("No").
				Value(&wantsSettings),
		).WithHideFunc(func() bool {
			return wantsSettings
		}),
		huh.NewGroup(
			huh.NewSelect[string]().
				Key(inputTypeKey).
				Title("Select input type").
				Options(
					shopifyInputOptions...,
				),
		).WithHideFunc(func() bool {
			return !wantsSettings
		}),
		huh.NewGroup(
			huh.NewInput().
				Key(idKey).
				Title("What is the id?").
				Validate(notEmptyValidator).
				Validate(snakeCaseValidator),
			huh.NewInput().
				Key(labelKey).
				Validate(notEmptyValidator).
				Title("What is the label?").
				Description("Label shows in the admin."),
		).WithHideFunc(func() bool {
			return !wantsSettings
		}),
	).WithTheme(ShopTheme())
}

// newIncludingBlocksForm asks for a section block and one of its settings.
// skipPrompt skips "Include section blocks?"; settingOnly asks only for a
// setting to add to the current block.
func newIncludingBlocksForm(skipPrompt, settingOnly bool) *huh.Form {
	wantsBlocks := skipPrompt
	inputsOnly := settingOnly

	return huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Key(responseKey).
				Title("Include section blocks?").
				Description("Blocks allow the admin to reorder in the editor.").
				Value(&wantsBlocks),
		).WithHideFunc(func() bool {
			return wantsBlocks
		}),
		huh.NewGroup(
			huh.NewInput().
				Key(nameKey).
				Title("Block name?").
				Description("You can name this whatever you'd like.").
				Validate(notEmptyValidator),
			huh.NewInput().
				Key(typeKey).
				Title("Type?").
				Description("Reminder - this is a type that you make up").
				Validate(notEmptyValidator).
				Validate(snakeCaseValidator),
		).WithHideFunc(func() bool {
			return !wantsBlocks || inputsOnly
		}),
		huh.NewGroup(
			huh.NewSelect[string]().
				Key(inputTypeKey).
				Title("Select input type").
				Options(
					shopifyInputOptions...,
				),
		).WithHideFunc(func() bool {
			return !wantsBlocks
		}),
		huh.NewGroup(
			huh.NewInput().
				Key(idKey).
				Title("What is the id?").
				Description("How do we want to reference this in our liquid code?").
				Validate(notEmptyValidator).
				Validate(snakeCaseValidator),
			huh.NewInput().
				Key(labelKey).
				Validate(notEmptyValidator).
				Description("This is what is shown in the admin.").
				Title("What is the label?"),
		).WithHideFunc(func() bool {
			return !wantsBlocks
		}),
	).WithTheme(ShopTheme())
}

// Confirmation for adding more inputs.
func newContinueSettingsForm() *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Key(continueKey).
				Title("Add more inputs?").
				Affirmative("Yes").
				Negative("No"),
		),
	).WithTheme(ShopTheme())
}
func newAdditionalBlocksForm() *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Key(continueKey).
				Title("Add more blocks?").
				Description("You can add more custom section blocks."),
		),
	).WithTheme(ShopTheme())
}

func newPreviewFileForm() *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Key(responseKey).
				Title("Ready to create the file?").
				Description("Yes, to create the file.\nNo, will return you to the start.").
				Affirmative("Yes, make the file now!").
				Negative("No, I'd like to try again..."),
		),
	).WithTheme(ShopTheme())
}

var shopifyParamOptions = []huh.Option[string]{
	huh.NewOption("String", "string"),
	huh.NewOption("Number", "number"),
	huh.NewOption("Boolean", "boolean"),
	huh.NewOption("Object", "object"),
}

// Collection of valid Shopify inputs.
var shopifyInputOptions = []huh.Option[string]{
	huh.NewOption("Text", "text"),
	huh.NewOption("Inline Rich Text", "inline_richtext"),
	huh.NewOption("Textarea", "textarea"),
	huh.NewOption("Rich Text", "richtext"),
	huh.NewOption("Number", "number"),
	huh.NewOption("Range Slider", "range"),
	huh.NewOption("Select", "select"),
	huh.NewOption("Radio Buttons", "radio"),
	huh.NewOption("Checkbox", "checkbox"),
	huh.NewOption("Color Picker", "color"),
	huh.NewOption("Background Color", "color_background"),
	huh.NewOption("Color Scheme", "color_scheme"),
	huh.NewOption("Font Picker", "font_picker"),
	huh.NewOption("Image Picker", "image_picker"),
	huh.NewOption("Video URL", "video_url"),
	huh.NewOption("URL", "url"),
	huh.NewOption("HTML", "html"),
	huh.NewOption("Liquid", "liquid"),
	huh.NewOption("Article", "article"),
	huh.NewOption("Blog", "blog"),
	huh.NewOption("Collection", "collection"),
	huh.NewOption("Collection List", "collection_list"),
	huh.NewOption("Product", "product"),
	huh.NewOption("Product List", "product_list"),
	huh.NewOption("Page", "page"),
	huh.NewOption("Link List", "link_list"),
}

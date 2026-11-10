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
	snakeCaseErr = errors.New("invalid liquid format, must be snake_case")
	emptyErr     = errors.New("field cannot be empty")
)

var snakeCaseValidator = func() func(s string) error {
	re := regexp.MustCompile("^[a-z_][a-z0-9_]*$")
	return func(s string) error {
		if re.MatchString(s) {
			return nil
		}
		return snakeCaseErr
	}
}()

var notEmptyValidator = func() func(s string) error {
	return func(s string) error {
		if len(s) <= 0 {
			return emptyErr
		}
		return nil
	}
}()

// User is asked to select the file type they wish to work with, give it a name
// and select an action.
var newFileSelector = func() *huh.Form {
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

// Pass in a func() bool { return true } to skip the first question.
// Pass another func() bool { return true } to skip the inital doc setup.
var newSnippetDocForm = func(opt ...func() bool) *huh.Form {

	var wantsParams bool
	var wantsMoreParams bool
	if len(opt) > 0 {
		wantsParams = opt[0]()
	}
	if len(opt) > 1 {
		wantsMoreParams = opt[1]()
	}

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

var newContinueSnippetParams = func() *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Key(continueKey).
				Title("Add more params?"),
		),
	).WithTheme(ShopTheme())
}

// Pass in a function that returns true to skip the first question.
var newSettingsForm = func(opt ...func() bool) *huh.Form {
	var wantsSettings bool
	if len(opt) > 0 {
		wantsSettings = opt[0]()
	}
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

// Available options:
//
//  1. func () bool { return true } - to skip the 'include section blocks question'
//  2. func () bool { return true } - to only fill in the input settings
var newIncludingBlocksForm = func(opt ...func() bool) *huh.Form {
	var wantsBlocks bool
	var inputsOnly bool
	if len(opt) > 0 {
		wantsBlocks = opt[0]()
	}
	if len(opt) > 1 {
		inputsOnly = opt[1]()
	}

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

// Confirmatiom for adding more inputs.
var newContinueSettingsForm = func() *huh.Form {
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
var newAdditionalBlocksForm = func() *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Key(continueKey).
				Title("Add more blocks?").
				Description("You can add more custom section blocks."),
		),
	).WithTheme(ShopTheme())
}

var newPreviewFileForm = func() *huh.Form {
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

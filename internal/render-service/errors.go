package renderservice

import "errors"

// Input errors. Callers can match these with errors.Is to show usage help.
var (
	ErrMissingName    = errors.New("a name is required (use -n)")
	ErrInvalidName    = errors.New("name must contain at least one letter or number")
	ErrInvalidSetting = errors.New("invalid setting")
	ErrInvalidBlock   = errors.New("invalid block")
)

// Environment errors: the command was run outside a theme root.
var (
	ErrMissingSectionsDir = errors.New("no 'sections' directory found; run this from the root of a Shopify theme")
	ErrMissingBlocksDir   = errors.New("no 'blocks' directory found; run this from the root of a Shopify theme")
	ErrMissingSnippetsDir = errors.New("no 'snippets' directory found; run this from the root of a Shopify theme")
)

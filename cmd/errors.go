package main

import "errors"

// ErrMissingSubcommand is returned by parent commands run without a subcommand.
var ErrMissingSubcommand = errors.New("missing subcommand")

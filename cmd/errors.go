package main

import "errors"

var (
	ErrMissingFlagOrCommand = errors.New("missing commands, flags, or operand")
)

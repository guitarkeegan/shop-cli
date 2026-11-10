// Package colors formats CLI status messages. Color is used only when
// output is a terminal and NO_COLOR is not set (https://no-color.org).
package colors

import (
	"os"

	"github.com/mattn/go-isatty"
)

const (
	reset = "\033[0m"
	bold  = "\033[1m"
	red   = "\033[38;2;239;106;50m" // #EF6A32
	green = "\033[38;2;144;173;67m" // #90AD43
)

var enabled = os.Getenv("NO_COLOR") == "" &&
	isatty.IsTerminal(os.Stdout.Fd()) && isatty.IsTerminal(os.Stderr.Fd())

func style(code, text string) string {
	if !enabled {
		return text
	}
	return bold + code + text + reset
}

// Error formats text as a bold red error message.
func Error(text string) string {
	return style(red, text)
}

// Success formats text as a bold green success message.
func Success(text string) string {
	return style(green, text)
}

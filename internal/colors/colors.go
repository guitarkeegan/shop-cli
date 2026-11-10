// Package colors formats CLI status messages.
package colors

const (
	reset = "\033[0m"
	bold  = "\033[1m"
	red   = "\033[38;2;239;106;50m" // #EF6A32
	green = "\033[38;2;144;173;67m" // #90AD43
)

// Error formats text as a bold red error message.
func Error(text string) string {
	return bold + red + text + reset
}

// Success formats text as a bold green success message.
func Success(text string) string {
	return bold + green + text + reset
}

package colors

import "fmt"

// ANSI color escape codes for Birds of Paradise theme
const (
	// Reset
	ColorReset = "\033[0m"

	// Birds of Paradise Colors (using RGB approximations for terminal)
	// Background: rich dark brown (#362725)
	ColorBackground = "\033[48;2;54;39;37m"

	// Text/Foreground: warm off-white (#E6E6E6)
	ColorForeground = "\033[38;2;230;230;230m"

	// Primary accent colors from the theme
	ColorOrange  = "\033[38;2;239;138;98m"  // Bright orange (#EF8A62)
	ColorYellow  = "\033[38;2;240;191;63m"  // Golden yellow (#F0BF3F)
	ColorRed     = "\033[38;2;239;106;50m"  // Coral red (#EF6A32)
	ColorGreen   = "\033[38;2;144;173;67m"  // Moss green (#90AD43)
	ColorBlue    = "\033[38;2;113;187;219m" // Sky blue (#71BBDB)
	ColorCyan    = "\033[38;2;19;108;90m"   // Deep teal (#136C5A)
	ColorPurple  = "\033[38;2;188;125;190m" // Soft purple (#BC7DBE)
	ColorMagenta = "\033[38;2;222;120;160m" // Pink magenta (#DE78A0)

	// Darker variants for backgrounds
	ColorOrangeBg  = "\033[48;2;239;138;98m"
	ColorYellowBg  = "\033[48;2;240;191;63m"
	ColorRedBg     = "\033[48;2;239;106;50m"
	ColorGreenBg   = "\033[48;2;144;173;67m"
	ColorBlueBg    = "\033[48;2;113;187;219m"
	ColorCyanBg    = "\033[48;2;19;108;90m"
	ColorPurpleBg  = "\033[48;2;188;125;190m"
	ColorMagentaBg = "\033[48;2;222;120;160m"

	// Dimmed variants for secondary text
	ColorDimOrange = "\033[38;2;179;103;73m"
	ColorDimYellow = "\033[38;2;180;143;47m"
	ColorDimRed    = "\033[38;2;179;79;37m"
	ColorDimGreen  = "\033[38;2;108;129;50m"
	ColorDimBlue   = "\033[38;2;84;140;164m"
	ColorDimCyan   = "\033[38;2;14;81;67m"

	// Text styles
	ColorBold      = "\033[1m"
	ColorDim       = "\033[2m"
	ColorItalic    = "\033[3m"
	ColorUnderline = "\033[4m"
)

// Helper functions for colored text
func Orange(text string) string {
	return fmt.Sprintf("%s%s%s", ColorOrange, text, ColorReset)
}

func Yellow(text string) string {
	return fmt.Sprintf("%s%s%s", ColorYellow, text, ColorReset)
}

func Red(text string) string {
	return fmt.Sprintf("%s%s%s", ColorRed, text, ColorReset)
}

func Green(text string) string {
	return fmt.Sprintf("%s%s%s", ColorGreen, text, ColorReset)
}

func Blue(text string) string {
	return fmt.Sprintf("%s%s%s", ColorBlue, text, ColorReset)
}

func Cyan(text string) string {
	return fmt.Sprintf("%s%s%s", ColorCyan, text, ColorReset)
}

func Purple(text string) string {
	return fmt.Sprintf("%s%s%s", ColorPurple, text, ColorReset)
}

func Magenta(text string) string {
	return fmt.Sprintf("%s%s%s", ColorMagenta, text, ColorReset)
}

// Dimmed versions for secondary text
func DimOrange(text string) string {
	return fmt.Sprintf("%s%s%s", ColorDimOrange, text, ColorReset)
}

func DimYellow(text string) string {
	return fmt.Sprintf("%s%s%s", ColorDimYellow, text, ColorReset)
}

func DimRed(text string) string {
	return fmt.Sprintf("%s%s%s", ColorDimRed, text, ColorReset)
}

func DimGreen(text string) string {
	return fmt.Sprintf("%s%s%s", ColorDimGreen, text, ColorReset)
}

func DimBlue(text string) string {
	return fmt.Sprintf("%s%s%s", ColorDimBlue, text, ColorReset)
}

func DimCyan(text string) string {
	return fmt.Sprintf("%s%s%s", ColorDimCyan, text, ColorReset)
}

// Background colored text
func OrangeBg(text string) string {
	return fmt.Sprintf("%s%s%s%s", ColorOrangeBg, ColorForeground, text, ColorReset)
}

func YellowBg(text string) string {
	return fmt.Sprintf("%s%s%s%s", ColorYellowBg, ColorForeground, text, ColorReset)
}

func RedBg(text string) string {
	return fmt.Sprintf("%s%s%s%s", ColorRedBg, ColorForeground, text, ColorReset)
}

func GreenBg(text string) string {
	return fmt.Sprintf("%s%s%s%s", ColorGreenBg, ColorForeground, text, ColorReset)
}

func BlueBg(text string) string {
	return fmt.Sprintf("%s%s%s%s", ColorBlueBg, ColorForeground, text, ColorReset)
}

func CyanBg(text string) string {
	return fmt.Sprintf("%s%s%s%s", ColorCyanBg, ColorForeground, text, ColorReset)
}

func PurpleBg(text string) string {
	return fmt.Sprintf("%s%s%s%s", ColorPurpleBg, ColorForeground, text, ColorReset)
}

func MagentaBg(text string) string {
	return fmt.Sprintf("%s%s%s%s", ColorMagentaBg, ColorForeground, text, ColorReset)
}

// Text styling helpers
func Bold(text string) string {
	return fmt.Sprintf("%s%s%s", ColorBold, text, ColorReset)
}

func Dim(text string) string {
	return fmt.Sprintf("%s%s%s", ColorDim, text, ColorReset)
}

func Italic(text string) string {
	return fmt.Sprintf("%s%s%s", ColorItalic, text, ColorReset)
}

func Underline(text string) string {
	return fmt.Sprintf("%s%s%s", ColorUnderline, text, ColorReset)
}

// Combined styling helpers
func BoldRed(text string) string {
	return fmt.Sprintf("%s%s%s%s", ColorBold, ColorRed, text, ColorReset)
}

func BoldOrange(text string) string {
	return fmt.Sprintf("%s%s%s%s", ColorBold, ColorOrange, text, ColorReset)
}

func BoldYellow(text string) string {
	return fmt.Sprintf("%s%s%s%s", ColorBold, ColorYellow, text, ColorReset)
}

func BoldGreen(text string) string {
	return fmt.Sprintf("%s%s%s%s", ColorBold, ColorGreen, text, ColorReset)
}

func BoldBlue(text string) string {
	return fmt.Sprintf("%s%s%s%s", ColorBold, ColorBlue, text, ColorReset)
}

// Error/Warning/Info helpers using Birds of Paradise colors
func Error(text string) string {
	return BoldRed(text)
}

func Warning(text string) string {
	return BoldOrange(text)
}

func Success(text string) string {
	return BoldGreen(text)
}

func Info(text string) string {
	return BoldBlue(text)
}

func Notice(text string) string {
	return BoldYellow(text)
}

// Usage examples:
// log.Fatalf("%s\n", Error("error: something went wrong"))
// fmt.Printf("%s: %s\n", Warning("Warning"), "This is a warning message")
// fmt.Printf("%s\n", Success("Build completed successfully!"))
// fmt.Printf("%s: %s\n", Info("Info"), "Additional information here")
// fmt.Printf("%s\n", Notice("Remember to run tests before deploying"))

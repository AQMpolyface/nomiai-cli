package elevenlab

import "fmt"

var themes = map[string]struct {
	PromptColor  string
	ErrorColor   string
	AIColor      string
	SuccessColor string
}{
	"purple_pink": {
		PromptColor:  "\033[35m", // Purple
		ErrorColor:   "\033[31m", // Red
		AIColor:      "\033[36m", // Cyan for AI prompt
		SuccessColor: "\033[32m", // Green for success messages
	},
	"light_blue": {
		PromptColor:  "\033[94m", // Light Blue
		ErrorColor:   "\033[31m", // Red
		AIColor:      "\033[36m", // Cyan for AI prompt
		SuccessColor: "\033[32m", // Green for success messages
	},
	"indigo_fuchsia": {
		PromptColor:  "\033[38;5;54m", // Indigo
		ErrorColor:   "\033[31m",      // Red
		AIColor:      "\033[36m",      // Cyan for AI prompt
		SuccessColor: "\033[32m",      // Green for success messages
	},
}

const ResetColor = "\033[0m" // ansi code to destroy the prowompt

func ChooseTheme() string {
	fmt.Println("Available themes:")
	for name := range themes {
		fmt.Println("-", name)
	}

	var selectedTheme string
	for {
		fmt.Print("Select a theme: ")
		fmt.Scanln(&selectedTheme)
		if _, exists := themes[selectedTheme]; exists {
			return selectedTheme
		}
		fmt.Println("Invalid theme. Please select again.")
	}
}

// gotta implement it somehow
func GetTheme(name string) (struct {
	PromptColor  string
	ErrorColor   string
	AIColor      string
	SuccessColor string
}, bool) {
	theme, exists := themes[name]
	return theme, exists
}

// Package styles implements the theming logic for slides
package styles

import (
	_ "embed"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

const (
	salmon = lipgloss.Color("#E8B4BC")
)

var (
	// Author is the style for the author text in the bottom-left corner of the
	// presentation.
	Author = lipgloss.NewStyle().Foreground(salmon).Align(lipgloss.Left).MarginLeft(2)
	// Date is the style for the date text in the bottom-left corner of the
	// presentation.
	Date = lipgloss.NewStyle().Faint(true).Align(lipgloss.Left).Margin(0, 1)
	// Page is the style for the pagination progress information text in the
	// bottom-right corner of the presentation.
	Page = lipgloss.NewStyle().Foreground(salmon).Align(lipgloss.Right).MarginRight(3)
	// Slide is the style for the slide.
	Slide = lipgloss.NewStyle().Padding(1)
	// Status is the style for the status bar at the bottom of the
	// presentation.
	Status = lipgloss.NewStyle().Padding(1)
	// Search is the style for the search input at the bottom-left corner of
	// the screen when searching is active.
	Search = lipgloss.NewStyle().Faint(true).Align(lipgloss.Left).MarginLeft(2)
)

var (
	// DefaultTheme is the default theme for the presentation.
	//go:embed theme.json
	DefaultTheme []byte
)

// JoinHorizontal joins two strings horizontally and fills the space in-between.
func JoinHorizontal(left, right string, width int) string {
	w := width - lipgloss.Width(right)
	return lipgloss.PlaceHorizontal(w, lipgloss.Left, left) + right
}

// JoinVertical joins two strings vertically and fills the space in-between.
func JoinVertical(top, bottom string, height int) string {
	h := height - lipgloss.Height(bottom)
	return lipgloss.PlaceVertical(h, lipgloss.Top, top) + bottom
}

// SelectTheme picks a glamour style or JSON file and always enables syntax highlighting.
func SelectTheme(theme string) glamour.TermRendererOption {
	var opt glamour.TermRendererOption

	switch {
	case theme == "ascii" || theme == "light" || theme == "dark" || theme == "notty" || theme == "pink" || theme == "dracula" || theme == "tokyo-night":
		opt = glamour.WithStandardStyle(theme)
	case strings.HasPrefix(theme, "http://") || strings.HasPrefix(theme, "https://"):
		resp, err := http.Get(theme)
		if err != nil {
			return getDefaultTheme()
		}
		defer resp.Body.Close()
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return getDefaultTheme()
		}
		opt = glamour.WithStylesFromJSONBytes(b)
	default:
		if _, err := os.Stat(theme); err == nil {
			opt = glamour.WithStylesFromJSONFile(theme)
		} else {
			opt = glamour.WithStylesFromJSONBytes(DefaultTheme)
		}
	}

	// Enable Chroma syntax highlighting
	return glamour.WithOptions(
		opt,
		glamour.WithChromaFormatter("terminal256"),
	)
}

func getDefaultTheme() glamour.TermRendererOption {
	if termenv.EnvNoColor() {
		return glamour.WithStandardStyle("notty")
	}
	if !termenv.HasDarkBackground() {
		return glamour.WithOptions(
			glamour.WithStandardStyle("light"),
			glamour.WithChromaFormatter("terminal256"),
		)
	}
	return glamour.WithOptions(
		glamour.WithStylesFromJSONBytes(DefaultTheme),
		glamour.WithChromaFormatter("terminal256"),
	)
}

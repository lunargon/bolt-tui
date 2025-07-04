package styles

import "github.com/charmbracelet/lipgloss"

// Styles defines UI styles
type Styles struct {
	Tab            lipgloss.Style
	ActiveTab      lipgloss.Style
	TabBorder      lipgloss.Style
	TabBorderRight lipgloss.Style
	Header         lipgloss.Style
	Cell           lipgloss.Style
	Selected       lipgloss.Style
	Help           lipgloss.Style
	Title          lipgloss.Style
}

var TableStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240")).
	BorderTop(true).
	BorderBottom(true).
	BorderRight(true).
	BorderLeft(true)

// DefaultStyles returns default styles
func DefaultStyles() Styles {
	tab := lipgloss.NewStyle().Padding(0, 2)
	activeTab := tab.Copy().Bold(true).Foreground(lipgloss.Color("#1E88E5"))
	tabBorder := lipgloss.NewStyle().Foreground(lipgloss.Color("#BBBBBB"))
	tabBorderRight := tabBorder.Copy().BorderRight(true)
	header := lipgloss.NewStyle().Bold(true).Padding(0, 1).Foreground(lipgloss.Color("#3949AB"))
	cell := lipgloss.NewStyle().Padding(0, 1)
	selected := lipgloss.NewStyle().Padding(0, 1).Background(lipgloss.Color("#1E88E5")).Foreground(lipgloss.Color("#FFFFFF"))
	help := lipgloss.NewStyle().Foreground(lipgloss.Color("#626262"))
	title := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#3949AB")).Padding(1, 2)

	return Styles{
		Tab:            tab,
		ActiveTab:      activeTab,
		TabBorder:      tabBorder,
		TabBorderRight: tabBorderRight,
		Header:         header,
		Cell:           cell,
		Selected:       selected,
		Help:           help,
		Title:          title,
	}
}

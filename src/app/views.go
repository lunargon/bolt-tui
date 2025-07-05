package app

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/lunargon/bolt-tui/src/app/styles"
)

// renderTitle renders the application title
func (m *Model) renderTitle() string {
	return m.styles.Title.Render("BoltDB TUI - " + m.db.Path)
}

// renderError renders error messages
func (m *Model) renderError() string {
	return fmt.Sprintf("Error: %v", m.err)
}

// renderMainContent renders the main content based on current state
func (m *Model) renderMainContent() string {
	switch m.state {
	case stateBuckets:
		return m.renderBucketsView()
	case stateCreateBucket:
		return m.renderCreateBucketView()
	case stateCreateKey:
		return m.renderCreateKeyView()
	case stateEditValue:
		return m.renderEditValueView()
	case stateEditBucket:
		return m.renderEditBucketView()
	case stateEditKey:
		return m.renderEditKeyView()
	case stateConfirmDelete:
		return m.renderConfirmDeleteView()
	case stateConfirmDeleteBucket:
		return m.renderConfirmDeleteBucketView()
	case stateSettings:
		return m.renderSettingsView()
	default:
		return ""
	}
}

// renderBucketsView renders the main buckets and keys view
func (m *Model) renderBucketsView() string {
	if len(m.buckets) == 0 {
		return "No buckets found. Press Ctrl+T to create a new bucket."
	}

	var s strings.Builder
	s.WriteString(m.renderTabs())
	s.WriteString("\n\n")
	s.WriteString(fmt.Sprintf("Display Mode: %s", m.displayMode.String()))
	s.WriteString("\n\n")
	s.WriteString(styles.TableStyle.Render(m.table.View()))
	return s.String()
}

// renderCreateBucketView renders the bucket creation view
func (m *Model) renderCreateBucketView() string {
	var s strings.Builder
	s.WriteString("Create new bucket:\n\n")
	s.WriteString(m.textInput.View())
	s.WriteString("\n\nPress Enter to create, Esc to cancel")
	return s.String()
}

// renderCreateKeyView renders the key creation view
func (m *Model) renderCreateKeyView() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintf("Create new key in bucket '%s':\n\n", m.currentBucket))
	s.WriteString(m.textInput.View())
	s.WriteString("\n\nPress Enter to create, Esc to cancel")
	return s.String()
}

// renderEditValueView renders the value editing view
func (m *Model) renderEditValueView() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintf("Edit value of key '%s' in bucket '%s':\n\n", m.currentKey, m.currentBucket))
	s.WriteString(m.textInput.View())
	s.WriteString("\n\nPress Enter to save, Esc to cancel")
	return s.String()
}

// renderEditBucketView renders the bucket name editing view
func (m *Model) renderEditBucketView() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintf("Edit bucket name '%s':\n\n", m.originalBucketName))
	s.WriteString(m.textInput.View())
	s.WriteString("\n\nPress Enter to save, Esc to cancel")
	return s.String()
}

// renderEditKeyView renders the key name editing view
func (m *Model) renderEditKeyView() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintf("Edit key name '%s' in bucket '%s':\n\n", m.originalKeyName, m.currentBucket))
	s.WriteString(m.textInput.View())
	s.WriteString("\n\nPress Enter to save, Esc to cancel")
	return s.String()
}

// renderConfirmDeleteView renders the key deletion confirmation view
func (m *Model) renderConfirmDeleteView() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintf("⚠️  DELETE CONFIRMATION ⚠️\n\n"))
	s.WriteString(fmt.Sprintf("Are you sure you want to delete key '%s' from bucket '%s'?\n\n", m.deleteKey, m.currentBucket))
	s.WriteString("This action cannot be undone.\n\n")
	s.WriteString("Press Enter to confirm deletion, Esc to cancel")
	return s.String()
}

// renderConfirmDeleteBucketView renders the bucket deletion confirmation view
func (m *Model) renderConfirmDeleteBucketView() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintf("⚠️  DELETE BUCKET CONFIRMATION ⚠️\n\n"))
	s.WriteString(fmt.Sprintf("Are you sure you want to delete bucket '%s'?\n\n", m.deleteBucket))
	s.WriteString("This will permanently delete ALL keys and values in this bucket.\n")
	s.WriteString("This action cannot be undone.\n\n")
	s.WriteString("Press Enter to confirm deletion, Esc to cancel")
	return s.String()
}

// renderSettingsView renders the settings view
func (m *Model) renderSettingsView() string {
	var s strings.Builder
	s.WriteString("Settings - Data Display Mode:\n\n")

	displayModes := []DisplayMode{DisplayString, DisplayBase64, DisplayBase58, DisplayHex}

	for i, mode := range displayModes {
		prefix := "  "
		if i == m.settingsSelection {
			prefix = "→ "
		}

		current := ""
		if mode == m.displayMode {
			current = " (current)"
		}

		s.WriteString(fmt.Sprintf("%s%s%s\n", prefix, mode.String(), current))
	}

	s.WriteString("\n\nUse ↑/↓ to navigate, Enter to select, Esc to cancel")
	return s.String()
}

// renderHelp renders the help section
func (m *Model) renderHelp() string {
	if m.showHelp {
		return "\n\n" + m.help.View(m.keyMap)
	}
	return "\n\n" + m.styles.Help.Render(" Press ? for help ")
}

func (m *Model) renderTabs() string {
	var renderedTabs []string

	for i, bucket := range m.buckets {
		var style lipgloss.Style
		if i == m.activeTab {
			style = m.styles.ActiveTab
		} else {
			style = m.styles.Tab
		}
		renderedTabs = append(renderedTabs, style.Render(bucket))
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
}

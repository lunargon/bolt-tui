package app

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lunargon/bolt-tui/src/app/styles"
	"github.com/lunargon/bolt-tui/src/bolt"
)

type state int

const (
	stateBuckets state = iota
	stateCreateBucket
	stateCreateKey
	stateEditValue
	stateEditBucket
	stateEditKey
	stateConfirmDelete
	stateConfirmDeleteBucket
)

// KeyMap defines keybindings
type KeyMap struct {
	Up           key.Binding
	Down         key.Binding
	Left         key.Binding
	Right        key.Binding
	NewTab       key.Binding
	New          key.Binding
	Delete       key.Binding
	DeleteBucket key.Binding
	Enter        key.Binding
	Esc          key.Binding
	Quit         key.Binding
	Help         key.Binding
	PrevTab      key.Binding
	NextTab      key.Binding
	SelectTab    key.Binding
	Edit         key.Binding
	EditBucket   key.Binding
}

// DefaultKeyMap returns default keybindings
func DefaultKeyMap() KeyMap {
	return KeyMap{
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "move up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "move down"),
		),
		Left: key.NewBinding(
			key.WithKeys("left", "h"),
			key.WithHelp("←/h", "move left"),
		),
		Right: key.NewBinding(
			key.WithKeys("right", "l"),
			key.WithHelp("→/l", "move right"),
		),
		NewTab: key.NewBinding(
			key.WithKeys("ctrl+t"),
			key.WithHelp("ctrl+t", "new tab"),
		),
		New: key.NewBinding(
			key.WithKeys("ctrl+n"),
			key.WithHelp("ctrl+n", "new key"),
		),
		Delete: key.NewBinding(
			key.WithKeys("ctrl+d"),
			key.WithHelp("ctrl+d", "delete"),
		),
		DeleteBucket: key.NewBinding(
			key.WithKeys("ctrl+r"),
			key.WithHelp("ctrl+r", "remove bucket"),
		),
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select/confirm"),
		),
		Esc: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "back"),
		),
		Quit: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("ctrl+c", "quit"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "toggle help"),
		),
		PrevTab: key.NewBinding(
			key.WithKeys("shift+tab"),
			key.WithHelp("shift+tab", "previous tab"),
		),
		NextTab: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "next tab"),
		),
		// Not implemented yet
		SelectTab: key.NewBinding(
			key.WithKeys("1", "2", "3", "4", "5", "6", "7", "8", "9"),
			key.WithHelp("1-9", "select tab"),
		),
		EditBucket: key.NewBinding(
			key.WithKeys("ctrl+b"),
			key.WithHelp("ctrl+b", "edit bucket name"),
		),
		Edit: key.NewBinding(
			key.WithKeys("ctrl+e"),
			key.WithHelp("ctrl+e", "edit bucket/key"),
		),
	}
}

type Model struct {
	db            *bolt.DB
	state         state
	buckets       []string
	activeTab     int
	table         table.Model
	currentBucket string
	currentKey    string
	// value              string
	textInput          textinput.Model
	help               help.Model
	keyMap             KeyMap
	styles             styles.Styles
	width              int
	height             int
	showHelp           bool
	err                error
	originalBucketName string // For storing original bucket name during editing
	originalKeyName    string // For storing original key name during editing
	deleteKey          string // For storing key name to be deleted
	deleteBucket       string // For storing bucket name to be deleted
	newlyCreatedBucket string // For tracking newly created bucket to make it active
}

func New(dbPath string) (*Model, error) {
	db := &bolt.DB{Path: dbPath}
	err := db.Open()
	if err != nil {
		return nil, err
	}

	// Initialize table
	columns := []table.Column{
		{Title: "Key", Width: 50},
		{Title: "Value", Width: 60},
	}
	rows := []table.Row{}

	tbl := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	// Set default styles
	s := styles.DefaultStyles()
	tbl.SetStyles(table.Styles{
		Header:   s.Header,
		Cell:     s.Cell,
		Selected: s.Selected,
	})

	// Initialize text input
	ti := textinput.New()
	ti.Placeholder = "Enter name..."
	ti.Focus()

	// Initialize help
	km := DefaultKeyMap()
	h := help.New()

	m := &Model{
		db:        db,
		state:     stateBuckets,
		table:     tbl,
		textInput: ti,
		help:      h,
		keyMap:    km,
		styles:    s,
		showHelp:  false,
	}

	return m, nil
}

func (m *Model) Init() tea.Cmd {
	return m.loadBuckets
}

func (m *Model) loadBuckets() tea.Msg {
	buckets, err := m.db.GetBuckets()
	if err != nil {
		return err
	}

	return bucketsLoadedMsg{buckets}
}

type bucketsLoadedMsg struct {
	buckets []string
}

// Load key and values
// TODO: Add more view: base64, base58 or hex
func (m *Model) loadKeysAndValues(bucket string) tea.Cmd {
	return func() tea.Msg {
		keys, err := m.db.GetKeysInBucket(bucket)
		if err != nil {
			return err
		}

		rows := make([]table.Row, len(keys))
		for i, k := range keys {
			value, err := m.db.GetValue(bucket, k)
			if err != nil {
				rows[i] = table.Row{k, "Error: " + err.Error()}
			} else {
				rows[i] = table.Row{k, string(value)}
			}
		}

		return keysLoadedMsg{rows}
	}
}

type keysLoadedMsg struct {
	rows []table.Row
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.Quit):
			m.db.Close()
			return m, tea.Quit

		case key.Matches(msg, m.keyMap.Help):
			m.showHelp = !m.showHelp
			return m, nil

		case key.Matches(msg, m.keyMap.Esc):
			if m.state == stateCreateBucket || m.state == stateCreateKey ||
				m.state == stateEditBucket || m.state == stateEditKey || m.state == stateEditValue ||
				m.state == stateConfirmDelete || m.state == stateConfirmDeleteBucket {
				m.state = stateBuckets
				return m, nil
			}
		case key.Matches(msg, m.keyMap.NewTab):
			if m.state == stateBuckets {
				m.state = stateCreateBucket
				m.textInput.SetValue("")
				return m, textinput.Blink
			} else if m.state == stateCreateBucket {
				m.state = stateBuckets
				return m, nil
			}
		case key.Matches(msg, m.keyMap.New):
			if len(m.buckets) > 0 {
				m.state = stateCreateKey
				m.textInput.SetValue("")
				return m, textinput.Blink
			}

		case key.Matches(msg, m.keyMap.Edit):
			if m.state == stateBuckets && len(m.buckets) > 0 && len(m.table.Rows()) > 0 {
				selectedRow := m.table.SelectedRow()
				if selectedRow != nil {
					// Edit the selected key
					m.originalKeyName = selectedRow[0]
					m.currentKey = m.originalKeyName
					m.textInput.SetValue(m.originalKeyName)
					m.state = stateEditKey
					return m, textinput.Blink
				}
			}

		case key.Matches(msg, m.keyMap.EditBucket):
			if m.state == stateBuckets && len(m.buckets) > 0 {
				// Edit the current bucket name
				m.originalBucketName = m.currentBucket
				m.textInput.SetValue(m.currentBucket)
				m.state = stateEditBucket
				return m, textinput.Blink
			}

		case key.Matches(msg, m.keyMap.Delete):
			if m.state == stateBuckets && len(m.buckets) > 0 && len(m.table.Rows()) > 0 {
				selectedRow := m.table.SelectedRow()
				if selectedRow != nil {
					// Store the key to be deleted and show confirmation
					m.deleteKey = selectedRow[0]
					m.state = stateConfirmDelete
					return m, nil
				}
			}

		case key.Matches(msg, m.keyMap.DeleteBucket):
			if m.state == stateBuckets && len(m.buckets) > 0 {
				// Store the bucket to be deleted and show confirmation
				m.deleteBucket = m.currentBucket
				m.state = stateConfirmDeleteBucket
				return m, nil
			}

		case key.Matches(msg, m.keyMap.PrevTab):
			if m.state == stateBuckets && len(m.buckets) > 0 {
				m.activeTab = max(0, m.activeTab-1)
				m.currentBucket = m.buckets[m.activeTab]
				return m, m.loadKeysAndValues(m.currentBucket)
			}

		case key.Matches(msg, m.keyMap.NextTab):
			if m.state == stateBuckets && len(m.buckets) > 0 {
				m.activeTab = min(len(m.buckets)-1, m.activeTab+1)
				m.currentBucket = m.buckets[m.activeTab]
				return m, m.loadKeysAndValues(m.currentBucket)
			}

		case key.Matches(msg, m.keyMap.Enter):
			if m.state == stateConfirmDelete {
				// Delete the key
				err := m.db.DeleteValue(m.currentBucket, m.deleteKey)
				if err != nil {
					m.err = err
					return m, nil
				}
				m.state = stateBuckets
				m.deleteKey = ""
				return m, m.loadKeysAndValues(m.currentBucket)
			} else if m.state == stateConfirmDeleteBucket {
				// Delete the bucket
				err := m.db.DeleteBucket(m.deleteBucket)
				if err != nil {
					m.err = err
					return m, nil
				}
				m.state = stateBuckets
				m.deleteBucket = ""
				// Reset activeTab if we deleted the current bucket
				if m.activeTab >= len(m.buckets)-1 {
					m.activeTab = max(0, len(m.buckets)-2)
				}
				return m, m.loadBuckets
			} else if m.state == stateCreateBucket {
				newBucket := m.textInput.Value()
				if newBucket != "" {
					err := m.db.CreateBucket(newBucket)
					if err != nil {
						m.err = err
						return m, nil
					}
					m.newlyCreatedBucket = newBucket
					m.state = stateBuckets
					return m, m.loadBuckets
				}
			} else if m.state == stateCreateKey && len(m.buckets) > 0 {
				newKey := m.textInput.Value()
				if newKey != "" {
					err := m.db.PutValue(m.currentBucket, newKey, []byte(""))
					if err != nil {
						m.err = err
						return m, nil
					}
					m.state = stateBuckets
					return m, m.loadKeysAndValues(m.currentBucket)
				}
			} else if m.state == stateEditValue && m.currentKey != "" {
				newValue := m.textInput.Value()
				if newValue != "" {
					err := m.db.PutValue(m.currentBucket, m.currentKey, []byte(newValue))
					if err != nil {
						m.err = err
						return m, nil
					}
					m.state = stateBuckets
					return m, m.loadKeysAndValues(m.currentBucket)
				}
			} else if m.state == stateEditBucket && m.originalBucketName != "" {
				newBucketName := m.textInput.Value()
				if newBucketName != "" && newBucketName != m.originalBucketName {
					err := m.db.RenameBucket(m.originalBucketName, newBucketName)
					if err != nil {
						m.err = err
						return m, nil
					}
					// Update the current bucket name to the new name
					m.currentBucket = newBucketName
					m.state = stateBuckets
					return m, m.loadBuckets
				} else if newBucketName == m.originalBucketName {
					// No change, just go back
					m.state = stateBuckets
					return m, nil
				}
			} else if m.state == stateEditKey && m.originalKeyName != "" {
				newKeyName := m.textInput.Value()
				if newKeyName != "" && newKeyName != m.originalKeyName {
					err := m.db.RenameKey(m.currentBucket, m.originalKeyName, newKeyName)
					if err != nil {
						m.err = err
						return m, nil
					}
					m.state = stateBuckets
					return m, m.loadKeysAndValues(m.currentBucket)
				} else if newKeyName == m.originalKeyName {
					// No change, just go back
					m.state = stateBuckets
					return m, nil
				}
			} else if m.state == stateBuckets && len(m.buckets) > 0 && len(m.table.Rows()) > 0 {
				// Get the selected row from the table
				selectedRow := m.table.SelectedRow()
				if selectedRow != nil {
					// The key is in the first column
					m.currentKey = selectedRow[0]
					// Get the current value to populate the text input
					value, err := m.db.GetValue(m.currentBucket, m.currentKey)
					if err != nil {
						m.err = err
						return m, nil
					}
					// Set the text input value to the current value
					m.textInput.SetValue(string(value))
					// Transition to edit value state
					m.state = stateEditValue
					return m, textinput.Blink
				}
			}

		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.table.SetWidth(msg.Width - 4)    // Account for margins
		m.table.SetHeight(msg.Height - 10) // Account for header, tabs, help, etc.
		return m, nil

	case bucketsLoadedMsg:
		m.buckets = msg.buckets
		if len(m.buckets) > 0 {
			// If we just created a new bucket, find it and make it active
			if m.newlyCreatedBucket != "" {
				for i, bucket := range m.buckets {
					if bucket == m.newlyCreatedBucket {
						m.activeTab = i
						m.newlyCreatedBucket = "" // Clear the flag
						break
					}
				}
			}

			// If we have a current bucket name (e.g., after renaming), find its new index
			if m.currentBucket != "" {
				for i, bucket := range m.buckets {
					if bucket == m.currentBucket {
						m.activeTab = i
						break
					}
				}
			}

			// Ensure activeTab is within bounds
			if m.activeTab >= len(m.buckets) {
				m.activeTab = len(m.buckets) - 1
			}
			if m.activeTab < 0 {
				m.activeTab = 0
			}

			m.currentBucket = m.buckets[m.activeTab]
			return m, m.loadKeysAndValues(m.currentBucket)
		}
		return m, nil

	case keysLoadedMsg:
		m.table.SetRows(msg.rows)
		return m, nil

	case error:
		m.err = msg
		return m, nil
	}

	var cmd tea.Cmd
	switch m.state {
	case stateBuckets:
		m.table, cmd = m.table.Update(msg)

	case stateCreateBucket, stateCreateKey, stateEditBucket, stateEditKey, stateEditValue:
		m.textInput, cmd = m.textInput.Update(msg)
	}

	return m, cmd
}

func (m *Model) View() string {
	var s strings.Builder

	// Title
	s.WriteString(m.styles.Title.Render("BoltDB TUI - " + m.db.Path))
	s.WriteString("\n\n")

	// Error message
	if m.err != nil {
		s.WriteString(fmt.Sprintf("Error: %v\n\n", m.err))
	}

	switch m.state {
	case stateBuckets:
		// Render tabs for buckets
		if len(m.buckets) > 0 {
			s.WriteString(m.renderTabs())
			s.WriteString("\n\n")

			// Render table for keys and values
			s.WriteString(styles.TableStyle.Render(m.table.View()))
		} else {
			s.WriteString("No buckets found. Press 'n' to create a new bucket.")
		}

	case stateCreateBucket:
		s.WriteString("Create new bucket:\n\n")
		s.WriteString(m.textInput.View())

	case stateCreateKey:
		s.WriteString(fmt.Sprintf("Create new key in bucket '%s':\n\n", m.currentBucket))
		s.WriteString(m.textInput.View())

	case stateEditValue:
		s.WriteString(fmt.Sprintf("Edit value of key '%s' in bucket '%s':\n\n", m.currentKey, m.currentBucket))
		s.WriteString(m.textInput.View())

	case stateEditBucket:
		s.WriteString(fmt.Sprintf("Edit bucket name '%s':\n\n", m.originalBucketName))
		s.WriteString(m.textInput.View())

	case stateEditKey:
		s.WriteString(fmt.Sprintf("Edit key name '%s' in bucket '%s':\n\n", m.originalKeyName, m.currentBucket))
		s.WriteString(m.textInput.View())

	case stateConfirmDelete:
		s.WriteString(fmt.Sprintf("Are you sure you want to delete key '%s' from bucket '%s'?\n\n", m.deleteKey, m.currentBucket))
		s.WriteString("Press Enter to confirm, Esc to cancel")

	case stateConfirmDeleteBucket:
		s.WriteString(fmt.Sprintf("Are you sure you want to delete bucket '%s'?\n\n", m.deleteBucket))
		s.WriteString("This will delete all keys and values in the bucket.\n")
		s.WriteString("Press Enter to confirm, Esc to cancel")
	}

	// Help
	if m.showHelp {
		s.WriteString("\n\n" + m.help.View(m.keyMap))
	} else {
		s.WriteString("\n\n" + m.styles.Help.Render(" Press ? for help "))
	}

	return lipgloss.NewStyle().Margin(1, 2).Render(s.String())
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

func (m *Model) Close() {
	if m.db != nil {
		m.db.Close()
	}
}

// Add these methods to implement the help.KeyMap interface

// ShortHelp returns keybindings to be shown in the short help view.
// It's part of the help.KeyMap interface.
func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Enter, k.Esc, k.NewTab, k.New,
		k.Edit, k.EditBucket, k.Delete, k.DeleteBucket, k.Help, k.Quit}
}

// FullHelp returns keybindings for the expanded help view.
// It's part of the help.KeyMap interface.
func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right}, // first column
		{k.Enter, k.Esc, k.NewTab, k.New, k.Edit, k.EditBucket, k.Delete, k.DeleteBucket}, // second column
		{k.PrevTab, k.NextTab, k.SelectTab},                                               // third column
		{k.Help, k.Quit},                                                                  // fourth column
	}
}

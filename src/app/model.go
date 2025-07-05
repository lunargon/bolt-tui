package app

import (
	"encoding/base64"
	"encoding/hex"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lunargon/bolt-tui/src/app/styles"
	"github.com/lunargon/bolt-tui/src/bolt"
	"github.com/mr-tron/base58"
)

type Model struct {
	db                 *bolt.DB
	state              state
	buckets            []string
	activeTab          int
	table              table.Model
	currentBucket      string
	currentKey         string
	textInput          textinput.Model
	help               help.Model
	keyMap             KeyMap
	styles             styles.Styles
	width              int
	height             int
	showHelp           bool
	err                error
	originalBucketName string      // For storing original bucket name during editing
	originalKeyName    string      // For storing original key name during editing
	deleteKey          string      // For storing key name to be deleted
	deleteBucket       string      // For storing bucket name to be deleted
	newlyCreatedBucket string      // For tracking newly created bucket to make it active
	displayMode        DisplayMode // Current display mode for data
	settingsSelection  int         // Current selection in settings
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

// Load key and values with support for different display modes
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
				displayValue := m.formatValue(value)
				rows[i] = table.Row{k, displayValue}
			}
		}

		return keysLoadedMsg{rows}
	}
}

type keysLoadedMsg struct {
	rows []table.Row
}

// formatValue formats the value according to the current display mode
func (m *Model) formatValue(value []byte) string {
	switch m.displayMode {
	case DisplayBase64:
		return base64.StdEncoding.EncodeToString(value)
	case DisplayBase58:
		return base58.Encode(value)
	case DisplayHex:
		return hex.EncodeToString(value)
	default:
		return string(value)
	}
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle global messages first
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m.handleWindowResize(msg)
	case bucketsLoadedMsg:
		return m.handleBucketsLoaded(msg)
	case keysLoadedMsg:
		return m.handleKeysLoaded(msg)
	case error:
		m.err = msg
		return m, nil
	}

	// Handle key messages
	if msg, ok := msg.(tea.KeyMsg); ok {
		// Handle global keys
		if cmd := m.handleGlobalKeys(msg); cmd != nil {
			return m, cmd
		}

		// Handle state-specific keys
		switch m.state {
		case stateBuckets:
			return m.handleBucketsState(msg)
		case stateCreateBucket:
			return m.handleCreateBucketState(msg)
		case stateCreateKey:
			return m.handleCreateKeyState(msg)
		case stateEditValue:
			return m.handleEditValueState(msg)
		case stateEditBucket:
			return m.handleEditBucketState(msg)
		case stateEditKey:
			return m.handleEditKeyState(msg)
		case stateConfirmDelete:
			return m.handleConfirmDeleteState(msg)
		case stateConfirmDeleteBucket:
			return m.handleConfirmDeleteBucketState(msg)
		case stateSettings:
			return m.handleSettingsState(msg)
		}
	}

	// Handle other input updates
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
	s.WriteString(m.renderTitle())
	s.WriteString("\n\n")

	// Error message
	if m.err != nil {
		s.WriteString(m.renderError())
		s.WriteString("\n\n")
	}

	// Main content based on state
	s.WriteString(m.renderMainContent())

	// Help
	s.WriteString(m.renderHelp())

	return lipgloss.NewStyle().Margin(1, 2).Render(s.String())
}

func (m *Model) Close() {
	if m.db != nil {
		m.db.Close()
	}
}

// ShortHelp returns keybindings to be shown in the short help view.
// It's part of the help.KeyMap interface.
func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Enter, k.Esc, k.NewTab, k.New,
		k.Edit, k.EditBucket, k.Delete, k.DeleteBucket, k.Settings, k.Help, k.Quit}
}

// FullHelp returns keybindings for the expanded help view.
// It's part of the help.KeyMap interface.
func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right}, // first column
		{k.Enter, k.Esc, k.NewTab, k.New, k.Edit, k.EditBucket, k.Delete, k.DeleteBucket}, // second column
		{k.PrevTab, k.NextTab, k.SelectTab, k.Settings},                                   // third column
		{k.Help, k.Quit}, // fourth column
	}
}

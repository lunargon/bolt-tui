package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/filepicker"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/lunargon/bolt-tui/src/app"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "bolt-tui",
	Short: "A TUI for viewing and managing BoltDB files",
	Long:  `A Terminal User Interface (TUI) for viewing and managing BoltDB files.`,
	Run: func(cmd *cobra.Command, args []string) {
		// If a file path is provided, open it directly
		filePath, _ := cmd.Flags().GetString("file")
		if filePath != "" {
			// Convert to absolute path if needed
			absPath, err := filepath.Abs(filePath)
			if err != nil {
				fmt.Printf("Error getting absolute path: %v\n", err)
				os.Exit(1)
			}

			// Check if file exists
			if _, err = os.Stat(absPath); os.IsNotExist(err) {
				fmt.Printf("File does not exist: %s\n", absPath)
				os.Exit(1)
			}

			// Open the app directly with the provided file
			appModel, err := app.New(absPath)
			if err != nil {
				fmt.Printf("Error opening database: %v\n", err)
				os.Exit(1)
			}

			p := tea.NewProgram(appModel, tea.WithAltScreen())
			if _, err := p.Run(); err != nil {
				fmt.Printf("Error running program: %v\n", err)
				os.Exit(1)
			}
			return
		}

		// Otherwise, show the file picker
		startDir, _ := cmd.Flags().GetString("dir")

		// Handle "." as a special case for current directory
		if startDir == "." || startDir == "" {
			// Use current directory
			currentDir, err := os.Getwd()
			if err != nil {
				fmt.Printf("Error getting current directory: %v\n", err)
				os.Exit(1)
			}
			startDir = currentDir
		} else {
			// Convert to absolute path if it's not already
			absPath, err := filepath.Abs(startDir)
			if err != nil {
				fmt.Printf("Error getting absolute path: %v\n", err)
				os.Exit(1)
			}
			startDir = absPath
		}

		// Create file picker model
		fp := filepicker.New()
		fp.Cursor = "->"
		fp.AllowedTypes = []string{".db"}
		fp.CurrentDirectory = startDir

		m := FilePickerModel{
			filepicker: fp,
		}

		p := tea.NewProgram(&m, tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			fmt.Printf("Error running program: %v\n", err)
			os.Exit(1)
		}
	},
}

// FilePickerModel represents the file picker state
type FilePickerModel struct {
	filepicker   filepicker.Model
	selectedFile string
	quitting     bool
	err          error
}

func (m *FilePickerModel) Init() tea.Cmd {
	return m.filepicker.Init()
}

func (m *FilePickerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.filepicker, cmd = m.filepicker.Update(msg)

	// Did the user select a file?
	if didSelect, path := m.filepicker.DidSelectFile(msg); didSelect {
		// Get the path of the selected file.
		m.selectedFile = path

		// Launch the main app with the selected file
		appModel, err := app.New(path)
		if err != nil {
			m.err = err
			return m, nil
		}

		return appModel, appModel.Init()
	}

	return m, cmd
}

func (m *FilePickerModel) View() string {
	if m.quitting {
		return ""
	}

	var s string
	if m.err != nil {
		s = fmt.Sprintf("\n  Error: %v\n\n%s\n", m.err, m.filepicker.View())
	} else {
		s = fmt.Sprintf("\n  Select a BoltDB file:\n\n%s\n", m.filepicker.View())
	}

	return s
}

// Execute executes the root command
func Execute() {
	rootCmd.PersistentFlags().StringP("file", "f", "", "Path to BoltDB file to open directly")
	rootCmd.PersistentFlags().StringP("dir", "d", "", "Starting directory for file picker (use '.' for current directory)")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

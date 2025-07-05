package app

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// handleGlobalKeys handles keys that work in all states
func (m *Model) handleGlobalKeys(msg tea.KeyMsg) tea.Cmd {
	switch {
	case key.Matches(msg, m.keyMap.Quit):
		m.db.Close()
		return tea.Quit
	case key.Matches(msg, m.keyMap.Help):
		m.showHelp = !m.showHelp
		return nil
	case key.Matches(msg, m.keyMap.Esc):
		if m.state != stateBuckets {
			m.state = stateBuckets
			return nil
		}
	}
	return nil
}

// handleBucketsState handles the main buckets view state
func (m *Model) handleBucketsState(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keyMap.NewTab):
		m.state = stateCreateBucket
		m.textInput.SetValue("")
		m.textInput.Focus()
		return m, textinput.Blink

	case key.Matches(msg, m.keyMap.New):
		if len(m.buckets) > 0 {
			m.state = stateCreateKey
			m.textInput.SetValue("")
			m.textInput.Focus()
			return m, textinput.Blink
		}

	case key.Matches(msg, m.keyMap.Settings):
		m.state = stateSettings
		m.settingsSelection = 0
		return m, nil

	case key.Matches(msg, m.keyMap.Edit):
		if len(m.buckets) > 0 && len(m.table.Rows()) > 0 {
			selectedRow := m.table.SelectedRow()
			if selectedRow != nil {
				m.originalKeyName = selectedRow[0]
				m.currentKey = m.originalKeyName
				m.textInput.SetValue(m.originalKeyName)
				m.textInput.Focus()
				m.state = stateEditKey
				return m, textinput.Blink
			}
		}

	case key.Matches(msg, m.keyMap.EditBucket):
		if len(m.buckets) > 0 {
			m.originalBucketName = m.currentBucket
			m.textInput.SetValue(m.currentBucket)
			m.textInput.Focus()
			m.state = stateEditBucket
			return m, textinput.Blink
		}

	case key.Matches(msg, m.keyMap.Delete):
		if len(m.buckets) > 0 && len(m.table.Rows()) > 0 {
			selectedRow := m.table.SelectedRow()
			if selectedRow != nil {
				m.deleteKey = selectedRow[0]
				m.state = stateConfirmDelete
				return m, nil
			}
		}

	case key.Matches(msg, m.keyMap.DeleteBucket):
		if len(m.buckets) > 0 {
			m.deleteBucket = m.currentBucket
			m.state = stateConfirmDeleteBucket
			return m, nil
		}

	case key.Matches(msg, m.keyMap.PrevTab):
		if len(m.buckets) > 0 {
			m.activeTab = max(0, m.activeTab-1)
			m.currentBucket = m.buckets[m.activeTab]
			return m, m.loadKeysAndValues(m.currentBucket)
		}

	case key.Matches(msg, m.keyMap.NextTab):
		if len(m.buckets) > 0 {
			m.activeTab = min(len(m.buckets)-1, m.activeTab+1)
			m.currentBucket = m.buckets[m.activeTab]
			return m, m.loadKeysAndValues(m.currentBucket)
		}

	case key.Matches(msg, m.keyMap.Enter):
		if len(m.buckets) > 0 && len(m.table.Rows()) > 0 {
			selectedRow := m.table.SelectedRow()
			if selectedRow != nil {
				m.currentKey = selectedRow[0]
				value, err := m.db.GetValue(m.currentBucket, m.currentKey)
				if err != nil {
					m.err = err
					return m, nil
				}
				m.textInput.SetValue(string(value))
				m.textInput.Focus()
				m.state = stateEditValue
				return m, textinput.Blink
			}
		}
	}

	return m, nil
}

// handleCreateBucketState handles bucket creation
func (m *Model) handleCreateBucketState(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch {
	case key.Matches(msg, m.keyMap.NewTab):
		m.state = stateBuckets
		return m, nil
	case key.Matches(msg, m.keyMap.Enter):
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
	}

	// Update text input with the message
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

// handleCreateKeyState handles key creation
func (m *Model) handleCreateKeyState(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch {
	case key.Matches(msg, m.keyMap.Enter):
		if len(m.buckets) > 0 {
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
		}
	}

	// Update text input with the message
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

// handleEditBucketState handles bucket name editing
func (m *Model) handleEditBucketState(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch {
	case key.Matches(msg, m.keyMap.Enter):
		if m.originalBucketName != "" {
			newBucketName := m.textInput.Value()
			if newBucketName != "" && newBucketName != m.originalBucketName {
				err := m.db.RenameBucket(m.originalBucketName, newBucketName)
				if err != nil {
					m.err = err
					return m, nil
				}
				m.currentBucket = newBucketName
				m.state = stateBuckets
				return m, m.loadBuckets
			} else if newBucketName == m.originalBucketName {
				m.state = stateBuckets
				return m, nil
			}
		}
	}

	// Update text input with the message
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

// handleEditKeyState handles key name editing
func (m *Model) handleEditKeyState(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch {
	case key.Matches(msg, m.keyMap.Enter):
		if m.originalKeyName != "" {
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
				m.state = stateBuckets
				return m, nil
			}
		}
	}

	// Update text input with the message
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

// handleEditValueState handles value editing
func (m *Model) handleEditValueState(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch {
	case key.Matches(msg, m.keyMap.Enter):
		if m.currentKey != "" {
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
		}
	}

	// Update text input with the message
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

// handleConfirmDeleteState handles key deletion confirmation
func (m *Model) handleConfirmDeleteState(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keyMap.Enter):
		err := m.db.DeleteValue(m.currentBucket, m.deleteKey)
		if err != nil {
			m.err = err
			return m, nil
		}
		m.state = stateBuckets
		m.deleteKey = ""
		return m, m.loadKeysAndValues(m.currentBucket)
	}
	return m, nil
}

// handleConfirmDeleteBucketState handles bucket deletion confirmation
func (m *Model) handleConfirmDeleteBucketState(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keyMap.Enter):
		err := m.db.DeleteBucket(m.deleteBucket)
		if err != nil {
			m.err = err
			return m, nil
		}
		m.state = stateBuckets
		m.deleteBucket = ""
		if m.activeTab >= len(m.buckets)-1 {
			m.activeTab = max(0, len(m.buckets)-2)
		}
		return m, m.loadBuckets
	}
	return m, nil
}

// handleWindowResize handles window resize events
func (m *Model) handleWindowResize(msg tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
	m.width = msg.Width
	m.height = msg.Height
	m.table.SetWidth(msg.Width - 4)
	m.table.SetHeight(msg.Height - 10)
	return m, nil
}

// handleBucketsLoaded handles the buckets loaded message
func (m *Model) handleBucketsLoaded(msg bucketsLoadedMsg) (tea.Model, tea.Cmd) {
	m.buckets = msg.buckets
	if len(m.buckets) > 0 {
		// If we just created a new bucket, find it and make it active
		if m.newlyCreatedBucket != "" {
			for i, bucket := range m.buckets {
				if bucket == m.newlyCreatedBucket {
					m.activeTab = i
					m.newlyCreatedBucket = ""
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
}

// handleKeysLoaded handles the keys loaded message
func (m *Model) handleKeysLoaded(msg keysLoadedMsg) (tea.Model, tea.Cmd) {
	m.table.SetRows(msg.rows)
	return m, nil
}

// handleSettingsState handles the settings view
func (m *Model) handleSettingsState(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keyMap.Up):
		if m.settingsSelection > 0 {
			m.settingsSelection--
		}
		return m, nil
	case key.Matches(msg, m.keyMap.Down):
		if m.settingsSelection < 3 { // 4 display modes (0-3)
			m.settingsSelection++
		}
		return m, nil
	case key.Matches(msg, m.keyMap.Enter):
		m.displayMode = DisplayMode(m.settingsSelection)
		m.state = stateBuckets
		m.textInput.Blur()
		// Reload current bucket with new display mode
		if len(m.buckets) > 0 {
			return m, m.loadKeysAndValues(m.currentBucket)
		}
		return m, nil
	}
	return m, nil
}

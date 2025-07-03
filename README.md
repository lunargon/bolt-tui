# Bolt TUI ( bolt-tui)

[![License](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/lunargon/bolt-tui)](https://goreportcard.com/report/github.com/lunargon/bolt-tui)

A Terminal User Interface (TUI) for viewing and managing BoltDB files. Built with Go and the Bubbletea framework.

## Features

- 🗂️ **Interactive File Picker**: Browse and select BoltDB files from your filesystem
- 📊 **Database Browser**: View all buckets and their key-value pairs
- ✏️ **Edit Operations**: Create, update, and delete buckets and keys
- 🎨 **Beautiful UI**: Clean, modern terminal interface with syntax highlighting
- ⌨️ **Keyboard Navigation**: Full keyboard support with intuitive shortcuts
- 🔍 **Help System**: Built-in help to guide you through available commands
- 📑 **Tab Navigation**: Organize your work with multiple tabs

## Installation

### Prerequisites

- Go 1.24.1 or higher

### Build From Source

```bash
git clone https://github.com/lunargon/bolt-tui.git
cd bolt-tui
go build -o bolt-tui cli/main.go
```

## Usage

### Basic Usage

Run the application to open the file picker:

```bash
./bolt-tui
```

### CLI Options

Open a specific BoltDB file directly:

```bash
./bolt-tui -f /path/to/your/database.db
```

Start file picker in a specific directory:

```bash
./bolt-tui -d /path/to/directory
```

Use current directory:

```bash
./bolt-tui -d .
```

### Keyboard Shortcuts

#### General Navigation
| Key | Action |
|-----|--------|
| `↑/k` | Move up |
| `↓/j` | Move down |
| `←/h` | Move left |
| `→/l` | Move right |
| `Enter` | Select/Confirm |
| `Esc` | Go back/Cancel |
| `Ctrl+c` | Quit |
| `?` | Toggle help |

#### Tab Management
| Key | Action |
|-----|--------|
| `Tab` | Next tab |
| `Shift+Tab` | Previous tab |
| `1-9` | Select tab by number ( Not implement this) |

#### Bucket Operations
| Key | Action |
|-----|--------|
| `Ctrl+t` | Create new bucket |
| `Ctrl+b` | Edit bucket name |
| `Ctrl+r` | Remove bucket |

#### Key-Value Operations
| Key | Action |
|-----|--------|
| `Ctrl+n` | Create new key |
| `Ctrl+e` | Edit key name |
| `Ctrl+d` | Delete key |
| `Enter` | Edit value (when key selected) |

## Project Structure

```
bolt-tui/
├── main.go              # Test entry point with file picker ( for testing)
├── cli/
│   └── main.go          # CLI entry point
├── src/
│   ├── app/             # TUI application logic
│   │   ├── model.go     # Main application model and UI
│   │   └── helper.go    # Helper functions
│   ├── bolt/            # BoltDB wrapper
│   │   └── bolt.go      # Database operations
│   └── cmd/             # CLI commands
│       └── main.go      # Cobra command definitions
├── seed/                # Database seeding utilities
│   └── seed.go
├── go.mod               # Go module definition
├── go.sum               # Go module checksums
└── README.md            # This file
```

## Dependencies

- **[BoltDB](https://github.com/boltdb/bolt)** - Embedded key/value database
- **[Bubbletea](https://github.com/charmbracelet/bubbletea)** - TUI framework
- **[Bubbles](https://github.com/charmbracelet/bubbles)** - TUI components
- **[Lipgloss](https://github.com/charmbracelet/lipgloss)** - Style definitions
- **[Cobra](https://github.com/spf13/cobra)** - CLI framework

## Development

### Building

Build the main CLI application:
```bash
go build -o bolt-tui cli/main.go
```

### Running Tests

Seed .db file for testing:
```bash
go run seed/seed.go
```

Run the main CLI application:
```bash
go run cli/main.go
```

Or run the test entry point (for testing purposes):
```bash
go run main.go
```

### Building for Different Platforms

```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o bolt-tui-linux

# macOS
GOOS=darwin GOARCH=amd64 go build -o bolt-tui-macos

# Windows
GOOS=windows GOARCH=amd64 go build -o bolt-tui-windows.exe
```

## Todo

- [ ] **Refactor**
- [ ] Update UI
- [ ] Add feature to jump tab with number
- [ ] Have switch to view `byte` value or `string` value

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feat/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feat/amazing-feature`)
5. Open a Pull Request


## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Built with [Charm](https://charm.sh/) TUI libraries
- Inspired by the need for a simple, interactive BoltDB browser
- Thanks to the Go community for excellent tooling and libraries

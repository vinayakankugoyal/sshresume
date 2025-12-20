# SSH Resume

An interactive resume/portfolio accessible via SSH, built with Go and [Charm](https://charm.sh/) libraries. Present your professional information in a beautiful terminal UI with a file-tree navigation system that users can access by simply SSH-ing to your server.

## Features

- 📂 **File-Tree Navigation** - Organize your content in a familiar directory structure with support for nested folders.
- 📄 **Split-View Interface** - Clean split-pane layout with a sidebar for navigation and a main viewport for content.
- 🎨 **Markdown Support** - Write content in markdown with syntax highlighting via [Glamour](https://github.com/charmbracelet/glamour).
- ⌨️ **Keyboard Navigation** - Intuitive vim-like keybindings for browsing and reading.
- 🔐 **SSH Server** - Self-hosted SSH server with automatic host key generation.

## Demo

Connect to see it in action:
```bash
ssh resume.fullydisfunctional.com
```

## Installation

### Prerequisites

- Go 1.25.3 or later
- SSH client (for testing)

### Build from Source

```bash
# Clone the repository
git clone https://github.com/vinayakankugoyal/sshresume.git
cd sshresume

# Install dependencies
make install

# Build the binary
make build
```

The binary will be available at `bin/sshresume`.

### Build for Linux AMD64

```bash
make build-linux-amd64
```

The binary will be available at `bin/linux/amd64/sshresume`.

## Usage

### Running the Server

```bash
# Run with default settings (localhost:23234)
./bin/sshresume

# Or use make
make run

# Development mode with go run
make dev
```

### Command Line Options

```bash
./bin/sshresume [options]

Options:
  -host string
        address to bind to (default "localhost")
  -port string
        port to bind to (default "23234")
  -config string
        path to config file (default "config.yaml")
```

## Configuration

The application is configured using a YAML file (default `config.yaml`).

### Setting Up Your Config File

1. Create a `config.yaml` or copy the example:
   ```bash
   cp config.example.yaml config.yaml
   ```

2. Point the `folder` key to the directory containing your markdown files:

   ```yaml
   # config.yaml
   
   # Folder containing your markdown files (supports nested directories)
   folder: "./content"
   ```

### Organizing Content

The UI will automatically generate a navigation tree based on the directory structure pointed to by `folder`.

- **Directories** become expandable nodes in the sidebar.
- **Markdown files (`.md`)** become selectable items.
- Hidden files/folders (starting with `.`) are ignored.
- Items are sorted with directories first, then files alphabetically.

**Example Structure:**
```
content/
├── 01-About.md
├── 02-Experience/
│   ├── Job1.md
│   └── Job2.md
├── 03-Projects/
│   ├── ProjectA.md
│   └── ProjectB.md
└── 04-Contact.md
```

## Keyboard Shortcuts

| Key | Action |
| :--- | :--- |
| `Tab` | Toggle focus between Sidebar and Content |
| `q` / `Ctrl+C` | Quit |

### Sidebar Navigation (When Focused)
| Key | Action |
| :--- | :--- |
| `j` / `↓` | Move cursor down |
| `k` / `↑` | Move cursor up |
| `Enter` / `Space` | Expand directory / Open file |

### Content Navigation (When Focused)
| Key | Action |
| :--- | :--- |
| `j` / `↓` | Scroll down |
| `k` / `↑` | Scroll up |
| `d` / `u` | Page down / Page up |

## Development

### Development Commands

```bash
# Format code
make fmt

# Run linter (requires golangci-lint)
make lint

# Clean build artifacts
make clean
```

## SSH Host Key

On first run, the server will generate an SSH host key at `.ssh/id_ed25519` if it doesn't exist. This key is used to identify your SSH server.

## Deployment

1. Build the binary for your target platform.
2. Copy the binary, your config file, and your content folder to your server.
3. Run the binary with appropriate host and port settings.
4. Configure your firewall to allow traffic on the SSH port.

### Example Systemd Service

```ini
[Unit]
Description=SSH Resume Server
After=network.target

[Service]
Type=simple
User=youruser
WorkingDirectory=/path/to/sshresume
ExecStart=/path/to/sshresume/bin/sshresume -host 0.0.0.0 -port 23234 -config /path/to/config.yaml
Restart=on-failure

[Install]
WantedBy=multi-user.target
```

## Technologies Used

- [Go](https://golang.org/)
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Glamour](https://github.com/charmbracelet/glamour) - Markdown rendering
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Styling
- [Wish](https://github.com/charmbracelet/wish) - SSH server
- [Bubbles](https://github.com/charmbracelet/bubbles) - TUI components

## License

This project is open source and available under the MIT License.

## Author

**Vinayak Goyal**
- Email: vinayaklovespizza@gmail.com
- GitHub: [github.com/vinayakankugoyal](https://github.com/vinayakankugoyal)
- LinkedIn: [linkedin.com/in/vinayakgoyal](https://linkedin.com/in/vinayakgoyal)
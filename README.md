# SSH Resume

An interactive resume/portfolio accessible via SSH, built with Go and [Charm](https://charm.sh/) libraries. Present your professional information in a beautiful terminal UI that users can access by simply SSH-ing to your server.

## Features

- 📄 **Interactive TUI** - Beautiful terminal user interface with tabbed navigation
- 🎨 **Markdown Support** - Write content in markdown with syntax highlighting via Glamour
- ⌨️ **Keyboard Navigation** - Navigate through tabs and scroll through content with keyboard shortcuts
- 🔐 **SSH Server** - Self-hosted SSH server with automatic host key generation

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

### Example

```bash
# Run on custom host and port with a specific config file
./bin/sshresume -host 0.0.0.0 -port 2222 -config ./myconfig.yaml
```

## Configuration

### Setting Up Your Config File

1. Create a `config.yaml` with your information:

```yaml
# Profile information displayed in the header
profile:
  name: "Your Name"
  email: "your.email@example.com"
  github: "github.com/yourusername"
  linkedin: "linkedin.com/in/yourusername"

# Tabs define the navigation structure
# Each tab has a name (displayed in UI) and file (markdown file to render)
tabs:
  - name: "Education"
    file: "education.md"
  - name: "Work"
    file: "work.md"
  - name: "Skills"
    file: "skills.md"
  - name: "Talks"
    file: "talks.md"
```

### Customizing Tabs

You can add, remove, or reorder tabs in the config file. Each tab requires:
- `name`: The display name shown in the UI
- `file`: Path to the markdown file (relative to where you run the binary, or absolute path)

Example with custom tabs:
```yaml
tabs:
  - name: "About"
    file: "content/about.md"
  - name: "Projects"
    file: "content/projects.md"
  - name: "Contact"
    file: "content/contact.md"
```

### Creating Content Files

Create markdown files referenced in your config. The content is rendered with [Glamour](https://github.com/charmbracelet/glamour) and supports standard markdown syntax

## Keyboard Shortcuts

When connected via SSH:

- `↑/↓` or scroll - Scroll through content
- `←/→` or `h/l` - Switch between tabs
- `Tab/Shift+Tab` - Switch between tabs
- `n/p` - Next/Previous tab
- `q` or `Ctrl+C` - Quit

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

1. Build the binary for your target platform
2. Copy the binary and content files to your server
3. Run the binary with appropriate host and port settings
4. Configure your firewall to allow traffic on the SSH port
5. (Optional) Set up a systemd service for automatic startup

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

- [Go](https://golang.org/) - Programming language
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Glamour](https://github.com/charmbracelet/glamour) - Markdown rendering
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Style definitions for terminal output
- [Wish](https://github.com/charmbracelet/wish) - SSH server toolkit
- [Bubbles](https://github.com/charmbracelet/bubbles) - TUI components (viewport)

## License

This project is open source and available under the MIT License.

## Author

**Vinayak Goyal**
- Email: vinayaklovespizza@gmail.com
- GitHub: [github.com/vinayakankugoyal](https://github.com/vinayakankugoyal)
- LinkedIn: [linkedin.com/in/vinayakgoyal](https://linkedin.com/in/vinayakgoyal)

## Acknowledgments

Built with the amazing [Charm](https://charm.sh/) suite of libraries for building terminal applications.

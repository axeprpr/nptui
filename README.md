# NPTUI - Netplan TUI

A terminal user interface (TUI) program for netplan, similar to nmtui, written in Go.

## Features

- 🖥️ User-friendly terminal interface (TUI)
- 🔌 Manage network interface configurations
- 🌐 Support for DHCP and static IP configuration
- 📡 Configure DNS servers
- ✅ Easy to use, similar to nmtui interface
- 🚀 Support for ARM64 and AMD64 architectures

## Requirements

- Linux operating system (with netplan)
- Root privileges (required to modify network configuration)
- netplan.io package

## Installation

### Install from DEB Package

#### AMD64 (x86_64)

```bash
sudo dpkg -i nptui-1.0.0-amd64.deb
```

#### ARM64

```bash
sudo dpkg -i nptui-1.0.0-arm64.deb
```

### Build from Source

#### Prerequisites

- Go 1.21 or higher

#### Build Steps

```bash
# Clone or enter project directory
cd nptui

# Install dependencies
make deps

# Build
make build

# Install to system
sudo make install
```

## Usage

Start the program (requires root privileges):

```bash
sudo nptui
```

### Main Features

1. **Edit Network Interfaces** - Configure network adapters

   - Select the network interface to configure
   - Choose DHCP or static IP configuration
   - Set IP address, gateway, DNS server

2. **Apply Configuration** - Apply changes

   - Save and apply netplan configuration

3. **Quit** - Exit the program

### Keyboard Shortcuts

- `Up/Down` or `↑/↓` - Navigate
- `Enter` - Select/Confirm
- `Tab` - Switch form fields
- `Esc` - Cancel/Return
- `q` - Quit (main menu)
- `b` - Back (interface list)

## Development

### Project Structure

```
nptui/
├── main.go           # Main program entry
├── go.mod            # Go module configuration
├── netplan/          # netplan configuration management
│   └── netplan.go
├── ui/               # TUI interface module
│   └── app.go
├── debian/           # Debian packaging configuration
│   ├── control-amd64
│   ├── control-arm64
│   ├── postinst
│   ├── postrm
│   └── copyright
├── Makefile          # Build script
└── README.md         # This document
```

### Build Targets

```bash
# Build for current architecture
make build

# Build AMD64 version
make build-amd64

# Build ARM64 version
make build-arm64

# Package AMD64 DEB
make deb-amd64

# Package ARM64 DEB
make deb-arm64

# Package all architectures
make deb-all

# Clean build files
make clean

# Run tests
make test
```

### Packaging Process

Generate DEB packages:

```bash
# Build DEB packages for all architectures
make deb-all

# Output files in build/ directory:
# - nptui-1.0.0-amd64.deb
# - nptui-1.0.0-arm64.deb
```

## Configuration Files

The program reads and modifies the following netplan configuration files:

- `/etc/netplan/*.yaml` - Existing configuration files
- `/etc/netplan/01-netcfg.yaml` - Default save location

## Important Notes

- Must be run with root privileges
- Configuration changes require application to take effect: `sudo netplan apply`
- It is recommended to backup existing configurations before making changes
- Incorrect network configuration may cause connection interruption

## License

MIT License

## Contributing

Issues and Pull Requests are welcome!

## Related Links

- [Netplan Documentation](https://netplan.io/)
- [tview Library](https://github.com/rivo/tview)

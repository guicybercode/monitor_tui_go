# SysTUI

<div align="center">

![SysTUI](https://img.shields.io/badge/SysTUI-Advanced%20System%20Monitor-blue?style=for-the-badge)
![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go)
![Rust](https://img.shields.io/badge/Rust-2021-000000?style=for-the-badge&logo=rust)
![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)

**A powerful terminal-based system monitoring and management tool for Linux**

</div>

## Overview

SysTUI is an advanced Text User Interface (TUI) application designed for Linux power users who prefer working directly in the terminal. Built with Go and Rust, it provides real-time system monitoring, process management, service control, network monitoring, package management, and efficient log analysis capabilities.

## Features

### 🎯 Core Functionality

- **📊 Interactive Dashboard** - Real-time monitoring of CPU, RAM, disk, and network usage
- **🔧 Process Management** - List, kill, and renice processes with an intuitive interface
- **⚙️ Systemd Service Control** - Start, stop, and restart services seamlessly
- **🌐 Network Monitor** - View active connections, open ports, and traffic per interface
- **📝 Config Editor** - Edit configuration files with syntax highlighting support
- **📦 Package Manager** - Visual interface for apt, dnf, or pacman with auto-detection
- **📄 Log Analysis** - Efficient log parsing powered by Rust for large files like `/var/log/syslog` and `journalctl`
- **📤 Report Export** - Generate detailed reports in JSON or Markdown format

### 🚀 Advanced Features

- **🔌 WebAssembly Plugin Support** - Extensible architecture via WASM plugins
- **🌍 Headless API Mode** - Expose metrics via REST API for integration with other tools
- **⚡ High Performance** - Rust-powered log parsing for handling large files efficiently
- **🎨 Beautiful TUI** - Modern interface built with Bubbletea and Lipgloss

## Installation

### Prerequisites

- Go 1.21 or later
- Rust toolchain (for building the log parser)
- Linux system with systemd
- Build tools (gcc, make)

### Build from Source

```bash
git clone https://github.com/guicybercode/systui.git
cd systui
make build
```

The binary will be available at `./systui`.

### Quick Start

```bash
./systui
```

For headless API mode:

```bash
./systui --headless --port 8080
```

## Usage

### TUI Mode

Launch SysTUI and use the following keyboard shortcuts:

- **1-6**: Switch between views (Dashboard, Processes, Services, Network, Packages, Logs)
- **j/k** or **↑/↓**: Navigate lists
- **d**: Kill selected process
- **s**: Start selected service
- **x**: Stop selected service
- **t**: Restart selected service
- **r**: Refresh current view
- **q** or **Ctrl+C**: Quit

### Headless API Mode

Start the API server:

```bash
./systui --headless --port 8080
```

Available endpoints:

- `GET /metrics` - Get system metrics (CPU, memory, disk)
- `GET /processes` - List all processes
- `GET /services` - List all systemd services
- `GET /network` - Get network statistics and connections
- `GET /report` - Generate full system report

Example:

```bash
curl http://localhost:8080/metrics
```

### Export Reports

Export functionality is available programmatically or can be integrated into the TUI. Reports include:

- CPU metrics and per-core usage
- Memory usage statistics
- Disk usage and partition information
- Process list with resource usage
- Service status
- Network interface statistics

## Architecture

```
SysTUI
├── Go Components (TUI & System Integration)
│   ├── Dashboard - Real-time metrics display
│   ├── Process Manager - Process listing and control
│   ├── Service Manager - Systemd integration
│   ├── Network Monitor - Connection and traffic monitoring
│   ├── Package Manager - Multi-distro package support
│   └── Log Viewer - Integration with Rust parser
│
├── Rust Components (Performance-Critical)
│   └── Log Parser - Efficient parsing with nom/regex
│       - Date filtering
│       - Severity filtering
│       - Regex pattern matching
│
├── WebAssembly Runtime
│   └── Plugin System - Extensible via WASM plugins
│
└── API Server
    └── REST API - Headless mode metrics endpoint
```

## Log Analysis

The Rust-powered log parser supports:

- **Large File Handling** - Efficiently processes files like `/var/log/syslog`
- **Date Filtering** - Filter logs by date range
- **Severity Filtering** - Filter by ERROR, WARN, INFO, DEBUG
- **Regex Search** - Advanced pattern matching

## Package Manager Support

SysTUI automatically detects your Linux distribution's package manager:

- **APT** - Debian/Ubuntu systems
- **DNF** - Fedora/RHEL systems
- **Pacman** - Arch Linux systems

## Development

### Project Structure

```
golang_project/
├── cmd/systui/          # Main entry point
├── internal/
│   ├── tui/             # TUI components
│   ├── system/          # System metrics collectors
│   ├── logparser/       # Go-Rust FFI bindings
│   ├── api/             # Headless API server
│   ├── plugins/wasm/    # WASM runtime
│   └── exports/         # Report export functionality
├── rust/                # Rust log parser module
└── plugins/example/     # Example WASM plugin
```

### Building

```bash
make build          # Build both Go and Rust components
make build-go       # Build only Go components
make build-rust     # Build only Rust components
make run            # Build and run
make clean          # Clean build artifacts
```

### Testing

```bash
make test           # Run all tests
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License.

## Acknowledgments

- Built with [Bubbletea](https://github.com/charmbracelet/bubbletea) and [Lipgloss](https://github.com/charmbracelet/lipgloss)
- System metrics powered by [gopsutil](https://github.com/shirou/gopsutil)
- Log parsing with [nom](https://github.com/Geal/nom) and [regex](https://github.com/rust-lang/regex)
- WebAssembly runtime by [wazero](https://github.com/tetratelabs/wazero)

---

**그들이 사도의 가르침을 받아 서로 교제하며 떡을 떼며 기도하기를 전혀 힘쓰니라** - 사도행전 2:42

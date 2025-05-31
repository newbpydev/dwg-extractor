# Go DWG Extractor

[![Go Version](https://img.shields.io/badge/Go-1.18+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Build Status](https://img.shields.io/badge/Build-Passing-brightgreen.svg)](https://github.com/remym/go-dwg-extractor)

A powerful command-line tool with a Terminal User Interface (TUI) for extracting data from DWG files. Convert DWG files to DXF format and interactively explore layers, blocks, text, and other CAD entities with advanced filtering, selection, and clipboard integration.

[Table of Contents](#table-of-contents)

[Table of Contents](#table-of-contents)
## Table of Contents

- [Description](#description)
- [Features](#features)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Usage](#usage)
- [Configuration](#configuration)
- [Build from Source](#build-from-source)
- [Examples](#examples)
- [Troubleshooting](#troubleshooting)
- [Contributing](#contributing)
- [License](#license)

## Description

Go DWG Extractor is a cross-platform application that provides both command-line and Terminal User Interface modes for working with DWG files. It leverages the ODA File Converter to convert DWG files to DXF format, then provides powerful tools for exploring and extracting data from the converted files.

The application features a modern TUI with keyboard shortcuts, layer filtering, entity selection, and clipboard integration, making it easy to work with CAD data in terminal environments.

## Features

- **Cross-platform support** - Works on Windows, Linux, and macOS
- **Terminal User Interface** - Modern, responsive TUI with keyboard navigation
- **DWG to DXF conversion** - Automatic conversion using ODA File Converter
- **Layer filtering** - Advanced search and filter capabilities for layers
- **Entity exploration** - Browse lines, circles, text, blocks, and polylines
- **Clipboard integration** - Copy selected data in multiple formats (text, CSV, JSON)
- **Keyboard shortcuts** - Efficient navigation with Ctrl+C, F1, Tab, and more
- **Multiple output formats** - Export data as text, CSV, or JSON
- **Error handling** - Comprehensive error reporting with recovery suggestions
- **Version information** - Built-in version tracking and build metadata

## Prerequisites

Before installing Go DWG Extractor, ensure you have the following:

### Required

- **Go 1.18+** - For building from source
- **ODA File Converter** - For DWG to DXF conversion
  - Download from: [https://www.opendesign.com/guestfiles/oda_file_converter](https://www.opendesign.com/guestfiles/oda_file_converter)
  - Install to default location or set `ODA_CONVERTER_PATH` environment variable

### Supported Platforms

- **Windows** - Windows 10/11 (amd64)
- **Linux** - Ubuntu 18.04+ or equivalent (amd64)
- **macOS** - macOS 10.15+ (amd64, arm64)

## Installation

### Option 1: Download Pre-built Binaries

1. Download the latest release for your platform from the [Releases](https://github.com/remym/go-dwg-extractor/releases) page
2. Extract the archive to your desired location
3. Add the executable to your PATH (optional)

### Option 2: Install from Source

```bash
# Clone the repository
git clone https://github.com/remym/go-dwg-extractor.git
cd go-dwg-extractor

# Build and install
make install
```

## Usage

Go DWG Extractor provides two main modes of operation:

### Command Line Mode

Extract data from DWG files and output to console:

```bash
# Extract data from a DWG file
./go-dwg-extractor extract -file sample.dwg

# Extract with custom output directory
./go-dwg-extractor extract -file sample.dwg -output ./output
```

### Terminal User Interface Mode

Launch the interactive TUI for exploring DWG data:

```bash
# Launch TUI with a specific DWG file
./go-dwg-extractor tui -file sample.dwg

# Launch TUI with sample data (for testing)
./go-dwg-extractor tui
```

### Version Information

```bash
# Show version information
./go-dwg-extractor version
```

### Help

```bash
# Show help information
./go-dwg-extractor help
```

## Configuration

### Environment Variables

The application can be configured using environment variables:

- **`ODA_CONVERTER_PATH`** - Path to ODA File Converter executable
  - Default: `C:\Program Files\ODA\ODAFileConverter 26.4.0\ODAFileConverter.exe` (Windows)
  - Example: `export ODA_CONVERTER_PATH="/usr/local/bin/ODAFileConverter"`

### Configuration File

The application automatically detects and validates the ODA File Converter installation. If the converter is not found in the default location, set the `ODA_CONVERTER_PATH` environment variable.

## Build from Source

### Using Make

```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Create distribution packages
make package

# Complete build pipeline
make all
```

### Using Makefile\n\nYou can also use the Makefile directly for various build tasks:\n\n```makefile\n# Available Makefile targets\nbuild:          # Build for current platform\nbuild-all:      # Build for all platforms\ntest:           # Run all tests\ntest-coverage:  # Run tests with coverage\npackage:        # Create distribution packages\nclean:          # Clean build artifacts\nhelp:           # Show available targets\n```

### Using Build Scripts

#### Linux/macOS (Bash)

```bash
# Run complete build pipeline
./scripts/build.sh all

# Build only
./scripts/build.sh build

# Clean previous builds
./scripts/build.sh clean
```

#### Windows (PowerShell)

```powershell
# Run complete build pipeline
.\scripts\build.ps1 all

# Build with specific version
.\scripts\build.ps1 all -Version "v1.0.0"

# Build only
.\scripts\build.ps1 build
```

### Manual Build

```bash
# Build with version information
go build -ldflags "-X main.version=v1.0.0 -X main.gitCommit=$(git rev-parse --short HEAD) -X main.buildTime=$(date -u '+%Y-%m-%dT%H:%M:%SZ')" -o go-dwg-extractor .
```

## Examples

### Basic Usage

```bash
# Extract layer information from a DWG file
./go-dwg-extractor extract -file architectural-plan.dwg

# Launch TUI to explore the file interactively
./go-dwg-extractor tui -file architectural-plan.dwg
```

### TUI Navigation

Once in the TUI mode, use these keyboard shortcuts:

- **Tab** - Switch between panes (search, layers, entities, details)
- **↑/↓** - Navigate within lists
- **Enter** - Select item or layer
- **Space** - Toggle layer visibility
- **Ctrl+C** - Copy selected items to clipboard
- **Ctrl+F** - Focus search input
- **F1** - Toggle help view
- **Escape** - Clear selection or go back
- **Ctrl+Q** - Quit application

### Clipboard Operations

1. Navigate to a layer or entity in the TUI
2. Press **Ctrl+C** to copy to clipboard
3. Data is copied in multiple formats (text, CSV, JSON)
4. Paste into your preferred text editor or spreadsheet

### Layer Filtering

In the TUI search box, you can use:

- **Text search**: Type layer name to filter
- **Status filter**: `on:true` or `on:false` to filter by visibility
- **Frozen filter**: `frozen:true` or `frozen:false` to filter by frozen status

## Troubleshooting

### Common Issues

#### ODA File Converter not found

**Error**: `ODA File Converter not found at default path`

**Solution**:
1. Download and install ODA File Converter from the official website
2. Set the `ODA_CONVERTER_PATH` environment variable:
   ```bash
   export ODA_CONVERTER_PATH="/path/to/ODAFileConverter"
   ```

#### Permission denied

**Error**: `Permission denied when accessing DWG file`

**Solution**:
1. Check file permissions: `chmod 644 your-file.dwg`
2. Ensure you have read access to the file
3. Run with appropriate user permissions

#### DWG file not supported

**Error**: `DWG file format not supported`

**Solution**:
1. Verify the file is a valid DWG file
2. Check if the DWG version is supported by ODA File Converter
3. Try converting the file to a newer DWG format using AutoCAD

#### Conversion failed

**Error**: `Failed to convert DWG to DXF`

**Solution**:
1. Check if the DWG file is corrupted
2. Verify ODA File Converter is properly installed
3. Ensure sufficient disk space for conversion
4. Check the application logs for detailed error information

### Debug Mode

For detailed debugging information, you can:

1. Check the application logs
2. Run with verbose output (if available)
3. Verify the ODA File Converter works independently

### Getting Help

If you encounter issues not covered here:

1. Check the [Issues](https://github.com/remym/go-dwg-extractor/issues) page
2. Create a new issue with:
   - Your operating system and version
   - Go DWG Extractor version (`./go-dwg-extractor version`)
   - Steps to reproduce the problem
   - Error messages or logs

## Contributing

We welcome contributions to Go DWG Extractor! Here's how you can help:

### Development Setup

1. Fork the repository
2. Clone your fork: `git clone https://github.com/yourusername/go-dwg-extractor.git`
3. Create a feature branch: `git checkout -b feature/your-feature-name`
4. Make your changes following the coding standards
5. Run tests: `make test`
6. Submit a pull request

### Coding Standards

- Follow Go best practices and conventions
- Write tests for new functionality (TDD approach)
- Maintain test coverage above 80%
- Use meaningful commit messages
- Update documentation as needed

### Testing

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run specific package tests
go test ./pkg/tui/ -v
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

**Go DWG Extractor** - Making CAD data accessible in terminal environments.

For more information, visit the [project repository](https://github.com/remym/go-dwg-extractor).

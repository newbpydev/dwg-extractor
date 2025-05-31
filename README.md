# Go DWG Extractor

A command-line tool with a Terminal User Interface (TUI) for extracting and managing data from DWG files.

## Features

- Convert DWG to DXF using ODA File Converter
- Extract layers, blocks, attributes, and text from DWG files
- Interactive TUI for data exploration
- Copy selected data to clipboard

## Prerequisites

- Go 1.18 or later
- ODA File Converter (for DWG to DXF conversion)

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/remym/go-dwg-extractor.git
   cd go-dwg-extractor
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Install the ODA File Converter from [ODA Teigha File Converter](https://www.opendesign.com/guestfiles/teighafileconverter)

## Configuration

Set the path to the ODA File Converter using the `ODA_CONVERTER_PATH` environment variable:

```bash
export ODA_CONVERTER_PATH="/path/to/FileConverterApp"  # Linux/macOS
# OR
set ODA_CONVERTER_PATH="C:\\Program Files\\ODA\\FileConverter\\FileConverterApp.exe"  # Windows
```

## Usage

```bash
go run main.go -file path/to/your/file.dwg
```

## Development

### Running Tests

```bash
go test ./...
```

### Project Structure

```
go-dwg-extractor/
├── cmd/                    # CLI commands
├── pkg/                    # Reusable packages
│   ├── config/             # Configuration management
│   ├── converter/          # DWG to DXF conversion
│   ├── dxfparser/          # DXF parsing logic
│   ├── tui/                # Terminal UI components
│   ├── clipboard/          # Clipboard interaction
│   └── data/               # Data structures
├── assets/                 # Sample files
├── go.mod
├── go.sum
└── README.md
```

## License

MIT

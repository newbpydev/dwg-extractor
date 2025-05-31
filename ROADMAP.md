# Project: Go DWG Extractor with TUI

**Objective:** Develop a command-line Go application with a Terminal User Interface (TUI) that allows users to specify a DWG file, view extracted data (layers, blocks, attributes, text, etc.), select specific items, and copy them to the clipboard. The application will use an external ODA File Converter for DWG to DXF conversion.

**Core Technologies:**
* **Go:** For the main application logic.
* **ODA File Converter:** External CLI tool for DWG to DXF conversion.
* **Go DXF Library:** `github.com/yofu/go-dxf` (as per the provided guide).
* **Go TUI Library:** A library like `github.com/rivo/tview` or `github.com/charmbracelet/bubbletea`. (The plan will be generic, the agent can choose).
* **Go Clipboard Library:** `github.com/atotto/clipboard` or similar.
* **Go Standard Libraries:** `os/exec`, `flag`, `fmt`, `log`, `encoding/json`, `testing`.

**Development Approach:** SCRUM with Test-Driven Development (TDD).

## 1. Project Setup and Initial Configuration (Sprint 0)

**Goal:** Establish the project structure, install necessary tools, and configure the environment.

**Tasks:**

1.  **[x] Task 0.1: Environment Setup**
    * [x] Install Go (version 1.18+).
    * [x] Install the ODA File Converter.
        * Installed at: `C:\Program Files\ODA\ODAFileConverter 26.4.0`
        * Note: Verify command-line conversion works:
            ```bash
            "C:\Program Files\ODA\ODAFileConverter 26.4.0\ODAFileConverter.exe" -i "path\to\sample.dwg" -o "path\to\output_dir" -f DXF -v ACAD2018
            ```
    * [x] Sample DWG file (`sample.dwg`) exists in the project root.
    * [ ] (Optional) Manually convert `sample.dwg` to `sample.dxf` using the ODA File Converter to have a reference DXF file.

2.  **[x] Task 0.2: Project Directory Structure**
    * Created the main project directory: `dwg-extractor`
    * Initialized Go module: `github.com/remym/go-dwg-extractor`
    * Created the following directory structure:
        ```
        go-dwg-extractor/
        ├── main.go                 // Main application entry point
        ├── cmd/                    // CLI argument parsing and TUI initialization
        │   └── root.go
        ├── pkg/
        │   ├── config/             // Configuration management
        │   │   └── config.go
        │   │   └── config_test.go
        │   ├── converter/          // DWG to DXF conversion logic
        │   │   └── converter.go
        │   │   └── converter_test.go
        │   ├── dxfparser/          // DXF parsing logic
        │   │   └── parser.go
        │   │   └── parser_test.go
        │   ├── tui/                  // Terminal User Interface components
        │   │   └── view.go
        │   │   └── (TUI tests might be more focused on underlying data models)
        │   ├── clipboard/          // Clipboard interaction
        │   │   └── clipboard.go
        │   │   └── clipboard_test.go
        │   └── data/               // Data structures for DXF entities
        │       └── models.go
        │       └── models_test.go
        ├── assets/                 // Sample files, bundled tools (optional)
        │   └── sample_files/
        │       └── sample.dwg
        │       └── sample.dxf      // Pre-converted reference
        │   └── oda_converter/      // (Optional: if bundling the converter)
        ├── go.mod
        ├── go.sum
        └── README.md
        ```

3.  **[x] Task 0.3: Install Initial Go Dependencies**
    * [x] DXF parsing library: `github.com/yofu/go-dxf`
    * [x] TUI library: `github.com/rivo/tview`
    * [x] Clipboard library: `github.com/atotto/clipboard`

4.  **[x] Task 0.4: Basic `main.go` and CLI Flag Parsing**
    * **TDD:**
        * [x] **Write failing test:** In `cmd/root_test.go`, test for parsing a `-file` flag.
        * [x] **Implement:** In `cmd/root.go`, use the `flag` package to accept a DWG file path (e.g., `-file="path/to/drawing.dwg"`).
        * [x] **Write failing test:** Test for error handling if the file flag is not provided.
        * [x] **Implement:** Add logic to handle missing file flag.
    * [x] In `main.go`, call the command execution logic from `cmd/root.go`.

## 2. Configuration Management (Sprint 1) - COMPLETED ✅

**Goal:** Implement robust configuration handling, especially for the ODA File Converter path.

**Tasks:**

1.  **[x] Task 1.1: Define Configuration Structure**
    * [x] In `pkg/config/config.go`, defined `AppConfig` struct with `ODAConverterPath` field.
    * [x] Added error types and validation logic.

2.  **[x] Task 1.2: Load Configuration**
    * **TDD:**
        * [x] **Write failing test:** In `pkg/config/config_test.go`, test loading from environment variable.
        * [x] **Implement:** Function to load from `ODA_CONVERTER_PATH` environment variable.
        * [x] **Write failing test:** Test default path when env var is not set.
        * [x] **Implement:** Default path fallback logic.
        * [x] **Write failing test:** Test path validation.
        * [x] **Implement:** Path validation to check if file exists and is a regular file.

3.  **[x] Task 1.3: Integrate Configuration Loading**
    * [x] Updated `cmd/root.go` to load and validate configuration at startup.
    * [x] Made configuration available to other packages via package variable.
    * [x] Updated tests to work with the new configuration system.

## 3. DWG to DXF Conversion (Sprint 2)

**Goal:** Implement the functionality to convert a given DWG file to DXF using the ODA File Converter.

**Tasks:**

1.  **[x] Task 2.1: Define Converter Interface and Struct**
    * [x] Created `DWGConverter` interface in `pkg/converter/converter.go` with `ConvertToDXF` method.
    * [x] Implemented `odaconverter` struct that holds the ODA converter path.
    * [x] Added `NewDWGConverter` constructor with input validation.
    * [x] Added comprehensive tests for the interface and implementation.

2.  **[ ] Task 2.2: Implement Conversion Logic**
    * **TDD:**
        * [ ] **Write failing test:** In `pkg/converter/converter_test.go`, mock `os/exec` or use a test script to simulate the ODA converter. Test successful conversion:
            * Input: valid DWG path, valid output directory.
            * Expected: path to the generated DXF file, no error.
        * [ ] **Implement:** The `ConvertToDXF` method.
            * Use `os/exec` to call the ODA File Converter CLI.
                * Command: `ODAConverterPath -i <dwgPath> -o <outputDir> -f DXF -v ACAD2018` (or a configurable DXF version).
            * Ensure the output DXF file name is predictable (e.g., same base name as DWG, but with `.dxf` extension).
            * Handle temporary output directory creation and cleanup if necessary.
        * [ ] **Write failing test:** Test for ODA converter command failure (e.g., converter not found, invalid DWG).
        * [ ] **Implement:** Error handling for `cmd.Run()` or `cmd.CombinedOutput()`. Capture and return stderr from the converter.
        * [ ] **Write failing test:** Test for input file not found.
        * [ ] **Implement:** Pre-check for DWG file existence.
        * [ ] **Write failing test:** Test for output directory creation failure.
        * [ ] **Implement:** Error handling for `os.MkdirAll`.

3.  **[ ] Task 2.3: Integrate Conversion into Main Flow**
    * In `cmd/root.go`, after parsing the DWG file path flag and loading config:
        * Instantiate the converter.
        * Call `ConvertToDXF` to convert the input DWG.
        * For now, log the path to the generated DXF file or an error message.
        * Ensure the temporary DXF file is cleaned up after processing (or placed in a designated temp area).

## 4. DXF Parsing (Sprint 3)

**Goal:** Implement parsing of the generated DXF file to extract relevant data.

**Tasks:**

1.  **[ ] Task 3.1: Define Data Structures for Extracted Information**
    * In `pkg/data/models.go`, define Go structs to represent the information you want to extract. Examples:
        ```go
        package data

        // Point defines a 2D or 3D point.
        type Point struct {
            X, Y, Z float64
        }

        // LayerInfo holds information about a DXF layer.
        type LayerInfo struct {
            Name    string
            Color   int
            IsOn    bool
            IsFrozen bool
            // Add other relevant properties
        }

        // AttributeInfo holds information about a block attribute.
        type AttributeInfo struct {
            Tag   string
            Value string
            // Add other relevant properties like position if needed
        }

        // BlockInfo holds information about a block instance (Insert entity).
        type BlockInfo struct {
            Name           string
            Layer          string
            InsertionPoint Point
            Attributes     []AttributeInfo
            // Add scale, rotation if needed
        }

        // TextInfo holds information about a Text entity.
        type TextInfo struct {
            Value          string
            Layer          string
            InsertionPoint Point
            Height         float64
            // Add other relevant properties
        }

        // LineInfo, CircleInfo, PolylineInfo etc. can be added as needed.

        // ExtractedData holds all data parsed from the DXF.
        type ExtractedData struct {
            DXFVersion string
            Layers     []LayerInfo
            Blocks     []BlockInfo
            Texts      []TextInfo
            // Add other entity lists (Lines, Circles, Polylines etc.)
        }
        ```
    * **TDD:**
        * [ ] **Write failing tests:** In `pkg/data/models_test.go`, write simple tests for these structs if they have any methods or complex initialization. (Often, these are just data holders).

2.  **[ ] Task 3.2: Implement DXF Parser Logic**
    * In `pkg/dxfparser/parser.go`, create a `Parser` struct/functions.
    * **TDD:**
        * [ ] **Write failing test:** In `pkg/dxfparser/parser_test.go`, test parsing a known, simple DXF file (use `assets/sample_files/sample.dxf`).
            * Input: path to `sample.dxf`.
            * Expected: correctly populated `data.ExtractedData` struct.
        * [ ] **Implement:** A function `ParseDXF(dxfPath string) (*data.ExtractedData, error)`.
            * Use `os.Open` to open the DXF file.
            * Use `github.com/yofu/go-dxf`'s `dxf.Parse(file)` method.
            * Iterate through `d.Header` to get version (`header.ACADVER`).
            * Iterate through `d.Tables.Layer.Layers` to extract `LayerInfo`.
            * Iterate through `d.Entities`:
                * Use a type switch (`switch e := ent.(type)`) to identify entity types:
                    * `*entities.Insert`: Extract `BlockInfo` and its `Attributes`.
                    * `*entities.Text`, `*entities.MText`: Extract `TextInfo`.
                    * `*entities.Line`, `*entities.Circle`, `*entities.LwPolyline`, etc.: Extract relevant geometric info if required by the client.
                * Populate your `data.ExtractedData` struct.
        * [ ] **Write failing test:** Test for DXF file not found.
        * [ ] **Implement:** Error handling for file opening.
        * [ ] **Write failing test:** Test for DXF parsing errors (malformed DXF).
        * [ ] **Implement:** Error handling for `dxf.Parse()`.

3.  **[ ] Task 3.3: Integrate DXF Parsing into Main Flow**
    * In `cmd/root.go`, after successful DWG to DXF conversion:
        * Call `dxfparser.ParseDXF()` with the path to the generated DXF file.
        * For now, print the extracted data to the console (e.g., using `fmt.Printf("%+v\n", extractedData)` or marshal to JSON and print) to verify.
        * Handle and log any errors from parsing.

## 5. Basic TUI Implementation (Sprint 4)

**Goal:** Display the extracted DXF data in a basic, non-interactive TUI.

**Tasks:**

1.  **[ ] Task 4.1: Choose and Set Up TUI Library**
    * Confirm the choice of TUI library (e.g., `tview`).
    * In `pkg/tui/view.go`, set up the basic TUI application structure.
        * Example with `tview`:
            ```go
            package tui

            import (
                "[github.com/yourusername/go-dwg-extractor/pkg/data](https://github.com/yourusername/go-dwg-extractor/pkg/data)" // Your data models
                "[github.com/rivo/tview](https://github.com/rivo/tview)"
            )

            type AppTUI struct {
                app      *tview.Application
                // Add other TUI components like lists, text views here
                // e.g., layersList *tview.List
                // e.g., detailsView *tview.TextView
            }

            func NewAppTUI() *AppTUI {
                // Initialize TUI components
                return &AppTUI{
                    app: tview.NewApplication(),
                }
            }

            func (t *AppTUI) Run(extractedData *data.ExtractedData) error {
                // Populate TUI components with extractedData
                // Set layout
                // return t.app.SetRoot( /* your root layout */, true).Run()
                return nil // Placeholder
            }
            ```

2.  **[ ] Task 4.2: Design TUI Layout**
    * Sketch a simple layout. For example:
        * A main list view to show categories (Layers, Blocks, Texts).
        * A detail view to show items within a selected category.
        * A status bar for messages or instructions.

3.  **[ ] Task 4.3: Display Extracted Data (Read-Only)**
    * **TDD (for data formatting logic, not direct TUI rendering):**
        * [ ] **Write failing test:** For functions that format `data.LayerInfo`, `data.BlockInfo`, etc., into strings suitable for TUI list items.
        * [ ] **Implement:** These formatting functions.
    * In `pkg/tui/view.go`:
        * Modify `Run` method to accept `*data.ExtractedData`.
        * Create TUI components (e.g., `tview.List` for categories, another `tview.List` or `tview.Table` for items).
        * Populate these components with the data from `extractedData`.
            * Example: A list showing "Layers", "Blocks", "Texts".
            * When "Layers" is (conceptually) selected, another list shows individual layer names.
    * Focus on displaying data first; interactivity comes next.

4.  **[ ] Task 4.4: Integrate TUI into Main Flow**
    * In `cmd/root.go`, after successful DXF parsing:
        * Instantiate `AppTUI`.
        * Call `tuiApp.Run(extractedData)`.
        * Handle any errors from the TUI application.

## 6. TUI Interactivity and Selection (Sprint 5)

**Goal:** Enable navigation, selection of data items in the TUI.

**Tasks:**

1.  **[ ] Task 5.1: Implement Navigation**
    * In `pkg/tui/view.go`:
        * Allow navigation between different panes/views (e.g., from category list to item list).
        * Enable keyboard navigation (up/down arrows) within lists.
        * When a category (e.g., "Layers") is selected in the main list, update the item list/details view to show all layers.
        * When an item (e.g., a specific layer) is selected, update a details pane to show more information about it (if applicable).

2.  **[ ] Task 5.2: Implement Item Selection Logic**
    * Allow users to select/deselect multiple items in the item list (e.g., using Spacebar).
    * Maintain a list of currently selected items internally in `AppTUI`.
    * Visually indicate selected items (e.g., change prefix, color).

3.  **[ ] Task 5.3: Data Display Refinements**
    * Improve how data is presented. For `BlockInfo`, you might show `BlockName` in the list, and `Attributes` in a separate detail view when the block is selected.
    * For `AttributeInfo`, display `Tag: Value`.
    * Consider what information is most useful for the client to see and copy.

## 7. Clipboard Functionality (Sprint 6)

**Goal:** Allow the user to copy selected data to the system clipboard.

**Tasks:**

1.  **[ ] Task 6.1: Implement Clipboard Interaction Logic**
    * In `pkg/clipboard/clipboard.go`:
        * Create a function `CopyToClipboard(text string) error`.
        * Use `github.com/atotto/clipboard`'s `clipboard.WriteAll(text)` method.
    * **TDD:**
        * [ ] **Write failing test:** In `pkg/clipboard/clipboard_test.go`. This might be tricky to test without actual clipboard interaction. You could test that the `clipboard.WriteAll` function is called with the correct string (requires mocking or a more abstract interface for the clipboard). Alternatively, test the formatting of the string to be copied.

2.  **[ ] Task 6.2: Format Selected Data for Clipboard**
    * In `pkg/tui/view.go` (or a helper function):
        * Create a function that takes the list of selected data items (e.g., `[]data.BlockInfo`, `[]data.TextInfo`).
        * Formats this data into a single string suitable for pasting into a text editor or spreadsheet.
            * Consider a simple format, e.g., one item per line, attributes tab-separated or similar.
            * `BlockName: MyBlock, Layer: Layer1, InsertionPoint: (10.0, 20.0), Attributes: [Tag1:Val1, Tag2:Val2]`
            * `Text: "Hello World", Layer: Notes, InsertionPoint: (5.0, 15.0)`
    * **TDD:**
        * [ ] **Write failing test:** For this formatting logic.
        * [ ] **Implement:** The formatting function.

3.  **[ ] Task 6.3: Integrate Clipboard into TUI**
    * In `pkg/tui/view.go`:
        * Add a key binding (e.g., 'c' or 'Ctrl+C').
        * When the key is pressed:
            * Get the currently selected items.
            * Format them into a string.
            * Call `clipboard.CopyToClipboard()` with the formatted string.
            * Display a status message in the TUI (e.g., "Selected items copied to clipboard!").

## 8. Error Handling, Refinements, and Packaging (Sprint 7)

**Goal:** Enhance error handling, refine the TUI, add help text, and prepare for distribution.

**Tasks:**

1.  **[ ] Task 7.1: Comprehensive Error Handling**
    * Review all parts of the application for error handling.
    * Ensure errors from file operations, conversion, parsing, and TUI are caught and displayed gracefully to the user (e.g., in a TUI status bar or modal dialog).
    * Provide informative error messages.

2.  **[ ] Task 7.2: TUI Enhancements**
    * Add a help view/popup showing key bindings (navigation, selection, copy, quit).
    * Improve visual styling if needed (colors, borders).
    * Ensure smooth quitting (e.g., 'q' or 'Ctrl+Q').

3.  **[ ] Task 7.3: Bundling ODA Converter (Optional but Recommended)**
    * **Decision:** Decide whether to bundle the ODA File Converter with your application or require users to install it separately.
    * **If bundling:**
        * Place the converter executables for different OSes (Windows, Linux, macOS) in `assets/oda_converter/`.
        * Modify the `pkg/config/config.go` to first check for a bundled converter (e.g., relative to the application executable) before checking environment variables or config files.
        * Go's `os.Executable()` can help find the path of the running binary.
        * This significantly improves ease of use for the client.

4.  **[ ] Task 7.4: Build and Packaging**
    * Create build scripts (e.g., `Makefile` or simple shell scripts) for compiling the Go application for different platforms (cross-compilation).
        * `GOOS=windows GOARCH=amd64 go build -o go-dwg-extractor.exe .`
        * `GOOS=linux GOARCH=amd64 go build -o go-dwg-extractor .`
        * `GOOS=darwin GOARCH=amd64 go build -o go-dwg-extractor .`
    * If bundling the converter, ensure the build script copies the converter alongside the Go executable into a distribution package (e.g., a ZIP file).

5.  **[ ] Task 7.5: README Documentation**
    * Update `README.md` with:
        * Project description.
        * Prerequisites (if ODA converter is not bundled).
        * Installation instructions.
        * Usage instructions (CLI flags, TUI navigation, key bindings).
        * Configuration options (environment variables, config file).
        * Troubleshooting common issues.

## 9. Testing and Quality Assurance (Ongoing & Final Sprint)

**Goal:** Ensure the application is robust, reliable, and meets requirements.

**Tasks:**

1.  **[ ] Task 8.1: Unit Test Coverage Review**
    * Ensure all critical logic in `pkg/` subdirectories has good unit test coverage.
    * `go test ./... -cover`

2.  **[ ] Task 8.2: Manual End-to-End Testing**
    * Test with various DWG files:
        * Simple and complex files.
        * Files with many layers, blocks, attributes.
        * Files with different DWG versions (if possible, though ODA converter handles this).
        * Empty or corrupted DWG files (to test error handling).
    * Test on target operating systems (Windows, Linux, macOS).
    * Verify TUI behavior, selection, and clipboard functionality.

3.  **[ ] Task 8.3: User Acceptance Testing (UAT) with Client (Simulated)**
    * If this were for a real client, this is where they would test. For the AI agent, ensure all specified requirements are met.

## Detailed Notes for the AI Agent:

* **TDD Process:**
    1.  **Red:** Write a small test for a specific piece of functionality. Ensure it fails because the functionality doesn't exist yet.
    2.  **Green:** Write the minimum amount of code to make the test pass.
    3.  **Refactor:** Improve the code (clarity, efficiency, remove duplication) while ensuring all tests still pass.
* **ODA Converter Path:** The `ODA_CONVERTER_PATH` environment variable is a good starting point. A config file (`~/.config/go-dwg-extractor/config.json` or `./config.json`) is a more permanent solution. Bundling is the most user-friendly.
* **TUI Library Choice:**
    * `tview`: Feature-rich, component-based, good for complex layouts.
    * `bubbletea`: Based on The Elm Architecture, good for state management, more functional approach.
    * Choose one and stick to its patterns. The plan is generic enough.
* **Error Handling:** Use `log.Fatalf` for unrecoverable errors at startup. In TUI, display errors in a status bar or modal. Return errors up the call stack.
* **DXF Data Extraction:** The client's specific needs determine which entities and attributes are crucial. The provided `data.ExtractedData` is a suggestion; expand it as needed. For example, if polyline coordinates are needed, add `PolylineInfo` and parse `*entities.LwPolyline` or `*entities.Polyline`.
* **File Paths:** Be careful with file path separators on different OSes. Use `filepath.Join` for constructing paths.
* **Temporary Files:** If the DXF file is temporary, ensure it's cleaned up using `defer os.Remove(dxfFilePath)` or similar. A dedicated temporary directory managed by the app is also a good idea. `os.MkdirTemp` can be useful.
* **Clipboard:** Test clipboard functionality thoroughly, as it can be platform-dependent.
* **User Experience (UX) in TUI:**
    * Keep the TUI clean and intuitive.
    * Provide clear instructions or a help screen.
    * Ensure responsive feedback (e.g., "Copied!" message).

This detailed plan should guide the AI agent through the development process, ensuring all requirements are met with a TDD approach.
This plan provides a structured, phased approach to building your Go DWG extractor. Remember that each "Sprint" is a time-boxed iteration, and the tasks within can be adjusted as development progresses. The key is to maintain the TDD cycle and focus on delivering working software incrementally.
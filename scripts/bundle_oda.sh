#!/bin/bash

# Bundle ODA File Converter Script
# This script helps bundle ODA File Converter executables for distribution

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to show usage
show_usage() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Bundle ODA File Converter executables for distribution"
    echo ""
    echo "Options:"
    echo "  -h, --help              Show this help message"
    echo "  -s, --source PATH       Source directory containing ODA converters"
    echo "  -d, --destination PATH  Destination directory (default: assets/oda_converter)"
    echo "  -p, --platform PLATFORM Platform to bundle (windows|linux|darwin|all)"
    echo "  -v, --verbose           Verbose output"
    echo ""
    echo "Examples:"
    echo "  $0 --source /path/to/oda --platform all"
    echo "  $0 -s ~/Downloads/ODA -p windows -d ./bundled"
    echo ""
    echo "Note: You must have the ODA File Converter installed or downloaded"
    echo "      from https://www.opendesign.com/guestfiles/oda_file_converter"
}

# Default values
SOURCE_DIR=""
DEST_DIR="assets/oda_converter"
PLATFORM="all"
VERBOSE=false

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            show_usage
            exit 0
            ;;
        -s|--source)
            SOURCE_DIR="$2"
            shift 2
            ;;
        -d|--destination)
            DEST_DIR="$2"
            shift 2
            ;;
        -p|--platform)
            PLATFORM="$2"
            shift 2
            ;;
        -v|--verbose)
            VERBOSE=true
            shift
            ;;
        *)
            print_error "Unknown option: $1"
            show_usage
            exit 1
            ;;
    esac
done

# Function to validate platform
validate_platform() {
    case $1 in
        windows|linux|darwin|all)
            return 0
            ;;
        *)
            print_error "Invalid platform: $1. Must be one of: windows, linux, darwin, all"
            exit 1
            ;;
    esac
}

# Function to find ODA converter executable
find_oda_executable() {
    local platform=$1
    local search_dir=$2
    
    case $platform in
        windows)
            find "$search_dir" -name "ODAFileConverter.exe" -type f 2>/dev/null | head -1
            ;;
        linux|darwin)
            find "$search_dir" -name "ODAFileConverter" -type f -executable 2>/dev/null | head -1
            ;;
    esac
}

# Function to copy ODA converter
copy_oda_converter() {
    local platform=$1
    local source_file=$2
    local dest_dir="$DEST_DIR/$platform"
    
    # Create destination directory
    mkdir -p "$dest_dir"
    
    # Determine target filename
    local target_file
    case $platform in
        windows)
            target_file="$dest_dir/ODAFileConverter.exe"
            ;;
        linux|darwin)
            target_file="$dest_dir/ODAFileConverter"
            ;;
    esac
    
    # Copy the file
    if [[ $VERBOSE == true ]]; then
        print_status "Copying $source_file to $target_file"
    fi
    
    cp "$source_file" "$target_file"
    
    # Make executable on Unix systems
    if [[ $platform != "windows" ]]; then
        chmod +x "$target_file"
    fi
    
    print_success "Bundled $platform ODA converter to $target_file"
}

# Function to bundle platform
bundle_platform() {
    local platform=$1
    
    print_status "Bundling ODA converter for $platform..."
    
    # Find the executable
    local oda_executable
    if [[ -n "$SOURCE_DIR" ]]; then
        oda_executable=$(find_oda_executable "$platform" "$SOURCE_DIR")
    else
        # Try common installation paths
        case $platform in
            windows)
                # Common Windows paths
                for path in \
                    "/c/Program Files/ODA/ODAFileConverter"* \
                    "/c/Program Files (x86)/ODA/ODAFileConverter"* \
                    "$HOME/Downloads/ODAFileConverter"*; do
                    if [[ -d "$path" ]]; then
                        oda_executable=$(find_oda_executable "$platform" "$path")
                        [[ -n "$oda_executable" ]] && break
                    fi
                done
                ;;
            linux)
                for path in \
                    "/usr/local/bin" \
                    "/usr/bin" \
                    "$HOME/bin" \
                    "$HOME/Downloads/ODAFileConverter"*; do
                    if [[ -d "$path" ]]; then
                        oda_executable=$(find_oda_executable "$platform" "$path")
                        [[ -n "$oda_executable" ]] && break
                    fi
                done
                ;;
            darwin)
                for path in \
                    "/usr/local/bin" \
                    "/Applications/ODAFileConverter"* \
                    "$HOME/Applications/ODAFileConverter"* \
                    "$HOME/Downloads/ODAFileConverter"*; do
                    if [[ -d "$path" ]]; then
                        oda_executable=$(find_oda_executable "$platform" "$path")
                        [[ -n "$oda_executable" ]] && break
                    fi
                done
                ;;
        esac
    fi
    
    if [[ -z "$oda_executable" ]]; then
        print_warning "ODA File Converter not found for $platform"
        print_warning "Please specify the source directory with --source option"
        return 1
    fi
    
    if [[ ! -f "$oda_executable" ]]; then
        print_error "ODA executable not found: $oda_executable"
        return 1
    fi
    
    # Copy the converter
    copy_oda_converter "$platform" "$oda_executable"
    return 0
}

# Main execution
print_status "Starting ODA File Converter bundling process..."

# Validate inputs
validate_platform "$PLATFORM"

if [[ -n "$SOURCE_DIR" && ! -d "$SOURCE_DIR" ]]; then
    print_error "Source directory does not exist: $SOURCE_DIR"
    exit 1
fi

# Create base destination directory
mkdir -p "$DEST_DIR"

# Bundle converters
success_count=0
total_count=0

if [[ "$PLATFORM" == "all" ]]; then
    platforms=("windows" "linux" "darwin")
else
    platforms=("$PLATFORM")
fi

for platform in "${platforms[@]}"; do
    total_count=$((total_count + 1))
    if bundle_platform "$platform"; then
        success_count=$((success_count + 1))
    fi
done

# Summary
echo ""
print_status "Bundling Summary:"
print_status "Successfully bundled: $success_count/$total_count platforms"

if [[ $success_count -gt 0 ]]; then
    echo ""
    print_success "ODA File Converter bundling completed!"
    print_status "Bundled converters are available in: $DEST_DIR"
    print_status "You can now build and distribute your application with bundled converters."
    echo ""
    print_status "Next steps:"
    print_status "1. Test the bundled converters: go test ./pkg/config/"
    print_status "2. Build your application: make build-all"
    print_status "3. Package for distribution: make package"
else
    echo ""
    print_error "No converters were successfully bundled."
    print_error "Please check the source paths and try again."
    exit 1
fi 
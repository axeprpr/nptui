#!/bin/bash
# Build script for nptui - builds both ARM and x86 DEB packages

set -e

echo "======================================"
echo "NPTUI Build Script"
echo "======================================"
echo ""

# Check if dpkg-deb is available
if ! command -v dpkg-deb &> /dev/null; then
    echo "Error: dpkg-deb is not installed."
    echo "Please install it: sudo apt-get install dpkg-dev"
    exit 1
fi

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed."
    echo "Please install Go 1.21 or later."
    exit 1
fi

# Get Go version
GO_VERSION=$(go version | awk '{print $3}')
echo "Using Go version: $GO_VERSION"
echo ""

# Clean previous builds
echo "Cleaning previous builds..."
make clean
echo ""

# Download dependencies
echo "Downloading dependencies..."
make deps
echo ""

# Build all DEB packages
echo "Building DEB packages for all architectures..."
make deb-all
echo ""

echo "======================================"
echo "Build completed successfully!"
echo "======================================"
echo ""
echo "Generated packages:"
ls -lh build/*.deb
echo ""
echo "To install:"
echo "  AMD64: sudo dpkg -i build/nptui-1.0.0-amd64.deb"
echo "  ARM64: sudo dpkg -i build/nptui-1.0.0-arm64.deb"
echo ""


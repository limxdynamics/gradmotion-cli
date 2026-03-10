#!/usr/bin/env bash
set -euo pipefail

VERSION="${VERSION:-dev}"
COMMIT="${COMMIT:-unknown}"
DATE="${DATE:-$(date -u +%Y-%m-%dT%H:%M:%SZ)}"

OUTPUT_DIR="dist"
rm -rf "$OUTPUT_DIR"
mkdir -p "$OUTPUT_DIR"

LDFLAGS="-s -w -X main.version=${VERSION} -X main.commit=${COMMIT} -X main.date=${DATE}"

echo "Building gm CLI ${VERSION}..."

# macOS (Intel)
echo "  → darwin/amd64"
GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="$LDFLAGS" -o "$OUTPUT_DIR/gm-darwin-amd64" ./cmd/gradmotion

# macOS (Apple Silicon)
echo "  → darwin/arm64"
GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="$LDFLAGS" -o "$OUTPUT_DIR/gm-darwin-arm64" ./cmd/gradmotion

# Linux (amd64)
echo "  → linux/amd64"
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="$LDFLAGS" -o "$OUTPUT_DIR/gm-linux-amd64" ./cmd/gradmotion

# Linux (arm64)
echo "  → linux/arm64"
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="$LDFLAGS" -o "$OUTPUT_DIR/gm-linux-arm64" ./cmd/gradmotion

# Windows (amd64)
echo "  → windows/amd64"
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="$LDFLAGS" -o "$OUTPUT_DIR/gm-windows-amd64.exe" ./cmd/gradmotion

# Windows (arm64)
echo "  → windows/arm64"
GOOS=windows GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="$LDFLAGS" -o "$OUTPUT_DIR/gm-windows-arm64.exe" ./cmd/gradmotion

echo "✓ Build complete. Binaries in $OUTPUT_DIR/"
ls -lh "$OUTPUT_DIR/"

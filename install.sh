#!/bin/bash
set -e

# Cleanup function
cleanup() {
  local exit_code=$?
  if [[ -d "tmpshopcli" ]]; then
    cd ..
    rm -rf tmpshopcli
  fi
  exit $exit_code
}

# Set trap before creating temp directory
trap cleanup EXIT ERR INT TERM

OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

# Match goreleaser archive names (x86_64, arm64)
case $ARCH in
    x86_64|amd64) ARCH="x86_64" ;;
    aarch64|arm64) ARCH="arm64" ;;
    *) echo "Error: unsupported architecture ${ARCH}"; exit 1 ;;
esac

echo "OS: ${OS}, ARCH: ${ARCH}"

DOWNLOAD_URL="https://github.com/guitarkeegan/shop-cli/releases/latest/download"
ARCHIVE_NAME="shop-cli_${OS}_${ARCH}.tar.gz"

mkdir -p tmpshopcli
cd tmpshopcli

echo "downloading ${ARCHIVE_NAME}..."
curl -fsSL -o "$ARCHIVE_NAME" "${DOWNLOAD_URL}/${ARCHIVE_NAME}" || {
  echo "Error: ${ARCHIVE_NAME} not found in latest release"
  exit 1
}

# Checksum file name includes the version, so look it up from the release page redirect
VERSION=$(curl -fsSLI -o /dev/null -w '%{url_effective}' \
  https://github.com/guitarkeegan/shop-cli/releases/latest | sed 's#.*/tag/v##')

curl -fsSL -o "checksums.txt" "${DOWNLOAD_URL}/shop-cli_${VERSION}_checksums.txt" || {
  echo "checksum not found!"
  exit 1
}

if command -v sha256sum &> /dev/null; then
  SHA_CMD="sha256sum"
else
  SHA_CMD="shasum -a 256"
fi

echo "Verifying checksum..."
grep "$ARCHIVE_NAME" checksums.txt | $SHA_CMD -c - || {
  echo "✗ Checksum verification failed!"
  echo "The downloaded file may be corrupted or tampered with."
  cd ..
  rm -rf tmpshopcli
  exit 1
}
echo "✓ Checksum verified"

echo "Extracting..."
tar -xzf "$ARCHIVE_NAME"

INSTALL_DIR="${HOME}/.local/bin"
mkdir -p "$INSTALL_DIR"

mv shop "$INSTALL_DIR"
cd ..
rm -rf tmpshopcli

echo "✓ Downloaded shop and installed at: ${INSTALL_DIR}/shop"

# Check if install dir is in PATH
if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
  echo ""
  echo "⚠️  WARNING: $INSTALL_DIR is not in your PATH"
  echo ""
  echo "Add this line to your shell config file:"
  
  # Detect shell and suggest the right config file
  if [[ "$SHELL" == *"zsh"* ]]; then
    echo "  echo 'export PATH=\"\$HOME/.local/bin:\$PATH\"' >> ~/.zshrc"
    echo "  source ~/.zshrc"
  elif [[ "$SHELL" == *"bash"* ]]; then
    echo "  echo 'export PATH=\"\$HOME/.local/bin:\$PATH\"' >> ~/.bashrc"
    echo "  source ~/.bashrc"
  else
    echo "  export PATH=\"\$HOME/.local/bin:\$PATH\""
    echo ""
    echo "  Add this to your shell's config file (~/.bashrc, ~/.zshrc, etc.)"
  fi
  echo ""
  echo "Or run this now (temporary, only for current session):"
  echo "  export PATH=\"\$HOME/.local/bin:\$PATH\""
else
  echo "✓ Installation complete! Run 'shop --help' to get started"
fi

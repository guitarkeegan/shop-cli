#!/bin/bash
set -e

if [[ -z $GITHUB_TOKEN ]]; then
  echo "GITHUB_TOKEN must be set!"
  exit 1
fi

if ! command -v jq &> /dev/null; then
  echo "Error: jq is required but not installed"
  echo ""
  echo "Install jq:"
  exit 1
fi

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

echo "found GITHUB_TOKEN..."

OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case $ARCH in
    x86_64) ARCH="amd64" ;;
    aarch64|arm64) ARCH="arm64" ;;
    armv7l) ARCH="armv7" ;;
    i386|i686) ARCH="386" ;;
esac

echo "OS: ${OS}, ARCH: ${ARCH}"

RELEASES=$(curl -sL \
  -H "Accept: application/vnd.github+json" \
  -H "Authorization: Bearer ${GITHUB_TOKEN}" \
  -H "X-GitHub-Api-Version: 2022-11-28" \
  https://api.github.com/repos/guitarkeegan/shop-cli/releases/latest)

echo "$RELEASES" | jq empty || {
  echo "invalid json response"
  exit 1
}

ARCHIVE_NAME="shop-cli_${OS}_${ARCH}.tar.gz"
echo "searching for ${ARCHIVE_NAME}..."

ASSET_ID=$(echo "$RELEASES" | jq -r \
  ".assets[] | select(.name == \"$ARCHIVE_NAME\") | .id")

if [[ -z $ASSET_ID ]]; then
  echo "Error: ${ARCHIVE_NAME} not found in repo"
  echo "Found archives:"
  echo "$RELEASES" | jq -r ".assets[].name"
  exit 1
fi

echo "Found asset ID: ${ASSET_ID}"

CHECKSUM_ASSET_ID=$(echo "$RELEASES" | jq -r \
  '.assets[] | select(.name | test("checksums.txt$")) | .id')

if [[ -z $CHECKSUM_ASSET_ID ]]; then
  echo "checksum not found!"
  exit 1
fi

echo "retrieved checksum asset ID"

mkdir -p tmpshopcli
cd tmpshopcli

curl -sL \
  -H "Accept: application/octet-stream" \
  -H "Authorization: Bearer ${GITHUB_TOKEN}" \
  -H "X-GitHub-Api-Version: 2022-11-28" \
  -o "$ARCHIVE_NAME" \
  "https://api.github.com/repos/guitarkeegan/shop-cli/releases/assets/${ASSET_ID}"

echo "downloaded $ARCHIVE_NAME"

# Download checksums
curl -sL \
  -H "Accept: application/octet-stream" \
  -H "Authorization: Bearer ${GITHUB_TOKEN}" \
  -H "X-GitHub-Api-Version: 2022-11-28" \
  -o "checksums.txt" \
  "https://api.github.com/repos/guitarkeegan/shop-cli/releases/assets/${CHECKSUM_ASSET_ID}"

echo "Verifying checksum..."
grep "$ARCHIVE_NAME" checksums.txt | sha256sum -c - || {
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

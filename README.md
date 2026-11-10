# Shop CLI

The purpose of this CLI tool is to make the creation of themes, just a bit easier.

## Install

```bash
curl -fsSL https://raw.githubusercontent.com/guitarkeegan/shop-cli/main/install.sh | bash
```

## Development

### Prerequisites

- Go 1.24.2 or later

### Building from source

```bash
go build -o shop cmd/main.go
```

### Usage

```bash
# Create a new section with inputs
shop section create -n section_name text,image,range

# Create a new block with inputs
shop block create -n block_name text,image,range
```

### Development commands

```bash
# Clear generated files
make clear           # Clear both sections and blocks directories
make clear-sections  # Clear only sections directory
make clear-blocks    # Clear only blocks directory
```

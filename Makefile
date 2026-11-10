# Makefile for clearing sections and blocks directories
.PHONY: clear clear-sections clear-blocks help clear-logs clear-snippets clear-src

# Default target - shows help
help:
	@echo "Available targets:"
	@echo "  clear          - Clear sections, blocks, and snippets directories"
	@echo "  clear-sections - Clear only the sections directory"
	@echo "  clear-snippets - Clear only the sections directory"
	@echo "  clear-blocks   - Clear only the blocks directory"
	@echo "  clear-logs     - Clear only the logs file"
	@echo "  clear-src      - Clear only the src directory"
	@echo "  help           - Show this help message"

# Clear both sections and blocks
clear: clear-sections clear-blocks clear-snippets clear-src
	@echo "✅ All directories cleared"

# Clear src/sections directory
clear-src:
	@echo "🗑️  Clearing src/sections directory..."
	@if [ -d "src/sections" ]; then \
		rm -rf src/sections/*; \
		echo "✅ src/sections directory cleared"; \
	else \
		echo "⚠️  src/sections directory does not exist"; \
	fi

# Clear sections directory
clear-sections:
	@echo "🗑️  Clearing sections directory..."
	@if [ -d "sections" ]; then \
		rm -rf sections/*; \
		echo "✅ Sections directory cleared"; \
	else \
		echo "⚠️  Sections directory does not exist"; \
	fi
	 
# Clear snippets directory
clear-snippets:
	@echo "🗑️  Clearing snippets directory..."
	@if [ -d "snippets" ]; then \
		rm -rf snippets/*; \
		echo "✅ Snippets directory cleared"; \
	else \
		echo "⚠️  Snippets directory does not exist"; \
	fi

# Clear blocks directory  
clear-blocks:
	@echo "🗑️  Clearing blocks directory..."
	@if [ -d "blocks" ]; then \
		rm -rf blocks/*; \
		echo "✅ Blocks directory cleared"; \
	else \
		echo "⚠️  Blocks directory does not exist"; \
	fi

# Clear logs file
clear-logs:
	@echo "🗑️  Clearing logs file..."
	@if [ -f "logs" ]; then \
		rm -f logs; \
		echo "✅ Logs file cleared"; \
	else \
		echo "⚠️  Logs file does not exist"; \
	fi

# Shop CLI

Shop is a CLI tool for Shopify theme developers that generates Liquid template files — sections, blocks, and snippets — with proper schema JSON, Liquid variable assignments, and documentation blocks.

All commands must be run from the root of a Shopify theme project (where `sections/`, `blocks/`, and `snippets/` directories exist). Every command accepts `-h`/`--help`; `shop --version` prints the version.

On success a command prints `Created <path>` and exits 0. On failure it prints the error to stderr and exits non-zero.

## Commands

### `shop section create`

Creates a section `.liquid` file in `sections/`.

```bash
shop section create -n "Section Name" "type:id:label,type:id:label"
shop section create -n "Hero" -b "Card,text:title:Title,image_picker:img:Image" "text:heading:Heading"
```

| Flag/Operand | Description |
|---|---|
| `-n`, `--name` | **Required.** Section name (used as schema name, handlized for filename). |
| `-t`, `--template` | Interactive: pick a prebuilt feature to copy into the theme (requires `sections/` and `src/`). |
| `-b`, `--blocks` | Add a block with settings. Repeatable. Format: `"Block Name,type:id:label,type:id:label"`. |
| `settings` (operand) | Optional comma-separated settings: `type:id:label`. |

Default settings are automatically appended: header, color_scheme, padding_top (range 0-120), padding_bottom (range 0-120).

### `shop block create`

Creates a theme block `.liquid` file in `blocks/`.

```bash
shop block create -n "Block Name" "type:id:label,type:id:label"
```

| Flag/Operand | Description |
|---|---|
| `-n`, `--name` | **Required.** Block name and filename. |
| `settings` (operand) | Optional comma-separated settings: `type:id:label`. |

### `shop snippet create`

Creates a snippet `.liquid` file in `snippets/` with a `{% doc %}` block.

```bash
shop snippet create -n "My Snippet" -d "Renders a card" -e '{% render "my-snippet", title: "Hi" %}' "title:string,image?:object"
```

| Flag/Operand | Description |
|---|---|
| `-n`, `--name` | **Required.** Snippet filename. |
| `-d`, `--description` | Description included in the `{% doc %}` block. |
| `-e`, `--example` | Example implementation included in `{% doc %}`. |
| `-t`, `--template` | Launch interactive file picker to prepopulate content from a `.liquid` file. |
| `params` (operand) | Comma-separated doc params. Formats: `name`, `name:type`, `name?:type` (optional). Types: `string`, `number`, `boolean`, `object`. |

### `shop examples pull`

Optional. Downloads example files from a GitHub repository into the user cache directory, skipping the download if nothing changed.

```bash
shop examples pull
shop examples pull -f   # force re-download, bypass SHA hash check
```

**Required env vars:** `SHOP_GITHUB_TOKEN`, `SHOP_GITHUB_OWNER`, `SHOP_GITHUB_REPO`, `SHOP_EXAMPLES_PATH`.

## Setting String Format

Settings are specified as `type:id:label` separated by commas.

```
text:title:Title,image_picker:image:Background Image,range:opacity:Opacity
```

### Valid Shopify Setting Types

`article`, `article_list`, `blog`, `checkbox`, `collection`, `collection_list`, `color`, `color_background`, `color_scheme`, `color_scheme_group`, `font_picker`, `header`, `html`, `image_picker`, `inline_richtext`, `link_list`, `liquid`, `metaobject`, `metaobject_list`, `number`, `page`, `paragraph`, `product`, `product_list`, `radio`, `range`, `richtext`, `select`, `text`, `text_alignment`, `textarea`, `url`, `video`, `video_url`

Special: `range` auto-adds min/max/step defaults. `header` and `paragraph` use the label as their `content` and have no `id`. Labels may contain colons; everything after the second colon is the label.

## Block String Format (for `-b` flag)

First value is the block name, remaining values are its settings. The block type is the name's handle:

```
"Block Name,type:id:label,type:id:label"
```

Example: `"Product Card,text:title:Title,image_picker:image:Image"` creates a block of type `product-card`.

## Naming

Names become Shopify-style handles for filenames: accents are removed, letters are lowercased, and any run of characters other than letters, digits and underscores becomes a hyphen ("Hero Banner" becomes `hero-banner.liquid`, "Café / Menu" becomes `cafe-menu.liquid`). Files are never overwritten: if one exists, a numeric suffix is appended (`-1`, `-2`, etc.).

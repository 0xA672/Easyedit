# Configuration

EasyEdit loads one TOML configuration file when the interactive editor starts. No configuration file is required: when the file is absent, EasyEdit uses its built-in defaults.

## Configuration file location

| Platform | Path |
| --- | --- |
| Linux / macOS | `~/.easyedit.toml` |
| Windows | `%APPDATA%\easyedit\config.toml` |

On Windows, if `APPDATA` is not set, EasyEdit falls back to `%USERPROFILE%\AppData\Roaming\easyedit\config.toml`.

The path is selected automatically. The current command-line interface does not provide a flag to select an alternate configuration file.

## Create a configuration file

Create the directory and file for your platform, then add the TOML settings you want to change. For example, on Linux or macOS:

```sh
cat > "$HOME/.easyedit.toml" <<'EOF'
tab_width = 4
show_line_num = true
backup = true
soft_wrap = true
auto_indent = true
undo_limit = 100
vim_mode = false
scroll_margin = 3
show_whitespace = false

[theme]
cursor_line_bg = "#2a2a2a"
linenum_cur_fg = "yellow"
whitespace_fg = "#555555"
EOF
```

On Windows PowerShell:

```powershell
$configDir = Join-Path $env:APPDATA "easyedit"
New-Item -ItemType Directory -Force $configDir | Out-Null
notepad (Join-Path $configDir "config.toml")
```

Paste the TOML example above into Notepad and save it. Restart EasyEdit after changing the file; settings are loaded when the editor starts.

## Editor options

These options belong at the top level of the TOML file, before any `[theme]` or `[keys]` table.

| Option | Type | Default | Description / current status |
| --- | --- | ---: | --- |
| `tab_width` | integer | `4` | Tab-stop width used for visual tab positioning and indentation calculations. Values below `1` are reset to `4`. |
| `show_line_num` | boolean | `true` | Show the line-number gutter. |
| `backup` | boolean | `true` | Create a sibling `.bak` file containing the previous file contents when saving an existing file. |
| `soft_wrap` | boolean | `true` | Enable soft wrapping. `Alt+W` toggles it for the current session. |
| `auto_indent` | boolean | `true` | Copy indentation from the previous line when inserting a newline. |
| `undo_limit` | integer | `100` | Maximum undo history size when opening a file. Values below `1` are reset to `100`. |
| `vim_mode` | boolean | `false` | Start in Vim-style Normal mode instead of the default Insert mode. |
| `scroll_margin` | integer | `3` | Keep this many visible lines between the cursor and the top/bottom edge while scrolling. Negative values are reset to `0`. |
| `show_whitespace` | boolean | `false` | Display whitespace markers in the editor view. |
| `show_mode` | boolean | `true` | Present in the configuration schema and example, but not currently wired to the UI. Changing it may have no visible effect. |
| `tab_to_spaces` | boolean | `true` | Present in the configuration schema. The current UI uses it in the status-bar label, but it is not yet a complete tab-input conversion setting. |

Booleans must be written as TOML `true` or `false`, without quotes.

## Theme

Theme settings go inside a `[theme]` table. Colors accept names understood by `tcell` (for example, `yellow`, `green`, or `blue`) and hexadecimal colors such as `"#2a2a2a"`.

| Option | Default | Description / current status |
| --- | --- | --- |
| `default_fg` | `"default"` | Default foreground color; currently not consistently applied by the UI. |
| `default_bg` | `"default"` | Default background color; currently not consistently applied by the UI. |
| `linenum_fg` | `"blue"` | Line-number foreground; currently not consistently applied by the UI. |
| `linenum_bg` | `"default"` | Line-number background; currently not consistently applied by the UI. |
| `linenum_cur_fg` | `"yellow"` | Current-line number foreground. |
| `statusbar_fg` | `"white"` | Status-bar foreground; currently not consistently applied by the UI. |
| `statusbar_bg` | `"blue"` | Status-bar background; currently not consistently applied by the UI. |
| `search_match_bg` | `"yellow"` | Search-match background; currently not consistently applied by the UI. |
| `bracket_match` | `"green"` | Bracket-match color; currently not consistently applied by the UI. |
| `selection_bg` | `"navy"` | Selection background; currently not consistently applied by the UI. |
| `cursor_line_bg` | `"#2a2a2a"` | Current-line background. |
| `whitespace_fg` | `"#555555"` | Whitespace-marker foreground. |

The configuration type contains the full palette above, but the current UI explicitly applies only `linenum_cur_fg`, `cursor_line_bg`, and `whitespace_fg`. The other palette entries are documented for completeness and should be treated as not yet fully implemented.

## Key bindings

The `[keys]` table exists in the configuration schema, but custom key bindings are **not currently implemented**. The values are reserved and changing them does not remap the editor's shortcuts.

The current default shortcuts include:

| Action | Shortcut |
| --- | --- |
| Save | `Ctrl+S` |
| Quit | `Ctrl+Q` |
| Search | `Ctrl+F` |
| Replace | `Ctrl+H` |
| Select all | `Ctrl+A` |
| Undo / Redo | `Ctrl+Z` / `Ctrl+Y` |
| Cut / Copy / Paste | `Ctrl+X` / `Ctrl+C` / `Ctrl+V` |
| Toggle soft wrap | `Alt+W` |
| Enter command mode | `:` |

Vim-style mode switching is enabled separately with `vim_mode = true`.

## Loading and troubleshooting

- A missing configuration file results in the built-in defaults being used.
- If the file cannot be read or TOML decoding fails, EasyEdit currently falls back without displaying a detailed configuration error. Check the TOML syntax if a setting appears to be ignored.
- Keep top-level options before `[theme]` or `[keys]`. In TOML, once a table such as `[keys]` begins, subsequent key/value pairs remain in that table until another table header appears.
- Numeric settings are normalized as described above: `tab_width < 1` becomes `4`, `undo_limit < 1` becomes `100`, and `scroll_margin < 0` becomes `0`.
- Changes are read at startup. Restart the editor after editing the configuration file.

> **Note about `easyedit.example.toml`:** in the current repository version, root-level editor options appear after the `[keys]` table. TOML therefore interprets those options as members of `[keys]`, rather than as top-level settings. The example in this guide uses the correct table order. The checked-in example file should be corrected to match it.

For installation instructions, see [Installation Guide](installation.md).

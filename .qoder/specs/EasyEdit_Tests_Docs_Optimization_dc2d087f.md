# EasyEdit: Tests, Documentation & Optimization

## Context

EasyEdit is a Go-based terminal text editor with solid architecture but uneven test coverage and documentation. The `command`, `config`, and `document` packages are well-tested and documented, while `ui` (8 tests covering only `autoPair`), `main` (zero tests), and `document/undo.go` (no dedicated tests) have significant gaps. Some exported functions in `main.go` have Chinese-only comments, and `ui/editor.go` lacks detailed doc comments on many methods/fields.

## Task  tests for formatting
- Any other exported or testable unexported helper functions

Keep tests unit-focused (no tcell screen required) by testing internal state methods.

## Task 2: Add tests for `main.go` stream mode

**Files**: `d:\Easyedit\main_test.go` (new file)

Add tests for:
- `parseScript()` — splitting by `;` and `\n`
- `applyRules()` — substitution with `/old/new/`, `/old/new/g`, edge cases (empty, no match, special chars)
- Stream mode end-to-end via `os/exec` or function extraction

## Task 3: Add dedicated tests for `document/undo.go`

**Files**: `d:\Easyedit\document\undo_test.go` (new file)

Add tests for:
- `NewUndoStack()` — verify limit
- `Push()` — basic insert/delete/replace recording
- Merge logic — consecutive adjacent inserts merge into one step
- `Undo()`/`Redo()` — round-trip operations
- Stack limit enforcement — oldest steps dropped when limit exceeded
- `Clear()` — resets both stacks
- Redo stack cleared on new push after undo

## Task 4: Enhance `clipboard` tests

**Files**: `d:\Easyedit\clipboard\clipboard_test.go`

Add tests for:
- `WriteAll` with empty string
- `WriteAll` with unicode/multibyte content
- `ReadAll` after multiple sequential writes (last-write-wins)

## Task 5: Add English documentation to `main.go`

**Files**: `d:\Easyedit\main.go`

- Replace Chinese-only function comments with English godoc-style comments for: `main()`, `runStreamMode()`, `parseScript()`, `applyRules()`, `printUsage()`
- Add doc comment to `Version` variable

## Task 6: Add documentation to `ui/editor.go`

**Files**: `d:\Easyedit\ui\editor.go`

- Add doc comments to `Mode` constants (Normal, Insert, Visual, Command, Search, Replace)
- Add doc comments to `Editor` struct and key exported fields
- Add doc comments to all exported methods: `NewEditor()`, `Run()`, `Quit()`, `OpenFile()`, `SaveFile()`
- Add brief field comments for the 30+ Editor struct fields grouped by concern

## Task 7: Add documentation to `clipboard` package

**Files**: `d:\Easyedit\clipboard\clipboard.go`, `clipboard_stub.go`, `clipboard_windows.go`, `clipboard_unix.go`

- Expand package-level doc comment explaining platform strategy
- Add doc comments to `ReadAll`, `WriteAll`, `Unsupported`

## Task 8: Minor optimizations

**Files**: Multiple

Apply safe, non-breaking improvements:
1. **`document/undo.go`**: Remove unused `Position` struct (dead code)
2. **`document/document.go`**: In `Indentation()`, avoid unnecessary string↔rune conversion
3. **`highlight/highlight.go`**: Add nil-check guard in `GetTokenStyle` for robustness
4. **`main.go`**: Extract magic number `16` (gap buffer capacity) into a named constant
5. **`config/config.go`**: Add a log warning (to stderr) when config file has parse errors instead of silent fallback

## Verification

1. Run `go build ./...` — must compile cleanly
2. Run `go test ./... -v` — all existing + new tests must pass
3. Run `go vet ./...` — no new warnings
4. Verify no behavioral changes to existing functionality

## Dependencies

- Tasks 1-4 (tests) are independent of each other and can run in parallel
- Tasks 5-7 (docs) are independent of each other and can run in parallel
- Task 8 (optimizations) should run after tests are in place to verify no regressions
- Verification runs last after all tasks complete

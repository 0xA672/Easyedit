// Package clipboard provides cross-platform clipboard access via OS commands.
// Falls back to no-op on unsupported platforms (the editor has its own internal
// clipboard as fallback).
package clipboard

// ReadAll reads text from the system clipboard.
func ReadAll() (string, error) {
	return readAll()
}

// WriteAll writes text to the system clipboard.
func WriteAll(text string) error {
	return writeAll(text)
}

// Unsupported reports whether clipboard operations are unavailable on this platform.
var Unsupported bool

package clipboard

import (
	"testing"
)

func TestReadAllReturnsEmptyOnUnsupported(t *testing.T) {
	if !Unsupported {
		t.Skip("clipboard is supported on this platform, skipping no-op test")
	}
	got, err := ReadAll()
	if err != nil {
		t.Fatalf("ReadAll() returned error: %v", err)
	}
	if got != "" {
		t.Fatalf("ReadAll() = %q, want empty string", got)
	}
}

func TestWriteAllNoErrorOnUnsupported(t *testing.T) {
	if !Unsupported {
		t.Skip("clipboard is supported on this platform, skipping no-op test")
	}
	if err := WriteAll("test"); err != nil {
		t.Fatalf("WriteAll() returned error: %v", err)
	}
}

func TestReadAllWriteAllRoundtrip(t *testing.T) {
	if Unsupported {
		t.Skip("clipboard unsupported on this platform")
	}
	// Save original clipboard content
	orig, _ := ReadAll()

	testStr := "EasyEdit clipboard test: 你好, 世界!"
	if err := WriteAll(testStr); err != nil {
		t.Fatalf("WriteAll() returned error: %v", err)
	}

	got, err := ReadAll()
	if err != nil {
		t.Fatalf("ReadAll() returned error: %v", err)
	}
	if got != testStr {
		t.Fatalf("ReadAll() = %q, want %q", got, testStr)
	}

	// Restore original
	if orig != "" {
		WriteAll(orig)
	}
}

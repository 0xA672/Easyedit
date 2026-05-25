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

	// Save the original clipboard content before probing, so we can restore it later.
	origBeforeProbe, _ := ReadAll()

	// Probe whether the clipboard is actually usable (e.g., daemon is running, tool is installed).
	const probe = "EasyEdit clipboard probe"
	if err := WriteAll(probe); err != nil {
		// Try to restore the original content before skipping.
		if origBeforeProbe != "" {
			_ = WriteAll(origBeforeProbe)
		}
		t.Skipf("WriteAll probe failed, clipboard may be unavailable: %v", err)
	}

	gotProbe, err := ReadAll()
	if err != nil {
		if origBeforeProbe != "" {
			_ = WriteAll(origBeforeProbe)
		}
		t.Skipf("ReadAll probe failed, clipboard may be unavailable: %v", err)
	}
	if gotProbe != probe {
		if origBeforeProbe != "" {
			_ = WriteAll(origBeforeProbe)
		}
		t.Skip("clipboard appears unavailable (roundtrip probe mismatched)")
	}

	// Restore the original clipboard content before running the actual roundtrip test.
	if origBeforeProbe != "" {
		_ = WriteAll(origBeforeProbe)
	}

	// Actual roundtrip test.
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

	// Restore the original clipboard content after the test.
	if orig != "" {
		_ = WriteAll(orig)
	}
}

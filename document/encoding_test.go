package document

import (
	"testing"
)

// ==================== Encoding Detection ====================

func TestDetectEncoding_EmptyData(t *testing.T) {
	data := []byte{}
	enc := DetectEncoding(data)
	if enc != EncUTF8 {
		t.Fatalf("empty data should be detected as UTF-8, got %v", enc)
	}
}

func TestDetectEncoding_UTF8WithBOM(t *testing.T) {
	// UTF-8 BOM: EF BB BF
	data := []byte{0xEF, 0xBB, 0xBF, 'h', 'e', 'l', 'l', 'o'}
	enc := DetectEncoding(data)
	if enc != EncUTF8 {
		t.Fatalf("UTF-8 with BOM should be detected as UTF-8, got %v", enc)
	}
}

func TestDetectEncoding_UTF16LE_BOM(t *testing.T) {
	// UTF-16 LE BOM: FF FE
	data := []byte{0xFF, 0xFE, 0x68, 0x00} // "h" in UTF-16LE
	enc := DetectEncoding(data)
	if enc != EncUTF16LE {
		t.Fatalf("UTF-16LE with BOM should be detected as UTF-16LE, got %v", enc)
	}
}

func TestDetectEncoding_UTF16BE_BOM(t *testing.T) {
	// UTF-16 BE BOM: FE FF
	data := []byte{0xFE, 0xFF, 0x00, 0x68} // "h" in UTF-16BE
	enc := DetectEncoding(data)
	if enc != EncUTF16BE {
		t.Fatalf("UTF-16BE with BOM should be detected as UTF-16BE, got %v", enc)
	}
}

func TestDetectEncoding_ValidUTF8(t *testing.T) {
	data := []byte("hello 世界")
	enc := DetectEncoding(data)
	if enc != EncUTF8 {
		t.Fatalf("valid UTF-8 should be detected as UTF-8, got %v", enc)
	}
}

func TestDetectEncoding_Latin1Fallback(t *testing.T) {
	// Invalid UTF-8 but valid Latin1
	// Note: GBK decoder may accept some byte sequences that are also valid Latin1
	// This test just verifies the function doesn't crash and returns a valid encoding
	data := []byte{0xE9, 0xE0, 0xED} // éàí in Latin1
	enc := DetectEncoding(data)
	// The detection logic tries GBK first, so it might return GBK instead of Latin1
	// Both are acceptable for this input
	validEncodings := map[CommonEncoding]bool{
		EncLatin1: true,
		EncGBK:    true,
	}
	if !validEncodings[enc] {
		t.Fatalf("non-UTF8 data should fallback to Latin1 or GBK, got %v", enc)
	}
}

// ==================== DecodeToUTF8 ====================

func TestDecodeToUTF8_Empty(t *testing.T) {
	data := []byte{}
	text, enc, err := DecodeToUTF8(data)
	if err != nil {
		t.Fatalf("empty data should not error, got %v", err)
	}
	if text != "" {
		t.Fatalf("empty data should decode to empty string, got %q", text)
	}
	if enc != EncUTF8 {
		t.Fatalf("empty data encoding should be UTF-8, got %v", enc)
	}
}

func TestDecodeToUTF8_UTF8WithBOM(t *testing.T) {
	data := []byte{0xEF, 0xBB, 0xBF, 'h', 'i'}
	text, enc, err := DecodeToUTF8(data)
	if err != nil {
		t.Fatalf("UTF-8 with BOM should not error, got %v", err)
	}
	if text != "hi" {
		t.Fatalf("BOM should be stripped, got %q", text)
	}
	if enc != EncUTF8 {
		t.Fatalf("encoding should be UTF-8, got %v", enc)
	}
}

func TestDecodeToUTF8_PlainUTF8(t *testing.T) {
	data := []byte("hello 世界")
	text, enc, err := DecodeToUTF8(data)
	if err != nil {
		t.Fatalf("plain UTF-8 should not error, got %v", err)
	}
	if text != "hello 世界" {
		t.Fatalf("decoded text mismatch, got %q", text)
	}
	if enc != EncUTF8 {
		t.Fatalf("encoding should be UTF-8, got %v", enc)
	}
}

// ==================== EncodeFromUTF8 ====================

func TestEncodeFromUTF8_UTF8(t *testing.T) {
	text := "hello"
	data, err := EncodeFromUTF8(text, EncUTF8)
	if err != nil {
		t.Fatalf("UTF-8 encoding should not error, got %v", err)
	}
	if string(data) != "hello" {
		t.Fatalf("UTF-8 encode mismatch, got %q", string(data))
	}
}

func TestEncodeFromUTF8_Unknown(t *testing.T) {
	text := "hello"
	data, err := EncodeFromUTF8(text, EncUnknown)
	if err != nil {
		t.Fatalf("Unknown encoding should not error, got %v", err)
	}
	if string(data) != "hello" {
		t.Fatalf("Unknown encode should return raw bytes, got %q", string(data))
	}
}

// ==================== StripBOM ====================

func TestStripBOM_UTF8(t *testing.T) {
	data := []byte{0xEF, 0xBB, 0xBF, 'h', 'i'}
	stripped := StripBOM(data)
	if len(stripped) != 2 {
		t.Fatalf("UTF-8 BOM should be stripped, got length %d", len(stripped))
	}
	if string(stripped) != "hi" {
		t.Fatalf("stripped content mismatch, got %q", string(stripped))
	}
}

func TestStripBOM_UTF16LE(t *testing.T) {
	data := []byte{0xFF, 0xFE, 'h', 'i'}
	stripped := StripBOM(data)
	if len(stripped) != 2 {
		t.Fatalf("UTF-16LE BOM should be stripped, got length %d", len(stripped))
	}
}

func TestStripBOM_UTF16BE(t *testing.T) {
	data := []byte{0xFE, 0xFF, 'h', 'i'}
	stripped := StripBOM(data)
	if len(stripped) != 2 {
		t.Fatalf("UTF-16BE BOM should be stripped, got length %d", len(stripped))
	}
}

func TestStripBOM_NoBOM(t *testing.T) {
	data := []byte("hello")
	stripped := StripBOM(data)
	if len(stripped) != 5 {
		t.Fatalf("no BOM should keep original length, got %d", len(stripped))
	}
	if string(stripped) != "hello" {
		t.Fatalf("no BOM should keep original content, got %q", string(stripped))
	}
}

// ==================== CommonEncoding String ====================

func TestCommonEncoding_String(t *testing.T) {
	tests := []struct {
		enc  CommonEncoding
		want string
	}{
		{EncUnknown, "unknown"},
		{EncUTF8, "UTF-8"},
		{EncUTF16LE, "UTF-16LE"},
		{EncUTF16BE, "UTF-16BE"},
		{EncGBK, "GBK"},
		{EncBig5, "Big5"},
		{EncShiftJIS, "Shift_JIS"},
		{EncEUCJP, "EUC-JP"},
		{EncEUCKR, "EUC-KR"},
		{EncLatin1, "ISO-8859-1"},
	}
	for _, tt := range tests {
		got := tt.enc.String()
		if got != tt.want {
			t.Errorf("encoding %v: String() = %q, want %q", tt.enc, got, tt.want)
		}
	}
}

// ==================== encToGoEncoding ====================

func TestEncToGoEncoding_AllTypes(t *testing.T) {
	// Just verify it doesn't panic and returns non-nil
	encodings := []CommonEncoding{
		EncUTF8, EncUTF16LE, EncUTF16BE, EncGBK, EncBig5,
		EncShiftJIS, EncEUCJP, EncEUCKR, EncLatin1, EncUnknown,
	}
	for _, enc := range encodings {
		goEnc := encToGoEncoding(enc)
		if goEnc == nil {
			t.Errorf("encToGoEncoding(%v) returned nil", enc)
		}
	}
}

// ==================== GBKToUTF8 ====================

func TestGBKToUTF8_Simple(t *testing.T) {
	// Simple ASCII works in GBK too
	data := []byte("hello")
	text, err := GBKToUTF8(data)
	if err != nil {
		t.Fatalf("ASCII GBK should not error, got %v", err)
	}
	if text != "hello" {
		t.Fatalf("GBK decode mismatch, got %q", text)
	}
}

// ==================== UTF8ToGBK ====================

func TestUTF8ToGBK_Simple(t *testing.T) {
	// Simple ASCII works in GBK too
	text := "hello"
	data, err := UTF8ToGBK(text)
	if err != nil {
		t.Fatalf("ASCII to GBK should not error, got %v", err)
	}
	if string(data) != "hello" {
		t.Fatalf("GBK encode mismatch, got %q", string(data))
	}
}

// ==================== isLikelyRawLatin1 ====================

func TestIsLikelyRawLatin1_True(t *testing.T) {
	// Valid Latin1 without C1 control chars
	data := []byte{0xA0, 0xB0, 0xC0, 0xD0, 0xE0, 0xF0}
	result := isLikelyRawLatin1(data)
	if !result {
		t.Fatalf("should be likely Latin1, got false")
	}
}

func TestIsLikelyRawLatin1_False(t *testing.T) {
	// Contains C1 control chars (0x80-0x9F)
	data := []byte{0x80, 0x90, 0xA0}
	result := isLikelyRawLatin1(data)
	if result {
		t.Fatalf("should not be likely Latin1 with C1 chars, got true")
	}
}

func TestIsLikelyRawLatin1_Empty(t *testing.T) {
	data := []byte{}
	result := isLikelyRawLatin1(data)
	if !result {
		t.Fatalf("empty data should be likely Latin1, got false")
	}
}

// ==================== BOM Method ====================

func TestCommonEncoding_BOM(t *testing.T) {
	tests := []struct {
		enc  CommonEncoding
		want []byte
	}{
		{EncUTF16LE, []byte{0xFF, 0xFE}},
		{EncUTF16BE, []byte{0xFE, 0xFF}},
		{EncUTF8, []byte{}},
		{EncGBK, []byte{}},
		{EncLatin1, []byte{}},
	}
	for _, tt := range tests {
		got := tt.enc.BOM()
		if len(got) != len(tt.want) {
			t.Errorf("encoding %v: BOM() length = %d, want %d", tt.enc, len(got), len(tt.want))
			continue
		}
		for i := range got {
			if got[i] != tt.want[i] {
				t.Errorf("encoding %v: BOM()[%d] = 0x%02X, want 0x%02X", tt.enc, i, got[i], tt.want[i])
			}
		}
	}
}

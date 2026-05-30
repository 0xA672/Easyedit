// Package document provides the data structures for editor documents,
// including encoding detection and conversion support.
package document

import (
	"unicode/utf8"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/encoding/unicode"
)

// CommonEncoding represents a human-readable encoding name.
type CommonEncoding int

const (
	EncUnknown CommonEncoding = iota
	EncUTF8
	EncUTF16LE
	EncUTF16BE
	EncGBK
	EncBig5
	EncShiftJIS
	EncEUCJP
	EncEUCKR
	EncLatin1
)

var encodingName = map[CommonEncoding]string{
	EncUnknown:  "unknown",
	EncUTF8:     "UTF-8",
	EncUTF16LE:  "UTF-16LE",
	EncUTF16BE:  "UTF-16BE",
	EncGBK:      "GBK",
	EncBig5:     "Big5",
	EncShiftJIS: "Shift_JIS",
	EncEUCJP:    "EUC-JP",
	EncEUCKR:    "EUC-KR",
	EncLatin1:   "ISO-8859-1",
}

func (e CommonEncoding) String() string {
	if s, ok := encodingName[e]; ok {
		return s
	}
	return "unknown"
}

// encToEncoding maps our common encoding to golang.org/x/text encoding.
func encToGoEncoding(e CommonEncoding) encoding.Encoding {
	switch e {
	case EncUTF8:
		return encoding.Nop // identity
	case EncUTF16LE:
		return unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
	case EncUTF16BE:
		return unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM)
	case EncGBK:
		return simplifiedchinese.GBK
	case EncBig5:
		return traditionalchinese.Big5
	case EncShiftJIS:
		return japanese.ShiftJIS
	case EncEUCJP:
		return japanese.EUCJP
	case EncEUCKR:
		return korean.EUCKR
	case EncLatin1:
		return charmap.ISO8859_1
	default:
		return charmap.ISO8859_1 // safe fallback
	}
}

// DetectEncoding detects the character encoding of raw byte data.
func DetectEncoding(data []byte) CommonEncoding {
	if len(data) == 0 {
		return EncUTF8
	}

	// 1. BOM detection
	if len(data) >= 3 && data[0] == 0xEF && data[1] == 0xBB && data[2] == 0xBF {
		return EncUTF8
	}
	if len(data) >= 2 && data[0] == 0xFF && data[1] == 0xFE {
		return EncUTF16LE
	}
	if len(data) >= 2 && data[0] == 0xFE && data[1] == 0xFF {
		return EncUTF16BE
	}

	// 2. Valid UTF-8 — most common case
	if utf8.Valid(data) {
		return EncUTF8
	}

	// 3. Try common CJK encodings in order of likelihood
	candidates := []CommonEncoding{EncGBK, EncShiftJIS, EncBig5, EncEUCJP, EncEUCKR, EncLatin1}
	for _, c := range candidates {
		if tryDecode(data, encToGoEncoding(c)) {
			return c
		}
	}

	return EncLatin1 // best effort
}

// tryDecode returns true if data can be decoded without error using the given encoding.
func tryDecode(data []byte, enc encoding.Encoding) bool {
	decoder := enc.NewDecoder()
	_, err := decoder.Bytes(data)
	return err == nil
}

// DecodeToUTF8 decodes raw bytes to a UTF-8 string.
// Returns the decoded string, detected encoding, and any error.
func DecodeToUTF8(data []byte) (string, CommonEncoding, error) {
	if len(data) == 0 {
		return "", EncUTF8, nil
	}

	enc := DetectEncoding(data)

	// Strip BOM bytes for UTF-16
	start := 0
	switch enc {
	case EncUTF16LE:
		if len(data) >= 2 && data[0] == 0xFF && data[1] == 0xFE {
			start = 2
		}
	case EncUTF16BE:
		if len(data) >= 2 && data[0] == 0xFE && data[1] == 0xFF {
			start = 2
		}
	case EncUTF8:
		if len(data) >= 3 && data[0] == 0xEF && data[1] == 0xBB && data[2] == 0xBF {
			start = 3
		}
	}

	if enc == EncUTF8 {
		return string(data[start:]), EncUTF8, nil
	}

	goEnc := encToGoEncoding(enc)
	decoder := goEnc.NewDecoder()
	out, err := decoder.Bytes(data[start:])
	if err != nil {
		return "", enc, err
	}
	return string(out), enc, nil
}

// EncodeFromUTF8 encodes a UTF-8 string back to the given encoding.
func EncodeFromUTF8(text string, enc CommonEncoding) ([]byte, error) {
	if enc == EncUTF8 || enc == EncUnknown {
		return []byte(text), nil
	}

	goEnc := encToGoEncoding(enc)
	encoder := goEnc.NewEncoder()
	out, err := encoder.Bytes([]byte(text))
	if err != nil {
		return nil, err
	}
	return out, nil
}

// StripBOM removes BOM bytes from data if present.
func StripBOM(data []byte) []byte {
	if len(data) >= 3 && data[0] == 0xEF && data[1] == 0xBB && data[2] == 0xBF {
		return data[3:]
	}
	if len(data) >= 2 && ((data[0] == 0xFF && data[1] == 0xFE) || (data[0] == 0xFE && data[1] == 0xFF)) {
		return data[2:]
	}
	return data
}

// GBKToUTF8 converts a GBK-encoded byte slice to a UTF-8 string.
func GBKToUTF8(data []byte) (string, error) {
	decoder := simplifiedchinese.GBK.NewDecoder()
	out, err := decoder.Bytes(data)
	if err != nil {
		return "", err
	}
	return string(out), nil
}

// UTF8ToGBK converts a UTF-8 string to a GBK-encoded byte slice.
func UTF8ToGBK(text string) ([]byte, error) {
	encoder := simplifiedchinese.GBK.NewEncoder()
	out, err := encoder.Bytes([]byte(text))
	if err != nil {
		return nil, err
	}
	return out, nil
}

// isLikelyRawLatin1 returns true if the data is likely just extended-ASCII
// (each byte is a valid single char) but not valid UTF-8 — used as a fast path.
func isLikelyRawLatin1(data []byte) bool {
	for _, b := range data {
		if b >= 0x80 && b < 0xA0 { // C1 control chars, uncommon in text files
			return false
		}
	}
	return true
}

// BOM returns an ASCII/UTF-8 BOM sequence for UTF-16 encodings (empty for others).
func (e CommonEncoding) BOM() []byte {
	switch e {
	case EncUTF16LE:
		return []byte{0xFF, 0xFE}
	case EncUTF16BE:
		return []byte{0xFE, 0xFF}
	default:
		return nil
	}
}

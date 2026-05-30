//go:build windows

package clipboard

import (
	"syscall"
	"unsafe"
)

var (
	user32                  = syscall.NewLazyDLL("user32.dll")
	kernel32                = syscall.NewLazyDLL("kernel32.dll")
	procOpenClipboard       = user32.NewProc("OpenClipboard")
	procCloseClipboard      = user32.NewProc("CloseClipboard")
	procGetClipboardData    = user32.NewProc("GetClipboardData")
	procSetClipboardData    = user32.NewProc("SetClipboardData")
	procEmptyClipboard      = user32.NewProc("EmptyClipboard")
	procGlobalAlloc         = kernel32.NewProc("GlobalAlloc")
	procGlobalLock          = kernel32.NewProc("GlobalLock")
	procGlobalUnlock        = kernel32.NewProc("GlobalUnlock")
	procGlobalSize          = kernel32.NewProc("GlobalSize")
	procMultiByteToWideChar = kernel32.NewProc("MultiByteToWideChar")
	procWideCharToMultiByte = kernel32.NewProc("WideCharToMultiByte")

	CF_UNICODETEXT = uintptr(13)
	GMEM_MOVEABLE  = uintptr(0x0002)
	GMEM_ZEROINIT  = uintptr(0x0040)
	GHND           = GMEM_MOVEABLE | GMEM_ZEROINIT

	CP_UTF8 = uint32(65001)
)

func readAll() (string, error) {
	r, _, _ := procOpenClipboard.Call(0)
	if r == 0 {
		return "", nil
	}
	defer procCloseClipboard.Call()

	h, _, _ := procGetClipboardData.Call(CF_UNICODETEXT)
	if h == 0 {
		return "", nil
	}

	p, _, _ := procGlobalLock.Call(h)
	if p == 0 {
		return "", nil
	}
	defer procGlobalUnlock.Call(h)

	// Get the size of the global handle in bytes
	size, _, _ := procGlobalSize.Call(h)
	if size == 0 {
		return "", nil
	}

	// Maximum number of uint16 elements (includes null terminator)
	maxChars := size / 2
	if maxChars == 0 {
		return "", nil
	}

	// Build a slice over the locked memory and find the null terminator
	u16Slice := unsafe.Slice((*uint16)(unsafe.Pointer(p)), maxChars)
	chars := 0
	for chars < len(u16Slice) {
		if u16Slice[chars] == 0 {
			break
		}
		chars++
	}
	if chars == 0 {
		return "", nil
	}

	// Convert UTF-16 to UTF-8
	size = uintptr(WideCharToMultiByte(CP_UTF8, 0, &u16Slice[0], chars, nil, 0, nil, nil))
	if size == 0 {
		return "", nil
	}

	buf := make([]byte, size)
	WideCharToMultiByte(CP_UTF8, 0, (*uint16)(unsafe.Pointer(p)), chars, &buf[0], int(size), nil, nil)
	return string(buf), nil
}

func writeAll(text string) error {
	r, _, _ := procOpenClipboard.Call(0)
	if r == 0 {
		return nil
	}
	defer procCloseClipboard.Call()
	procEmptyClipboard.Call()

	// Convert UTF-8 to UTF-16
	utf16 := UTF16FromString(text)

	// Allocate global memory
	h, _, _ := procGlobalAlloc.Call(GHND, uintptr(len(utf16)*2+2))
	if h == 0 {
		return nil
	}

	p, _, _ := procGlobalLock.Call(h)
	if p != 0 {
		for i, v := range utf16 {
			*(*uint16)(unsafe.Pointer(p + uintptr(i*2))) = v
		}
		*(*uint16)(unsafe.Pointer(p + uintptr(len(utf16)*2))) = 0 // null terminator
		procGlobalUnlock.Call(h)
	}

	procSetClipboardData.Call(CF_UNICODETEXT, h)
	// After SetClipboardData, the clipboard owns the memory; don't free it
	return nil
}

func UTF16FromString(s string) []uint16 {
	// Trim null in input (paranoid safety check)
	for i := 0; i < len(s); i++ {
		if s[i] == 0 {
			s = s[:i]
			break
		}
	}
	if len(s) == 0 {
		return nil
	}
	// Convert UTF-8 → UTF-16.
	// Go strings are NOT null-terminated, so pass the exact byte count.
	n := MultiByteToWideChar(CP_UTF8, 0, s, len(s), nil, 0)
	if n == 0 {
		return nil
	}
	buf := make([]uint16, n)
	MultiByteToWideChar(CP_UTF8, 0, s, len(s), &buf[0], n)
	return buf
}

func MultiByteToWideChar(codePage uint32, dwFlags uint32, lpMultiByteStr string, cbMultiByte int, lpWideCharStr *uint16, cchWideChar int) int {
	lpMultiByteStrBytes := []byte(lpMultiByteStr)
	if lpWideCharStr != nil && cchWideChar > 0 {
		r, _, _ := procMultiByteToWideChar.Call(
			uintptr(codePage),
			uintptr(dwFlags),
			uintptr(unsafe.Pointer(&lpMultiByteStrBytes[0])),
			uintptr(cbMultiByte),
			uintptr(unsafe.Pointer(lpWideCharStr)),
			uintptr(cchWideChar),
		)
		return int(r)
	}
	r, _, _ := procMultiByteToWideChar.Call(
		uintptr(codePage),
		uintptr(dwFlags),
		uintptr(unsafe.Pointer(&lpMultiByteStrBytes[0])),
		uintptr(cbMultiByte),
		0,
		0,
	)
	return int(r)
}

func WideCharToMultiByte(codePage uint32, dwFlags uint32, lpWideCharStr *uint16, cchWideChar int, lpMultiByteStr *byte, cbMultiByte int, lpDefaultChar *byte, lpUsedDefaultChar *bool) int {
	if lpMultiByteStr != nil && cbMultiByte > 0 {
		r, _, _ := procWideCharToMultiByte.Call(
			uintptr(codePage),
			uintptr(dwFlags),
			uintptr(unsafe.Pointer(lpWideCharStr)),
			uintptr(cchWideChar),
			uintptr(unsafe.Pointer(lpMultiByteStr)),
			uintptr(cbMultiByte),
			0,
			0,
		)
		return int(r)
	}
	r, _, _ := procWideCharToMultiByte.Call(
		uintptr(codePage),
		uintptr(dwFlags),
		uintptr(unsafe.Pointer(lpWideCharStr)),
		uintptr(cchWideChar),
		0,
		0,
		0,
		0,
	)
	return int(r)
}

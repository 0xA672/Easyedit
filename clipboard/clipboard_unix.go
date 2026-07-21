//go:build darwin || dragonfly || freebsd || (linux && !android) || netbsd || openbsd || (solaris && !illumos)

package clipboard

import (
	"os/exec"
	"strings"
)

// getCopyCommand returns a command to copy text to clipboard.
// Uses hardcoded trusted binaries to prevent command injection.
func getCopyCommand() *exec.Cmd {
	switch {
	case hasCmd("pbcopy"):
		return exec.Command("pbcopy")
	case hasCmd("wl-copy"):
		return exec.Command("wl-copy")
	case hasCmd("xclip"):
		return exec.Command("xclip", "-in", "-selection", "clipboard")
	case hasCmd("xsel"):
		return exec.Command("xsel", "--input", "--clipboard")
	case hasCmd("termux-clipboard-set"):
		return exec.Command("termux-clipboard-set")
	}
	Unsupported = true
	return nil
}

// getPasteCommand returns a command to paste text from clipboard.
// Uses hardcoded trusted binaries to prevent command injection.
func getPasteCommand() *exec.Cmd {
	switch {
	case hasCmd("pbpaste"):
		return exec.Command("pbpaste")
	case hasCmd("wl-paste"):
		return exec.Command("wl-paste", "--no-newline")
	case hasCmd("xclip"):
		return exec.Command("xclip", "-out", "-selection", "clipboard")
	case hasCmd("xsel"):
		return exec.Command("xsel", "--output", "--clipboard")
	case hasCmd("termux-clipboard-get"):
		return exec.Command("termux-clipboard-get")
	}
	Unsupported = true
	return nil
}

func readAll() (string, error) {
	if Unsupported {
		return "", nil
	}
	cmd := getPasteCommand()
	if cmd == nil {
		return "", nil
	}
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimRight(string(out), "\r\n"), nil
}

func writeAll(text string) error {
	if Unsupported {
		return nil
	}
	cmd := getCopyCommand()
	if cmd == nil {
		return nil
	}
	in, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	if _, err := in.Write([]byte(text)); err != nil {
		return err
	}
	if err := in.Close(); err != nil {
		return err
	}
	return cmd.Wait()
}

func hasCmd(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

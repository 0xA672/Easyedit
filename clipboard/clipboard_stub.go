//go:build aix || android || illumos || ios || js || wasip1 || plan9

package clipboard

func readAll() (string, error) {
	return "", nil
}

func writeAll(_ string) error {
	return nil
}

func init() {
	Unsupported = true
}

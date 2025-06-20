//go:build !windows
// +build !windows

package metrics

import "fmt"

func GetWindowsVersion() (string, error) {
	return "", fmt.Errorf("GetWindowsVersion is only supported on Windows")
}

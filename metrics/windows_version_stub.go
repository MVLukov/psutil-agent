//go:build !windows
// +build !windows

package metrics

import "fmt"

func GetWindowsVersion() ([]string, error) {
	return []string{}, fmt.Errorf("GetWindowsVersion is only supported on Windows")
}

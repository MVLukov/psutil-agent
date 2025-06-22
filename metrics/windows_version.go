//go:build windows
// +build windows

package metrics

import (
	"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows/registry"
)

type OSVERSIONINFOEX struct {
	dwOSVersionInfoSize uint32
	dwMajorVersion      uint32
	dwMinorVersion      uint32
	dwBuildNumber       uint32
	dwPlatformId        uint32
	szCSDVersion        [128]uint16
	wServicePackMajor   uint16
	wServicePackMinor   uint16
	wSuiteMask          uint16
	wProductType        byte
	wReserved           byte
}

var (
	modntdll          = syscall.NewLazyDLL("ntdll.dll")
	procRtlGetVersion = modntdll.NewProc("RtlGetVersion")

	modkernel32        = syscall.NewLazyDLL("kernel32.dll")
	procGetProductInfo = modkernel32.NewProc("GetProductInfo")
)

func GetWindowsVersion() ([]string, error) {
	version, _, err := getWindowsVersion()
	if err != nil {
		fmt.Println("Failed to get version:", err)
		return []string{}, err
	}

	edition, err := getWindowsEditionFromRegistry()
	if err != nil {
		fmt.Println("Failed to get edition:", err)
		return []string{}, err
	}

	return []string{version, edition}, nil
}

func getWindowsVersion() (string, uint32, error) {
	var info OSVERSIONINFOEX
	info.dwOSVersionInfoSize = uint32(unsafe.Sizeof(info))

	r, _, err := procRtlGetVersion.Call(uintptr(unsafe.Pointer(&info)))
	if r != 0 {
		return "", 0, err
	}

	build := info.dwBuildNumber

	switch {
	case info.dwMajorVersion == 10 && build >= 22000:
		return "Windows 11", build, nil
	case info.dwMajorVersion == 10:
		return "Windows 10", build, nil
	case info.dwMajorVersion == 6 && info.dwMinorVersion == 3:
		return "Windows 8.1", build, nil
	case info.dwMajorVersion == 6 && info.dwMinorVersion == 2:
		return "Windows 8", build, nil
	case info.dwMajorVersion == 6 && info.dwMinorVersion == 1:
		return "Windows 7", build, nil
	default:
		return fmt.Sprintf("Unknown (major=%d minor=%d build=%d)",
			info.dwMajorVersion, info.dwMinorVersion, build), build, nil
	}
}

func getWindowsEditionFromRegistry() (string, error) {
	key, err := registry.OpenKey(registry.LOCAL_MACHINE,
		`SOFTWARE\Microsoft\Windows NT\CurrentVersion`, registry.QUERY_VALUE)
	if err != nil {
		return "", err
	}
	defer key.Close()

	editionID, _, err := key.GetStringValue("EditionID")
	if err != nil {
		return "", err
	}

	return editionID, nil
}

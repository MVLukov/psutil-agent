//go:build windows
// +build windows

package metrics

import (
	"syscall"
	"unsafe"
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

func GetWindowsVersion() (string, error) {
	mod := syscall.NewLazyDLL("ntdll.dll")
	proc := mod.NewProc("RtlGetVersion")

	var info OSVERSIONINFOEX
	info.dwOSVersionInfoSize = uint32(unsafe.Sizeof(info))

	r, _, err := proc.Call(uintptr(unsafe.Pointer(&info)))
	if r != 0 {
		return "", err
	}

	major := info.dwMajorVersion
	minor := info.dwMinorVersion
	build := info.dwBuildNumber

	version := "Unknown Windows version"

	switch {
	case major == 6 && minor == 1:
		version = "Windows 7"
	case major == 6 && minor == 2:
		version = "Windows 8"
	case major == 6 && minor == 3:
		version = "Windows 8.1"
	case major == 10 && build < 22000:
		version = "Windows 10"
	case major == 10 && build >= 22000:
		version = "Windows 11"
	case major == 6 && minor == 0:
		version = "Windows Vista or Server 2008"
	case major == 6 && minor == 1 && info.wProductType != 1:
		version = "Windows Server 2008 R2"
	case major == 10 && info.wProductType != 1 && build < 22000:
		version = "Windows Server 2016/2019"
	case major == 10 && info.wProductType != 1 && build >= 22000:
		version = "Windows Server 2022"
	}

	return version, nil
}

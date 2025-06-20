//go:build windows
// +build windows

package metrics

import (
	"fmt"
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

var (
	modntdll          = syscall.NewLazyDLL("ntdll.dll")
	procRtlGetVersion = modntdll.NewProc("RtlGetVersion")

	modkernel32        = syscall.NewLazyDLL("kernel32.dll")
	procGetProductInfo = modkernel32.NewProc("GetProductInfo")
)

func GetWindowsVersion() ([]string, error) {
	var info OSVERSIONINFOEX
	info.dwOSVersionInfoSize = uint32(unsafe.Sizeof(info))

	r, _, err := procRtlGetVersion.Call(uintptr(unsafe.Pointer(&info)))
	if r != 0 {
		return []string{}, err
	}

	major := info.dwMajorVersion
	minor := info.dwMinorVersion
	build := info.dwBuildNumber
	productType := info.wProductType

	// Get edition (Home, Pro, etc.)
	var productTypeCode uint32
	procGetProductInfo.Call(
		uintptr(major),
		uintptr(minor),
		uintptr(info.wServicePackMajor),
		uintptr(info.wServicePackMinor),
		uintptr(unsafe.Pointer(&productTypeCode)),
	)

	edition := mapProductType(productTypeCode)

	version := "Unknown Windows version"

	switch {
	case major == 6 && minor == 0:
		if productType == 1 {
			version = "Windows Vista"
		} else {
			version = "Windows Server 2008"
		}
	case major == 6 && minor == 1:
		if productType == 1 {
			version = "Windows 7"
		} else {
			version = "Windows Server 2008 R2"
		}
	case major == 6 && minor == 2:
		if productType == 1 {
			version = "Windows 8"
		} else {
			version = "Windows Server 2012"
		}
	case major == 6 && minor == 3:
		if productType == 1 {
			version = "Windows 8.1"
		} else {
			version = "Windows Server 2012 R2"
		}
	case major == 10 && build < 22000:
		if productType == 1 {
			version = "Windows 10"
		} else {
			version = "Windows Server 2016/2019"
		}
	case major == 10 && build >= 22000:
		if productType == 1 {
			version = "Windows 11"
		} else {
			version = "Windows Server 2022"
		}
	default:
		version = fmt.Sprintf("Unknown (major=%d minor=%d build=%d)", major, minor, build)
	}

	return []string{version, edition}, nil
}

func mapProductType(code uint32) string {
	switch code {
	case 0x00000006:
		return "Business"
	case 0x00000010:
		return "Business N"
	case 0x00000012:
		return "Cluster Server"
	case 0x00000008:
		return "Datacenter Server"
	case 0x0000000C:
		return "Datacenter Server Core"
	case 0x0000004F:
		return "Education"
	case 0x00000050:
		return "Education N"
	case 0x00000048:
		return "Enterprise"
	case 0x0000001B:
		return "Enterprise N"
	case 0x0000000A:
		return "Home Basic"
	case 0x0000002A:
		return "Home Basic N"
	case 0x00000003:
		return "Home Premium"
	case 0x0000001A:
		return "Home Premium N"
	case 0x00000005:
		return "Home Starter"
	case 0x0000002F:
		return "Home Starter N"
	case 0x00000013:
		return "Home Server"
	case 0x00000065:
		return "IoT Enterprise"
	case 0x0000002C:
		return "Professional"
	case 0x00000045:
		return "Professional N"
	case 0x00000067:
		return "Professional Workstation"
	case 0x00000068:
		return "Professional Workstation N"
	case 0x0000000E:
		return "Server Standard"
	case 0x0000000D:
		return "Server Standard Core"
	case 0x00000018:
		return "Ultimate"
	case 0x0000001C:
		return "Ultimate N"
	default:
		return "Unknown Edition"
	}
}

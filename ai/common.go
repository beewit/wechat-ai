package ai

import (
	"github.com/lxn/win"
	"syscall"
	"unsafe"
)

const (
	MOD_ALT      = 0x0001
	MOD_CONTROL  = 0x0002
	MOD_NOREPEAT = 0x400
	MOD_SHIFT    = 0x0004
	MOD_WIN      = 0x0008
)

const HOT_KEY_CTRL_SHFT_H = "HOT_KEY_CTRL_SHFT_H"

func IntPtr(n int) uintptr {
	return uintptr(n)
}

func StrPtr(s string) uintptr {
	return uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(s)))
}

func StrUint16(s string) *uint16 {
	return syscall.StringToUTF16Ptr(s)
}

func HotKey_Register(hWnd win.HWND, key int) bool {
	r, _, err := syscall.Syscall6(RegisterHotKey, 4,
		uintptr(hWnd),
		uintptr(StrPtr(HOT_KEY_CTRL_SHFT_H)),
		uintptr(MOD_CONTROL|MOD_SHIFT),
		uintptr(IntPtr(key)),
		0,
		0)
	errStr := err.Error()
	if errStr != "The operation completed successfully." {
		println(errStr)
	}
	return r != 0
}

func HotKey_Unregister(hWnd win.HWND, id string) bool {
	r, _, err := syscall.Syscall(RegisterHotKey, 2,
		uintptr(hWnd),
		uintptr(StrPtr(id)),
		0)
	errStr := err.Error()
	if errStr != "The operation completed successfully." {
		println(errStr)
	}
	return r != 0
}

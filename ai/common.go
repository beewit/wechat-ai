package ai

import (
	"errors"
	"fmt"
	"syscall"
	"time"
	"unsafe"

	"github.com/beewit/wechat-ai/enum"
	"github.com/lxn/win"
)

const (
	MOD_ALT      = 0x0001
	MOD_CONTROL  = 0x0002
	MOD_NOREPEAT = 0x400
	MOD_SHIFT    = 0x0004
	MOD_WIN      = 0x0008
)

const HOT_KEY_CTRL_SHFT_H = "HOT_KEY_CTRL_SHFT_H"

var (
	// Library
	libuser32 = win.MustLoadLibrary("user32.dll")

	// Functions
	Mouse = win.MustGetProcAddress(libuser32, "mouse_event")
	keyBD = win.MustGetProcAddress(libuser32, "keybd_event")

	CloseClipboard   = win.MustGetProcAddress(libuser32, "CloseClipboard")
	EmptyClipboard   = win.MustGetProcAddress(libuser32, "EmptyClipboard")
	OpenClipboard    = win.MustGetProcAddress(libuser32, "OpenClipboard")
	SetClipboardData = win.MustGetProcAddress(libuser32, "SetClipboardData")

	GetSystemMetrics = win.MustGetProcAddress(libuser32, "GetSystemMetrics")
	ClientToScreen   = win.MustGetProcAddress(libuser32, "ClientToScreen")
	RegisterHotKey   = win.MustGetProcAddress(libuser32, "RegisterHotKey")
	UnregisterHotKey = win.MustGetProcAddress(libuser32, "UnregisterHotKey")
)

func KeybdEven(bVk, bScan int) (uintptr, uintptr, error) {
	r, r2, err := syscall.Syscall(keyBD, 3, uintptr(bVk), 0, uintptr(bScan))
	println(err.Error())
	return r, r2, err
}

//鼠标左键操作
func MouseClick() {
	MouseEvent(win.MOUSEEVENTF_LEFTDOWN)
	MouseEvent(win.MOUSEEVENTF_LEFTUP)
}

//鼠标右键操作
func MouseRightClick() {
	MouseEvent(win.MOUSEEVENTF_RIGHTDOWN)
	MouseEvent(win.MOUSEEVENTF_RIGHTUP)
}

func MouseScroll(count int) {
	for i := 0; i < count; i++ {
		_, _, err := syscall.Syscall(Mouse, 3, uintptr(win.MOUSEEVENTF_WHEEL), 0, 120)
		errStr := err.Error()
		if errStr != "The operation completed successfully." {
			println(errStr)
		}
		time.Sleep(time.Millisecond)
	}
}

func MouseEvent(me int) (uintptr, uintptr, error) {
	r, r2, err := syscall.Syscall(Mouse, 3, uintptr(me), 0, 0)
	errStr := err.Error()
	if errStr != "The operation completed successfully." {
		println(errStr)
	}
	return r, r2, err
}

//键盘Control+V
func KeydbCV() {
	KeybdEven(win.VK_CONTROL, 0)
	KeybdEven(enum.VK_V, 0)
	KeybdEven(win.VK_CONTROL, win.KEYEVENTF_KEYUP)
	KeybdEven(enum.VK_V, win.KEYEVENTF_KEYUP)
}

//键盘Control+A
func KeydbCA() {
	KeybdEven(win.VK_CONTROL, 0)
	KeybdEven(enum.VK_A, 0)
	KeybdEven(win.VK_CONTROL, win.KEYEVENTF_KEYUP)
	KeybdEven(enum.VK_A, win.KEYEVENTF_KEYUP)
}

//键盘Control+Alt+W
func KeydbCSW() {
	KeybdEven(win.VK_CONTROL, 0)
	KeybdEven(win.VK_MENU, 0)
	KeybdEven(enum.VK_W, 0)
	KeybdEven(win.VK_CONTROL, win.KEYEVENTF_KEYUP)
	KeybdEven(win.VK_MENU, win.KEYEVENTF_KEYUP)
	KeybdEven(enum.VK_W, win.KEYEVENTF_KEYUP)
}

func KeydbKey(key int) {
	KeybdEven(key, 0)
	KeybdEven(key, win.KEYEVENTF_KEYUP)
}

//键盘BackSpace
func KeydbBack() {
	KeybdEven(win.VK_BACK, 0)
	KeybdEven(win.VK_BACK, win.KEYEVENTF_KEYUP)
}

func KeydbEnter() {
	KeybdEven(win.VK_RETURN, 0)
	time.Sleep(time.Second)
	KeybdEven(win.VK_RETURN, win.KEYEVENTF_KEYUP)
}

func GetScreen() (int32, int32) {
	x, _, err := syscall.Syscall(GetSystemMetrics, 1, win.SM_CXSCREEN, 0, 0)
	errStr := err.Error()
	if errStr != "The operation completed successfully." {
		println(errStr)
	}
	y, _, err := syscall.Syscall(GetSystemMetrics, 1, win.SM_CYSCREEN, 0, 0)
	errStr = err.Error()
	if errStr != "The operation completed successfully." {
		println(errStr)
	}
	return int32(x), int32(y)
}

func GetClientToScreen(hWnd win.HWND, rect *win.POINT) bool {
	ret, _, _ := syscall.Syscall(ClientToScreen, 2,
		uintptr(hWnd),
		uintptr(unsafe.Pointer(rect)),
		0)
	return ret != 0
}

func ForegroundWindow(winClass, winTitle string) (h win.HWND) {
	h = win.FindWindow(StrUint16(winClass), StrUint16(winTitle))
	win.SetForegroundWindow(h)
	return
}

func FindWindow(winClass, winTitle string, foreground bool, start func()) (h win.HWND, size enum.Size, rect win.RECT, err error) {
	h = win.FindWindow(StrUint16(winClass), StrUint16(winTitle))
	//激活窗体
	//win.MustGetProcAddress(libuser32, "SwitchToThisWindow")
	//win.ShowWindow(h, win.SW_SHOW)
	//win.SetWindowPos(h, win.HWND_TOP, 0, 0, 0, 0, win.SWP_NOSIZE)
	if start != nil {
		start()
	}
	flog := win.SetForegroundWindow(h)
	if !flog {
		err = errors.New("查找【" + winTitle + "】窗体失败")
		return
	}
	if foreground {
		//激活窗体
		flog = win.SetForegroundWindow(h)
		if !flog {
			err = errors.New("激活【" + winTitle + "】窗体失败")
			return
		}
	}
	var baseRect win.RECT
	flog = win.GetClientRect(h, &baseRect)
	if !flog {
		err = errors.New("获取【" + winTitle + "】窗体窗口大小失败")
		return
	}
	size.Width = baseRect.Right
	size.Height = baseRect.Bottom
	println(fmt.Sprintf("【"+winTitle+"】窗体大小：width：%v，height：%v", size.Width, size.Height))
	flog = win.GetWindowRect(h, &rect)
	if !flog {
		err = errors.New("查找【" + winTitle + "】窗体坐标失败")
		return
	}
	println(fmt.Sprintf("查找【"+winTitle+"】坐标：TOP：%v，Left：%v，Bottom：%v，Right：%v", rect.Top, rect.Left, rect.Bottom, rect.Right))
	return
}

func SetCursorPos(hwnd win.HWND, x, y int32) {
	var point win.POINT
	point.X = x
	point.Y = y
	GetClientToScreen(hwnd, &point)
	println(fmt.Sprintf("原始坐标：x：%v，y：%v，坐标：x：%v，y：%v", x, y, point.X, point.Y))
	win.SetCursorPos(point.X, point.Y)
}

func IntPtr(n int) uintptr {
	return uintptr(n)
}

func StrPtr(s string) uintptr {
	return uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(s)))
}

func StrUint16(s string) *uint16 {
	return syscall.StringToUTF16Ptr(s)
}

func HotKeyRegister(hWnd win.HWND, key int) bool {
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

func HotKeyUnregister(hWnd win.HWND, id string) bool {
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

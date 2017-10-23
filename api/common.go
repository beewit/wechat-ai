package api

import (
	"fmt"
	"syscall"

	"github.com/lxn/win"

	"time"
	"unsafe"

	"github.com/beewit/wechat-ai/enum"
	"github.com/pkg/errors"
	//	"net/url"
	//	"math/rand"
)

var (
	// Library
	libuser32 = win.MustLoadLibrary("user32.dll")

	// Functions
	mouse = win.MustGetProcAddress(libuser32, "mouse_event")
	keybd = win.MustGetProcAddress(libuser32, "keybd_event")

	CloseClipboard   = win.MustGetProcAddress(libuser32, "CloseClipboard")
	EmptyClipboard   = win.MustGetProcAddress(libuser32, "EmptyClipboard")
	OpenClipboard    = win.MustGetProcAddress(libuser32, "OpenClipboard")
	SetClipboardData = win.MustGetProcAddress(libuser32, "SetClipboardData")

	GetSystemMetrics = win.MustGetProcAddress(libuser32, "GetSystemMetrics")
)

//鼠标左键操作
func mouseClick() {
	mouse_event(win.MOUSEEVENTF_LEFTDOWN)
	mouse_event(win.MOUSEEVENTF_LEFTUP)
}

//鼠标右键操作
func mouseRightClick() {
	mouse_event(win.MOUSEEVENTF_RIGHTDOWN)
	mouse_event(win.MOUSEEVENTF_RIGHTUP)
}

func mouseScroll(count int) {
	for i := 0; i < count; i++ {
		_, _, err := syscall.Syscall(mouse, 3, uintptr(win.MOUSEEVENTF_WHEEL), 0, 120)
		println(err.Error())
		time.Sleep(time.Millisecond)
	}
}

func mouse_event(me int) (uintptr, uintptr, error) {
	r, r2, err := syscall.Syscall(mouse, 3, uintptr(me), 0, 0)
	println(err.Error())
	return r, r2, err
}

//键盘Control+V
func keydbCV() {
	keybd_even(win.VK_CONTROL, 0)
	keybd_even(enum.VK_V, 0)
	keybd_even(win.VK_CONTROL, win.KEYEVENTF_KEYUP)
	keybd_even(enum.VK_V, win.KEYEVENTF_KEYUP)
}

//键盘Control+A
func keydbCA() {
	keybd_even(win.VK_CONTROL, 0)
	keybd_even(enum.VK_A, 0)
	keybd_even(win.VK_CONTROL, win.KEYEVENTF_KEYUP)
	keybd_even(enum.VK_A, win.KEYEVENTF_KEYUP)
}

//键盘BackSpace
func keydbBack() {
	keybd_even(win.VK_BACK, 0)
	keybd_even(win.VK_BACK, win.KEYEVENTF_KEYUP)
}

func keydbEnter() {
	keybd_even(win.VK_RETURN, 0)
	time.Sleep(time.Second)
	keybd_even(win.VK_RETURN, win.KEYEVENTF_KEYUP)
}

func getScreen() (int32, int32) {
	x, r2, err := syscall.Syscall(GetSystemMetrics, 1, win.SM_CXSCREEN, 0, 0)
	println(x, r2, err.Error())
	y, r2, err := syscall.Syscall(GetSystemMetrics, 1, win.SM_CYSCREEN, 0, 0)
	println(y, r2, err.Error())
	return int32(x), int32(y)
}

func findWindow(winClass, winTitle string, foreground bool) (h win.HWND, size enum.Size, rect win.RECT, err error) {
	h = win.FindWindow(StrUint16(winClass), StrUint16(winTitle))
	//激活窗体
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
	fmt.Printf("【"+winTitle+"】窗体大小：width：%v，height：%v", size.Width, size.Height)
	flog = win.GetWindowRect(h, &rect)
	if !flog {
		err = errors.New("查找【" + winTitle + "】窗体坐标失败")
		return
	}
	fmt.Printf("查找【"+winTitle+"】坐标：TOP：%v，Left：%v，Bottom：%v，Right：%v", rect.Top, rect.Left, rect.Bottom, rect.Right)
	return
}

func Wechat(title string, off *enum.Offset) (err error) {
	//启动复制
	_, size, rect, err := findWindow("Chrome_WidgetWin_1", title, true)
	if err != nil {
		println(err.Error())
		return
	}
	webLeft := int32(off.Left) + 60
	webTop := int32(off.Top) + 100
	win.SetCursorPos(webLeft, webTop)
	time.Sleep(time.Second)
	//右键
	mouseRightClick()
	time.Sleep(time.Second)
	//鼠标移动到右键菜单复制
	win.SetCursorPos(webLeft+30, webTop+60)
	time.Sleep(time.Second)
	mouseClick()

	//屏幕分辨率
	x, y := getScreen()
	fmt.Printf("分辨率：x：%v,y:%v", x, y)
	//找到窗体句柄
	_, size, rect, err = findWindow("WeChatMainWndForPC", "微信测试版", true)
	if err != nil {
		println(err.Error())
		return
	}

	//最大化
	//win.PostMessage(wechatWin, win.WM_SYSCOMMAND, win.SC_MAXIMIZE, 0)

	//选择文件传输助手
	win.SetCursorPos(rect.Left+33, rect.Top+92)
	time.Sleep(time.Second)
	mouseClick()
	//mouseClick()
	time.Sleep(time.Second)
	win.SetCursorPos(rect.Left+190, rect.Top+80)
	time.Sleep(time.Second)
	mouseScroll(300)
	time.Sleep(time.Second)
	mouseClick()

	//鼠标定位会话编辑框
	win.SetCursorPos(rect.Left+360, rect.Top+size.Height-60)
	time.Sleep(time.Second)
	//var point win.POINT
	//win.GetCursorPos(&point)
	//println(fmt.Sprintf("%v", point))
	//win.SendMessage(0, win.MOUSEEVENTF_LEFTDOWN, 0, 0)
	//选择会话编辑框
	mouseClick()
	time.Sleep(time.Second)
	keydbCA()
	time.Sleep(time.Second)
	keydbBack()
	time.Sleep(time.Second)
	keydbCV()
	////发送消息
	//keydbEnter()
	//time.Sleep(time.Second)
	////鼠标定位二维码
	//win.SetCursorPos(rect.Left+size.Width-110, rect.Top+size.Height-200)
	time.Sleep(time.Second)
	//点击二维码，进入图片查看器
	mouseClick()
	time.Sleep(time.Millisecond)
	mouseClick()
	time.Sleep(time.Second)
	//获取图片二维码窗口信息
	_, _, rect, err = findWindow("ImagePreviewWnd", "图片查看器", false)
	if err != nil {
		println(err.Error())
		return
	}
	time.Sleep(time.Second)
	//鼠标移动到图片二维码中心位置
	imgLeft := rect.Left + 120
	imgTop := rect.Top + 200
	win.SetCursorPos(imgLeft, imgTop)
	time.Sleep(time.Second)
	//右键
	mouseRightClick()
	time.Sleep(time.Second)
	//鼠标移动到右键菜单识别二维码
	win.SetCursorPos(imgLeft+70, imgTop+90)
	time.Sleep(time.Second)
	//点击识别二维码
	mouseClick()
	time.Sleep(time.Second * 3)
	//查找加群窗口
	_, size, rect, err = findWindow("CefWebViewWnd", "微信测试版", false)
	if err != nil {
		println(err.Error())
		return
	}
	//鼠标移动到加群按钮
	win.SetCursorPos(rect.Left+size.Width/2, rect.Top+370)
	//点击按钮加群
	mouseClick()
	return nil
}

// SetText sets the current text data of the clipboard.
func SetText(s string) error {
	_, _, err2 := syscall.Syscall(OpenClipboard, 0, 0, 0, 0)
	println(err2.Error())
	println("1----------")

	utf16, err := syscall.UTF16FromString(s)
	if err != nil {
		return err
	}

	hMem := win.GlobalAlloc(win.GMEM_MOVEABLE, uintptr(len(utf16)*2))
	if hMem == 0 {
		return errors.New("GlobalAlloc")
	}

	p := win.GlobalLock(hMem)
	if p == nil {
		return errors.New("GlobalLock()")
	}

	win.MoveMemory(p, unsafe.Pointer(&utf16[0]), uintptr(len(utf16)*2))

	win.GlobalUnlock(hMem)

	if 0 == win.SetClipboardData(win.CF_UNICODETEXT, win.HANDLE(hMem)) {
		// We need to free hMem.
		defer win.GlobalFree(hMem)

		return errors.New("SetClipboardData")
	}
	return nil
}

func keybd_even(bVk, bScan int) (uintptr, uintptr, error) {
	r, r2, err := syscall.Syscall(keybd, 3, uintptr(bVk), 0, uintptr(bScan))
	println(err.Error())
	return r, r2, err
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

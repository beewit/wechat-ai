package ai

import (
	"fmt"
	"time"
	"github.com/beewit/wechat-ai/enum"
	//	"net/url"
	//	"math/rand"
)

func Wechat(title string, off *enum.Offset) (err error) {
	//启动复制
	chromeWin, size, _, err := FindWindow("Chrome_WidgetWin_1", title, true, nil)
	if err != nil {
		println(err.Error())
		return
	}

	SetCursorPos(chromeWin, int32(off.Left)+150, int32(off.Top)+150)

	//【1】不同电脑上定位准确，转换屏幕坐标到客户坐标
	//webLeft := int32(off.Left) + 150
	//webTop := int32(off.Top) + 150
	//
	//var point win.POINT
	//point.X = webLeft
	//point.Y = webTop
	//
	//GetClientToScreen(chromeWin, &point)
	//win.SetCursorPos(point.X, point.Y)

	//【2】不同电脑上定位不准确，未转换屏幕坐标到客户坐标
	//win.SetCursorPos(webLeft, webTop)

	time.Sleep(time.Second)
	//右键
	MouseRightClick()
	time.Sleep(time.Second)
	//快捷键Y
	KeydbKey(enum.VK_Y)
	time.Sleep(time.Second)
	////鼠标移动到右键菜单复制
	//win.SetCursorPos(point.X+30, point.Y+60)
	//time.Sleep(time.Second)
	//mouseClick()

	//屏幕分辨率
	x, y := GetScreen()
	println(fmt.Sprintf("分辨率：x：%v,y:%v", x, y))
	time.Sleep(time.Second * 2)
	//找到窗体句柄
	wechatWin, size, _, err := FindWindow("WeChatMainWndForPC", "微信", true, func() {
		KeydbCSW()
	})
	if err != nil {
		wechatWin, size, _, err = FindWindow("WeChatMainWndForPC", "微信测试版", true, nil)
		if err != nil {
			println(err.Error())
			return
		}
	}
	//最大化
	//win.PostMessage(wechatWin, win.WM_SYSCOMMAND, win.SC_MAXIMIZE, 0)

	//选择文件传输助手

	SetCursorPos(wechatWin, 33, 95)
	//win.SetCursorPos(rect.Left+33, rect.Top+92)
	time.Sleep(time.Second)
	MouseClick()
	//mouseClick()
	time.Sleep(time.Second)
	SetCursorPos(wechatWin, 190, 85)
	//win.SetCursorPos(rect.Left+190, rect.Top+80)
	time.Sleep(time.Second)
	MouseScroll(300)
	time.Sleep(time.Second)
	MouseClick()

	//鼠标定位会话编辑框
	if x < 1550 && x > 1500 {
		println("【定位会话框位置】470")
		SetCursorPos(wechatWin, 470, size.Height-65)
	} else {
		SetCursorPos(wechatWin, 380, size.Height-65)
		println("【定位会话框位置】380")
	}
	//win.SetCursorPos(rect.Left+360, rect.Top+size.Height-60)
	time.Sleep(time.Second)
	//var point win.POINT
	//win.GetCursorPos(&point)
	//println(fmt.Sprintf("%v", point))
	//win.SendMessage(0, win.MOUSEEVENTF_LEFTDOWN, 0, 0)
	//选择会话编辑框
	MouseClick()
	time.Sleep(time.Second)
	KeydbCA()
	time.Sleep(time.Second)
	KeydbBack()
	time.Sleep(time.Second)
	KeydbCV()
	////发送消息
	//keydbEnter()
	//time.Sleep(time.Second)
	////鼠标定位二维码
	//win.SetCursorPos(rect.Left+size.Width-110, rect.Top+size.Height-200)
	time.Sleep(time.Second)
	//点击二维码，进入图片查看器
	MouseClick()
	time.Sleep(time.Millisecond)
	MouseClick()
	time.Sleep(time.Second)
	//获取图片二维码窗口信息
	imageWin, _, _, err := FindWindow("ImagePreviewWnd", "图片查看器", false, nil)
	if err != nil {
		println(err.Error())
		return
	}
	time.Sleep(time.Second)
	//鼠标移动到图片二维码中心位置
	//imgLeft := rect.Left + 120
	//imgTop := rect.Top + 200
	SetCursorPos(imageWin, 120, 200)
	//win.SetCursorPos(imgLeft, imgTop)
	time.Sleep(time.Second)
	//右键
	MouseRightClick()
	time.Sleep(time.Second)
	//鼠标移动到右键菜单识别二维码
	SetCursorPos(imageWin, 190, 290)
	//win.SetCursorPos(imgLeft+70, imgTop+90)
	time.Sleep(time.Second)
	//点击识别二维码
	MouseClick()
	time.Sleep(time.Second * 3)
	//查找加群窗口
	cefWebViewWndWin, size, _, err := FindWindow("CefWebViewWnd", "微信", false, nil)
	if err != nil {
		cefWebViewWndWin, size, _, err = FindWindow("CefWebViewWnd", "微信测试版", false, nil)
		if err != nil {
			println(err.Error())
			return
		}
	}
	//鼠标移动到加群按钮
	SetCursorPos(cefWebViewWndWin, size.Width/2, 370)
	//win.SetCursorPos(rect.Left+size.Width/2, rect.Top+370)
	time.Sleep(time.Second * 2)
	//点击按钮加群
	MouseClick()
	return nil
}

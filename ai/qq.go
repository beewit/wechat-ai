package ai

import (
	"github.com/beewit/beekit/utils"
	"time"
	"github.com/lxn/win"
	"github.com/beewit/wechat-ai/enum"
	"github.com/pkg/errors"
	"github.com/beewit/beekit/utils/convert"
)

func QQLogin(qq int64, pwd string) (err error) {
	err = utils.StartQQ()
	if err != nil {
		return
	}
	var qqWin win.HWND
	qqWin, _, err = getQQWin("QQ", 1)
	if err != nil {
		return
	}
	time.Sleep(time.Second * 2)
	SetCursorPos(qqWin, 230, 285)
	time.Sleep(time.Second * 2)
	MouseClick()
	time.Sleep(time.Second * 2)
	KeydbStr(convert.ToString(qq))
	SetCursorPos(qqWin, 230, 310)
	time.Sleep(time.Second * 2)
	MouseClick()
	time.Sleep(time.Second * 2)
	KeydbBack()
	time.Sleep(time.Second * 2)
	KeydbStr(pwd)
	time.Sleep(time.Second * 2)
	SetCursorPos(qqWin, 230, 370)
	time.Sleep(time.Second * 2)
	MouseClick()
	var size enum.Size
	qqWin, size, err = getQQWin("QQ", 1)
	if err != nil {
		return
	}
	if size.Height < 520 {
		err = errors.New("登录失败")
		return
	}
	return
}

//获取登录QQ界面win.HWND，检查次数5次
func getQQWin(winTitle string, count int) (qqWin win.HWND, size enum.Size, err error) {
	time.Sleep(time.Second * 5)
	qqWin, size, _, err = FindWindow("TXGuiFoundation", winTitle, false, nil)
	if err != nil {
		if count < 6 {
			count++
			getQQWin(winTitle, count)
		}
	}
	return
}

func AddQQGroup(qq int64, remark string) (err error) {
	var qqWin, addGroupWin win.HWND
	var size enum.Size
	qqWin, size, err = getQQWin("QQ", 1)
	if err != nil {
		return
	}
	SetCursorPos(qqWin, 56, size.Height-26)
	MouseClick()
	time.Sleep(time.Second * 3)
	qqWin, _, err = getQQWin("查找", 1)
	if err != nil {
		return
	}
	//找群
	time.Sleep(time.Second * 2)
	SetCursorPos(qqWin, 334, 47)
	time.Sleep(time.Second * 2)
	MouseClick()
	//输入QQ群号
	time.Sleep(time.Second * 2)
	SetCursorPos(qqWin, 175, 110)
	time.Sleep(time.Second * 2)
	MouseClick()
	time.Sleep(time.Second * 2)
	KeydbCA()
	time.Sleep(time.Second * 2)
	KeydbStr(convert.ToString(qq))
	time.Sleep(time.Second * 2)
	SetCursorPos(qqWin, 650, 110)
	time.Sleep(time.Second * 2)
	MouseClick()
	//加群
	time.Sleep(time.Second * 5)
	SetCursorPos(qqWin, 270, 345)
	time.Sleep(time.Second * 2)
	MouseClick()
	time.Sleep(time.Second * 3)
	//加群弹窗
	addGroupWin, _, err = getQQWin("添加群", 1)
	if err != nil {
		return
	}
	//输入验证信息
	time.Sleep(time.Second * 2)
	SetCursorPos(addGroupWin, 210, 100)
	time.Sleep(time.Second * 2)
	SetClipboard(remark)
	time.Sleep(time.Second * 2)
	KeydbCV()
	time.Sleep(time.Second * 2)
	SetCursorPos(addGroupWin, 345, 345)
	time.Sleep(time.Second * 2)
	MouseClick()
	time.Sleep(time.Second * 2)
	SetCursorPos(addGroupWin, 421, 348)
	time.Sleep(time.Second * 2)
	MouseClick()
	return
}

func AddQQFriend(qq int64, remark string) (err error) {
	var qqWin, addFriendWin win.HWND
	qqWin, _, err = getQQWin("QQ", 1)
	if err != nil {
		return
	}
	SetCursorPos(qqWin, 56, 686)
	MouseClick()
	time.Sleep(time.Second * 3)
	qqWin, _, err = getQQWin("查找", 1)
	if err != nil {
		return
	}
	//找人
	time.Sleep(time.Second * 2)
	SetCursorPos(qqWin, 233, 47)
	time.Sleep(time.Second * 2)
	MouseClick()
	//输入QQ号
	time.Sleep(time.Second * 2)
	SetCursorPos(qqWin, 100, 106)
	time.Sleep(time.Second * 2)
	MouseClick()
	time.Sleep(time.Second * 2)
	KeydbCA()
	time.Sleep(time.Second * 2)
	KeydbStr(convert.ToString(qq))
	time.Sleep(time.Second * 5)
	SetCursorPos(qqWin, 654, 120)
	time.Sleep(time.Second * 2)
	MouseClick()
	//加人
	time.Sleep(time.Second * 2)
	SetCursorPos(qqWin, 127, 312)
	time.Sleep(time.Second * 2)
	MouseClick()
	time.Sleep(time.Second * 3)
	// - 添加好友
	addFriendWin, _, _, err = ThisWindow() // FindWindow( "工蜂小智 - 添加好友", false, nil)
	if err != nil {
		return
	}
	//输入验证信息
	time.Sleep(time.Second * 2)
	SetCursorPos(addFriendWin, 210, 100)
	time.Sleep(time.Second * 2)
	MouseClick()
	time.Sleep(time.Second * 2)
	KeydbCA()
	time.Sleep(time.Second * 2)
	SetClipboard(remark)
	time.Sleep(time.Second * 2)
	KeydbCV()
	time.Sleep(time.Second * 2)
	SetCursorPos(addFriendWin, 345, 345)
	time.Sleep(time.Second * 2)
	MouseClick()
	time.Sleep(time.Second * 2)
	MouseClick()
	time.Sleep(time.Second * 2)
	SetCursorPos(addFriendWin, 421, 348)
	time.Sleep(time.Second * 2)
	MouseClick()
	return
}

func QQ() (err error) {
	err = utils.StartQQ()
	if err != nil {
		return
	}
	time.Sleep(time.Second * 3)
	var qqWin, addGroupWin win.HWND
	qqWin, _, err = getQQWin("QQ", 1)
	if err != nil {
		return
	}
	SetCursorPos(qqWin, 230, 285)
	time.Sleep(time.Second * 2)
	MouseClick()
	time.Sleep(time.Second * 2)
	KeydbStr("3240033436")
	SetCursorPos(qqWin, 230, 310)
	time.Sleep(time.Second * 2)
	MouseClick()
	time.Sleep(time.Second * 2)
	KeydbBack()
	time.Sleep(time.Second * 2)
	KeydbStr("13696433488love")
	time.Sleep(time.Second * 2)
	SetCursorPos(qqWin, 230, 370)
	time.Sleep(time.Second * 2)
	MouseClick()
	time.Sleep(time.Second * 3)
	qqWin, _, err = getQQWin("QQ", 1)
	if err != nil {
		return
	}
	time.Sleep(time.Second * 2)
	SetCursorPos(qqWin, 56, 686)
	MouseClick()
	time.Sleep(time.Second * 3)
	qqWin, _, err = getQQWin("查找", 1)
	if err != nil {
		return
	}
	//找群
	time.Sleep(time.Second * 2)
	SetCursorPos(qqWin, 334, 47)
	time.Sleep(time.Second * 2)
	MouseClick()
	//输入QQ群号
	time.Sleep(time.Second * 2)
	SetCursorPos(qqWin, 175, 110)
	time.Sleep(time.Second * 2)
	MouseClick()
	time.Sleep(time.Second * 2)
	KeydbStr("553147171")
	time.Sleep(time.Second * 2)
	SetCursorPos(qqWin, 650, 110)
	time.Sleep(time.Second * 2)
	MouseClick()
	//加群
	time.Sleep(time.Second * 2)
	SetCursorPos(qqWin, 270, 345)
	time.Sleep(time.Second * 2)
	MouseClick()
	time.Sleep(time.Second * 3)
	//加群弹窗
	addGroupWin, _, err = getQQWin("添加群", 1)
	if err != nil {
		return
	}
	//输入验证信息
	time.Sleep(time.Second * 2)
	SetCursorPos(addGroupWin, 210, 100)
	time.Sleep(time.Second * 2)
	SetClipboard("你好啊")
	time.Sleep(time.Second * 2)
	KeydbCV()
	time.Sleep(time.Second * 2)
	SetCursorPos(addGroupWin, 345, 345)
	time.Sleep(time.Second * 2)
	MouseClick()
	return
}

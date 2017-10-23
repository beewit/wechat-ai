package main

import (
	"github.com/beewit/wechat-client/send"
	"github.com/beewit/wechat-client/enum"
	"github.com/beewit/beekit/utils"
	"regexp"
	"strings"
	"fmt"
	"time"
	"os"
)

//import (
//	"github.com/beewit/wechat-client/api"
//	"github.com/beewit/wechat-client/enum"
//	"github.com/sclevine/agouti"
//	"time"
//)

var (
	uuid       string
	err        error
	loginMap   send.LoginMap
	contactMap map[string]send.User
	groupMap   map[string][]send.User /* 关键字为key的，群组数组 */
)

func main() {
	//file, err := os.Open("sam.jpg")
	//if err != nil {
	//	println(err.Error())
	//	return
	//}
	////base64.StdEncoding()
	//r, r2, err := syscall.Syscall(OpenClipboard, 0, 0, 0, 0)
	//println(r, r2, err.Error())
	//r, r2, err = syscall.Syscall(GetClipboardFormatName, 0, 0, 0, 0)
	//println(r, r2, err.Error())

	//driver := agouti.ChromeDriver(agouti.ChromeOptions("args", []string{
	//	"--user-data-dir=ChromeUserData",
	//	"--gpu-process",
	//	"--start-maximized",
	//	"--disable-infobars",
	//	"--app=http://weixinqun.com/group?id=1036320",
	//	"--webkit-text-size-adjust"}))
	//driver.Start()
	//page, err := driver.NewPage()
	//if err != nil {
	//	println("Failed to open page.")
	//	return
	//}
	////page.Navigate("http://weixinqun.com/group?id=1036320")
	//page.RunScript(`$(".checkCode span:eq(1)").mouseover()`, nil, nil)
	//time.Sleep(time.Second * 3)
	//var of *enum.Offset
	//page.RunScript(`return  $(".shiftcode:eq(1) img").offset()`, nil, &of)
	//title, err := page.Title()
	//if err != nil {
	//	println("error:" + err.Error())
	//	return
	//}
	//if of != nil {
	//	api.Wechat(title, of)
	//}

	/* 从微信服务器获取UUID */
	uuid, err := send.GetUUIDFromWX()
	if err != nil {
		panic(err)
	}
	/* 根据UUID获取二维码 */
	err = send.DownloadImagIntoDir(enum.QRCODE_URL+uuid, "qrcode")
	if err != nil {
		panic(err)
	}
	utils.Open("qrcode/qrcode.jpg")
	/* 轮询服务器判断二维码是否扫过暨是否登陆了 */
	for {
		fmt.Println("正在验证登陆... ...")
		status, msg := send.CheckLogin(uuid)

		if status == 200 {
			fmt.Println("登陆成功,处理登陆信息...")
			loginMap, err := send.ProcessLoginInfo(msg)
			if err != nil {
				panic(err)
			}

			fmt.Println("登陆信息处理完毕,正在初始化微信...")
			err = send.InitWX(&loginMap)
			if err != nil {
				panic(err)
			}

			fmt.Println("初始化完毕,通知微信服务器登陆状态变更...")
			err = send.NotifyStatus(&loginMap)
			if err != nil {
				panic(err)
			}

			fmt.Println("通知完毕,本次登陆信息：")
			fmt.Println(enum.SKey + "\t\t" + loginMap.BaseRequest.SKey)
			fmt.Println(enum.PassTicket + "\t\t" + loginMap.PassTicket)
			break
		} else if status == 201 {
			fmt.Println("请在手机上确认")
		} else if status == 408 {
			fmt.Println("请扫描二维码")
		} else {
			fmt.Println(msg)
		}
	}

	fmt.Println("开始获取联系人信息...")
	contactMap, err := send.GetAllContact(&loginMap)
	if err != nil {
		panic(err)
	}
	fmt.Printf("成功获取 %d个 联系人信息,开始整理群组信息...\n", len(contactMap))

	groupMap = send.MapGroupInfo(contactMap)
	groupSize := 0
	for _, v := range groupMap {
		groupSize += len(v)
	}
	fmt.Printf("整理完毕，共有 %d个 群组是焦点群组，它们是：\n", groupSize)
	for key, v := range groupMap {
		fmt.Println(key)
		for _, user := range v {
			fmt.Println("========>" + user.NickName)
		}
	}

	fmt.Println("开始监听消息响应...")
	var retcode, selector int64
	regAt := regexp.MustCompile(`^@.*@.*丁丁.*$`) /* 群聊时其他人说话时会在前面加上@XXX */
	regGroup := regexp.MustCompile(`^@@.+`)
	regAd := regexp.MustCompile(`(朋友圈|点赞)+`)
	secretCode := "cmd"
	for {
		retcode, selector, err = send.SyncCheck(&loginMap)
		if err != nil {
			fmt.Println(retcode, selector)
			panic(err)
			if retcode == 1101 {
				fmt.Println("帐号已在其他地方登陆，程序将退出。")
				os.Exit(2)
			}
			continue
		}

		if retcode == 0 && selector != 0 {
			fmt.Printf("selector=%d,有新消息产生,准备拉取...\n", selector)
			wxRecvMsges, err := send.WebWxSync(&loginMap)
			panic(err)

			for i := 0; i < wxRecvMsges.MsgCount; i++ {
				if wxRecvMsges.MsgList[i].MsgType == 1 {
					/* 普通文本消息 */
					fmt.Println(
						contactMap[wxRecvMsges.MsgList[i].FromUserName].NickName+":",
						wxRecvMsges.MsgList[i].Content)

					if regGroup.MatchString(wxRecvMsges.MsgList[i].FromUserName) && regAt.MatchString(wxRecvMsges.MsgList[i].Content) {
						/* 有人在群里@我，发个消息回答一下 */
						wxSendMsg := send.WxSendMsg{}
						wxSendMsg.Type = 1
						wxSendMsg.Content = "我是丁丁编写的微信机器人，我已帮你通知我的主人，请您稍等片刻，他会跟您联系"
						wxSendMsg.FromUserName = wxRecvMsges.MsgList[i].ToUserName
						wxSendMsg.ToUserName = wxRecvMsges.MsgList[i].FromUserName
						wxSendMsg.LocalID = fmt.Sprintf("%d", time.Now().Unix())
						wxSendMsg.ClientMsgId = wxSendMsg.LocalID

						/* 加点延时，避免消息次序混乱，同时避免微信侦察到机器人 */
						time.Sleep(time.Second)

						go send.SendMsg(&loginMap, wxSendMsg)
					} else if (!regGroup.MatchString(wxRecvMsges.MsgList[i].FromUserName)) && regAd.MatchString(wxRecvMsges.MsgList[i].Content) {
						/* 有人私聊我，并且内容含有「朋友圈」、「点赞」等敏感词，则回复 */
						wxSendMsg := send.WxSendMsg{}
						wxSendMsg.Type = 1
						wxSendMsg.Content = "我是丁丁编写的微信机器人，我已经感应到一些敏感字符，我的主人对微商、朋友圈集赞、砍价等活动不感兴趣，请不要再发这些请求给他了，Sorry！"
						wxSendMsg.FromUserName = wxRecvMsges.MsgList[i].ToUserName
						wxSendMsg.ToUserName = wxRecvMsges.MsgList[i].FromUserName
						wxSendMsg.LocalID = fmt.Sprintf("%d", time.Now().Unix())
						wxSendMsg.ClientMsgId = wxSendMsg.LocalID

						time.Sleep(time.Second)

						go send.SendMsg(&loginMap, wxSendMsg)
					} else if (!regGroup.MatchString(wxRecvMsges.MsgList[i].FromUserName)) && strings.EqualFold(wxRecvMsges.MsgList[i].Content, secretCode) {
						/* 有人私聊我，并且内容是密语，输出加群菜单 */
						wxSendMsg := send.WxSendMsg{}
						wxSendMsg.Type = 1
						wxSendMsg.Content = fmt.Sprintf("我是丁丁编写的微信群聊助手，我为您提供了以下分组关键词：\n\n%s\n您可以输入上方关键词或者其索引号获取更详细的群聊目录，比如您可以输入\"编程\"或者\"1-0\"，我会为您细化系统内所有与计算机编程相关的领域，您可以进一步选择加入该领域的微信群聊。", e.GetFatherKeywordsStr())
						wxSendMsg.FromUserName = wxRecvMsges.MsgList[i].ToUserName
						wxSendMsg.ToUserName = wxRecvMsges.MsgList[i].FromUserName
						wxSendMsg.LocalID = fmt.Sprintf("%d", time.Now().Unix())
						wxSendMsg.ClientMsgId = wxSendMsg.LocalID

						time.Sleep(time.Second)

						go send.SendMsg(&loginMap, wxSendMsg)
					} else if !regGroup.MatchString(wxRecvMsges.MsgList[i].FromUserName) {
						/* 有人私聊我，依次判断是否是父级目录结构，输出子目录 */
						count := 0
						content := wxRecvMsges.MsgList[i].Content
						for _, v := range e.GetFocusGroupKeywords() {
							reg := regexp.MustCompile("^(" + strings.ToLower(v.FatherName) + ")$")
							if reg.MatchString(strings.ToLower(content)) {
								/* 判断为父级目录 */
								description, keywords, exampleStr := e.GetChildKeywordsInfo(v.FatherName)
								wxSendMsg := send.WxSendMsg{}
								wxSendMsg.Type = 1
								wxSendMsg.Content = fmt.Sprintf("%s\n\n%s\n您可以输入以上关键词或者其索引号，我会为您寻找系统内的微信群组并邀请您加入，比如您可以输入%s", description, exampleStr, keywords)
								wxSendMsg.FromUserName = wxRecvMsges.MsgList[i].ToUserName
								wxSendMsg.ToUserName = wxRecvMsges.MsgList[i].FromUserName
								wxSendMsg.LocalID = fmt.Sprintf("%d", time.Now().Unix())
								wxSendMsg.ClientMsgId = wxSendMsg.LocalID

								time.Sleep(time.Second)

								go send.SendMsg(&loginMap, wxSendMsg)
								break
							}
							count++
						}

						if count != len(send.GetFocusGroupKeywords()) {
							continue
						}

						/* 依次判断是否为子目录 */
						for key, groupUsers := range groupMap {
							reg := regexp.MustCompile("^(" + strings.ToLower(key) + ")$")
							if reg.MatchString(strings.ToLower(content)) {
								/* 判断为子目录 */
								for _, user := range groupUsers {
									time.Sleep(time.Second)
									go send.InviteMember(&loginMap, wxRecvMsges.MsgList[i].FromUserName, user.UserName)
								}

								break
							}
						}
					}
				}
			}
		}
	}
}

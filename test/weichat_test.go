package test

import (
	"testing"
	"github.com/beewit/wechat-ai/send"
	"fmt"
	"github.com/beewit/beekit/utils/convert"
	"github.com/beewit/wechat-ai/enum"
	"os"
	"encoding/base64"
	"strings"
	"io/ioutil"
	"encoding/json"
	"time"
	"net/http"
	"github.com/beewit/beekit/redis"
)

func TestTimUninx(t *testing.T) {
	println(fmt.Sprintf("%d", time.Now().UnixNano()/1000000))
}

func TestJson(t *testing.T) {
	b, _ := ioutil.ReadFile("userinfo.json")
	var initInfo *send.InitInfo
	err := json.Unmarshal(b, &initInfo)
	if err != nil {
		t.Error(err.Error())
		return
	}
	printlnInfo(initInfo)
}

func printlnInfo(initInfo *send.InitInfo) {
	if initInfo != nil {
		for _, v := range initInfo.AllContactList {
			for _, vv := range v.MemberList {
				println("【" + v.NickName + "】" + vv.UserName)
			}
		}
	}
}

func printlnInfo2(loginMap *send.LoginMap) {
	initInfo := loginMap.InitInfo
	if initInfo != nil {
		for _, v := range initInfo.AllContactList {
			for _, vv := range v.MemberList {
				println("【" + v.NickName + "】" + vv.UserName)
				if v.NickName == "工蜂小智内测" {
					vu := send.VerifyUser{}
					vu.Value = vv.UserName
					b, err := send.AddUser(loginMap, "你好，我是工蜂小智助手", []send.VerifyUser{vu})
					if err != nil {
						println(err.Error())
					} else {
						println(b.BaseResponse.Ret)
					}
				}
			}
		}
	}
}

func TestAddUser(t *testing.T) {
	var timestamp int64 = time.Now().UnixNano() / 1000000
	urlMap := map[string]string{}
	urlMap[enum.R] = fmt.Sprintf("%d", ^(int32)(timestamp))

	var br send.BaseRequest
	err := json.Unmarshal([]byte(`{"Uin":3374717726,"Sid":"WsrVCY8kERN6yF/6","Skey":"@crypt_4ffa22b9_adeb7efadd694c33864662a5e68ce7a9","DeviceID":"e270484016143139"}`), &br)
	if err != nil {
		println(err.Error())
		return
	}
	vu := send.VerifyUser{}
	vu.Value = "@aa99026aca9686d9b6446db973c1c1342efe2fc5bd4bf12aeadeb34a503edbb8"

	wxAddUser := send.WxAddUser{}
	wxAddUser.SKey = "@crypt_4ffa22b9_adeb7efadd694c33864662a5e68ce7a9"
	wxAddUser.VerifyContent = "我是执手并肩看天下"
	wxAddUser.SceneListCount = 1
	wxAddUser.SceneList = []int{33}
	wxAddUser.Opcode = 2
	wxAddUser.BaseRequest = br
	wxAddUser.VerifyUserList = []send.VerifyUser{vu}
	wxAddUser.VerifyUserListSize = 1
	jsonBytes, err := json.Marshal(wxAddUser)
	if err != nil {
		return
	}
	println(string(jsonBytes))
	// TODO: 发送微信消息时暂不处理返回值
	rep, err := http.Post(enum.WEB_WX_VERIFY_USER+send.GetURLParams(urlMap), enum.JSON_HEADER, strings.NewReader(string(jsonBytes)))
	if err != nil {
		return
	}
	bts, err := ioutil.ReadAll(rep.Body)
	if err != nil {
		println(err.Error())
	} else {
		println(string(bts))
	}
}

func TestLoginWX(t *testing.T) {
	var err error
	var SendStatusMsg string
	var LoginMap *send.LoginMap
	var UUid string
	wlm, err := redis.Cache.GetString("wechat_login_map")
	if wlm == "" {
		/* 从微信服务器获取UUID */
		UUid, err = send.GetUUIDFromWX()
		if err != nil {
			t.Error("GetUUIDFromWX Error：" + err.Error())

			return
		}
		/* 根据UUID获取二维码 */
		base64Img, err := send.DownloadImage(enum.QRCODE_URL + UUid)
		if err != nil {
			t.Error("DownloadImage Error：" + err.Error())
			return
		}
		//解压
		dist, _ := base64.StdEncoding.DecodeString(strings.Replace(base64Img, "data:image/jpeg;base64,", "", -1))
		//写入新文件
		f, err := os.OpenFile("output.jpg", os.O_RDWR|os.O_CREATE, os.ModePerm)
		defer f.Close()
		f.Write(dist)
		if err != nil {
			t.Error("output.jpg Error：" + err.Error())
			return
		}
		for {
			SendStatusMsg = "【" + UUid + "】正在验证登陆... ..."
			t.Log(SendStatusMsg)
			status, msg := send.CheckLogin(UUid)
			if status == 200 {
				SendStatusMsg = "登陆成功,处理登陆信息..."
				t.Log(SendStatusMsg)
				LoginMap, err = send.ProcessLoginInfo(msg)
				if err != nil {
					SendStatusMsg = "错误：登陆成功,处理登陆信息...，error：" + err.Error()
					t.Log(SendStatusMsg)
					return
				}
				SendStatusMsg = "登陆信息处理完毕,正在初始化微信..."
				t.Log(SendStatusMsg)
				err = send.InitWX(LoginMap)
				if err != nil {
					if err != nil {
						SendStatusMsg = "错误：登陆信息处理完毕,正在初始化微信...，error：" + err.Error()
						t.Log(SendStatusMsg)
						return
					}
				}
				SendStatusMsg = "初始化完毕,通知微信服务器登陆状态变更..."
				t.Log(SendStatusMsg)
				err = send.NotifyStatus(LoginMap)
				if err != nil {
					panic(err)
				}
				SendStatusMsg = "通知完毕,本次登陆信息获取成功"
				t.Log(SendStatusMsg)
				//fmt.Println(enum.SKey + "\t\t" + loginMap.BaseRequest.SKey)
				//fmt.Println(enum.PassTicket + "\t\t" + loginMap.PassTicket)
				break
			} else if status == 201 {
				SendStatusMsg = "请在手机上确认登录"
				t.Log(SendStatusMsg)
			} else if status == 408 {
				SendStatusMsg = "请扫描登录二维码"
				t.Log(SendStatusMsg)
			} else {
				SendStatusMsg = fmt.Sprintf("未知情况，返回状态码：%d", status)
				t.Log(SendStatusMsg)
			}
		}
		redis.Cache.SetAndExpire("wechat_login_map", convert.ToObjStr(LoginMap), 60*60*60)
	} else {
		json.Unmarshal([]byte(wlm), &LoginMap)
	}
	SendStatusMsg = "开始获取联系人信息..."
	t.Log(SendStatusMsg)
	ContactMap, err := send.GetAllContact(LoginMap)
	if err != nil {
		SendStatusMsg = "错误：开始获取联系人信息...，error：" + err.Error()
		t.Log(SendStatusMsg)
	}
	ss := convert.ToObjStr(ContactMap)

	t.Log("联系人信息" + ss)
	SendStatusMsg = "【" + convert.ToString(len(ContactMap)) + "】准备群发消息..."
	t.Log(SendStatusMsg)

	printlnInfo2(LoginMap)
}

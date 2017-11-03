package smartQQ

import (
	"time"
	"net/http"
	"fmt"
	"net/url"
	"github.com/beewit/beekit/utils"
	"strings"
	"github.com/pkg/errors"
	"regexp"
	"io/ioutil"
	"os"
	"encoding/base64"
	"github.com/beewit/beekit/utils/convert"
	"encoding/json"
	"github.com/beewit/wechat-ai/smartWechat"
)

var (
	loginUrl              = "https://ui.ptlogin2.qq.com/cgi-bin/login?daid=164&target=self&style=16&mibao_css=m_webqq&appid=501004106&enable_qlogin=0&no_verifyimg=1&s_url=http%3A%2F%2Fw.qq.com%2Fproxy.html&f_url=loginerroralert&strong_login=1&login_state=10&t=20131024001"
	qrShowUrl             = "https://ssl.ptlogin2.qq.com/ptqrshow?appid=501004106&e=2&l=M&s=3&d=72&v=4&t=0.1509281804550%v&daid=164&pt_3rd_aid=0"
	qrLoginUrl            = "https://ssl.ptlogin2.qq.com/ptqrlogin?ptqrtoken={ptqrtoken}&webqq_type=10&remember_uin=1&login2qq=1&aid=501004106&u1=http%3A%2F%2Fw.qq.com%2Fproxy.html%3Flogin2qq%3D1%26webqq_type%3D10&ptredirect=0&ptlang=2052&daid=164&from_ui=1&pttype=1&dumy=&fp=loginerroralert&action=0-0-4090&mibao_css=m_webqq&t=1&g=1&js_type=0&js_ver=10231&login_sig=&pt_randsalt=2"
	getUserFriendsUrl     = "http://s.web2.qq.com/api/get_user_friends2"
	getVFWebQQUrl         = "http://s.web2.qq.com/api/getvfwebqq?ptwebqq=&clientid=53999199&psessionid=&t=%s"
	login2Url             = "http://d1.web2.qq.com/channel/login2"
	testLoginUrl          = "http://d1.web2.qq.com/channel/get_online_buddies2"
	poll2Url              = "http://d1.web2.qq.com/channel/poll2"
	refererProxyUrl       = "http://s.web2.qq.com/proxy.html?v=20130916001&callback=1&id=1"
	refererProxy2Url      = "http://s.web2.qq.com/proxy.html?v=20130916001&callback=1&id=2"
	getGroupNameListMask2 = "http://s.web2.qq.com/api/get_group_name_list_mask2"
	getGroupInfoUrl       = "http://s.web2.qq.com/api/get_group_info_ext2"
	getFriendInfoUrl      = "http://s.web2.qq.com/api/get_friend_info2"
	sendMsgUrl            = "http://d1.web2.qq.com/channel/send_buddy_msg2"
	sendQunMsgUrl         = "http://d1.web2.qq.com/channel/send_qun_msg2"
)

type QQClient struct {
	Login              Login
	UserInfo           map[int64]UserInfo
	FriendsMap         FriendsMap
	GroupInfo          map[int64]GroupInfo
	QtQrShowUrl        string
	LoginQrCode        string
	PtQrToken          string
	PtWebQQ            string
	VFWebQQ            string
	PSessionId         string
	Cookies            []*http.Cookie
	QrCodeFilePath     string
	LoginCacheFilePath string
	TimeOut            time.Duration
	HeadMap            map[string]string
	Status             bool
	StatusMessage      string
	LoginCheck         bool
	LoginCheckTimeOut  time.Duration
	MsgID              int
}

type Login struct {
	Url      string
	QQ       int64
	Nickname string
	Status   bool
	Desc     string
}

type QQResponse struct {
	Result  Result `json:"result,omitempty"`
	RetCode int    `json:"retcode,omitempty"`
}

type QQResponseObj struct {
	Result  interface{} `json:"result,omitempty"`
	RetCode int         `json:"retcode,omitempty"`
}

type QQResponsePoll struct {
	Result  []Poll `json:"result,omitempty"`
	RetCode int    `json:"retcode,omitempty"`
}

type Result struct {
	VFWebQQ string `json:"vfwebqq,omitempty"`
	FriendsListInfo
	Group
	Login2
	GroupInfoExt2
	UserInfo
}

type FriendsListInfo struct {
	Categories []Categories `json:"categories,omitempty"`
	Friends    []Friends    `json:"friends,omitempty"`
	Info       []UserInfo   `json:"info,omitempty"`
	MarkNames  []MarkNames  `json:"marknames,omitempty"`
	VipInfo    []VipInfo    `json:"vipinfo,omitempty"`
}

type FriendsMap struct {
	Categories map[int]Categories  `json:"categories,omitempty"`
	Friends    map[int64]Friends   `json:"friends,omitempty"`
	Info       map[int64]UserInfo  `json:"info,omitempty"`
	MarkNames  map[int64]MarkNames `json:"marknames,omitempty"`
	VipInfo    map[int64]VipInfo   `json:"vipinfo,omitempty"`
}

//分组
type Categories struct {
	Index int    `json:"index,omitempty"`
	Name  string `json:"name,omitempty"`
	Sort  int    `json:"sort,omitempty"`
}

//朋友
type Friends struct {
	Categories int   `json:"categories,omitempty"`
	Flag       int   `json:"flag,omitempty"`
	Uin        int64 `json:"uin,omitempty"`
}

//用户详情
type UserInfo struct {
	Allow           int      `json:"allow,omitempty"`
	Birthday        Birthday `json:"birthday,omitempty"`
	Blood           int      `json:"blood,omitempty"`
	City            string   `json:"city,omitempty"`
	ClientType      int      `json:"client_type,omitempty"`
	College         string   `json:"college,omitempty"`
	Constel         int      `json:"constel,omitempty"`
	Country         string   `json:"country,omitempty"`
	Email           string   `json:"email,omitempty"`
	Face            int      `json:"face,omitempty,omitempty"`
	Gender          string   `json:"gender,omitempty"`
	HomePage        string   `json:"homepage,omitempty"`
	Mobile          string   `json:"mobile,omitempty"`
	Nick            string   `json:"nick,omitempty,omitempty"`
	OCCupation      string   `json:"occupation,omitempty"`
	Personal        string   `json:"personal,omitempty"`
	Phone           string   `json:"phone,omitempty"`
	Province        string   `json:"province,omitempty"`
	ShengXiao       int      `json:"shengxiao,omitempty"`
	Stat            int      `json:"stat,omitempty"`
	Uin             int64    `json:"uin,omitempty,omitempty"`
	Flag            int64    `json:"flag,omitempty,omitempty"`
	IsVip           int      `json:"vip_info,omitempty"`
	MarkName        string   `json:"-"`
	VipLevel        int      `json:"-"`
	IsFriend        bool     `json:"-"`
	CategoriesIndex int      `json:"-"`
	CategoriesName  string   `json:"-"`
}

type Birthday struct {
	day   int `json:"day,omitempty"`
	month int `json:"month,omitempty"`
	year  int `json:"year,omitempty"`
}

//标记备注
type MarkNames struct {
	MarkName string `json:"markname,omitempty"`
	Type     int    `json:"type,omitempty"`
	Uin      int64  `json:"uin,omitempty"`
}

//Vip
type VipInfo struct {
	IsVip    int   `json:"is_vip,omitempty"`
	Uin      int64 `json:"u,omitempty"`
	VipLevel int   `json:"vip_level,omitempty"`
}

//群组
type Group struct {
	GMarkList []interface{} `json:"gmarklist,omitempty"`
	GMaskList []interface{} `json:"gmasklist,omitempty"`
	GNameList []GroupInfo   `json:"gnamelist,omitempty"`
}

//群信息
type GroupInfoExt2 struct {
	Cards   []Cards    `json:"cards,omitempty"`
	GInfo   GroupInfo  `json:"ginfo,omitempty"`
	Info    []UserInfo `json:"minfo,omitempty"`
	VipInfo []VipInfo  `json:"vipinfo,omitempty"`
	//stats用处不大，忽略
}

//群名片
type Cards struct {
	Card  string    `json:"card,omitempty"`
	Uin   int64     `json:"muin,omitempty"`
	GInfo GroupInfo `json:"ginfo,omitempty"`
}

//群资料
type GroupInfo struct {
	Class      int                 `json:"class,omitempty"`
	Code       int64               `json:"code,omitempty"`
	CreateTime int64               `json:"createtime,omitempty"`
	Face       int                 `json:"face,omitempty"`
	FingerMemo string              `json:"fingermemo,omitempty"`
	Flag       int64               `json:"flag,omitempty"`
	GId        int64               `json:"gid,omitempty"`
	Level      int                 `json:"level,omitempty"`
	Memo       string              `json:"memo,omitempty"`
	Name       string              `json:"name,omitempty"`
	Option     int                 `json:"option,omitempty"`
	Owner      int64               `json:"owner,omitempty"`
	Cards      map[int64]Cards     `json:"-"`
	UserInfo   map[int64]UserInfo  `json:"-"`
	MarkNames  map[int64]MarkNames `json:"-"`
	//minfo 用处不大，忽略
}

type Login2 struct {
	PSessionId string `json:"psessionid,omitempty"`
	Uin        int64  `json:"uin,omitempty"`
	VFWebQQ    string `json:"vfwebqq,omitempty"`
	UserState  int    `json:"user_state,omitempty"`
}

//心跳
type Poll struct {
	PollType string    `json:"poll_type,omitempty"`
	Value    PollValue `json:"value,omitempty"`
}

type PollValue struct {
	Content   []interface{} `json:"content,omitempty"`
	FromUin   int64         `json:"from_uin,omitempty"`
	GroupCode int64         `json:"group_code,omitempty"`
	MsgID     int64         `json:"msg_id,omitempty"`
	MsgType   int           `json:"msg_type,omitempty"`
	SendUin   int64         `json:"send_uin,omitempty"`
	Time      int64         `json:"time,omitempty"`
	ToUin     int64         `json:"to_uin,omitempty"`
}

func NewQQClient(qq *QQClient) *QQClient {
	qq.UserInfo = make(map[int64]UserInfo)
	qq.GroupInfo = make(map[int64]GroupInfo)
	cookies := []*http.Cookie{}
	cookies = append(cookies, &http.Cookie{Name: "RK", Value: "OfeLBai4FB"})
	cookies = append(cookies, &http.Cookie{Name: "pgv_pvi", Value: "911366144"})
	cookies = append(cookies, &http.Cookie{Name: "pgv_info", Value: "OfeLBai4FB"})
	cookies = append(cookies, &http.Cookie{Name: "pgv_info", Value: "ssid pgv_pvid=1051433466"})
	cookies = append(cookies, &http.Cookie{Name: "ptcz", Value: "ad3bf14f9da2738e09e498bfeb93dd9da7540dea2b7a71acfb97ed4d3da4e277"})
	cookies = append(cookies, &http.Cookie{Name: "qrsig", Value: "hJ9GvNx*oIvLjP5I5dQ19KPa3zwxNI62eALLO*g2JLbKPYsZIRsnbJIxNe74NzQQ"})
	if qq.TimeOut <= 0 {
		qq.TimeOut = time.Duration(30 * time.Second)
	}
	if qq.HeadMap == nil {
		qq.HeadMap = map[string]string{
			"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36",
			"DNT":        "1",
		}
	}
	if qq.Cookies == nil || len(qq.Cookies) <= 0 {
		qq.Cookies = cookies
	}
	qq.MsgID = utils.NewRandom().Number(7)
	qq.Status = true
	return qq
}

func uid4444(str string) string {
	skey := []byte(str)
	e := 0
	for i, n := 0, len(str); n > i; i++ {
		e += (e << 5) + int(skey[i])
	}
	return fmt.Sprint(2147483647 & e)
}

func Hash(x int64, K string) string {
	N := make([]int64, len(K))
	for T := 0; T < len(K); T++ {
		N[T%4] ^= int64(K[T])
	}
	U := "ECOK"
	V := make([]int64, 4)
	V[0] = ((x >> 24) & 255) ^ int64(U[0])
	V[1] = ((x >> 16) & 255) ^ int64(U[1])
	V[2] = ((x >> 8) & 255) ^ int64(U[2])
	V[3] = ((x >> 0) & 255) ^ int64(U[3])
	U2 := make([]int64, 8)
	var str string
	for T := 0; T < 8; T++ {
		if T%2 == 0 {
			U2[T] = N[T>>1]
		} else {
			U2[T] = V[T>>1]
		}
		str += convert.ToString(U2[T])
	}

	N1 := "0123456789ABCDEF"
	var V1 string
	for T := 0; T < len(U2); T++ {
		V1 += string(N1[((U2[T] >> 4) & 15)])
		V1 += string(N1[((U2[T] >> 0) & 15)])
	}
	return V1
}

func Start(qq *QQClient) (newQQ *QQClient, err error) {
	newQQ = qq
	flog := true
	if newQQ.LoginCacheFilePath != "" {
		bts, _ := ioutil.ReadFile(newQQ.LoginCacheFilePath)
		if bts != nil && string(bts) != "" {
			json.Unmarshal(bts, &qq)
		}
		rep, err := newQQ.TestLogin()
		if err == nil && rep.RetCode == 0 {
			flog = false
		}
	}
	if flog {
		newQQ = NewQQClient(&QQClient{LoginCacheFilePath: newQQ.LoginCacheFilePath, QrCodeFilePath: newQQ.QrCodeFilePath})
		_, err = newQQ.PtqrShow()
		if err != nil {
			return
		}
		newQQ, err = newQQ.CheckLogin(nil)
		if err != nil {
			return
		}
	}
	_, err = newQQ.GetFriends()
	if err != nil {
		return
	}
	_, err = newQQ.GetGroup()
	if err != nil {
		return
	}
	if newQQ.QrCodeFilePath != "" {
		var f *os.File
		f, err = os.OpenFile(newQQ.QrCodeFilePath, os.O_RDWR|os.O_CREATE, os.ModePerm)
		defer f.Close()
		f.Write([]byte(convert.ToObjStr(qq)))
		if err != nil {
			return
		}
	}
	return
}

func (qq *QQClient) setQrShowUrl() {
	qq.QtQrShowUrl = fmt.Sprintf(qrShowUrl, utils.NewRandom().Number(4))
}

func (qq *QQClient) getQrLoginUrl() string {
	cookie := qq.GetCookie("qrsig")
	if cookie == nil {
		return ""
	}
	qq.PtQrToken = uid4444(cookie.Value)
	return strings.Replace(qrLoginUrl, "{ptqrtoken}", qq.PtQrToken, -1)
}

func (qq *QQClient) PtqrShow() (*QQClient, error) {
	qq.setQrShowUrl()
	base64Img, cookies, err := qq.DownloadImageCookie(qq.QtQrShowUrl)
	if err != nil {
		return qq, err
	}
	qq.StatusMessage = "扫码登录QQ"
	qq.LoginQrCode = base64Img
	qq.Cookies = cookies
	if qq.QrCodeFilePath != "" {
		//解压
		dist, _ := base64.StdEncoding.DecodeString(strings.Replace(base64Img, "data:image/jpeg;base64,", "", -1))
		//写入新文件
		f, err := os.OpenFile(qq.QrCodeFilePath, os.O_RDWR|os.O_CREATE, os.ModePerm)
		defer f.Close()
		f.Write(dist)
		if err != nil {
			return qq, err
		}
	}
	return qq, nil
}

func (qq *QQClient) CheckLogin(backFunc func(newQQ *QQClient, err error)) (newQQ *QQClient, err error) {
	newQQ = qq
	newQQ.Login.Status = false
	//每次重新调用关闭上次调用程序
	if newQQ.LoginCheck {
		newQQ.LoginCheck = false
		time.Sleep(time.Second)
		newQQ.CheckLogin(backFunc)
	}
	newQQ.LoginCheck = true
	defer func() {
		if backFunc != nil {
			backFunc(newQQ, err)
		}
		newQQ.LoginCheck = false
	}()
	if newQQ.TimeOut <= 0 {
		//10分钟超时
		newQQ.TimeOut = time.Minute * 10
	}
	var timeOut time.Duration
	for {
		if newQQ.LoginCheck {
			println("检测登录中...")
			qrLoginUrl := newQQ.getQrLoginUrl()
			if qrLoginUrl == "" {
				err = errors.New("获取【qrsig】cookie失败")
				return
			}
			var resp *http.Response
			var bodyBytes []byte
			resp, bodyBytes, err = newQQ.Get(newQQ.getQrLoginUrl())
			if err != nil {
				return
			}
			result := string(bodyBytes)
			println("登录检测，返回结果：", result)
			if strings.Contains(result, "二维码未失效") {
				newQQ.Login.Desc = "二维码未失效"
			} else if strings.Contains(result, "二维码认证中") {
				newQQ.Login.Desc = "二维码认证中"
			} else if strings.Contains(result, "二维码已失效") {
				newQQ.Login.Desc = "二维码已失效"
				_, err = newQQ.PtqrShow()
				if err != nil {
					return
				}
			} else if strings.Contains(result, "登录成功") {
				c := newQQ.GetNewCookie(resp.Cookies(), "ptwebqq")
				if c != nil {
					newQQ.PtWebQQ = c.Value
					println("ptwebqq：", c.Value)
				}
				newQQ.UpdateCookie(resp.Cookies())
				newQQ.Login.Nickname = newQQ.GetNickName(result)
				newQQ.Login.Url = newQQ.GetLoginUrl(result)
				q := newQQ.GetQQ(result)
				if q != 0 {
					newQQ.Login.QQ = q
				}

				_, err = newQQ.GetVFWebQQ()
				if err != nil {
					return
				}
				_, err = newQQ.Login2()
				if err != nil {
					return
				}
				newQQ.Login.Desc = "初始化朋友"
				newQQ.GetFriends()
				newQQ.Login.Desc = "初始化群组"
				newQQ.GetGroup()
				newQQ.Login.Desc = "登录成功"
				newQQ.Login.Status = true
				return
			} else {
				newQQ.Login.Desc = fmt.Sprintf("获取二维码扫描状态时出错，ERROR：%s", result)
				err = errors.New(newQQ.Login.Desc)
				return
			}
			timeOut = time.Second * time.Duration(3)
			if timeOut > newQQ.TimeOut {
				err = errors.New("登录校验超时，请重新扫码登录")
				return
			}
			time.Sleep(time.Second * 3)
		} else {
			err = errors.New("已关闭登录校验")
			return
		}
	}
}

func (qq *QQClient) GetVFWebQQ() (qqRes QQResponse, err error) {
	var resp *http.Response
	var bodyBytes []byte
	resp, bodyBytes, err = qq.HttpRequestGet(qq.Login.Url, nil, nil)
	if err != nil {
		return
	}
	if len(resp.Cookies()) > 0 {
		qq.UpdateCookie(resp.Cookies())
	}
	qq.Get(refererProxyUrl)

	head := map[string]string{
		"Host":    "s.web2.qq.com",
		"Origin":  "http://s.web2.qq.com",
		"Referer": refererProxyUrl,
	}
	resp, bodyBytes, err = qq.HttpRequestGet(fmt.Sprintf(getVFWebQQUrl, GetTimeUinx()), head, nil)
	defer resp.Body.Close()
	if err != nil {
		return
	}
	err = json.Unmarshal(bodyBytes, &qqRes)
	if err != nil {
		return
	}
	if qqRes.RetCode == 0 {
		qq.VFWebQQ = qqRes.Result.VFWebQQ
	}
	return
}

func (qq *QQClient) Login2() (qqRes QQResponse, err error) {
	var bodyBytes []byte
	r := `r={"ptwebqq":"","clientid":53999199,"psessionid":"","status":"online"}`
	head := map[string]string{
		"Host":    "d1.web2.qq.com",
		"Origin":  "http://d1.web2.qq.com",
		"Referer": "http://d1.web2.qq.com/proxy.html?v=20151105001&callback=1&id=2",
	}
	_, bodyBytes, err = qq.HttpRequestPost(login2Url, head, r)
	if err != nil {
		return
	}
	err = json.Unmarshal(bodyBytes, &qqRes)
	if err != nil {
		return
	}
	if qqRes.RetCode == 0 {
		qq.PSessionId = qqRes.Result.PSessionId
	}
	return
}

func (qq *QQClient) TestLogin() (qqRes QQResponseObj, err error) {
	u, _ := url.Parse(".qq.com")
	jar := new(smartWechat.Jar)
	jar.SetCookies(u, qq.Cookies)
	var bodyBytes []byte

	head := map[string]string{
		"Cookie":  qq.GetCookieStr(),
		"Host":    "d1.web2.qq.com",
		"Referer": "http://d1.web2.qq.com/proxy.html?v=20151105001&callback=1&id=2",
	}
	pars := map[string]interface{}{
		"vfwebqq":    qq.VFWebQQ,
		"clientid":   "53999199",
		"psessionid": qq.PSessionId,
		"t":          GetTimeUinx(),
	}
	_, bodyBytes, err = qq.HttpRequestGet(testLoginUrl, head, pars)
	if err != nil {
		println(err.Error())
		return
	}
	err = json.Unmarshal(bodyBytes, &qqRes)
	if err != nil {
		return
	}
	return
}

func (qq *QQClient) GetFriends() (qqRes QQResponse, err error) {
	u, _ := url.Parse(".qq.com")
	jar := new(smartWechat.Jar)
	jar.SetCookies(u, qq.Cookies)
	var bodyBytes []byte
	head := map[string]string{
		"Host":    "s.web2.qq.com",
		"Origin":  "http://s.web2.qq.com",
		"Referer": refererProxyUrl,
	}
	r := `r={"vfwebqq":"` + qq.VFWebQQ + `","hash":"` + Hash(qq.Login.QQ, qq.PtWebQQ) + `"}`
	_, bodyBytes, err = qq.HttpRequestPost(getUserFriendsUrl, head, r)
	if err != nil {
		return
	}
	err = json.Unmarshal(bodyBytes, &qqRes)
	if err != nil {
		return
	}
	if qqRes.RetCode == 0 {
		qq.FriendsMap = qq.ConvertFriendsListToMap(&qqRes.Result.FriendsListInfo)
	}
	return
}

func (qq *QQClient) GetFriendInfo(uin int64) (qqRes QQResponse, err error) {
	var bodyBytes []byte
	head := map[string]string{
		"Host":    "s.web2.qq.com",
		"Referer": refererProxyUrl,
	}
	parMap := map[string]interface{}{
		"tuin":       uin,
		"vfwebqq":    qq.VFWebQQ,
		"psessionid": qq.PSessionId,
		"clientid":   "53999199",
		"t":          GetTimeUinx(),
	}
	_, bodyBytes, err = qq.HttpRequestGet(getFriendInfoUrl, head, parMap)
	if err != nil {
		return
	}
	err = json.Unmarshal(bodyBytes, &qqRes)
	if err != nil {
		return
	}
	if qqRes.RetCode == 0 {

	}
	return
}

func (qq *QQClient) GetGroup() (qqRes QQResponse, err error) {
	var bodyBytes []byte
	r := `r={"vfwebqq":"` + qq.VFWebQQ + `","hash":"` + Hash(qq.Login.QQ, qq.PtWebQQ) + `"}`
	head := map[string]string{
		"Host":    "s.web2.qq.com",
		"Origin":  "http://s.web2.qq.com",
		"Referer": refererProxyUrl,
	}
	_, bodyBytes, err = qq.HttpRequestPost(getGroupNameListMask2, head, r)
	if err != nil {
		return
	}
	err = json.Unmarshal(bodyBytes, &qqRes)
	if err != nil {
		return
	}
	if qqRes.RetCode == 0 {
		qq.UpdateGroupToMap(&qqRes.Result.Group)
	}
	return
}

func (qq *QQClient) GetGroupInfo(groupCode int64) (qqRes QQResponse, err error) {
	var bodyBytes []byte
	head := map[string]string{
		"Host":    "s.web2.qq.com",
		"Referer": refererProxyUrl,
	}
	parMap := map[string]interface{}{
		"gcode":   groupCode,
		"vfwebqq": qq.VFWebQQ,
		"t":       GetTimeUinx(),
	}
	_, bodyBytes, err = qq.HttpRequestGet(getGroupInfoUrl, head, parMap)
	if err != nil {
		return
	}
	err = json.Unmarshal(bodyBytes, &qqRes)
	if err != nil {
		return
	}
	if qqRes.RetCode == 0 {
		qq.UpdateGroupInfo(&qqRes.Result.GroupInfoExt2)
	}
	return
}

func (qq *QQClient) Poll2(poll2 func(qq *QQClient, result QQResponsePoll)) (res QQResponsePoll, err error) {
	println("Poll2 -->  qq.Status && qq.Login.Status", qq.Status, qq.Login.Status)
	var flog bool
	for qq.Status && qq.Login.Status {
		res, flog, err = qq.Poll()
		if err != nil {
			println(err.Error())
			if flog {
				//post请求异常，继续下一次
				continue
			}
		}
		if res.RetCode == 0 {
			poll2(qq, res)
			time.Sleep(time.Second * 2)
		} else {
			//有可能心跳连接失败需要重新登录
			return
		}
	}
	return
}

func (qq *QQClient) Poll() (res QQResponsePoll, flog bool, err error) {
	var bodyBytes []byte
	r := `r={"ptwebqq":"` + qq.PtWebQQ + `","clientid":53999199,"psessionid":"` + qq.PSessionId + `","key":""}`
	head := map[string]string{
		"Host":    "d1.web2.qq.com",
		"Origin":  "http://d1.web2.qq.com",
		"Referer": "http://d1.web2.qq.com/proxy.html?v=20151105001&callback=1&id=2",
	}
	println(r)
	_, bodyBytes, err = qq.HttpRequestPost(poll2Url, head, r)
	if err != nil {
		println(err.Error())
		flog = true
		return
	}
	err = json.Unmarshal(bodyBytes, &res)
	if err != nil {
		println(err.Error())
		return
	}
	return
}

func (qq *QQClient) SendMsg(toUin int64, content string) (res QQResponse, err error) {
	qq.MsgID++
	var bodyBytes []byte
	r := `r={"to":` + convert.ToString(toUin) + `,"content":"[\"` + content + `\",[\"font\",{\"name\":\"宋体\",\"size\":10,\"style\":[0,0,0],\"color\":\"000000\"}]]","face":96,"clientid":53999199,"msg_id":` + convert.ToString(qq.MsgID) + `,"psessionid":"` + qq.PSessionId + `"}`
	println(r)
	head := map[string]string{
		"Host":    "d1.web2.qq.com",
		"Origin":  "http://d1.web2.qq.com",
		"Referer": "http://d1.web2.qq.com/cfproxy.html?v=20151105001&callback=1",
	}
	_, bodyBytes, err = qq.HttpRequestPost(sendMsgUrl, head, r)
	if err != nil {
		return
	}
	err = json.Unmarshal(bodyBytes, &res)
	if err != nil {
		return
	}
	return
}

func (qq *QQClient) SendQunMsg(toGroupUin int64, content string) (res QQResponse, err error) {
	qq.MsgID++
	var bodyBytes []byte
	r := `r={"group_uin":` + convert.ToString(toGroupUin) + `,"content":"[\"` + content + `\",[\"font\",{\"name\":\"宋体\",\"size\":10,\"style\":[0,0,0],\"color\":\"000000\"}]]","face":594,"clientid":53999199,"msg_id":` + convert.ToString(qq.MsgID) + `,"psessionid":"` + qq.PSessionId + `"}`
	println(r)
	head := map[string]string{
		"Host":    "d1.web2.qq.com",
		"Origin":  "http://d1.web2.qq.com",
		"Referer": "http://d1.web2.qq.com/cfproxy.html?v=20151105001&callback=1",
	}
	_, bodyBytes, err = qq.HttpRequestPost(sendQunMsgUrl, head, r)
	if err != nil {
		return
	}
	err = json.Unmarshal(bodyBytes, &res)
	if err != nil {
		return
	}
	return
}

func (qq *QQClient) ConvertFriendsListToMap(userInfo *FriendsListInfo) FriendsMap {
	if userInfo != nil && len(userInfo.Friends) > 0 {
		friendsMap := FriendsMap{}

		if len(userInfo.Categories) > 0 {
			categoriesMap := map[int]Categories{}
			for _, v := range userInfo.Categories {
				categoriesMap[v.Index] = v
			}
			friendsMap.Categories = categoriesMap
		}

		if len(userInfo.Friends) > 0 {
			friendMap := map[int64]Friends{}
			for _, v := range userInfo.Friends {
				friendMap[v.Uin] = v
			}
			friendsMap.Friends = friendMap
		}

		if len(userInfo.Info) > 0 {
			infoMap := map[int64]UserInfo{}
			for _, v := range userInfo.Info {
				infoMap[v.Uin] = v
			}
			friendsMap.Info = infoMap
		}

		if len(userInfo.MarkNames) > 0 {
			markNamesMap := map[int64]MarkNames{}
			for _, v := range userInfo.MarkNames {
				markNamesMap[v.Uin] = v
			}
			friendsMap.MarkNames = markNamesMap
		}
		if len(userInfo.VipInfo) > 0 {
			vipMap := map[int64]VipInfo{}
			for _, v := range userInfo.VipInfo {
				vipMap[v.Uin] = v
			}
			friendsMap.VipInfo = vipMap
		}
		//添加到用户中【添加关系为好友】
		info := UserInfo{}
		for _, v := range friendsMap.Friends {
			tu := qq.UserInfo[v.Uin]
			if tu.Uin <= 0 {
				info = UserInfo{}
			} else {
				info = tu
			}
			info.Uin = v.Uin
			info.IsFriend = true
			if friendsMap.VipInfo != nil && len(userInfo.VipInfo) > 0 {
				info.IsVip = friendsMap.VipInfo[v.Uin].IsVip
				info.VipLevel = friendsMap.VipInfo[v.Uin].VipLevel
			}
			if friendsMap.Info != nil && len(userInfo.Info) > 0 {
				info.Face = friendsMap.Info[v.Uin].Face
				info.Nick = friendsMap.Info[v.Uin].Nick
				info.Flag = friendsMap.Info[v.Uin].Flag
			}
			if friendsMap.MarkNames != nil && len(userInfo.MarkNames) > 0 {
				info.MarkName = friendsMap.MarkNames[v.Uin].MarkName
			}
			if friendsMap.Categories != nil && len(userInfo.Categories) > 0 {
				c := friendsMap.Categories[friendsMap.Friends[v.Uin].Categories]
				info.CategoriesIndex = c.Index
				info.CategoriesName = c.Name
			}
			qq.UserInfo[v.Uin] = info
			println(convert.ToObjStr(info))
		}
		return friendsMap
	}
	return FriendsMap{}
}

func (qq *QQClient) UpdateGroupToMap(group *Group) {
	if group != nil && len(group.GNameList) > 0 {
		g := GroupInfo{}
		groupMap := map[int64]GroupInfo{}
		for _, v := range group.GNameList {
			tg := qq.GroupInfo[v.Code]
			if tg.Code <= 0 {
				g = GroupInfo{}
			} else {
				g = tg
			}
			g.Code = v.Code
			g.Flag = v.Flag
			g.GId = v.GId
			g.Name = v.Name
			groupMap[v.GId] = g
			qq.GroupInfo[v.GId] = g
		}
	}
}

func (qq *QQClient) UpdateGroupInfo(ge *GroupInfoExt2) {
	if ge != nil {
		g := GroupInfo{}
		g.Class = ge.GInfo.Class
		g.Code = ge.GInfo.Code
		g.CreateTime = ge.GInfo.CreateTime
		g.Face = ge.GInfo.Face
		g.FingerMemo = ge.GInfo.FingerMemo
		g.Flag = ge.GInfo.Flag
		g.GId = ge.GInfo.GId
		g.Level = ge.GInfo.Level
		g.Memo = ge.GInfo.Memo
		g.Name = ge.GInfo.Name
		g.Option = ge.GInfo.Option
		g.Owner = ge.GInfo.Owner
		g.UserInfo = qq.CovertUserInfoArrayToMap(ge.Info)
		g.MarkNames = ge.GInfo.MarkNames
		qq.GroupInfo[ge.GInfo.GId] = g
		if g.UserInfo != nil {
			for _, v := range g.UserInfo {
				qq.UpdateUserInfo(v)
			}
		}
	}
}

func (qq *QQClient) CovertUserInfoArrayToMap(info []UserInfo) map[int64]UserInfo {
	if len(info) <= 0 {
		return nil
	}
	uMap := map[int64]UserInfo{}
	for i := 0; i < len(info); i++ {
		uMap[info[i].Uin] = info[i]
	}
	return uMap
}

func (qq *QQClient) UpdateUserInfo(info UserInfo) {
	println("UpdateUserInfo：", convert.ToObjStr(info))
	qq.UserInfo[info.Uin] = info
}

func (qq *QQClient) GetNickName(result string) string {
	if regName := regexp.MustCompile(`'登录成功！', '(.*)\'\)`).FindAllStringSubmatch(result, -1); len(regName) == 1 {
		return regName[0][1]
	}
	return ""
}

func (qq *QQClient) GetLoginUrl(result string) string {
	if regUrl := regexp.MustCompile(`ptuiCB\(\'0\',\'0\',\'([^\']+)\'`).FindAllStringSubmatch(result, -1); len(regUrl) == 1 {
		return regUrl[0][1]
	}
	return ""
}

func (qq *QQClient) GetQQ(result string) int64 {
	if regName := regexp.MustCompile(`&uin=(.*)&service=`).FindAllStringSubmatch(result, -1); len(regName) == 1 {
		return convert.MustInt64(regName[0][1])
	}
	return 0
}

func (qq *QQClient) UpdateCookie(cookies []*http.Cookie) {
	for _, cookie := range cookies {
		qq.Cookies = append(qq.Cookies, cookie)
	}
}

func (qq *QQClient) GetCookie(name string) *http.Cookie {
	for _, cookie := range qq.Cookies {
		if cookie.Name == name {
			return cookie
		}
	}
	return nil
}

func (qq *QQClient) GetCookieStr() string {
	var str string
	for _, cookie := range qq.Cookies {
		if cookie.Value != "" {
			str += cookie.Name + "=" + cookie.Value + ";"
		}
	}
	return str
}

func (qq *QQClient) GetNewCookie(cookies []*http.Cookie, name string) *http.Cookie {
	for _, cookie := range cookies {
		if cookie.Name == name {
			return cookie
		}
	}
	return nil
}

func (qq *QQClient) HttpRequestPost(url string, head map[string]string, body interface{}) (resp *http.Response, bodyBytes []byte, err error) {
	if head == nil {
		head = map[string]string{}
	}
	head["Content-Type"] = "application/x-www-form-urlencoded"
	head["Cookie"] = qq.GetCookieStr()
	return qq.HttpRequest("POST", url, head, body)
}
func (qq *QQClient) HttpRequestGet(url string, head map[string]string, parMap map[string]interface{}) (resp *http.Response, bodyBytes []byte, err error) {
	if head == nil {
		head = map[string]string{}
	}
	head["Content-Type"] = "utf-8"
	head["Cookie"] = qq.GetCookieStr()
	return qq.HttpRequest("GET", qq.GetURLParams(url, parMap), head, nil)
}

func (qq *QQClient) Get(url string) (resp *http.Response, bodyBytes []byte, err error) {
	head := map[string]string{}
	head["Content-Type"] = "utf-8"
	head["Cookie"] = qq.GetCookieStr()
	return qq.HttpRequest("GET", url, head, nil)
}

func (qq *QQClient) HttpRequest(method string, url string, head map[string]string, body interface{}) (resp *http.Response, bodyBytes []byte, err error) {
	client := &http.Client{Timeout: qq.TimeOut, CheckRedirect: RedirectPolicyFunc}
	var bodyStr string
	switch val := body.(type) {
	case string:
		bodyStr = string(val)
	default:
		var jsonBytes []byte
		jsonBytes, err = json.Marshal(val)
		if err != nil {
			return
		}
		bodyStr = string(jsonBytes)
	}
	req, err := http.NewRequest(method, url, strings.NewReader(bodyStr))
	if err != nil {
		return
	}
	for k, v := range qq.HeadMap {
		req.Header.Set(k, v)
	}
	if head != nil {
		for k, v := range head {
			req.Header.Set(k, v)
		}
	}
	resp, err = client.Do(req)
	if err != nil {
		return
	}
	bodyBytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	println(url, "【返回结果】", string(bodyBytes))
	return
}

func (qq *QQClient) GetURLParams(thisUrl string, m map[string]interface{}) string {
	if m == nil {
		return thisUrl
	}
	var result = ""
	if !strings.Contains(thisUrl, "?") {
		result = "?"
	}
	for key, value := range m {
		if key != "" && value != "" {
			result += fmt.Sprintf("%s=%s&", key, url.QueryEscape(convert.ToString(value)))
		}
	}
	if result == "" {
		return thisUrl
	}
	return thisUrl + result[:len(result)-1]
}

func (qq *QQClient) SetJar(domain string, cookie []*http.Cookie) (jar *smartWechat.Jar) {
	u, _ := url.Parse(domain)
	jar.SetCookies(u, cookie)
	return
}

/* 下载URL指向的JPG base64*/
func (qq *QQClient) DownloadImageCookie(url string) (string, []*http.Cookie, error) {
	println("DownloadImageCookie  url：", url)
	resp, err := http.Get(url)
	if err != nil {
		return "", nil, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", nil, err
	}
	base64Img := "data:image/jpeg;base64," + base64.StdEncoding.EncodeToString(b)
	return base64Img, resp.Cookies(), err
}

func GetTimeUinx() string {
	return fmt.Sprintf("%d", time.Now().UnixNano()/1000000)
}

func RedirectPolicyFunc(req *http.Request, via []*http.Request) error {
	return http.ErrUseLastResponse
}

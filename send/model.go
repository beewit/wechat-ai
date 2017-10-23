package send

import (
	"net/url"
	"net/http"
	"fmt"
	"encoding/xml"
	"time"
	"math/rand"
)

/*
 * <error>
 *   <ret>0</ret>
 *   <message></message>
 *   <skey>@crypt_3aaab8d5_aa9febb1c57122a4569c1b1dc4772eac</skey>
 *   <wxsid>vjqCszEkQQw9jep1</wxsid>
 *   <wxuin>154158775</wxuin>
 *   <pass_ticket>wbFO7Vqg%2BpADuIcrQPDM1e0KjmNvgsH8jYAEoT0FtSY%3D</pass_ticket>
 *   <isgrayscale>1</isgrayscale>
 * </error>
 */
type LoginCallbackXMLResult struct {
	XMLName     xml.Name `xml:"error"` /* 根节点定义 */
	Ret         string   `xml:"ret"`
	Message     string   `xml:"message"`
	SKey        string   `xml:"skey"`
	WXSid       string   `xml:"wxsid"`
	WXUin       string   `xml:"wxuin"`
	PassTicket  string   `xml:"pass_ticket"`
	IsGrayscale string   `xml:"isgrayscale"`
}

type BaseRequest struct {
	Uin      string `json:"Uin"`
	Sid      string `json:"Sid"`
	SKey     string `json:"Skey"`
	DeviceID string `json:"DeviceID"`
}

/* 微信初始化时返回的大JSON，选择性地提取一些关键数据 */
type InitInfo struct {
	User     User             `json:"User"`
	SyncKeys SyncKeysJsonData `json:"SyncKey"`
}

/* 微信获取所有联系人列表时返回的大JSON */
type ContactList struct {
	MemberCount int    `json:"MemberCount"`
	MemberList  []User `json:"MemberList"`
}

/* 微信通用User结构，可根据需要扩展 */
type User struct {
	Uin        int64  `json:"Uin"`
	UserName   string `json:"UserName"`
	NickName   string `json:"NickName"`
	RemarkName string `json:"RemarkName"`
	Sex        int8   `json:"Sex"`
	Province   string `json:"Province"`
	City       string `json:"City"`
}

type SyncKeysJsonData struct {
	Count    int       `json:"Count"`
	SyncKeys []SyncKey `json:"List"`
}

type SyncKey struct {
	Key int64 `json:"Key"`
	Val int64 `json:"Val"`
}

/* 设计一个构造成字符串的结构体方法 */
func (sks SyncKeysJsonData) ToString() string {
	resultStr := ""

	for i := 0; i < sks.Count; i++ {
		resultStr = resultStr + fmt.Sprintf("%d_%d|", sks.SyncKeys[i].Key, sks.SyncKeys[i].Val)
	}

	return resultStr[:len(resultStr)-1]
}

/* 微信消息对象 */
type WxRecvMsges struct {
	MsgCount int              `json:"AddMsgCount"`
	MsgList  []WxRecvMsg      `json:"AddMsgList"`
	SyncKeys SyncKeysJsonData `json:"SyncKey"`
}

/* 微信接受消息对象元素 */
type WxRecvMsg struct {
	MsgId        string `json:"MsgId"`
	FromUserName string `json:"FromUserName"`
	ToUserName   string `json:"ToUserName"`
	MsgType      int    `json:"MsgType"`
	Content      string `json:"Content"`
	CreateTime   int64  `json:"CreateTime"`
}

/**
 * "Type":1,
 * "Content":"1",
 * "FromUserName":"@9499e6e8dfd2c1020ecb6cc727982bef",
 * "ToUserName":"@9499e6e8dfd2c1020ecb6cc727982bef",
 * "LocalID":"15046739462870976",
 * "ClientMsgId":"15046739462870976"
 * 微信发送消息对象元素
 */
type WxSendMsg struct {
	Type         int    `json:"Type"`
	Content      string `json:"Content"`
	FromUserName string `json:"FromUserName"`
	ToUserName   string `json:"ToUserName"`
	LocalID      string `json:"LocalID"`
	ClientMsgId  string `json:"ClientMsgId"`
}


/* 获取联系人列表时需要带入Cookie信息，实现CookieJar接口 */
type Jar struct {
	cookies []*http.Cookie
}

func (jar *Jar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	jar.cookies = cookies
}
func (jar *Jar) Cookies(u *url.URL) []*http.Cookie {
	return jar.cookies
}

type LoginMap struct {
	PassTicket  string
	BaseRequest BaseRequest /* 将涉及登陆有关的验证数据封装成对象 */

	SelfNickName string
	SelfUserName string

	SyncKeys   SyncKeysJsonData /* 同步消息时需要验证的Keys */
	SyncKeyStr string           /* Keys组装成的字符串 */

	Cookies []*http.Cookie /* 微信相关API需要用到的Cookies */
}


var (
	uuid       string
	err        error
	loginMap   LoginMap
	contactMap map[string]User
	groupMap   map[string][]User /* 关键字为key的，群组数组 */
)

/**
 * 这个文件用来维护焦点微信群的关键字的正则表达式
 * TODO:后期使用数据库维护
 */

type FocusGroupKeywords struct {
	FatherName   string
	Description  string
	ExampleStr   string
	ChildrenName []string
}

var focusGroupKeywords []FocusGroupKeywords

func init() {
	focusGroupKeywords = []FocusGroupKeywords{}
	/*

	focusGroupKeywords = append(focusGroupKeywords, FocusGroupKeywords{
		FatherName:  "1-0|编程",
		Description: "这些群组可以快速帮你结识相关领域志同道合的爱好者，程序猿并不是一个人在战斗！",
		ExampleStr:  `"Golang"或者其索引号"1-1"，我会邀请您加入专业探讨Golang技术的微信群。`,
		ChildrenName: []string{"1-1|Golang", "1-2|Java", "1-3|Python", "1-4|Nodejs", "1-5|Qt",
			"1-6|SQL|数据库",
			"1-7|Angular", "1-8|Vue", "1-9|React", "1-10|jQuery",
			"1-11|Linux", "1-12|Android", "1-13|IOS", "1-14|全栈"}})

	*/
}
func GetFocusGroupKeywords() []FocusGroupKeywords {
	return focusGroupKeywords
}

/* 获取所有关键词 */
func GetFocusGroupKeywordChildren() []string {
	keywords := []string{}

	for _, v := range focusGroupKeywords {
		for _, str := range v.ChildrenName {
			keywords = append(keywords, str)
		}
	}

	return keywords
}

/* 得到关键词父级目录 */
func GetFatherKeywordsStr() string {
	keywordsStr := ""

	for _, v := range focusGroupKeywords {
		keywordsStr += v.FatherName + "\n"
	}

	return keywordsStr
}

/* 返回分组介绍和所有分组关键词的组装 */
func GetChildKeywordsInfo(fatherName string) (string, string, string) {
	keywordsStr := ""
	var index int

	for index = 0; index < len(focusGroupKeywords); index++ {
		if focusGroupKeywords[index].FatherName == fatherName {
			break
		}
	}

	if index == len(focusGroupKeywords) {
		return "", "", ""
	}

	for _, str := range focusGroupKeywords[index].ChildrenName {
		keywordsStr += str + "\n"
	}

	return focusGroupKeywords[index].Description, focusGroupKeywords[index].ExampleStr, keywordsStr
}




/**
 * 有序(或者无序)地从一个map中按照index的顺序构造URL中的params
 * 加上有序的目的是为了防止有些环境下需要params根据key的ASC大小排序后进行签名加密
 */
func GetURLParams(values ...interface{}) string {
	var result = "?"
	if len(values) == 1 {
		maap := values[0].(map[string]string)
		for key, value := range maap {
			if key != "" && value != "" {
				result += fmt.Sprintf("%s=%s&", key, url.QueryEscape(value))
			}
		}
	} else if len(values) == 2 {
		index := values[1].([]string)
		maap := values[0].(map[string]string)
		for _, key := range index {
			if key != "" && maap[key] != "" {
				result += fmt.Sprintf("%s=%s&", key, url.QueryEscape(maap[key]))
			}
		}
	}

	return result[:len(result)-1]
}

/**
 *  生成随机字符串
 *  index：取随机序列的前index个
 *  0-9:10
 *  0-9a-z:10+24
 *  0-9a-zA-Z:10+24+24
 *  length：需要生成随机字符串的长度
 */
func GetRandomString(index int, length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(index)])
	}
	return string(result)
}
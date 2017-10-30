package enum

const (
	APPID     = "appid"
	FUN       = "fun"
	Lang      = "lang"
	LangValue = "zh_CN"
	TimeStamp = "_"
	UUID      = "uuid"
	R         = "r"

	/* 以下信息会存储在loginMap中 */
	Ret          = "ret"
	Message      = "message"
	SKey         = "skey"
	WXSid        = "wxsid"
	WXUin        = "wxuin"
	PassTicket   = "pass_ticket"
	IsGrayscale  = "isgrayscale"
	DeviceID     = "DeviceID"
	SelfUserName = "UserName"
	SelfNickName = "NickName"
	SyncKeyStr   = "synckeystr"

	Sid         = "sid"
	Uin         = "uin"
	DeviceId    = "deviceid"
	SyncKey     = "synckey"
	BaseRequest = "BaseRequest"

	WECHAT_RESPONE_NORMAL           = 0    //正常
	WECHAT_RESPONE_LOGIN_ERR        = -14  // ticket 错误
	WECHAT_RESPONE_LOGIN_OUT        = 1100 //退出未登录
	WECHAT_RESPONE_LOGIN_OTHERWHERE = 1101 //其它地方登陆
	WECHAT_RESPONE_MOBILE_LOGIN_OUT = 1102 //移动端退出
	WECHAT_RESPONE_FREQUENTLY       = 1205 //操作频繁
)

var (
	uuidParaEnum = map[string]string{
		APPID:     "wx782c26e4c19acffb",
		FUN:       "new",
		Lang:      LangValue,
		TimeStamp: ""}

	loginParaEnum = map[string]string{
		"loginicon": "true",
		"tip":       "0",
		UUID:        "",
		R:           "",
		TimeStamp:   ""}

	initParaEnum = map[string]string{
		R:          "",
		Lang:       LangValue,
		PassTicket: ""}
)

func GetUUIDParaEnum() map[string]string {
	return uuidParaEnum
}

func GetLoginParaEnum() map[string]string {
	return loginParaEnum
}

func GetInitParaEnum() map[string]string {
	return initParaEnum
}

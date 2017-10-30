package smartqq

import (
	"time"
	"net/http"
	"fmt"
	"github.com/beewit/wechat-ai/send"
	"net/url"
	"github.com/beewit/beekit/utils"
	"strings"
	"github.com/pkg/errors"
	"regexp"
	"io/ioutil"
	"os"
	"encoding/base64"
	"github.com/beewit/beekit/utils/convert"
)

var (
	qrShowUrl  = "https://ssl.ptlogin2.qq.com/ptqrshow?appid=501004106&e=2&l=M&s=3&d=72&v=4&t=0.1509281804550%v&daid=164&pt_3rd_aid=0"
	qrLoginUrl = "https://ssl.ptlogin2.qq.com/ptqrlogin?u1=http%3A%2F%2Fw.qq.com%2Fproxy.html&ptqrtoken={ptqrtoken}&ptredirect=0&h=1&t=1&g=1&from_ui=1" +
		"&ptlang=2052&action=0-0-1509278827677&js_ver=10231&js_type=1&login_sig=nN7dHna67PR7COUx5nsLSBQAsZOsDWdM*QWgBVE9htplwGgYqyPdDNDyPYG2HZwr&" +
		"pt_uistyle=40&aid=501004106&daid=164&mibao_css=m_webqq&"
)

type QQClient struct {
	Login          Login
	QtQrShowUrl    string
	LoginQrCode    string
	PtQrToken      string
	Cookies        []*http.Cookie
	QrCodeFilePath string
}

type Login struct {
	Url      string
	QQ       string
	Nickname string
	Status   bool
	Desc     string
}

func NewQQClient() *QQClient {
	qqClient := new(QQClient)
	return qqClient
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

	println(convert.ToObjStr(U2))

	N1 := "0123456789ABCDEF"
	var V1 uint8
	for T := 0; T < len(U2); T++ {
		V1 += N1[((U2[T]>>4)&15)]
		println(N1[(U2[T]>>4)&15])
		V1 += N1[((U2[T] >> 0) & 15)]
		println(V1)
	}
	println(string(V1))
	return string(V1)
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
	base64Img, cookies, err := send.DownloadImageCookie(qq.QtQrShowUrl)
	if err != nil {
		return qq, err
	}
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

func (qq *QQClient) CheckLogin() (newQQ *QQClient, err error) {
	newQQ = qq
	newQQ.Login.Status = false
	for {
		println("检测登录中...")
		u, _ := url.Parse("ptlogin2.qq.com")
		jar := new(send.Jar)
		jar.SetCookies(u, newQQ.Cookies)
		qrLoginUrl := newQQ.getQrLoginUrl()
		if qrLoginUrl == "" {
			err = errors.New("获取【qrsig】cookie失败")
			return
		}
		var resp *http.Response
		var bodyBytes []byte
		resp, err = newQQ.Get(jar, newQQ.getQrLoginUrl())
		if err != nil {
			return
		}
		newQQ.UpdateCookie(resp.Cookies())
		bodyBytes, err = ioutil.ReadAll(resp.Body)
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
			newQQ.Login.Status = true
			newQQ.Login.Desc = "登录成功"
			newQQ.Login.Nickname = newQQ.GetNickName(result)
			newQQ.Login.Url = newQQ.GetLoginUrl(result)
			qqCookie := qq.GetCookie("superuin")
			if qqCookie != nil {
				newQQ.Login.QQ = strings.Replace(qqCookie.Value, "o0", "", 1)
			}
			return
		} else {
			newQQ.Login.Desc = fmt.Sprintf("获取二维码扫描状态时出错，ERROR：%s", result)
			err = errors.New(newQQ.Login.Desc)
			return
		}
		time.Sleep(time.Second * 3)
	}
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

func (qq *QQClient) UpdateCookie(cookies []*http.Cookie) {
	for _, cookie := range cookies {
		if qq.GetCookie(cookie.Name) == nil {
			qq.Cookies = append(qq.Cookies, cookie)
		}
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

func (qq *QQClient) Get(jar *send.Jar, getUrl string) (resp *http.Response, err error) {
	client := &http.Client{Jar: jar}
	resp, err = client.Get(getUrl)
	return
}

func (qq *QQClient) SetJar(domain string, cookie []*http.Cookie) (jar *send.Jar) {
	u, _ := url.Parse(domain)
	jar.SetCookies(u, cookie)
	return
}

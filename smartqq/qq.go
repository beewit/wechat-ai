package smartqq

import (
	"github.com/beewit/wechat-ai/send"
	"strings"
	"os"
	"encoding/base64"
	"fmt"
	"net/url"
	"net/http"
	"github.com/beewit/beekit/utils/convert"
	"io/ioutil"
	"github.com/robertkrimen/otto"
	"github.com/beewit/beekit/utils"
	"time"
)

func main() {
	//https://ssl.ptlogin2.qq.com/ptqrshow?appid=501004106&e=2&l=M&s=3&d=72&v=4&t=0.6267585117576875&daid=164&pt_3rd_aid=0
	//https://ssl.ptlogin2.qq.com/ptqrshow?appid=501004106&e=2&l=M&s=3&d=72&v=4&t=0.15092818045507936&daid=164&pt_3rd_aid=0
	qrUrl := fmt.Sprintf("https://ssl.ptlogin2.qq.com/ptqrshow?appid=501004106&e=2&l=M&s=3&d=72&v=4&t=0.1509281804550%v&daid=164&pt_3rd_aid=0", utils.NewRandom().Number(4))
	println("请求Url:", qrUrl)
	base64Img, cookides, err := send.DownloadImageCookie(qrUrl)

	if err != nil {
		println("DownloadImage Error：" + err.Error())
		return
	}
	//解压
	dist, _ := base64.StdEncoding.DecodeString(strings.Replace(base64Img, "data:image/jpeg;base64,", "", -1))
	//写入新文件
	f, err := os.OpenFile("output.jpg", os.O_RDWR|os.O_CREATE, os.ModePerm)
	defer f.Close()
	f.Write(dist)
	if err != nil {
		println("output.jpg Error：" + err.Error())
		return
	}

	var qrsig string
	for i := 0; i < len(cookides); i++ {
		if cookides[i].Name == "qrsig" {
			qrsig = cookides[i].Value
		}
	}
	vm := otto.New()
	vm.Run(`
			function uid4444(t) {
			for (var e = 0, i = 0, n = t.length; n > i; ++i)
			e += (e << 5) + t.charCodeAt(i);
			return 2147483647 & e
		}
	`)
	value, _ := vm.Call("uid4444", nil, qrsig)
	val, _ := value.ToString()
	for {
		checkLogin(cookides, val)
		time.Sleep(time.Second * 2)
	}
}



func checkLogin(cookides []*http.Cookie, val string) (err error, cookie []*http.Cookie) {
	u, _ := url.Parse("ptlogin2.qq.com")
	jar := new(send.Jar)
	jar.SetCookies(u, cookides)
	println(convert.ToObjStr(cookides))
	println("解码后的值：", val)
	client := &http.Client{Jar: jar}
	url := "https://ssl.ptlogin2.qq.com/ptqrlogin?u1=http%3A%2F%2Fw.qq.com%2Fproxy.html&ptqrtoken=" + val + "&ptredirect=0&h=1&t=1&g=1&from_ui=1&ptlang=2052&action=0-0-1509278827677&js_ver=10231&js_type=1&login_sig=nN7dHna67PR7COUx5nsLSBQAsZOsDWdM*QWgBVE9htplwGgYqyPdDNDyPYG2HZwr&pt_uistyle=40&aid=501004106&daid=164&mibao_css=m_webqq&"
	println("检测登录", url)
	var resp *http.Response
	var bodyBytes []byte

	resp, err = client.Get(url)
	if err != nil {
		println(err.Error())
		return
	} else {
		bodyBytes, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			println(err.Error())
			return
		} else {
			println("bodyBytes", string(bodyBytes))
			cookie = resp.Cookies()
			return
		}
	}
}

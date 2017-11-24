package test

import (
	"time"
	"strings"
	"io/ioutil"
	"testing"
	"net/http"
	"github.com/beewit/wechat-ai/smartQQ"
	"github.com/beewit/beekit/utils/convert"
	"fmt"
	"encoding/json"
	"github.com/beewit/beekit/utils"
)

func TestAddQQ(t *testing.T) {
	client := &http.Client{Timeout: time.Second * time.Duration(30)}
	req, _ := http.NewRequest("POST", "http://ti.qq.com/mqqbase/cgi/qqrecommend/report?dwPageId=103&dwEntranceId=1031&dwExposeCnt=0&ddwUin=834979464&dwActionID=3&dwExposeTime=0&alghBuffer=200150|140006", nil)
	req.Header.Set("Cookie", "RK=uXnOl7jPkO; _qpsvr_localtk=tk5112; pgv_pvid=6382990204; p_uin=o2458208514; pt4_token=UzNrnGj2liegkz1S1DpnicWno4WarSWtjz43ZEnHOqw_; p_skey=zcS*c8Dr7vJf303SDTUrDmj-57Z05wcPylA7HKozV*U_; pt2gguin=o2458208514; uin=o2458208514; skey=@CtNhvCdei; ptisp=ctc; ptcz=6bcc6e291b0e5603b3a0c1cbea25387ef38863f6ddfa4d46627c7c1abb145800")

	req.Header.Set("Referer", "http://id.qq.com/possiblev3/possible.html?ver=7&frienduin=294477044&qqver=5545")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.2; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) QQ/8.9.5.22062 Chrome/43.0.2357.134 Safari/537.36 QBCore/3.43.716.400 QQBrowser/9.0.2524.400")

	resp, err := client.Do(req)
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	defer resp.Body.Close()
	println(string(bodyBytes))

}

func TestAddQQ2(t *testing.T) {
	client := &http.Client{Timeout: time.Second * time.Duration(30)}
	req, _ := http.NewRequest("GET", "http://ti.qq.com/mqqbase/cgi/qqrecommend/people?startpos=0&uinnum=15&relationuin=294477044&v=1509951924214&ldw=839709853", nil)
	req.Header.Set("Cookie", "RK=uXnOl7jPkO; _qpsvr_localtk=tk5112; pgv_pvid=6382990204; p_uin=o2458208514; pt4_token=UzNrnGj2liegkz1S1DpnicWno4WarSWtjz43ZEnHOqw_; p_skey=zcS*c8Dr7vJf303SDTUrDmj-57Z05wcPylA7HKozV*U_; pt2gguin=o2458208514; uin=o2458208514; skey=@CtNhvCdei; ptisp=ctc; ptcz=6bcc6e291b0e5603b3a0c1cbea25387ef38863f6ddfa4d46627c7c1abb145800")

	req.Header.Set("Referer", "http://id.qq.com/possiblev3/possible.html?ver=7&frienduin=294477044&qqver=5545")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.2; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) QQ/8.9.5.22062 Chrome/43.0.2357.134 Safari/537.36 QBCore/3.43.716.400 QQBrowser/9.0.2524.400")

	resp, err := client.Do(req)
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	defer resp.Body.Close()
	println(string(bodyBytes))

}

func TestGroupSearch(t *testing.T) {
	client := &http.Client{Timeout: time.Second * time.Duration(30)}
	req, _ := http.NewRequest("POST", "http://qun.qq.com/cgi-bin/group_search/pc_group_search", strings.NewReader(`k=%6A%61%76%61&n=24&st=1&iso=0&src=1&v=5503&bkn=839709853&isRecommend=false&city_id=0&from=1&keyword=%6A%61%76%61&sort=1&wantnum=24&page=0&ldw=839709853`))
	req.Header.Set("Cookie", "confirmuin=0; ptdrvs=H5ft9EaWLFnH4osg5uFyvxtAUUnnSs0a*oJvcYeR1xM_; ptvfsession=bf12354337e36b416b42458018a1d9c854f413f55a02f3096e1902a87a09f044935b3adc5d06202f273b340dd92ece587a30df9d9464953d; ptisp=ctc; pt2gguin=o2458208514; uin=o2458208514; skey=@CtNhvCdei; superuin=o2458208514; supertoken=3683791768; superkey=cztker8yeceVXVibYX8OPyVg7IOQVE1vkEKnORgWlvo_; pt_recent_uins=b1adf2ad0d85c3250886b780d157797145950c603b22b580a160977c10ff0971d20f91e2be1438a0eb32d1a84d78f818ade9ec8cd0dfd6ce; RK=kXnGh7jelO; ptnick_2458208514=e689bfe8afbae4b880e697b6e8aa93e8a880; pgv_pvi=8719185920; pgv_si=s1542224896; _qpsvr_localtk=0.05683227197119056;")
	req.Header.Set("Host", "qun.qq.com")
	req.Header.Set("Referer", "http://qun.qq.com/cgi-bin/group_search/pc_group_search")
	req.Header.Set("User-Agent", "Mozilla/4.0 (compatible; MSIE 8.0; Windows NT 5.1)")

	resp, err := client.Do(req)
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	defer resp.Body.Close()
	println(string(bodyBytes))
}

func TestGroupSearch2(t *testing.T) {
	client := &http.Client{Timeout: time.Second * time.Duration(30)}
	req, _ := http.NewRequest("POST", "http://qun.qq.com/cgi-bin/group_search/pc_group_search", strings.NewReader(`k=%6A%61%76%61&n=24&st=1&iso=0&src=1&v=5503&bkn=839709853&isRecommend=false&city_id=0&from=1&keyword=%6A%61%76%61&sort=1&wantnum=24&page=0&ldw=839709853`))
	req.Header.Set("Cookie", "confirmuin=0; ptdrvs=H5ft9EaWLFnH4osg5uFyvxtAUUnnSs0a*oJvcYeR1xM_; ptvfsession=bf12354337e36b416b42458018a1d9c854f413f55a02f3096e1902a87a09f044935b3adc5d06202f273b340dd92ece587a30df9d9464953d; ptisp=ctc; pt2gguin=o2458208514; uin=o2458208514; skey=@CtNhvCdei; superuin=o2458208514; supertoken=3683791768; superkey=cztker8yeceVXVibYX8OPyVg7IOQVE1vkEKnORgWlvo_; pt_recent_uins=b1adf2ad0d85c3250886b780d157797145950c603b22b580a160977c10ff0971d20f91e2be1438a0eb32d1a84d78f818ade9ec8cd0dfd6ce; RK=kXnGh7jelO; ptnick_2458208514=e689bfe8afbae4b880e697b6e8aa93e8a880; pgv_pvi=8719185920; pgv_si=s1542224896; _qpsvr_localtk=0.05683227197119056;")
	req.Header.Set("Host", "qun.qq.com")
	req.Header.Set("Referer", "http://qun.qq.com/cgi-bin/group_search/pc_group_search")
	req.Header.Set("User-Agent", "Mozilla/4.0 (compatible; MSIE 8.0; Windows NT 5.1)")

	resp, err := client.Do(req)
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	defer resp.Body.Close()
	println(string(bodyBytes))
}

func TestSearchQun(t *testing.T) {
	qq, err := smartQQ.Start(&smartQQ.QQClient{LoginCacheFilePath: "qqLogin.json", QrCodeFilePath: "qrcode.jpg"})
	if err != nil {
		println(err.Error())
		return
	}
	rep, _, err := qq.HttpRequestGet(strings.Replace(strings.Replace("http://ptlogin2.qun.qq.com/check_sig?pttype=2&uin={uid}&service=jump&nodirect=0&ptsigx={ptsigx}&s_url=http%3A%2F%2Fqun.qq.com%2Fmanage.html&f_url=&ptlang=2052&ptredirect=100&aid=1000101&daid=73&j_later=0&low_login_hour=0&regmaster=0&pt_login_type=2&pt_aid=715030901&pt_aaid=0&pt_light=0&pt_3rd_aid=0", "{uid}", convert.ToString(qq.Login.QQ), -1), "{ptsigx}", qq.PtSigX, -1), nil, nil)
	cks := rep.Cookies()
	println(convert.ToObjStr(cks))
	qq.HttpRequestPost(
		"http://qun.qq.com/cgi-bin/qun_mgr/get_group_list",
		map[string]string{"Referer": "http://qun.qq.com/member.html"},
		strings.NewReader("bkn=775849843"),
	)
}

func TestGroupSearchQun(t *testing.T) {
	qq, err := smartQQ.Start(&smartQQ.QQClient{LoginCacheFilePath: "qqLogin.json", QrCodeFilePath: "qrcode.jpg"})
	if err != nil {
		println(err.Error())
		return
	}
	sKey := qq.GetCookie("skey").Value
	uin := qq.GetCookie("uin").Value
	println(fmt.Sprintf("uin=%s; skey=%s", qq.GetCookie("uin").Value, sKey))
	client := &http.Client{Timeout: time.Second * time.Duration(30)}
	cityId := 0
	btnHash := btnHash(sKey)
	keyword := "美女"
	bnk := fmt.Sprintf(`k=&n=8&st=1&iso=1&src=1&v=4903&bkn=%s&isRecommend=false&city_id=%d&from=1&newSearch=true&keyword=%s&sort=0&wantnum=24&page=0&ldw=%s`,
		btnHash, cityId, keyword, btnHash)
	req, _ := http.NewRequest("POST", "http://qun.qq.com/cgi-bin/group_search/pc_group_search", strings.NewReader(bnk))
	req.Header.Set("Cookie", fmt.Sprintf("uin=%s; skey=%s", uin, sKey))
	req.Header.Set("Host", "qun.qq.com")
	req.Header.Set("Referer", "http://find.qq.com/index.html?version=1&im_version=5545&width=910&height=610&search_target=0")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.87 UBrowser/6.2.3831.3 Safari/537.36")

	resp, err := client.Do(req)
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	defer resp.Body.Close()
	println("【搜索群】", string(bodyBytes))
}

func TestGetGroupList(t *testing.T) {
	skey := "@Bl3yTiwVz"
	uin := "o2458208514"
	pskey := "sIcZS8p*5SOMK8iEu20jFGc*frCzuWB7p-2mslDlENk_"
	client := &http.Client{Timeout: time.Second * time.Duration(30)}
	req, _ := http.NewRequest("POST", "http://qun.qq.com/cgi-bin/qun_mgr/get_group_list", strings.NewReader(fmt.Sprintf("bkn=%s", btnHash(skey))))
	req.Header.Set("Cookie", fmt.Sprintf("uin=%s; skey=%s; p_uin=%s; p_skey=%s", uin, skey, uin, pskey))
	req.Header.Set("Host", "qun.qq.com")
	req.Header.Set("Referer", "http://qun.qq.com/member.html")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.87 UBrowser/6.2.3831.3 Safari/537.36")

	resp, err := client.Do(req)
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	defer resp.Body.Close()
	println("【查询群】", string(bodyBytes))
	var gl smartQQ.GroupList2
	json.Unmarshal(bodyBytes, &gl)
	println(convert.ToObjStr(gl))
}

func TestGetFriendList(t *testing.T) {
	skey := "@7Efu2eM9y"
	uin := "o0294477044"
	pskey := "JW545cG4-ydFPJeFxUAS15kWZHqAISg*K7yP*HeyL34_"
	client := &http.Client{Timeout: time.Second * time.Duration(30)}
	req, _ := http.NewRequest("POST", "http://qun.qq.com/cgi-bin/qun_mgr/get_friend_list", strings.NewReader(fmt.Sprintf("bkn=%s", btnHash(skey))))
	req.Header.Set("Cookie", fmt.Sprintf("uin=%s; skey=%s; p_uin=%s; p_skey=%s", uin, skey, uin, pskey))
	req.Header.Set("Host", "qun.qq.com")
	req.Header.Set("Referer", "http://qun.qq.com/member.html")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.87 UBrowser/6.2.3831.3 Safari/537.36")

	resp, err := client.Do(req)
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	defer resp.Body.Close()
	println("【查询好友】", string(bodyBytes))

	var flist smartQQ.FriendList2
	err = json.Unmarshal(bodyBytes, &flist)
	if err != nil {
		println(err.Error())
	}
	println(convert.ToObjStr(flist))
}

func TestSearchGroupMembers(t *testing.T) {
	skey := "@Bl3yTiwVz"
	uin := "o2458208514"
	pskey := "sIcZS8p*5SOMK8iEu20jFGc*frCzuWB7p-2mslDlENk_"
	client := &http.Client{Timeout: time.Second * time.Duration(30)}
	req, _ := http.NewRequest("POST", "http://qinfo.clt.qq.com/cgi-bin/qun_info/get_group_members_new", strings.NewReader(fmt.Sprintf("gc=11862108&st=0&end=5000&sort=0&bkn=%s", btnHash(skey))))
	req.Header.Set("Cookie", fmt.Sprintf("uin=%s; skey=%s; p_uin=%s; p_skey=%s", uin, skey, uin, pskey))
	req.Header.Set("Host", "qun.qq.com")
	req.Header.Set("Origin", "http://qun.qq.com")
	req.Header.Set("Referer", "http://qun.qq.com/member.html")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.87 UBrowser/6.2.3831.3 Safari/537.36")

	resp, err := client.Do(req)
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	defer resp.Body.Close()
	println("【查询群成员】", string(bodyBytes))
}

func btnHash(e string) string {
	n := int64(5381)
	i := len(e)
	for r := 0; r < i; r++ {
		n += (n << 5) + int64(e[r])
	}
	return convert.ToString(n & 2147483647)
}

func TestBtnHash(t *testing.T) {
	println(btnHash("@7KapkyMz0"))
}

func TestMac(t *testing.T) {
	//b0:25:aa:17:91:a0
	//eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.e30.ri553VU8FhX7M9DX2cbdFKvsXofDxHWcE3kXhPg7vOAd7f91fc14aee950d71abbf3d360f0633
	//eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.e30.r1QhVVk20WkJPjagz_ENT5gee2XZEvrEQRpChcKx3i4d7f91fc14aee950d71abbf3d360f0633
	println(utils.GetMac())
}

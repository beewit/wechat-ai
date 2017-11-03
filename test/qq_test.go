package test

import (
	"testing"
	"fmt"
	"regexp"
	"github.com/beewit/wechat-ai/smartQQ"
	"github.com/beewit/beekit/utils/convert"
	"io/ioutil"
	"encoding/json"
	"net/http"
	"time"
	"strings"
	"os"
)

func TestMap(t *testing.T) {

}

func TestStartQQ(t *testing.T) {
	qq, err := smartQQ.Start(&smartQQ.QQClient{LoginCacheFilePath: "qqLogin.json", QrCodeFilePath: "qrcode.jpg"})
	if err != nil {
		println(err.Error())
		return
	}
	if qq != nil && qq.GroupInfo != nil {
		for _, v := range qq.GroupInfo {
			reg, err := qq.GetGroupInfo(v.Code)
			if err != nil {
				println(fmt.Sprintf("【%s】加载群信息失败，ERROR：%s", v.Name, err.Error()))
			} else {
				println(fmt.Sprintf("【%s】加载群信息结果：%s", v.Name, convert.ToObjStr(reg)))
			}
			time.Sleep(time.Second * 3)
		}
	}

	qq.SendMsg(2829015118, "你好哦！")

	go func() {
		pollResult, err := qq.Poll2(func(qq *smartQQ.QQClient, result smartQQ.QQResponsePoll) {
			if len(result.Result) > 0 && len(result.Result[0].Value.Content) > 0 {
				var message string
				if result.Result[0].PollType == "group_message" {
					group := qq.GroupInfo[result.Result[0].Value.GroupCode]
					if group.GId > 0 {
						message = " 【群消息 - " + group.Name + "】 "
					}
				}
				sendUser := qq.FriendsMap.Info[result.Result[0].Value.SendUin]
				if sendUser.Uin > 0 {
					message += "   -   发送人《" + qq.FriendsMap.Info[result.Result[0].Value.SendUin].Nick + "》"
				}
				for i := 0; i < len(result.Result[0].Value.Content); i++ {
					if i > 0 {
						message += convert.ToObjStr(result.Result[0].Value.Content[i])
					}
				}
				println("您有新消息了哦！ ==>> ", message)
			}
		})
		if err != nil {
			t.Error("Poll2", err.Error())
			return
		}
		println(convert.ToObjStr(pollResult))
	}()

	time.Sleep(time.Second * 20)
}

func TestCheckSig(t *testing.T) {
	url := `{"result":[{"client_type":1,"status":"online","uin":2545321045},{"client_type":7,"status":"online","uin":2441435397},{"client_type":7,"status":"online","uin":3923104280}],"retcode":0}`
	var qqRes smartQQ.QQResponseObj
	err := json.Unmarshal([]byte(url), &qqRes)

	if err != nil {
		println(err.Error())
	}
	println(qqRes.RetCode)
}

func TestChar(t *testing.T) {
	client := &http.Client{Timeout: time.Second * time.Duration(30)}
	req, _ := http.NewRequest("POST", "http://d1.web2.qq.com/channel/send_buddy_msg2", strings.NewReader(`r={"to":2289681300,"content":"[\"sadf222\",[\"font\",{\"name\":\"宋体\",\"size\":10,\"style\":[0,0,0],\"color\":\"000000\"}]]","face":594,"clientid":53999199,"msg_id":43860002,"psessionid":"8368046764001d636f6e6e7365727665725f77656271714031302e3133332e34312e383400001ad00000066b026e040015808a206d0000000a406172314338344a69526d0000002859185d94e66218548d1ecb1a12513c86126b3afb97a3c2955b1070324790733ddb059ab166de6857"}`))
	req.Header.Set("Cookie", "RK=FWH2yaPqNd; tvfe_boss_uuid=e5cd66a63db6c808; _qpsvr_localtk=0.4244528137550774; pgv_pvi=1155806208; pgv_si=s7372022784; luin=o0294477044; lskey=000100004f911f2883f0bdfa406fed6fc34b1824f87412aac66c127d6fab16d21498e25c9f4b3eb797718d3e; ptui_loginuin=294477044; FTN5K=aa9fc15c; rv2=80B8EC1E02F46C6133A5AA8D23E8087F021BF1157AAC9C1B74; property20=33049670F316BCEEE0AF407FE81799CD73229E1F509FCD8789B2ECFA7B6515E4D9CA0B4110DE0405; qqmusic_uin=; qqmusic_key=; qqmusic_fromtag=; o_cookie=294477044; pgv_info=ssid=s7269520736; pgv_pvid=1432444209; ptisp=ctc; ptcz=2364cbbd95bd23a29569dd728d75bf743f40bb1f1c865f7da0d55d8ae0819f98; uin=o2458208514; skey=@aFn4kl249; pt2gguin=o2458208514; p_uin=o2458208514; pt4_token=8NBl7D3inbjkj0MTm0mHLcrdd8ecBv4eoDEVD21lW*A_; p_skey=6gIRjWGUvt*lKTMCPVFkkmddOv7YnjDHP30jL54w6EU_")
	req.Header.Set("Host", "d1.web2.qq.com")
	req.Header.Set("Origin", "http://d1.web2.qq.com")
	req.Header.Set("Referer", "http://d1.web2.qq.com/cfproxy.html?v=20151105001&callback=1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36")

	resp, err := client.Do(req)
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	defer resp.Body.Close()
	println(string(bodyBytes))
}

func TestFried(t *testing.T) {
	//10bb5430d1c1d83901320d7a737fa3d44bb2f167a6aca5ac798298ffe60adfca4ae3d584e5a725df
	client := &http.Client{Timeout: time.Second * time.Duration(30)}
	req, _ := http.NewRequest("POST", "http://s.web2.qq.com/api/get_user_friends2", strings.NewReader(`r={"vfwebqq":"10bb5430d1c1d83901320d7a737fa3d44bb2f167a6aca5ac798298ffe60adfca4ae3d584e5a725df","hash":"54D75BC600065649"}`))
	req.Header.Set("Cookie", "RK=FWH2yaPqNd; tvfe_boss_uuid=e5cd66a63db6c808; _qpsvr_localtk=0.4244528137550774; pgv_pvi=1155806208; pgv_si=s7372022784; luin=o0294477044; lskey=000100004f911f2883f0bdfa406fed6fc34b1824f87412aac66c127d6fab16d21498e25c9f4b3eb797718d3e; ptui_loginuin=294477044; FTN5K=aa9fc15c; rv2=80B8EC1E02F46C6133A5AA8D23E8087F021BF1157AAC9C1B74; property20=33049670F316BCEEE0AF407FE81799CD73229E1F509FCD8789B2ECFA7B6515E4D9CA0B4110DE0405; qqmusic_uin=; qqmusic_key=; qqmusic_fromtag=; o_cookie=294477044; pgv_info=ssid=s7269520736; pgv_pvid=1432444209; ptisp=ctc; ptcz=2364cbbd95bd23a29569dd728d75bf743f40bb1f1c865f7da0d55d8ae0819f98; uin=o2458208514; skey=@Bd39Ijhg1; pt2gguin=o2458208514; p_uin=o2458208514; pt4_token=eEL1cK5Q3Q0HnjA01JSmY1lVUSflkyQPWTKQTnCUbO0_; p_skey=StoA*ueMWys7p*xz3jEjeiJ2ZvFCktjDFDtbJYFKBNo_; ptwebqq=0d1a1c8deb2579fc200601245d1e76c31a0fe2ca61a044a7cc80e5eafb311601")
	req.Header.Set("Host", "s.web2.qq.com")
	req.Header.Set("Origin", "http://s.web2.qq.com")
	req.Header.Set("Referer", "http://s.web2.qq.com/proxy.html?v=20130916001&callback=1&id=1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	defer resp.Body.Close()
	println(string(bodyBytes))
}

func TestGroup(t *testing.T) {
	//10bb5430d1c1d83901320d7a737fa3d44bb2f167a6aca5ac798298ffe60adfca4ae3d584e5a725df
	client := &http.Client{Timeout: time.Second * time.Duration(30)}
	req, _ := http.NewRequest("POST", "http://s.web2.qq.com/api/get_group_name_list_mask2", strings.NewReader(`r={"vfwebqq":"10bb5430d1c1d83901320d7a737fa3d44bb2f167a6aca5ac798298ffe60adfca4ae3d584e5a725df","hash":"54D75BC600065649"}`))
	req.Header.Set("Cookie", "RK=FWH2yaPqNd; tvfe_boss_uuid=e5cd66a63db6c808; _qpsvr_localtk=0.4244528137550774; pgv_pvi=1155806208; pgv_si=s7372022784; luin=o0294477044; lskey=000100004f911f2883f0bdfa406fed6fc34b1824f87412aac66c127d6fab16d21498e25c9f4b3eb797718d3e; ptui_loginuin=294477044; FTN5K=aa9fc15c; rv2=80B8EC1E02F46C6133A5AA8D23E8087F021BF1157AAC9C1B74; property20=33049670F316BCEEE0AF407FE81799CD73229E1F509FCD8789B2ECFA7B6515E4D9CA0B4110DE0405; qqmusic_uin=; qqmusic_key=; qqmusic_fromtag=; o_cookie=294477044; pgv_info=ssid=s7269520736; pgv_pvid=1432444209; ptisp=ctc; ptcz=2364cbbd95bd23a29569dd728d75bf743f40bb1f1c865f7da0d55d8ae0819f98; uin=o2458208514; skey=@Bd39Ijhg1; pt2gguin=o2458208514; p_uin=o2458208514; pt4_token=eEL1cK5Q3Q0HnjA01JSmY1lVUSflkyQPWTKQTnCUbO0_; p_skey=StoA*ueMWys7p*xz3jEjeiJ2ZvFCktjDFDtbJYFKBNo_; ptwebqq=0d1a1c8deb2579fc200601245d1e76c31a0fe2ca61a044a7cc80e5eafb311601")
	req.Header.Set("Host", "s.web2.qq.com")
	req.Header.Set("Origin", "http://s.web2.qq.com")
	req.Header.Set("Referer", "http://s.web2.qq.com/proxy.html?v=20130916001&callback=1&id=1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36")

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", "148")
	resp, err := client.Do(req)
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	defer resp.Body.Close()
	println(string(bodyBytes))
}

func TestLoginQQ(t *testing.T) {
	var qq *smartQQ.QQClient
	bts, _ := ioutil.ReadFile("qqLogin.json")
	if bts != nil && string(bts) != "" {
		json.Unmarshal(bts, &qq)
	} else {
		qq = smartQQ.NewQQClient(qq)
		qq.QrCodeFilePath = "out.jpg"
		_, err := qq.PtqrShow()
		if err != nil {
			t.Error(err.Error())
			return
		}
		_, err = qq.CheckLogin(nil)
		if err != nil {
			println("错误：", err.Error())
			t.Error(err.Error())
			return
		}
		println(convert.ToObjStr(qq))
	}
	f, err := os.OpenFile("qqLogin.json", os.O_RDWR|os.O_CREATE, os.ModePerm)
	defer f.Close()
	f.Write([]byte(convert.ToObjStr(qq)))
	if err != nil {
		t.Error("OpenFile", err.Error())
		return
	}

	pollResult, err := qq.Poll2(func(qq *smartQQ.QQClient, result smartQQ.QQResponsePoll) {
		if len(result.Result) > 0 && len(result.Result[0].Value.Content) > 0 {
			var message string
			if result.Result[0].PollType == "group_message" {
				message = " 【群消息 - " + qq.GroupInfo[result.Result[0].Value.GroupCode].Name + "】 "
			}
			message += "   -   发送人《" + qq.FriendsMap.Info[result.Result[0].Value.SendUin].Nick + "》"
			for i := 0; i < len(result.Result[0].Value.Content); i++ {
				if i > 0 {
					message += convert.ToObjStr(result.Result[0].Value.Content[i])
				}
			}
			println("您有新消息了哦！ ==>> ", message)
		}
	})
	if err != nil {
		t.Error("Poll2", err.Error())
		return
	}
	println(convert.ToObjStr(pollResult))
}

func TestHash(t *testing.T) {
	str := smartQQ.Hash(convert.MustInt64("2458208514"), "")
	println(str)

}

func TestQQToken(t *testing.T) {
	qrsig := "5xY*M*lybU7ffUlOL5uBWEi0DZwsEK32r2mRfoba02F5k9ekyZc*HFntB162i2Ty"
	println(uid4444(qrsig))
}

func TestPtuiCB(t *testing.T) {
	str := `ptuiCB('0','0','http://ptlogin2.web2.qq.com/check_sig?pttype=1&uin=294477044&service=ptqrlogin&nodirect=0&ptsigx=c998030f792bc3ac8b7b7de6f5270d71d068148a3d8fb4fc017f0fb009df05291bc37035c6bc721549fe6e8694ea197a441115070cae6b25e0098fc5cdb0bde7&s_url=http%3A%2F%2Fw.qq.com%2Fproxy.html&f_url=&ptlang=2052&ptredirect=100&aid=501004106&daid=164&j_later=0&low_login_hour=0&regmaster=0&pt_login_type=3&pt_aid=0&pt_aaid=16&pt_light=0&pt_3rd_aid=0','0','登录成功！', '承諾，\/aiq一世的誓言')`
	//regexp_image_status := regexp.MustCompile(`'\d+'`)
	//code := regexp_image_status.FindAllString(str, 1)[0]
	//println(code)
	if reg_sig := regexp.MustCompile(`ptuiCB\(\'0\',\'0\',\'([^\']+)\'`).FindAllStringSubmatch(str, -1); len(reg_sig) == 1 {
		println("结果：", reg_sig[0][1])
	} else {
		fmt.Println("Check Sig Err:")
		return
	}
	//`ptuiCB\(\'0\',\'0\',\'\w+\',\'([^\']+)\'`

}

func uid4444(str string) string {
	skey := []byte(str)
	e := 0
	for i, n := 0, len(str); n > i; i++ {
		e += (e << 5) + int(skey[i])
	}

	return fmt.Sprint(2147483647 & e)
}

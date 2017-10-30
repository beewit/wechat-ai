package test

import (
	"testing"
	"fmt"
	"regexp"
	"github.com/beewit/wechat-ai/smartqq"
	"github.com/beewit/beekit/utils/convert"
)

func TestLoginQQ(t *testing.T) {
	qq := smartqq.NewQQClient()
	qq.QrCodeFilePath = "out.jpg"
	u := qq.GetLoginUrl(`ptuiCB('0','0','http://ptlogin2.web2.qq.com/check_sig?pttype=1&uin=294477044&service=ptqrlogin&nodirect=0&ptsigx=6f51f90df3b70b663c21d351d7415ef63db8c79207dbff149ffae7b1297c2f702af77ec5fb8ee02331e04c760d798de6091b7503901750bc16ae9ad4342cacb1&s_url=http%3A%2F%2Fw.qq.com%2Fproxy.html&f_url=&ptlang=2052&ptredirect=100&aid=501004106&daid=164&j_later=0&low_login_hour=0&regmaster=0&pt_login_type=3&pt_aid=0&pt_aaid=16&pt_light=0&pt_3rd_aid=0','0','登录成功！', '承諾，\/aiq一世的誓言')`)
	println(u)
	//_, err := qq.PtqrShow()
	//if err != nil {
	//	t.Error(err.Error())
	//}
	//_, err = qq.CheckLogin()
	//if err != nil {
	//	t.Error(err.Error())
	//}
	//println(convert.ToObjStr(qq))
}

func TestHash(t *testing.T) {
	str := smartqq.Hash(convert.MustInt64("294477044"), "7de29efeba3a479ea50dee3c321a5db7650c3d1826253aefd87bc9a56c738f55d0aed70282be76cb")
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

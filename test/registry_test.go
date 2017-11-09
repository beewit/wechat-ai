package test

import (
	"testing"
	"fmt"
	"golang.org/x/sys/windows/registry"
	"os/exec"
	"github.com/beewit/wechat-ai/ai"
	"strconv"
	"github.com/beewit/beekit/utils/convert"
	"github.com/lxn/win"
)

func TestRegistry(t *testing.T) {
	k, err := registry.OpenKey(registry.CURRENT_USER, `Software\\Tencent\\PlatForm_Type_List\\1`, registry.QUERY_VALUE)
	if err != nil {
		k, err = registry.OpenKey(registry.LOCAL_MACHINE, `Software\\Tencent\\PlatForm_Type_List\\1`, registry.QUERY_VALUE)
		if err != nil {
			t.Fatal(err)
		}
	}
	defer k.Close()

	s, _, err := k.GetStringValue("TypePath")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%q\n", s)
	println(s)
}

func TestOpen(t *testing.T) {
	CallEXE("D:\\QQ\\Info\\Bin\\QQ.exe")
}

func CallEXE(strGameName string) {
	cmd := exec.Command(strGameName)
	err := cmd.Start()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	println("执行成功")
}

func TestChar2(t *testing.T) {
	keys := "13696433488love"
	for _, v := range keys {
		println("OLD", v, "NEW", int(v), "char", string(v))
	}

}

func TestFZ(t *testing.T) {
	h := win.GetForegroundWindow()
	println(h)
	win.SetForegroundWindow(h)
}

func TestStart(t *testing.T) {
	err := ai.QQLogin(3240033436, "13696433488love")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	ai.AddQQGroup(643413680,"旅游一下了")
	ai.AddQQFriend(2818707257,"嗨喽")
	//ai.KeydbStr("13696433488loveLoveMm*?{}[]")/.,<>()
}

//二进制转十六进制
func btox(b string) int32 {
	base, _ := strconv.ParseInt(b, 2, 10)
	println(strconv.FormatInt(base, 16))
	return convert.MustInt32(strconv.FormatInt(base, 16))
}

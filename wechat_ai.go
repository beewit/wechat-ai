package main

import (
	"github.com/sclevine/agouti"
	"time"
)

func main() {
	Driver := agouti.ChromeDriver(agouti.ChromeOptions("args", []string{
		"--user-data-dir=ChromeUserData",
		"--gpu-process",
		"--start-maximized",
		"--disable-infobars",
		"--app=http://www.baidu.com",
		"--webkit-text-size-adjust"}))
	Driver.Start()
	var err error
	Page, err := Driver.NewPage()
	if err != nil {
		println("Failed to open page.")
		return
	}
	time.Sleep(time.Second * 3)
	Page.Navigate("http://www.jd.com")
}

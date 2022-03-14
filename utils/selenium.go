package utils

import (
	"fmt"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

const (
	//设置常量 分别设置chromedriver.exe的地址和本地调用端口
	seleniumPath = `C:\Program Files\Google\Chrome\Application\chromedriver.exe`
	port         = 9515
)

var (
	chromeCaps = chrome.Capabilities{
		Prefs: map[string]interface{}{ // 禁止加载图片，加快渲染速度
			"profile.managed_default_content_settings.images": 2,
		},
		Path: "",
		Args: []string{
			// "--headless",
			"--start-maximized",
			"--window-size=1920x1080",
			"--no-sandbox",
			"--user-agent=Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.77 Safari/537.36",
			"--disable-gpu",
			"--disable-impl-side-painting",
			"--disable-gpu-sandbox",
			"--disable-accelerated-2d-canvas",
			"--disable-accelerated-jpeg-decoding",
			"--test-type=ui",
			"--ignore-certificate-errors",
		},
	}
	//设置selenium服务的选项,设置为空。根据需要设置。
	ops = make([]selenium.ServiceOption, 0)
	//设置浏览器兼容性，设置浏览器名称为chrome
	caps = selenium.Capabilities{"browserName": "chrome"}
)

func NewService()(*selenium.Service, error){
	return selenium.NewChromeDriverService(seleniumPath, port, ops...)
}

func NewWindow()(selenium.WebDriver, error){
	caps.AddChrome(chromeCaps)
	return selenium.NewRemote(caps, fmt.Sprintf("http://127.0.0.1:%v/wd/hub", port))
}
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/martinlindhe/notify"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

const (
	mainURL          = "https://eclass.kunsan.ac.kr/index.jsp?sso=ok"
	seleniumPath     = "selenium-server-standalone.jar"
	chromeDriverPath = "chromedriver.exe"
	port             = 8080
)

func main() {
	opts := []selenium.ServiceOption{
		selenium.GeckoDriver(chromeDriverPath), // Specify the path to GeckoDriver in order to use Firefox.
		selenium.Output(os.Stderr),             // Output debug information to STDERR.
	}

	service, _ := selenium.NewSeleniumService(seleniumPath, port, opts...)
	defer service.Stop()

	caps := selenium.Capabilities{"browserName": "chrome"}

	chromeCaps := chrome.Capabilities{
		Path: "",
		Args: []string{
			"--no-sandbox",
			"--user-agent=Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_2) AppleWebKit/604.4.7 (KHTML, like Gecko) Version/11.0.2 Safari/604.4.7",
		},
	}

	caps.AddChrome(chromeCaps)

	wd, _ := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
	defer wd.Quit()

	wd.Get(mainURL)

	//loginID, _ := wd.FindElement(selenium.ByCSSSelector, "#id")
	//loginPW, _ := wd.FindElement(selenium.ByCSSSelector, "#pw")
	//
	//loginID.SendKeys("ID")
	//loginPW.SendKeys("PW")
	//
	//loginBtnElem, _ := wd.FindElement(selenium.ByXPATH, "//*[@id=\"loginForm\"]/fieldset/p[2]/a")
	//
	//loginBtnElem.Click()

	for {
		fmt.Println("대기 중...")
		windows, _ := wd.WindowHandles()

		for len(windows) == 2 {
			fmt.Println("찾음")

			wd.SwitchWindow(windows[1])
			_ = wd.SwitchFrame("contentsCheckForm")

			windowTitle, _ := wd.Title()

			if windowTitle != "Contents Viewer" {
				break
			}

			fmt.Println("AFK 버튼 기다리는 중...")

			afkBtnElem, notAfk := wd.FindElement(selenium.ByXPATH, "/html/body/form/div/div[2]/ul/li/a")
			if notAfk == nil {
				fmt.Printf("AFK 버튼 찾음")
				afkBtnElem.Click()

				notify.Alert("Auto AFK Clicker", "Clicked", "AFK button clicked", "")
			}

			windows, _ = wd.WindowHandles()

			time.Sleep(1 * time.Second)
		}

		time.Sleep(1 * time.Second)
	}
}

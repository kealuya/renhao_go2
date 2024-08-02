package tray_cmd

import (
	"fmt"
	"github.com/getlantern/systray/example/icon"
	"github.com/gookit/goutil/timex"
	"log"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/getlantern/systray"
)

func TrayCmd() {
	systray.Run(onReady, onExit)
	//systray.Run(onReadyComplex, onExit)

}

func onReady() {
	systray.SetIcon(getIcon("./tray_cmd/icon.ico"))
	systray.SetTitle("任浩的go程序")
	systray.SetTooltip("Pretty awesome超级棒")

	myState := systray.AddMenuItem("start...", "State")
	systray.AddSeparator()
	openBrowser := systray.AddMenuItem("openBrowser", "打开浏览器")
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit", "Quit the application")
	mQuit.SetIcon(getIcon("./tray_cmd/quit.ico"))

	go func() {
		<-mQuit.ClickedCh
		fmt.Println("Quit clicked")
		systray.Quit()
		fmt.Println("Quit finished")
	}()
	go func() {
		<-openBrowser.ClickedCh
		fmt.Println("openBrowser clicked")
		openBrowserFunc("https://www.baidu.com")
		fmt.Println("openBrowser finished")
	}()
	ticker := time.NewTicker(1 * time.Second)
	go func() {
		for range ticker.C {
			s := timex.Now().Format("2006-01-02 15:04:05")
			myState.SetTitle(s)
		}
	}()
}

func onExit() {
	// 清理操作
}

func getIcon(s string) []byte {
	b, err := os.ReadFile(s)
	if err != nil {
		log.Fatal(err)
	}
	return b
}
func openBrowserFunc(url string) {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "rundll32"
		args = []string{"url.dll,FileProtocolHandler", url}
	case "darwin":
		cmd = "open"
		args = []string{url}
	case "linux":
		cmd = "xdg-open"
		args = []string{url}
	default:
		fmt.Printf("unsupported platform: %s\n", runtime.GOOS)
		return
	}

	err := exec.Command(cmd, args...).Start()
	if err != nil {
		fmt.Printf("failed to open browser: %v\n", err)
	}
}
func onReadyComplex() {
	systray.SetTemplateIcon(icon.Data, icon.Data)
	systray.SetTitle("Awesome App")
	systray.SetTooltip("Lantern")
	mQuitOrig := systray.AddMenuItem("Quit", "Quit the whole app")
	go func() {
		<-mQuitOrig.ClickedCh
		fmt.Println("Requesting quit")
		systray.Quit()
		fmt.Println("Finished quitting")
	}()

	// We can manipulate the systray in other goroutines
	go func() {
		systray.SetTemplateIcon(icon.Data, icon.Data)
		systray.SetTitle("Awesome App")
		systray.SetTooltip("Pretty awesome棒棒嗒")
		mChange := systray.AddMenuItem("Change Me", "Change Me")
		mAllowRemoval := systray.AddMenuItem("Allow removal", "macOS only: allow removal of the icon when cmd is pressed")
		mChecked := systray.AddMenuItemCheckbox("Unchecked", "Check Me", true)
		mEnabled := systray.AddMenuItem("Enabled", "Enabled")
		// Sets the icon of a menu item. Only available on Mac.
		mEnabled.SetTemplateIcon(icon.Data, icon.Data)

		systray.AddMenuItem("Ignored", "Ignored")

		subMenuTop := systray.AddMenuItem("SubMenuTop", "SubMenu Test (top)")
		subMenuMiddle := subMenuTop.AddSubMenuItem("SubMenuMiddle", "SubMenu Test (middle)")
		subMenuBottom := subMenuMiddle.AddSubMenuItemCheckbox("SubMenuBottom - Toggle Panic!", "SubMenu Test (bottom) - Hide/Show Panic!", false)
		subMenuBottom2 := subMenuMiddle.AddSubMenuItem("SubMenuBottom - Panic!", "SubMenu Test (bottom)")

		mUrl := systray.AddMenuItem("Open UI", "my home")
		mQuit := systray.AddMenuItem("退出", "Quit the whole app")

		// Sets the icon of a menu item. Only available on Mac.
		mQuit.SetIcon(icon.Data)

		systray.AddSeparator()
		mToggle := systray.AddMenuItem("Toggle", "Toggle the Quit button")
		shown := true
		toggle := func() {
			if shown {
				subMenuBottom.Check()
				subMenuBottom2.Hide()
				mQuitOrig.Hide()
				mEnabled.Hide()
				shown = false
			} else {
				subMenuBottom.Uncheck()
				subMenuBottom2.Show()
				mQuitOrig.Show()
				mEnabled.Show()
				shown = true
			}
		}

		for {
			select {
			case <-mChange.ClickedCh:
				mChange.SetTitle("I've Changed")
			case <-mAllowRemoval.ClickedCh:

			case <-mChecked.ClickedCh:
				if mChecked.Checked() {
					mChecked.Uncheck()
					mChecked.SetTitle("Unchecked")
				} else {
					mChecked.Check()
					mChecked.SetTitle("Checked")
				}
			case <-mEnabled.ClickedCh:
				mEnabled.SetTitle("Disabled")
				mEnabled.Disable()
			case <-mUrl.ClickedCh:

			case <-subMenuBottom2.ClickedCh:
				panic("panic button pressed")
			case <-subMenuBottom.ClickedCh:
				toggle()
			case <-mToggle.ClickedCh:
				toggle()
			case <-mQuit.ClickedCh:
				systray.Quit()
				fmt.Println("Quit2 now...")
				return
			}
		}
	}()
}

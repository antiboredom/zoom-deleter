package main

import (
	"fmt"
	// "github.com/emersion/go-autostart"
	"github.com/getlantern/systray"
	"github.com/mitchellh/go-homedir"
	"os"
	"path"
	"runtime"
	"time"
	"zoom.lav.io/zoom_deleter/v2/icon"
)

func main() {
	onExit := func() {
	}

	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetIcon(icon.Data)
	// systray.SetTitle("Zoom Deleter")
	systray.SetTooltip("Zoom Deleter")
	mInfo := systray.AddMenuItem("Status: Active", "Status: Active")
	mInfo.Disable()

	systray.AddSeparator()

	mQuitOrig := systray.AddMenuItem("Quit", "Quit Zoom Deleter")
	go func() {
		<-mQuitOrig.ClickedCh
		fmt.Println("Requesting quit")
		systray.Quit()
		fmt.Println("Finished quitting")
	}()

	deleteTicker := time.NewTicker(1 * time.Second)
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				return
			case t := <-deleteTicker.C:
				fmt.Println("Tick at", t)
				deleter()
			}
		}
	}()
}

func deleter() {
	switch osName := runtime.GOOS; osName {
	case "darwin":
		e := os.RemoveAll("/Applications/zoom.us.app")
		if e != nil {

		}
	case "linux":
		fmt.Println("Not implemented yet for Linux!")
	case "windows":
		userHome, _ := homedir.Dir()
		zoomPath := path.Join(userHome, "AppData\\Roaming\\Zoom")
		_ = os.RemoveAll(zoomPath)
		_ = os.RemoveAll("C:\\Program Files (x86)\\Zoom")
		// Zoom\uninstall\Installer.exe /uninstall
	}
}

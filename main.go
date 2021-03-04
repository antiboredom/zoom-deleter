package main

import (
	"fmt"
	"github.com/emersion/go-autostart"
	"github.com/getlantern/systray"
	"github.com/mitchellh/go-homedir"
	"os"
	"path"
	"runtime"
	"time"
	"zoom.lav.io/zoom_deleter/v2/icon"
)

//os.Executable()

func main() {
	onExit := func() {
	}

	systray.Run(onReady, onExit)
}

func onReady() {
	executable, _ := os.Executable()
	app := &autostart.App{
		Name:        "zoomdeleter",
		DisplayName: "Zoom Deleter",
		Exec:        []string{executable},
	}

	systray.SetIcon(icon.Data)
	// systray.SetTitle("Zoom Deleter")
	systray.SetTooltip("Zoom Deleter")
	mInfo := systray.AddMenuItem("Status: Active", "Status: Active")
	mInfo.Disable()

	runningOnStart := app.IsEnabled()
	runningMessage := "Run At Start"
	mStart := systray.AddMenuItemCheckbox(runningMessage, runningMessage, runningOnStart)

	systray.AddSeparator()

	mQuitOrig := systray.AddMenuItem("Quit", "Quit Zoom Deleter")

	deleteTicker := time.NewTicker(5 * time.Second)
	done := make(chan bool)

	for {
		select {
		case <-done:
			return
		case t := <-deleteTicker.C:
			fmt.Println("Tick at", t)
			deleter()
		case <-mStart.ClickedCh:
			if mStart.Checked() {
				mStart.Uncheck()
				fmt.Println("Unchecking")
				if err := app.Disable(); err != nil {
					fmt.Println(err)
				}
			} else {
				mStart.Check()
				fmt.Println("Checking")
				if err := app.Enable(); err != nil {
					fmt.Println(err)
				}
			}
		case <-mQuitOrig.ClickedCh:
			systray.Quit()
			return
		}
	}

	// deleteTicker := time.NewTicker(5 * time.Second)
	// done := make(chan bool)
	//
	// go func() {
	// 	for {
	// 		select {
	// 		case <-done:
	// 			return
	// 		case t := <-deleteTicker.C:
	// 			fmt.Println("Tick at", t)
	// 			deleter()
	// 		}
	// 	}
	// }()
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

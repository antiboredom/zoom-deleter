package main

import (
	"fmt"
	"github.com/emersion/go-autostart"
	"github.com/getlantern/systray"
	"github.com/mitchellh/go-homedir"
	"os"
	"path"
	"path/filepath"
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

	deleteTicker := time.NewTicker(10 * time.Second)
	done := make(chan bool)

	for {
		select {
		case <-done:
			return
		case <-deleteTicker.C:
			deleter()
		case <-mStart.ClickedCh:
			if mStart.Checked() {
				mStart.Uncheck()
				if err := app.Disable(); err != nil {
					fmt.Println(err)
				}
			} else {
				mStart.Check()
				if err := app.Enable(); err != nil {
					fmt.Println(err)
				}
			}
		case <-mQuitOrig.ClickedCh:
			systray.Quit()
			return
		}
	}
}

func deleter() {
	switch osName := runtime.GOOS; osName {
	case "darwin":
		os.RemoveAll("/Applications/zoom.us.app")
		os.RemoveAll("/Applications/Microsoft Teams.app")
		goToMeetings, _ := filepath.Glob("/Applications/GoToMeeting*.app")
		for _, goToMeeting := range goToMeetings {
			os.RemoveAll(goToMeeting)
		}
	case "linux":
		fmt.Println("Not implemented yet for Linux!")
	case "windows":
		userHome, _ := homedir.Dir()

		zoomPath := path.Join(userHome, "AppData\\Roaming\\Zoom")
		os.RemoveAll(zoomPath)
		os.RemoveAll("C:\\Program Files (x86)\\Zoom")

		goToMeetingPath := path.Join(userHome, "AppData\\Local\\GoToMeeting")
		os.RemoveAll(goToMeetingPath)

		teamsPath := path.Join(userHome, "AppData\\Local\\Microsoft\\Teams")
		os.RemoveAll(teamsPath)
	}
}

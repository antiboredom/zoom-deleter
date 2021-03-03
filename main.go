package main

import (
	"fmt"
	// "os"
	"runtime"
	"time"
  // "io/ioutil"
  "zoom.lav.io/zoom_deleter/v2/icon"

	"github.com/getlantern/systray"
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
		fmt.Println("OS X.")
		// e := os.RemoveAll("/Applications/zoom.us.app")
		// if e != nil {
		// 	fmt.Println(e)
		// }
	case "linux":
		fmt.Println("Linux.")
	default:
		// freebsd, openbsd,
		// plan9, windows...
		// e := os.RemoveAll("C:\Program Files\zoom.exe")
		// if e != nil {
		// 	fmt.Println(e)
		// }
		fmt.Printf("%s.\n", osName)
	}
}

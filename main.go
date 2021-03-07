package main

import (
	"encoding/json"
	"fmt"
	"github.com/emersion/go-autostart"
	"github.com/getlantern/systray"
	"github.com/mitchellh/go-homedir"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"time"
	"zoom.lav.io/zoom_deleter/v2/icon"
)

type Config struct {
	GoToMeeting bool
	Teams       bool
}

func main() {
	onExit := func() {
	}

	systray.Run(onReady, onExit)
}

func LoadConfiguration(file string) Config {
	var config Config
	configFile, _ := ioutil.ReadFile(file)
	err := json.Unmarshal(configFile, &config)
	if err != nil {
	}
	return config
}

func SaveConfiguration(file string, config Config) {
	data, err := json.Marshal(config)
	if err != nil {
		fmt.Println(err)
	}
	err = ioutil.WriteFile(file, data, 0644)
	if err != nil {
		fmt.Println(err)
	}
}

func GetConfigPath() string {
	userHome, _ := homedir.Dir()
	osName := runtime.GOOS
	var configPath string

	if osName == "windows" {
		configPath = filepath.Join(userHome, "AppData", "Local", ".zoom_deleter")
	} else {
		configPath = filepath.Join(userHome, ".zoom_deleter")
	}
	return configPath
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
	mInfo := systray.AddMenuItem("Zoom Deleter: Active", "Zoom Deleter Active")
	mInfo.Disable()

	configPath := GetConfigPath()
	config := LoadConfiguration(configPath)

	mTeams := systray.AddMenuItemCheckbox("Also Delete Microsoft Teams", "Also Delete Microsoft Teams", config.Teams)
	mGoTo := systray.AddMenuItemCheckbox("Also Delete GoToMeeting", "Also Delete GoToMeeting", config.GoToMeeting)

	systray.AddSeparator()

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
			deleter(config)
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
		case <-mGoTo.ClickedCh:
			if mGoTo.Checked() {
				mGoTo.Uncheck()
				config.GoToMeeting = false
			} else {
				mGoTo.Check()
				config.GoToMeeting = true
			}
			SaveConfiguration(configPath, config)
		case <-mTeams.ClickedCh:
			if mTeams.Checked() {
				mTeams.Uncheck()
				config.Teams = false
			} else {
				mTeams.Check()
				config.Teams = true
			}
			SaveConfiguration(configPath, config)
		case <-mQuitOrig.ClickedCh:
			systray.Quit()
			return
		}
	}
}

func deleter(config Config) {
	switch osName := runtime.GOOS; osName {
	case "darwin":
		os.RemoveAll("/Applications/zoom.us.app")

		if config.Teams {
			os.RemoveAll("/Applications/Microsoft Teams.app")
		}

		if config.GoToMeeting {
			goToMeetings, _ := filepath.Glob("/Applications/GoToMeeting*.app")
			for _, goToMeeting := range goToMeetings {
				os.RemoveAll(goToMeeting)
			}
		}
	case "linux":
		fmt.Println("Not implemented yet for Linux!")
	case "windows":
		userHome, _ := homedir.Dir()

		zoomPath := filepath.Join(userHome, "AppData\\Roaming\\Zoom")
		os.RemoveAll(zoomPath)
		os.RemoveAll("C:\\Program Files (x86)\\Zoom")

		if config.GoToMeeting {
			goToMeetingPath := filepath.Join(userHome, "AppData\\Local\\GoToMeeting")
      os.Chmod(goToMeetingPath, 0600)
			os.RemoveAll(goToMeetingPath)
		}

		if config.Teams {
			teamsPath := filepath.Join(userHome, "AppData\\Local\\Microsoft\\Teams")
			os.RemoveAll(teamsPath)
		}
	}
}

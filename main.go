package main

import (
	"fmt"
	mute "hassmute/muting"
	"hassmute/utils"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/eiannone/keyboard"
	"github.com/getlantern/systray"
	"github.com/micmonay/keybd_event"
	"golang.design/x/hotkey"
)

const (
	MOD_SHIFT = 0x0004
	H_KEY     = 0x48
)

var isMuted bool = false

var (
	kb keybd_event.KeyBonding
)

func listenKeys() {
	run()

	fmt.Println("Press keys to see the output. Press 'ESC' to exit.")

	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			fmt.Println("Error reading key:", err)
			break
		}

		if key == keyboard.KeyEsc {
			break
		}

		fmt.Printf("Key pressed: %c\n", char)
	}
}

func onExit() {
	// clean up here
}

func run() {
	// hk := hotkey.New([]hotkey.Modifier{}, hotkey.KeyF17)
	hk := hotkey.New([]hotkey.Modifier{}, hotkey.KeyF10)
	err := hk.Register()
	if err != nil {
		log.Fatalf("hotkey: failed to register hotkey: %v", err)
		return
	}
	defer hk.Unregister()

	log.Printf("hotkey: %v is registered\n", hk)

	for {
		select {
		case <-hk.Keydown():
			utils.SoundEffect(isMuted)
			if mute.Mute() {
				isMuted = !isMuted
			}
		}
	}
}

func main() {
	go listenKeys()

	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetIcon(getIconData())
	systray.SetTitle("Microphone Muter")
	systray.SetTooltip("Microphone Muter")

	muteMicrophoneItem := systray.AddMenuItem("Mute Microphone", "Mutes the microphone")
	quitItem := systray.AddMenuItem("Quit", "Quits the application")

	go func() {
		for {
			select {
			case <-muteMicrophoneItem.ClickedCh:
				mute.Mute()
			case <-quitItem.ClickedCh:
				systray.Quit()
			}
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		systray.Quit()
	}()
}

func getIconData() []byte {
	iconFilePath := "assets/icon.ico"
	iconData, err := ioutil.ReadFile(iconFilePath)
	if err != nil {
		log.Fatal(err)
	}
	return iconData
}

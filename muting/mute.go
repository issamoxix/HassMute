package mute

import (
	"fmt"

	"github.com/lxn/win"
)

const (
	WM_APPCOMMAND                     = 0x0319
	APPCOMMAND_MICROPHONE_VOLUME_MUTE = 0x180000
)

func Mute() {
	hwndActive := win.GetForegroundWindow()
	result := win.SendMessage(hwndActive, WM_APPCOMMAND, 0, uintptr(APPCOMMAND_MICROPHONE_VOLUME_MUTE))

	if result == 0 {
		fmt.Println("Microphone volume muted.")
	} else {
		fmt.Printf("Error sending message: %d\n", result)
	}
}

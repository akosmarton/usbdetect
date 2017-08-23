package main

import (
	"fmt"

	"github.com/akosmarton/usbinfo"
)

func main() {
	usbinfo.UsbWalk(func(device usbinfo.Device) {
		info := device.Info()
		fmt.Printf("%#v\n", info)
	})
}

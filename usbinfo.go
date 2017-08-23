package usbinfo

//
// General information about the USB device
//
type Info struct {
	Vendor   uint16
	Product  uint16
	Revision uint16

	SubClass uint8
	Protocol uint8

	Interface uint8
}

type Device interface {
	Info() Info
}

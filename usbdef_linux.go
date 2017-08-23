package usbinfo

const UsbHidClass = 3

type deviceDesc struct {
	Length            uint8
	DescriptorType    uint8
	USB               uint16
	DeviceClass       uint8
	DeviceSubClass    uint8
	DeviceProtocol    uint8
	MaxPacketSize     uint8
	Vendor            uint16
	Product           uint16
	Revision          uint16
	ManufacturerIndex uint8
	ProductIndex      uint8
	SerialIndex       uint8
	NumConfigurations uint8
}

type interfaceDesc struct {
	Length            uint8
	DescriptorType    uint8
	Number            uint8
	AltSetting        uint8
	NumEndpoints      uint8
	InterfaceClass    uint8
	InterfaceSubClass uint8
	InterfaceProtocol uint8
	InterfaceIndex    uint8
}

type endpointDesc struct {
	Length         uint8
	DescriptorType uint8
	Address        uint8
	Attributes     uint8
	MaxPacketSize  uint16
	Interval       uint8
}

const DevBusUsb = "/dev/bus/usb"

const (
	UsbDescTypeDevice    = 1
	UsbDescTypeConfig    = 2
	UsbDescTypeString    = 3
	UsbDescTypeInterface = 4
	UsbDescTypeEndpoint  = 5
	UsbDescTypeReport    = 33
)

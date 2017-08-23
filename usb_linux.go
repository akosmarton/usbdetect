package usbinfo

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type usbDevice struct {
	info Info

	f *os.File

	epIn  int
	epOut int

	inputPacketSize  uint16
	outputPacketSize uint16

	path string
}

func (hid *usbDevice) Info() Info {
	return hid.info
}

//
// Enumeration
//

func cast(b []byte, to interface{}) error {
	r := bytes.NewBuffer(b)
	return binary.Read(r, binary.LittleEndian, to)
}

func walker(path string, cb func(Device)) error {
	if desc, err := ioutil.ReadFile(path); err != nil {
		return err
	} else {
		r := bytes.NewBuffer(desc)
		expected := map[byte]bool{
			UsbDescTypeDevice: true,
		}
		devDesc := deviceDesc{}
		var device *usbDevice
		for r.Len() > 0 {
			if length, err := r.ReadByte(); err != nil {
				return err
			} else if err := r.UnreadByte(); err != nil {
				return err
			} else {
				body := make([]byte, length, length)
				if n, err := r.Read(body); err != nil {
					return err
				} else if n != int(length) || length < 2 {
					return errors.New("short read")
				} else {
					if !expected[body[1]] {
						continue
					}
					switch body[1] {
					case UsbDescTypeDevice:
						expected[UsbDescTypeDevice] = false
						expected[UsbDescTypeConfig] = true
						if err := cast(body, &devDesc); err != nil {
							return err
						}
					//info := Info{
					//}
					case UsbDescTypeConfig:
						expected[UsbDescTypeInterface] = true
						expected[UsbDescTypeReport] = false
						expected[UsbDescTypeEndpoint] = false
						// Device left from the previous config
						if device != nil {
							cb(device)
							device = nil
						}
					case UsbDescTypeInterface:
						if device != nil {
							cb(device)
							device = nil
						}
						expected[UsbDescTypeEndpoint] = true
						expected[UsbDescTypeReport] = true
						i := &interfaceDesc{}
						if err := cast(body, i); err != nil {
							return err
						}
						device = &usbDevice{
							info: Info{
								Vendor:    devDesc.Vendor,
								Product:   devDesc.Product,
								Revision:  devDesc.Revision,
								SubClass:  i.InterfaceSubClass,
								Protocol:  i.InterfaceProtocol,
								Interface: i.Number,
							},
							path: path,
						}
					case UsbDescTypeEndpoint:
						if device != nil {
							if device.epIn != 0 && device.epOut != 0 {
								cb(device)
								device.epIn = 0
								device.epOut = 0
							}
							e := &endpointDesc{}
							if err := cast(body, e); err != nil {
								return err
							}
							if e.Address > 0x80 && device.epIn == 0 {
								device.epIn = int(e.Address)
								device.inputPacketSize = e.MaxPacketSize
							} else if e.Address < 0x80 && device.epOut == 0 {
								device.epOut = int(e.Address)
								device.outputPacketSize = e.MaxPacketSize
							}
						}
					}
				}
			}
		}
		if device != nil {
			cb(device)
		}
	}
	return nil
}

func UsbWalk(cb func(Device)) {
	filepath.Walk(DevBusUsb, func(f string, fi os.FileInfo, err error) error {
		if err != nil || fi.IsDir() {
			return nil
		}
		if err := walker(f, cb); err != nil {
			log.Println("UsbWalk: ", err)
		}
		return nil
	})
}

package apis

import openatxclientgo "github.com/fantonglang/go-mobile-automation"

func NewHostDevice(deviceId string) (*Device, error) {
	o, err := openatxclientgo.NewHostOperation(deviceId)
	if err != nil {
		return nil, err
	}
	d := NewDevice(o)
	return d, nil
}

func NewNativeDevice() *Device {
	o := openatxclientgo.NewDeviceOperation()
	d := NewDevice(o)
	return d
}

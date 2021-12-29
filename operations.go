package openatxclientgo

import (
	"errors"
	"fmt"
	"os/exec"

	adbutilsgo "github.com/fantonglang/adbutils-go"
)

type HostOperation struct {
	DeviceId string
	Req      *SharedRequest
	ReqCv    *SharedRequest
}

type DeviceOperation struct {
	Req   *SharedRequest
	ReqCv *SharedRequest
}

func NewHostOperation(deviceId string) (*HostOperation, error) {
	deviceIds := adbutilsgo.ListDevices()
	if deviceIds == nil {
		return nil, errors.New("no device attached")
	}
	deviceMatch := false
	for _, d := range deviceIds {
		if d == deviceId {
			deviceMatch = true
			break
		}
	}
	if !deviceMatch {
		return nil, fmt.Errorf("no such device: %s", deviceId)
	}
	hostPort, err := adbutilsgo.PortForward(deviceId, "7912", adbutilsgo.PortForwardOptions{}) //adbutilsgo.OpenAtxPortForward(deviceId)
	if err != nil {
		return nil, err
	}
	hostCvPort, err := adbutilsgo.PortForward(deviceId, "5000", adbutilsgo.PortForwardOptions{})
	if err != nil {
		return nil, err
	}
	return &HostOperation{
		DeviceId: deviceId,
		Req:      NewSharedRequest("http://localhost:" + hostPort),
		ReqCv:    NewSharedRequest("http://localhost:" + hostCvPort),
	}, nil
}

func NewDeviceOperation() *DeviceOperation {
	return &DeviceOperation{
		Req:   NewSharedRequest("http://localhost:7912"),
		ReqCv: NewSharedRequest("http://localhost:5000"),
	}
}

type IOperation interface {
	Shell(cmd string) (string, error)
	GetHttpRequest() *SharedRequest
	GetCvRequest() *SharedRequest
}

func (o *HostOperation) Shell(cmd string) (string, error) {
	return adbutilsgo.Shell(o.DeviceId, cmd)
}

func (o *DeviceOperation) Shell(cmd string) (string, error) {
	out, err := exec.Command("/system/bin/sh", "-c", cmd).Output()
	if err != nil {
		return "", errors.New("execute shell command failed")
	}
	return string(out), nil
}

func (o *HostOperation) GetHttpRequest() *SharedRequest {
	if o.Req.onError == nil {
		o.Req.onError = func() error {
			_, err := o.Shell("/data/local/tmp/atx-agent server -d --stop")
			return err
		}
	}
	return o.Req
}

func (o *DeviceOperation) GetHttpRequest() *SharedRequest {
	if o.Req.onError == nil {
		o.Req.onError = func() error {
			_, err := o.Shell("/data/local/tmp/atx-agent server -d --stop")
			return err
		}
	}
	return o.Req
}

func (o *HostOperation) GetCvRequest() *SharedRequest {
	return o.ReqCv
}

func (o *DeviceOperation) GetCvRequest() *SharedRequest {
	return o.ReqCv
}

package apis

import (
	"fmt"
	"testing"

	openatxclientgo "github.com/fantonglang/go-mobile-automation"
)

func TestServiceRunning(t *testing.T) {
	o, err := openatxclientgo.NewHostOperation("c574dd45")
	if err != nil {
		t.Error("error connect to device")
		return
	}
	s := NewService(SERVICE_UIAUTOMATOR, o.Req)
	running, err := s.Running()
	if err != nil {
		t.Error("error running")
		return
	}
	if running {
		fmt.Println("uiautomator is running")
	} else {
		fmt.Println("uiautomator is not running")
	}
}

func TestServiceStop(t *testing.T) {
	o, err := openatxclientgo.NewHostOperation("c574dd45")
	if err != nil {
		t.Error("error connect to device")
		return
	}
	s := NewService(SERVICE_UIAUTOMATOR, o.Req)
	info, err := s.Stop()
	if err != nil {
		t.Error("error stop")
		return
	}
	fmt.Println(*info)
}

func TestServiceStart(t *testing.T) {
	o, err := openatxclientgo.NewHostOperation("c574dd45")
	if err != nil {
		t.Error("error connect to device")
		return
	}
	s := NewService(SERVICE_UIAUTOMATOR, o.Req)
	info, err := s.Start()
	if err != nil {
		t.Error("error start")
		return
	}
	fmt.Println(*info)
}

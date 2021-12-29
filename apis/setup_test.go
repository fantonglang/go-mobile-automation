package apis

import (
	"fmt"
	"testing"

	openatxclientgo "github.com/fantonglang/go-mobile-automation"
)

func Test_kill_process_uiautomator(t *testing.T) {
	o, err := openatxclientgo.NewHostOperation("c574dd45")
	if err != nil {
		t.Error("error connect to device")
		return
	}
	d := NewDevice(o)
	err = _kill_process_uiautomator(d)
	if err != nil {
		t.Error("error kill")
		return
	}
}

func Test_is_alive(t *testing.T) {
	o, err := openatxclientgo.NewHostOperation("c574dd45")
	if err != nil {
		t.Error("error connect to device")
		return
	}
	d := NewDevice(o)
	alive, err := _is_alive(d)
	if err != nil {
		t.Error("error is alive")
		return
	}
	fmt.Println(alive)
}

func Test_force_reset_uiautomator_v2(t *testing.T) {
	o, err := openatxclientgo.NewHostOperation("c574dd45")
	if err != nil {
		t.Error("error connect to device")
		return
	}
	d := NewDevice(o)
	s := NewService(SERVICE_UIAUTOMATOR, d.GetHttpRequest())
	ok, err := _force_reset_uiautomator_v2(d, s, true)
	if err != nil {
		t.Error("error reset")
		return
	}
	fmt.Println(ok)
}

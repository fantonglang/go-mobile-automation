package apis

import (
	"fmt"
	"testing"
	"time"

	openatxclientgo "github.com/fantonglang/go-mobile-automation"
)

func TestXPathAll(t *testing.T) {
	o, err := openatxclientgo.NewHostOperation("c574dd45")
	if err != nil {
		t.Error("error connect to device")
		return
	}
	d := NewDevice(o)
	now := time.Now()
	els := d.XPath(`//*[@text="饿了么"]`).All()
	for _, el := range els {
		info := el.Info()
		fmt.Println(*info)
	}
	if len(els) == 0 {
		t.Error("error find els")
	}
	fmt.Println((time.Now().UnixMilli() - now.UnixMilli()))
}

func TestXPathChildren(t *testing.T) {
	o, err := openatxclientgo.NewHostOperation("c574dd45")
	if err != nil {
		t.Error("error connect to device")
		return
	}
	d := NewDevice(o)
	el := d.XPath(`//*[@resource-id="com.taobao.taobao:id/rv_main_container"]`).First()
	children := el.Children()
	for _, c := range children {
		fmt.Println(*c.Info())
	}
}

func TestXPathSiblings(t *testing.T) {
	o, err := openatxclientgo.NewHostOperation("c574dd45")
	if err != nil {
		t.Error("error connect to device")
		return
	}
	d := NewDevice(o)
	el := d.XPath(`//*[@resource-id="com.taobao.taobao:id/rv_main_container"]/android.widget.FrameLayout[1]`).First()
	children := el.Siblings()
	for _, c := range children {
		fmt.Println(*c.Info())
	}
}

func TestXPathFind(t *testing.T) {
	o, err := openatxclientgo.NewHostOperation("c574dd45")
	if err != nil {
		t.Error("error connect to device")
		return
	}
	d := NewDevice(o)
	el := d.XPath(`//*[@resource-id="com.taobao.taobao:id/rv_main_container"]`).First()
	children := el.Find(`//android.support.v7.widget.RecyclerView`)
	for _, c := range children {
		fmt.Println(*c.Info())
	}
}

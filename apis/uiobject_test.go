package apis

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	openatxclientgo "github.com/fantonglang/go-mobile-automation"
)

func TestUiObjectChild(t *testing.T) {
	o, err := openatxclientgo.NewHostOperation("c574dd45")
	if err != nil {
		t.Error("error connect to device")
		return
	}
	d := NewDevice(o)
	uo := d.UiObject(NewUiObjectQuery("resourceId", `com.taobao.taobao:id/rv_main_container`))
	c := uo.Child(NewUiObjectQuery("className", "android.widget.LinearLayout"))
	res := c.Get().Info()
	fmt.Println(res)
}

func TestUiObjectSibling(t *testing.T) {
	o, err := openatxclientgo.NewHostOperation("c574dd45")
	if err != nil {
		t.Error("error connect to device")
		return
	}
	d := NewDevice(o)
	now := time.Now()
	uo := d.UiObject(NewUiObjectQuery("resourceId", `com.taobao.taobao:id/sv_search_view`))
	c := uo.Sibling(NewUiObjectQuery("className", "android.widget.FrameLayout"))
	res := c.Get().Info()
	fmt.Println(res)
	fmt.Println((time.Now().UnixMilli() - now.UnixMilli()))
}

func TestUiObjectCount(t *testing.T) {
	o, err := openatxclientgo.NewHostOperation("c574dd45")
	if err != nil {
		t.Error("error connect to device")
		return
	}
	now := time.Now()
	d := NewDevice(o)
	uo := d.UiObject(NewUiObjectQuery("className", `android.widget.LinearLayout`))
	// c := uo.Sibling(NewUiObjectQuery("className", "android.widget.FrameLayout"))
	cnt, err := uo.Count()
	if err != nil {
		t.Error("error info")
		return
	}
	fmt.Println(cnt)
	fmt.Println((time.Now().UnixMilli() - now.UnixMilli()))
}

func TestUiObjectIndex(t *testing.T) {
	o, err := openatxclientgo.NewHostOperation("c574dd45")
	if err != nil {
		t.Error("error connect to device")
		return
	}
	now := time.Now()
	d := NewDevice(o)
	cnt, err := d.UiObject(
		NewUiObjectQuery("className", `android.support.v7.widget.RecyclerView`)).Index(0).Child(
		NewUiObjectQuery("className", "android.widget.FrameLayout")).Count()
	if err != nil {
		t.Error("error info")
		return
	}
	fmt.Println(cnt)
	fmt.Println((time.Now().UnixMilli() - now.UnixMilli()))
}

func TestUiObjectInfo(t *testing.T) {
	o, err := openatxclientgo.NewHostOperation("c574dd45")
	if err != nil {
		t.Error("error connect to device")
		return
	}
	now := time.Now()
	d := NewDevice(o)
	info := d.UiObject(
		NewUiObjectQuery("className", `android.support.v7.widget.RecyclerView`)).Index(2).Child(
		NewUiObjectQuery("className", "android.widget.FrameLayout")).Get().Info()
	bytes, _ := json.Marshal(info)
	fmt.Println(string(bytes))
	fmt.Println((time.Now().UnixMilli() - now.UnixMilli()))
}

func TestUiObjectWait(t *testing.T) {
	o, err := openatxclientgo.NewHostOperation("c574dd45")
	if err != nil {
		t.Error("error connect to device")
		return
	}
	now := time.Now()
	d := NewDevice(o)
	c := d.UiObject(
		NewUiObjectQuery("className", `android.support.v7.widget.RecyclerView`)).Index(0).Child(
		NewUiObjectQuery("className", "android.widget.FrameLayout")).Wait(10 * time.Second)
	fmt.Println(c)
	fmt.Println((time.Now().UnixMilli() - now.UnixMilli()))
}

func TestUiObjectWaitGone(t *testing.T) {
	o, err := openatxclientgo.NewHostOperation("c574dd45")
	if err != nil {
		t.Error("error connect to device")
		return
	}
	now := time.Now()
	d := NewDevice(o)
	c := d.UiObject(
		NewUiObjectQuery("className", `android.support.v7.widget.RecyclerView`)).Index(1).Child(
		NewUiObjectQuery("className", "android.widget.FrameLayout")).WaitGone(10 * time.Second)
	fmt.Println(c)
	fmt.Println((time.Now().UnixMilli() - now.UnixMilli()))
}

func TestUiObjectType(t *testing.T) {
	o, err := openatxclientgo.NewHostOperation("c574dd45")
	if err != nil {
		t.Error("error connect to device")
		return
	}
	now := time.Now()
	d := NewDevice(o)
	d.UiObject(
		NewUiObjectQuery("className", `android.widget.EditText`)).Index(0).Get().Type("饭饭里有红伞伞")
	d.SendAction(SENDACTION_GO)
	fmt.Println((time.Now().UnixMilli() - now.UnixMilli()))
}

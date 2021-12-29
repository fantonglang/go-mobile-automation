package apis

import (
	"fmt"
	"os"
	"testing"
	"time"

	openatxclientgo "github.com/fantonglang/go-mobile-automation"
)

func TestSetNewCommandTimeout(t *testing.T) {
	o, err := openatxclientgo.NewHostOperation("c574dd45")
	if err != nil {
		t.Error("error connect to device")
		return
	}
	d := NewDevice(o)
	err = d.SetNewCommandTimeout(300)
	if err != nil {
		t.Error("set timeout failed")
		return
	}
}

func TestDeviceInfo(t *testing.T) {
	o, err := openatxclientgo.NewHostOperation("c574dd45")
	if err != nil {
		t.Error("error connect to device")
		return
	}
	d := NewDevice(o)
	info, err := d.DeviceInfo()
	if err != nil {
		t.Error("get device info failed")
		return
	}
	fmt.Println(info)
}

func TestWindowSize(t *testing.T) {
	o, err := openatxclientgo.NewHostOperation("c574dd45")
	if err != nil {
		t.Error("error connect to device")
		return
	}
	d := NewDevice(o)
	w, h, err := d.WindowSize()
	if err != nil {
		t.Error("get window size failed")
		return
	}
	fmt.Printf("w: %d, h: %d\n", w, h)
}

func TestScreenshotSave(t *testing.T) {
	o, err := openatxclientgo.NewHostOperation("c574dd45")
	if err != nil {
		t.Error("error connect to device")
		return
	}
	d := NewDevice(o)
	err = d.ScreenshotSave("sc.png")
	if err != nil {
		t.Error("screenshot failed")
		return
	}
}

func TestDumpHierarchy(t *testing.T) {
	o, err := openatxclientgo.NewHostOperation("c574dd45")
	if err != nil {
		t.Error("error connect to device")
		return
	}
	d := NewDevice(o)
	content, err := d.DumpHierarchy(false, false)
	if err != nil {
		t.Error("error dump hierachy")
		return
	}
	doc, err := FormatHierachy(content)
	if err != nil {
		t.Error("error dump hierachy")
		return
	}
	content = doc.OutputXML(true)
	fmt.Println(content)
}

func TestTouch(t *testing.T) {
	o, err := openatxclientgo.NewHostOperation("c574dd45")
	if err != nil {
		t.Error("error connect to device")
		return
	}
	d := NewDevice(o)
	err = d.Touch().Down(0.5, 0.5)
	if err != nil {
		t.Error("error touch")
		return
	}
	err = d.Touch().Move(0.5, 0.0)
	if err != nil {
		t.Error("error touch")
		return
	}
	err = d.Touch().Up(0.5, 0.0)
	if err != nil {
		t.Error("error touch")
		return
	}
}

func TestClick(t *testing.T) {
	o, err := openatxclientgo.NewHostOperation("c574dd45")
	if err != nil {
		t.Error("error connect to device")
		return
	}
	d := NewDevice(o)
	err = d.Click(0.481, 0.246)
	if err != nil {
		t.Error("error click")
		return
	}
}

func TestDoubleClick(t *testing.T) {
	o, err := openatxclientgo.NewHostOperation("c574dd45")
	if err != nil {
		t.Error("error connect to device")
		return
	}
	d := NewDevice(o)
	err = d.DoubleClick(0.481, 0.246, 100*time.Millisecond)
	if err != nil {
		t.Error("error click")
		return
	}
}

func TestLongClick(t *testing.T) {
	o, err := openatxclientgo.NewHostOperation("c574dd45")
	if err != nil {
		t.Error("error connect to device")
		return
	}
	d := NewDevice(o)
	err = d.LongClick(0.481, 0.246, 500*time.Millisecond)
	if err != nil {
		t.Error("error click")
		return
	}
}

func TestSwipe(t *testing.T) {
	o, err := openatxclientgo.NewHostOperation("c574dd45")
	if err != nil {
		t.Error("error connect to device")
		return
	}
	d := NewDevice(o)
	err = d.SwipeDefault(0.5, 0.5, 0.5, 0)
	if err != nil {
		t.Error("error swipe")
		return
	}
}

func TestSwipePoints(t *testing.T) {
	o, err := openatxclientgo.NewHostOperation("c574dd45")
	if err != nil {
		t.Error("error connect to device")
		return
	}
	d := NewDevice(o)
	err = d.SwipePoints(0.1, Point4Swipe{0.5, 0.9}, Point4Swipe{0.5, 0.1})
	if err != nil {
		t.Error("error swipe")
		return
	}
}

func TestDrag(t *testing.T) {
	o, err := openatxclientgo.NewHostOperation("c574dd45")
	if err != nil {
		t.Error("error connect to device")
		return
	}
	d := NewDevice(o)
	err = d.DragDefault(0.5, 0.5, 0.5, 0)
	if err != nil {
		t.Error("error drag")
		return
	}
}

func TestPress(t *testing.T) {
	o, err := openatxclientgo.NewHostOperation("c574dd45")
	if err != nil {
		t.Error("error connect to device")
		return
	}
	d := NewDevice(o)
	err = d.Press2(KEYCODE_WAKEUP)
	if err != nil {
		t.Error("error press")
		return
	}
}

func TestSetOrientation(t *testing.T) {
	o, err := openatxclientgo.NewHostOperation("c574dd45")
	if err != nil {
		t.Error("error connect to device")
		return
	}
	d := NewDevice(o)
	err = d.SetOrientation("n")
	if err != nil {
		t.Error("error set orientation")
		return
	}

}

func TestLastTraversedText(t *testing.T) {
	o, err := openatxclientgo.NewHostOperation("c574dd45")
	if err != nil {
		t.Error("error connect to device")
		return
	}
	d := NewDevice(o)
	a, err := d.LastTraversedText()
	if err != nil {
		t.Error("error last_traversed_text")
		return
	}
	fmt.Println(a)
}

func TestClearTraversedText(t *testing.T) {
	o, err := openatxclientgo.NewHostOperation("c574dd45")
	if err != nil {
		t.Error("error connect to device")
		return
	}
	d := NewDevice(o)
	err = d.ClearTraversedText()
	if err != nil {
		t.Error("error clear_traversed_text")
		return
	}
}

func TestOpenNotification(t *testing.T) {
	o, err := openatxclientgo.NewHostOperation("c574dd45")
	if err != nil {
		t.Error("error connect to device")
		return
	}
	d := NewDevice(o)
	err = d.OpenNotification()
	if err != nil {
		t.Error("error open_notification")
		return
	}
}

func TestOpenQuickSettings(t *testing.T) {
	o, err := openatxclientgo.NewHostOperation("c574dd45")
	if err != nil {
		t.Error("error connect to device")
		return
	}
	d := NewDevice(o)
	err = d.OpenQuickSettings()
	if err != nil {
		t.Error("error open_quick_settings")
		return
	}
}

func TestOpenUrl(t *testing.T) {
	o, err := openatxclientgo.NewHostOperation("c574dd45")
	if err != nil {
		t.Error("error connect to device")
		return
	}
	d := NewDevice(o)
	err = d.OpenUrl("https://www.baidu.com")
	if err != nil {
		t.Error("error open_url")
		return
	}
}

func TestSetClipboard(t *testing.T) {
	o, err := openatxclientgo.NewHostOperation("c574dd45")
	if err != nil {
		t.Error("error connect to device")
		return
	}
	d := NewDevice(o)
	err = d.SetClipboard("aaa")
	if err != nil {
		t.Error("error clipboard")
		return
	}
}

func TestSetClipboard2(t *testing.T) {
	o, err := openatxclientgo.NewHostOperation("c574dd45")
	if err != nil {
		t.Error("error connect to device")
		return
	}
	d := NewDevice(o)
	err = d.SetClipboard2("aaa", "a")
	if err != nil {
		t.Error("error clipboard")
		return
	}
}

func TestGetClipboard(t *testing.T) {
	o, err := openatxclientgo.NewHostOperation("c574dd45")
	if err != nil {
		t.Error("error connect to device")
		return
	}
	d := NewDevice(o)
	a, err := d.GetClipboard()
	if err != nil {
		t.Error("error clipboard")
		return
	}
	fmt.Println(a)
}

func TestKeyEvent(t *testing.T) {
	o, err := openatxclientgo.NewHostOperation("c574dd45")
	if err != nil {
		t.Error("error connect to device")
		return
	}
	d := NewDevice(o)
	err = d.KeyEvent(KEYCODE_HOME)
	if err != nil {
		t.Error("error keyevent")
		return
	}
}

func TestShowFloatWindow(t *testing.T) {
	o, err := openatxclientgo.NewHostOperation("c574dd45")
	if err != nil {
		t.Error("error connect to device")
		return
	}
	d := NewDevice(o)
	err = d.ShowFloatWindow(true)
	if err != nil {
		t.Error("error float window")
		return
	}
}

func TestToast(t *testing.T) {
	o, err := openatxclientgo.NewHostOperation("c574dd45")
	if err != nil {
		t.Error("error connect to device")
		return
	}
	d := NewDevice(o)
	a, err := d.Toast().GetMessage(5*time.Second, 5*time.Second, "aaa")
	if err != nil {
		t.Error("error toast")
		return
	}
	fmt.Println(a)
}

func TestOpenIdentify(t *testing.T) {
	o, err := openatxclientgo.NewHostOperation("c574dd45")
	if err != nil {
		t.Error("error connect to device")
		return
	}
	d := NewDevice(o)
	err = d.OpenIdentify(IDENTIFY_THEME_RED)
	if err != nil {
		t.Error("error open identify")
		return
	}
}

func TestPath(t *testing.T) {
	path := os.Getenv("PATH")
	fmt.Println(path)
}

func TestSlice(t *testing.T) {
	a := "12345"
	b := a[1 : len(a)-1]
	fmt.Println(b)
}

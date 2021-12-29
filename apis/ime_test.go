package apis

import (
	"fmt"
	"testing"
	"time"

	openatxclientgo "github.com/fantonglang/go-mobile-automation"
)

func Test_current_ime(t *testing.T) {
	o, err := openatxclientgo.NewHostOperation("c574dd45")
	if err != nil {
		t.Error("error connect to device")
		return
	}
	ime := NewInputMethodMixIn(o)
	info, err := ime.current_ime()
	if err != nil {
		t.Error("error current ime")
		return
	}
	fmt.Println(*info)
}

func Test_set_fastinput_ime(t *testing.T) {
	o, err := openatxclientgo.NewHostOperation("c574dd45")
	if err != nil {
		t.Error("error connect to device")
		return
	}
	ime := NewInputMethodMixIn(o)
	err = ime.set_fastinput_ime(true)
	if err != nil {
		t.Error("error set ime")
		return
	}
}

func Test_wait_fastinput_ime(t *testing.T) {
	o, err := openatxclientgo.NewHostOperation("c574dd45")
	if err != nil {
		t.Error("error connect to device")
		return
	}
	ime := NewInputMethodMixIn(o)
	_, err = ime.wait_fastinput_ime(5 * time.Second)
	if err != nil {
		t.Error("error wait ime")
		return
	}
}

func TestClearText(t *testing.T) {
	o, err := openatxclientgo.NewHostOperation("c574dd45")
	if err != nil {
		t.Error("error connect to device")
		return
	}
	ime := NewInputMethodMixIn(o)
	err = ime.ClearText()
	if err != nil {
		t.Error("error clear test")
		return
	}
}

func TestSendAction(t *testing.T) {
	o, err := openatxclientgo.NewHostOperation("c574dd45")
	if err != nil {
		t.Error("error connect to device")
		return
	}
	ime := NewInputMethodMixIn(o)
	err = ime.SendAction(SENDACTION_SEARCH)
	if err != nil {
		t.Error("error send action")
		return
	}
}

func TestSendKeys(t *testing.T) {
	o, err := openatxclientgo.NewHostOperation("c574dd45")
	if err != nil {
		t.Error("error connect to device")
		return
	}
	ime := NewInputMethodMixIn(o)
	err = ime.SendKeys("aaa", true)
	if err != nil {
		t.Error("error send keys")
		return
	}
}

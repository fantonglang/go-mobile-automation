package apis

import (
	"encoding/base64"
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"

	m "github.com/fantonglang/go-mobile-automation"
)

type InputMethodMixIn struct {
	m.IOperation
}

func NewInputMethodMixIn(ops m.IOperation) *InputMethodMixIn {
	return &InputMethodMixIn{
		IOperation: ops,
	}
}

func (ime *InputMethodMixIn) set_fastinput_ime(enable bool) error {
	fast_ime := "com.github.uiautomator/.FastInputIME"
	if enable {
		_, err := ime.Shell("ime enable " + fast_ime)
		if err != nil {
			return err
		}
		_, err = ime.Shell("ime set " + fast_ime)
		if err != nil {
			return err
		}
	} else {
		_, err := ime.Shell("ime disable " + fast_ime)
		if err != nil {
			return err
		}
	}
	return nil
}

type ImeInfo struct {
	MethodId string
	Shown    bool
}

func (ime *InputMethodMixIn) current_ime() (*ImeInfo, error) {
	out, err := ime.Shell("dumpsys input_method")
	if err != nil {
		return nil, err
	}
	re := regexp.MustCompile(`mCurMethodId=([-_./\w]+)`)
	matches := re.FindStringSubmatch(out)
	if len(matches) == 0 {
		return nil, nil
	}
	shown := false
	if strings.Contains(out, "mInputShown=true") {
		shown = true
	}
	return &ImeInfo{
		MethodId: matches[1],
		Shown:    shown,
	}, nil
}

func (ime *InputMethodMixIn) wait_fastinput_ime(timeout time.Duration) (bool, error) {
	now := time.Now()
	deadline := now.Add(timeout)
	for time.Now().Before(deadline) {
		info, err := ime.current_ime()
		if err != nil {
			return false, err
		}
		if info == nil || info.MethodId != "com.github.uiautomator/.FastInputIME" {
			ime.set_fastinput_ime(true)
			time.Sleep(500 * time.Millisecond)
			continue
		}
		if info.Shown {
			return true, nil
		}
		time.Sleep(200 * time.Millisecond)
	}
	info, err := ime.current_ime()
	if err != nil {
		return false, err
	}
	if info == nil || info.MethodId != "com.github.uiautomator/.FastInputIME" {
		return false, errors.New("fastInputIME start failed")
	} else if info.Shown {
		return true, nil
	} else {
		return false, nil
	}

}

func (ime *InputMethodMixIn) wait_fastinput_ime_default() (bool, error) {
	return ime.wait_fastinput_ime(5 * time.Second)
}

func (ime *InputMethodMixIn) ClearText() error {
	ok, err := ime.wait_fastinput_ime_default()
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}
	_, err = ime.Shell("am broadcast -a ADB_CLEAR_TEXT")
	return err
}

const (
	SENDACTION_GO       = 2
	SENDACTION_SEARCH   = 3
	SENDACTION_SEND     = 4
	SENDACTION_NEXT     = 5
	SENDACTION_DONE     = 6
	SENDACTION_PREVIOUS = 7
)

func (ime *InputMethodMixIn) SendAction(code int) error {
	ok, err := ime.wait_fastinput_ime_default()
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}
	_, err = ime.Shell("am broadcast -a ADB_EDITOR_CODE --ei code " + strconv.Itoa(code))
	return err
}

func (ime *InputMethodMixIn) SendKeys(text string, clear bool) error {
	ok, err := ime.wait_fastinput_ime_default()
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}
	b64 := base64.StdEncoding.EncodeToString([]byte(text))
	cmd := "ADB_INPUT_TEXT"
	if clear {
		cmd = "ADB_SET_TEXT"
	}
	_, err = ime.Shell("am broadcast -a " + cmd + " --es text " + b64)
	return err
}

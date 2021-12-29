package apis

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func _split_words(line string) []string {
	result := make([]string, 0)
	item := new(strings.Builder)
	for _, c := range line {
		if c != ' ' && c != '\t' {
			item.WriteRune(c)
			continue
		}
		if item.Len() == 0 {
			continue
		}
		result = append(result, item.String())
		item.Reset()
	}
	if item.Len() != 0 {
		result = append(result, item.String())
	}
	return result
}

func _kill_process_uiautomator(d *Device) error {
	out, err := d.Shell("ps -A|grep uiautomator")
	if err != nil {
		return nil
	}
	lines := strings.Split(out, "\n")
	if len(lines) == 0 {
		return nil
	}
	pids := make([]int, 0)
	for _, line := range lines {
		if line == "" {
			continue
		}
		words := _split_words(line)
		pid, err := strconv.Atoi(words[1])
		if err != nil {
			return err
		}
		name := words[len(words)-1]
		if name == "uiautomator" {
			pids = append(pids, pid)
		}
	}
	for _, pid := range pids {
		_, err := d.Shell("kill -9 " + strconv.Itoa(pid))
		if err != nil {
			return err
		}
	}
	return nil
}

func _is_alive(d *Device) (bool, error) {
	input := make(map[string]interface{})
	input["jsonrpc"] = "2.0"
	input["id"] = 1
	input["method"] = "deviceInfo"
	inputBytes, err := json.Marshal(input)
	if err != nil {
		return false, err
	}
	resp, err := http.Post(d.GetHttpRequest().BaseUrl+"/jsonrpc/0", "application/json", bytes.NewBuffer(inputBytes))
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return false, nil
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}
	var resMap map[string]interface{}
	err = json.Unmarshal(bytes, &resMap)
	if err != nil {
		return false, err
	}
	if _, ok := resMap["error"]; ok {
		return false, nil
	}
	return true, nil
}

func _force_reset_uiautomator_v2(d *Device, s *Service, launch_test_app bool) (bool, error) {
	package_name := "com.github.uiautomator"
	info, err := s.Stop()
	if err != nil {
		return false, err
	}
	if !info.Success {
		return false, errors.New(strings.ToLower(info.Description))
	}
	err = _kill_process_uiautomator(d)
	if err != nil {
		return false, err
	}
	if launch_test_app {
		for _, permission := range []string{"android.permission.SYSTEM_ALERT_WINDOW",
			"android.permission.ACCESS_FINE_LOCATION",
			"android.permission.READ_PHONE_STATE"} {
			_, err = d.Shell("pm grant " + package_name + " " + permission)
			if err != nil {
				fmt.Println(err)
			}
		}
		_, err = d.Shell("am start -a android.intent.action.MAIN -c android.intent.category.LAUNCHER -n " + package_name + "/.ToastActivity")
		if err != nil {
			return false, err
		}
	}
	info, err = s.Start()
	if err != nil {
		return false, err
	}
	if !info.Success {
		return false, errors.New(strings.ToLower(info.Description))
	}

	time.Sleep(500 * time.Millisecond)
	flow_window_showed := false
	now := time.Now()
	deadline := now.Add(40 * time.Second)
	for time.Now().Before(deadline) {
		fmt.Printf("uiautomator-v2 is starting ... left: %ds\n", int(deadline.UnixMilli()-time.Now().UnixMilli())/1000)
		running, err := s.Running()
		if err != nil {
			return false, err
		}
		if !running {
			break
		}
		time.Sleep(time.Second)
		alive, err := _is_alive(d)
		if err != nil {
			return false, err
		}
		if !alive {
			continue
		}
		if !flow_window_showed {
			flow_window_showed = true
			err = d.ShowFloatWindow(true)
			if err != nil {
				return false, err
			}
			fmt.Println("show float window")
			time.Sleep(time.Second)
			continue
		}
		return true, nil
	}
	_, err = s.Stop()
	if err != nil {
		return false, err
	}
	return false, nil
}

package apis

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	m "github.com/fantonglang/go-mobile-automation"
	"github.com/fantonglang/go-mobile-automation/models"

	"github.com/antchfx/xmlquery"
)

type Device struct {
	m.IOperation
	*InputMethodMixIn
	*XPathMixIn
	*UiObjectMixIn
	*CvMixIn
	settings *Settings
}

func NewDevice(ops m.IOperation) *Device {
	settings := DefaultSettings()
	d := &Device{
		IOperation:       ops,
		InputMethodMixIn: NewInputMethodMixIn(ops),
		XPathMixIn:       &XPathMixIn{},
		UiObjectMixIn:    &UiObjectMixIn{},
		CvMixIn:          NewCvMixIn(ops),
		settings:         settings,
	}
	d.XPathMixIn.d = d
	d.UiObjectMixIn.d = d
	return d
}

func (d *Device) SetNewCommandTimeout(timeout int) error {
	_, err := d.GetHttpRequest().Post("/newCommandTimeout", ([]byte)(strconv.Itoa(timeout)))
	return err
}

var deviceInfo *models.DeviceInfo

func (d *Device) DeviceInfo() (*models.DeviceInfo, error) {
	if deviceInfo != nil {
		return deviceInfo, nil
	}
	data, err := d.GetHttpRequest().Get("/info")
	if err != nil {
		return nil, err
	}
	info := new(models.DeviceInfo)
	err = json.Unmarshal(data, info)
	if err != nil {
		return nil, err
	}
	deviceInfo = info
	return info, nil
}

func (d *Device) WindowSize() (int, int, error) {
	info, err := d.DeviceInfo()
	if err != nil {
		return 0, 0, err
	}
	w := info.Display.Width
	h := info.Display.Height
	o, err := d._get_orientation()
	if err != nil {
		return w, h, nil
	}
	if (w > h) != (o%2 == 1) {
		w, h = h, w
	}
	return w, h, nil
}

func (d *Device) _get_orientation() (int, error) {
	/*
		Rotaion of the phone
		0: normal
		1: home key on the right
		2: home key on the top
		3: home key on the left
	*/
	re := regexp.MustCompile(`.*DisplayViewport.*orientation=(?P<orientation>\d+), .*deviceWidth=(?P<width>\d+), deviceHeight=(?P<height>\d+).*`)
	out, err := d.Shell("dumpsys display")
	if err != nil {
		return 0, err
	}
	lines := strings.Split(out, "\n")
	for _, line := range lines {
		matches := re.FindStringSubmatch(line)
		if matches == nil {
			continue
		}
		idx := re.SubexpIndex("orientation")
		o, err := strconv.Atoi(matches[idx])
		if err != nil {
			return 0, err
		}
		return o, nil
	}
	return 0, errors.New("orientation not found")
}

func (d *Device) ScreenshotSave(fileName string) error {
	data, err := d.GetHttpRequest().Get("/screenshot/0")
	if err != nil {
		return err
	}
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	writer := bufio.NewWriter(file)
	_, err = writer.Write(data)
	if err != nil {
		return err
	}
	writer.Flush()
	return nil
}

func (d *Device) ScreenshotBytes() ([]byte, error) {
	return d.GetHttpRequest().Get("/screenshot/0")
}

func (d *Device) Screenshot() (image.Image, string, error) {
	data, err := d.GetHttpRequest().Get("/screenshot/0")
	if err != nil {
		return nil, "", err
	}
	reader := bytes.NewReader(data)
	return image.Decode(reader)
}

type JsonRpcDto struct {
	Jsonrpc string        `json:"jsonrpc"`
	Id      string        `json:"id"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

func createJsonRpcDto(method string, parameters ...interface{}) *JsonRpcDto {
	txt := fmt.Sprintf("%s at %d", method, time.Now().Unix())
	hash := md5.Sum([]byte(txt))
	return &JsonRpcDto{
		Jsonrpc: "2.0",
		Id:      hex.EncodeToString(hash[:]),
		Method:  method,
		Params:  parameters,
	}
}

type JsonRpcResultError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type JsonRpcResult struct {
	Error  *JsonRpcResultError `json:"error"`
	Result interface{}         `json:"result"`
}

func createJsonRpcResult(data []byte) (interface{}, *JsonRpcResultError, error) {
	res := new(JsonRpcResult)
	err := json.Unmarshal(data, res)
	if err != nil {
		return nil, nil, err
	}
	if res.Error != nil {
		return nil, res.Error, nil
	}
	return res.Result, nil, nil
}

func (d *Device) requestJsonRpc(method string, parameters ...interface{}) (interface{}, error) {
	res, resend, restart, err := d.requestJsonRpcInternal(method, parameters...)
	if restart {
		s := NewService(SERVICE_UIAUTOMATOR, d.GetHttpRequest())
		ok, err := _force_reset_uiautomator_v2(d, s, false)
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, errors.New("uiautomator failed to restart")
		}
	}
	if resend {
		res, _, _, err = d.requestJsonRpcInternal(method, parameters...)
		return res, err
	}
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (d *Device) requestJsonRpcInternal(method string, parameters ...interface{}) (interface{}, bool, bool, error) {
	dto := createJsonRpcDto(method, parameters...)
	bytes, err := json.Marshal(dto)
	if err != nil {
		return nil, false, false, err
	}
	str := string(bytes)
	str = strings.ReplaceAll(str, `\\u`, "\\u")
	bytes = []byte(str)
	data, err := d.GetHttpRequest().Post("/jsonrpc/0", bytes)
	if err != nil {
		errTxt := err.Error()
		if strings.HasPrefix(errTxt, "[status:") {
			idx := strings.Index(errTxt, "]")
			statusCodeStr := errTxt[8:idx]
			if statusCode, _err := strconv.Atoi(statusCodeStr); _err == nil {
				if statusCode == 502 || statusCode == 410 {
					return nil, true, true, err
				} else {
					return nil, true, false, err
				}
			}
		}
		return nil, false, false, err
	}
	res, e, err := createJsonRpcResult(data)
	if err != nil {
		return nil, false, false, err
	}
	if e != nil {
		if e_data, e_data_ok := e.Data.(string); e_data_ok && strings.Contains(e_data, "UiAutomation not connected") {
			return nil, true, true, errors.New(strings.ToLower(e.Message))
		} else {
			return nil, true, false, errors.New(strings.ToLower(e.Message))
		}
	}
	return res, false, false, nil
}

func formatXML(data []byte) ([]byte, error) {
	b := &bytes.Buffer{}
	decoder := xml.NewDecoder(bytes.NewReader(data))
	encoder := xml.NewEncoder(b)
	encoder.Indent("", "  ")
	for {
		token, err := decoder.Token()
		if err == io.EOF {
			encoder.Flush()
			return b.Bytes(), nil
		}
		if err != nil {
			return nil, err
		}
		err = encoder.EncodeToken(token)
		if err != nil {
			return nil, err
		}
	}
}

func (d *Device) DumpHierarchy(compressed bool, pretty bool) (string, error) {
	// 965 content = self.jsonrpc.dumpWindowHierarchy(compressed, None)
	result, err := d.requestJsonRpc("dumpWindowHierarchy", compressed, nil)
	if err != nil {
		return "", err
	}
	content, ok := result.(string)
	if !ok || content == "" {
		return "", errors.New("dump hierarchy is empty")
	}
	if pretty {
		_xml, err := formatXML([]byte(content))
		if err != nil {
			return content, nil
		}
		return string(_xml), nil
	}
	return content, nil
}

func (d *Device) DumpHierarchyDefault() (string, error) {
	return d.DumpHierarchy(false, false)
}

func FormatHierachy(content string) (*xmlquery.Node, error) {
	doc, err := xmlquery.Parse(strings.NewReader(content))
	if err != nil {
		return nil, err
	}
	els, err := xmlquery.QueryAll(doc, "//node")
	if err != nil {
		return nil, err
	}
	for _, t := range els {
		if len(t.Attr) == 0 {
			continue
		}
		for aidx, a := range t.Attr {
			if a.Name.Local == "class" {
				cls := a.Value
				t.Data = strings.ReplaceAll(cls, "$", "-")
				t.Attr = append(t.Attr[:aidx], t.Attr[aidx+1:]...)
				break
			}
		}
	}
	return doc, nil
}

func (d *Device) ImplicitlyWait(to time.Duration) {
	d.settings.ImplicitlyWait(to)
}

func (d *Device) pos_rel2abs(fast bool) (func(float32, float32) (int, int, error), error) {
	var _width, _height int
	if fast {
		_w, _h, err := d.WindowSize()
		if err != nil {
			return nil, err
		}
		_width = _w
		_height = _h
	}
	getSize := func() (int, int, error) {
		if fast {
			return _width, _height, nil
		} else {
			return d.WindowSize()
		}
	}
	return func(x, y float32) (int, int, error) {
		if x < 0 || y < 0 {
			return 0, 0, errors.New("坐标值不能为负")
		}
		var w, h int
		if x < 1 || y < 1 {
			_w, _h, err := getSize()
			if err != nil {
				return 0, 0, err
			}
			w = _w
			h = _h
		} else {
			return int(x), int(y), nil
		}
		_x := int(x)
		_y := int(y)
		if x < 1 {
			_x = (int)(x * float32(w))
		}
		if y < 1 {
			_y = (int)(y * float32(h))
		}
		return _x, _y, nil
	}, nil
}

type Touch struct {
	d *Device
}

func (t *Touch) Down(x float32, y float32) error {
	rel2abs, err := t.d.pos_rel2abs(t.d.settings.FastRel2Abs)
	if err != nil {
		return err
	}
	_x, _y, err := rel2abs(x, y)
	if err != nil {
		return err
	}
	_, err = t.d.requestJsonRpc("injectInputEvent", 0, _x, _y, 0)
	// fmt.Println(a)
	if err != nil {
		return err
	}
	return nil
}

func (t *Touch) Move(x float32, y float32) error {
	rel2abs, err := t.d.pos_rel2abs(t.d.settings.FastRel2Abs)
	if err != nil {
		return err
	}
	_x, _y, err := rel2abs(x, y)
	if err != nil {
		return err
	}
	_, err = t.d.requestJsonRpc("injectInputEvent", 2, _x, _y, 0)
	if err != nil {
		return err
	}
	return nil
}

func (t *Touch) Up(x float32, y float32) error {
	rel2abs, err := t.d.pos_rel2abs(t.d.settings.FastRel2Abs)
	if err != nil {
		return err
	}
	_x, _y, err := rel2abs(x, y)
	if err != nil {
		return err
	}
	_, err = t.d.requestJsonRpc("injectInputEvent", 1, _x, _y, 0)
	if err != nil {
		return err
	}
	return nil
}

func (d *Device) Touch() *Touch {
	return &Touch{
		d: d,
	}
}

func (d *Device) Click(x float32, y float32) error {
	rel2abs, err := d.pos_rel2abs(d.settings.FastRel2Abs)
	if err != nil {
		return err
	}
	_x, _y, err := rel2abs(x, y)
	if err != nil {
		return err
	}
	delayAfter := d.settings.operation_delay("click")
	defer delayAfter()
	_, err = d.requestJsonRpc("click", _x, _y)
	if err != nil {
		return err
	}
	return nil
}

func (d *Device) Tap(x, y int) error {
	_, err := d.Shell(fmt.Sprintf("input tap %d %d", x, y))
	return err
}

func (d *Device) DoubleClick(x float32, y float32, duration time.Duration) error {
	t := d.Touch()
	err := t.Down(x, y)
	if err != nil {
		return err
	}
	err = t.Up(x, y)
	if err != nil {
		return err
	}
	time.Sleep(duration)
	err = d.Click(x, y)
	if err != nil {
		return err
	}
	return nil
}

func (d *Device) DoubleClickDefault(x float32, y float32) error {
	return d.DoubleClick(x, y, 100*time.Millisecond)
}

func (d *Device) LongClick(x, y float32, duration time.Duration) error {
	t := d.Touch()
	err := t.Down(x, y)
	if err != nil {
		return err
	}
	time.Sleep(duration)
	err = t.Up(x, y)
	if err != nil {
		return err
	}
	return nil
}

func (d *Device) LongClickDefault(x, y float32) error {
	return d.LongClick(x, y, 500*time.Millisecond)
}

func (d *Device) Swipe(fx, fy, tx, ty float32, seconds float64) error {
	rel2abs, err := d.pos_rel2abs(d.settings.FastRel2Abs)
	if err != nil {
		return err
	}
	_fx, _fy, err := rel2abs(fx, fy)
	if err != nil {
		return err
	}
	_tx, _ty, err := rel2abs(tx, ty)
	if err != nil {
		return err
	}
	steps := int(math.Max(2, seconds*200))
	delayAfter := d.settings.operation_delay("swipe")
	defer delayAfter()
	_, err = d.requestJsonRpc("swipe", _fx, _fy, _tx, _ty, steps)
	return err
}

func (d *Device) SwipeDefault(fx, fy, tx, ty float32) error {
	return d.Swipe(fx, fy, tx, ty, 0.275)
}

type Point4Swipe struct {
	X float32
	Y float32
}

func (d *Device) SwipePoints(seconds float64, points ...Point4Swipe) error {
	if points == nil || len(points) == 1 || len(points) == 0 {
		return nil
	}
	rel2abs, err := d.pos_rel2abs(d.settings.FastRel2Abs)
	if err != nil {
		return err
	}
	ppoints := make([]int, 0)
	for _, p := range points {
		x, y, err := rel2abs(p.X, p.Y)
		if err != nil {
			return err
		}
		ppoints = append(ppoints, x, y)
	}
	steps := int(math.Max(2, seconds*200))
	_, err = d.requestJsonRpc("swipePoints", ppoints, steps)
	return err
}

func (d *Device) SwipePointsDefault(points ...Point4Swipe) error {
	return d.SwipePoints(0.275, points...)
}

func (d *Device) Drag(sx, sy, ex, ey float32, seconds float64) error {
	rel2abs, err := d.pos_rel2abs(d.settings.FastRel2Abs)
	if err != nil {
		return err
	}
	_fx, _fy, err := rel2abs(sx, sy)
	if err != nil {
		return err
	}
	_tx, _ty, err := rel2abs(ex, ey)
	if err != nil {
		return err
	}
	steps := int(math.Max(2, seconds*200))
	delayAfter := d.settings.operation_delay("drag")
	defer delayAfter()
	_, err = d.requestJsonRpc("drag", _fx, _fy, _tx, _ty, steps)
	return err
}

func (d *Device) DragDefault(sx, sy, ex, ey float32) error {
	return d.Drag(sx, sy, ex, ey, 0.275)
}

const (
	VSK_HOME        = "home"
	VSK_BACK        = "back"
	VSK_LEFT        = "left"
	VSK_RIGHT       = "right"
	VSK_UP          = "up"
	VSK_DOWN        = "down"
	VSK_CENTER      = "center"
	VSK_MENU        = "menu"
	VSK_SEARCH      = "search"
	VSK_ENTER       = "enter"
	VSK_DELETE      = "delete"
	VSK_DEL         = "del"
	VSK_RECENT      = "recent" //recent apps
	VSK_VOLUME_UP   = "volume_up"
	VSK_VOLUME_DOWN = "volume_down"
	VSK_VOLUME_MUTE = "volume_mute"
	VSK_CAMERA      = "camera"
	VSK_POWER       = "power"
)

func (d *Device) Press(key string) error {
	delayAfter := d.settings.operation_delay("press")
	defer delayAfter()
	_, err := d.requestJsonRpc("pressKey", key)
	return err
}

func (d *Device) Press2(key int) error {
	delayAfter := d.settings.operation_delay("press")
	defer delayAfter()
	_, err := d.requestJsonRpc("pressKeyCode", key)
	return err
}

func (d *Device) Press2WithMeta(key int, meta int) error {
	delayAfter := d.settings.operation_delay("press")
	defer delayAfter()
	_, err := d.requestJsonRpc("pressKeyCode", key, meta)
	return err
}

func (d *Device) SetOrientation(orient string) error {
	switch orient {
	case "n":
		orient = "natural"
	case "l":
		orient = "left"
	case "u":
		orient = "upsidedown"
	case "r":
		orient = "right"
	default:
		return nil
	}
	_, err := d.requestJsonRpc("setOrientation", orient)
	return err
}

func (d *Device) LastTraversedText() (interface{}, error) {
	return d.requestJsonRpc("getLastTraversedText")
}

func (d *Device) ClearTraversedText() error {
	_, err := d.requestJsonRpc("clearLastTraversedText")
	return err
}

func (d *Device) OpenNotification() error {
	_, err := d.requestJsonRpc("openNotification")
	return err
}

func (d *Device) OpenQuickSettings() error {
	_, err := d.requestJsonRpc("openQuickSettings")
	return err
}

func (d *Device) OpenUrl(url string) error {
	if url == "" {
		return nil
	}
	_, err := d.Shell("am start -a android.intent.action.VIEW -d " + url)
	return err
}

func (d *Device) GetClipboard() (string, error) {
	txt, err := d.requestJsonRpc("getClipboard")
	if err != nil {
		return "", err
	}
	return txt.(string), nil
}

func (d *Device) SetClipboard(text string) error {
	_, err := d.requestJsonRpc("setClipboard", nil, text)
	return err
}

func (d *Device) SetClipboard2(text string, label string) error {
	_, err := d.requestJsonRpc("setClipboard", label, text)
	return err
}

func (d *Device) KeyEvent(v int) error {
	_, err := d.Shell("input keyevent " + strconv.Itoa(v))
	return err
}

func (d *Device) ShowFloatWindow(show bool) error {
	arg := strings.ToLower(strconv.FormatBool(show))
	_, err := d.Shell("am start -n com.github.uiautomator/.ToastActivity -e showFloatWindow " + arg)
	return err
}

type Toast struct {
	d *Device
}

func (t *Toast) GetMessage(waitTimeout time.Duration, cacheTimeout time.Duration, defaultMessage string) (string, error) {
	now := time.Now()
	deadline := now.Add(waitTimeout)
	for time.Now().Before(deadline) {
		_msg, err := t.d.requestJsonRpc("getLastToast", cacheTimeout.Milliseconds())
		if err != nil {
			return "", err
		}
		msg, ok := _msg.(string)
		if ok {
			return msg, nil
		}
		time.Sleep(500 * time.Millisecond)
	}
	return defaultMessage, nil
}

func (t *Toast) Reset() error {
	_, err := t.d.requestJsonRpc("clearLastToast")
	return err
}

func (t *Toast) Show(text string, duration time.Duration) error {
	_, err := t.d.requestJsonRpc("makeToast", text, duration.Milliseconds())
	return err
}

func (d *Device) Toast() *Toast {
	return &Toast{
		d: d,
	}
}

const (
	IDENTIFY_THEME_BLACK = "black"
	IDENTIFY_THEME_RED   = "red"
)

func (d *Device) OpenIdentify(theme string) error {
	_, err := d.Shell("am start -W -n com.github.uiautomator/.IdentifyActivity -e theme " + theme)
	return err
}

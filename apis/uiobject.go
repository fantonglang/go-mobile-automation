package apis

import (
	"encoding/json"
	"image"
	"image/draw"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type UiObjectMixIn struct {
	d *Device
}

type UiObject struct {
	d     *Device
	query map[string]interface{}
}

type UiObjectQuery struct {
	key string
	val interface{}
}

type UiElement struct {
	parent *UiObject
	info   temp_info
}

var UIOBJECT_FIELDS map[string]int

func init() {
	UIOBJECT_FIELDS = map[string]int{
		"text":                  0x01,
		"textContains":          0x02,
		"textMatches":           0x04,
		"textStartsWith":        0x08,
		"className":             0x10,
		"classNameMatches":      0x20,
		"description":           0x40,
		"descriptionContains":   0x80,
		"descriptionMatches":    0x0100,
		"descriptionStartsWith": 0x0200,
		"checkable":             0x0400,
		"checked":               0x0800,
		"clickable":             0x1000,
		"longClickable":         0x2000,
		"scrollable":            0x4000,
		"enabled":               0x8000,
		"focusable":             0x010000,
		"focused":               0x020000,
		"selected":              0x040000,
		"packageName":           0x080000,
		"packageNameMatches":    0x100000,
		"resourceId":            0x200000,
		"resourceIdMatches":     0x400000,
		"index":                 0x800000,
		"instance":              0x01000000,
	}
}

func NewUiObjectQuery(key string, val interface{}) *UiObjectQuery {
	if _, ok := UIOBJECT_FIELDS[key]; ok {
		if str, _ok := val.(string); _ok {
			a := strconv.QuoteToASCII(str)
			a = strings.ReplaceAll(a, `"`, "")
			val = a
		}
		return &UiObjectQuery{
			key: key,
			val: val,
		}
	}
	return nil
}

func (u *UiObjectMixIn) UiObject(queries ...*UiObjectQuery) *UiObject {
	m := make(map[string]interface{})
	mask := 0
	for _, q := range queries {
		if q == nil {
			continue
		}
		m[q.key] = q.val
		_m := UIOBJECT_FIELDS[q.key]
		mask |= _m
	}
	if len(m) == 0 {
		return nil
	}
	m["mask"] = mask
	m["childOrSibling"] = make([]string, 0)
	m["childOrSiblingSelector"] = make([]string, 0)
	return &UiObject{
		d:     u.d,
		query: m,
	}
}

type temp_bound_type struct {
	Bottom int `json:"bottom"`
	Left   int `json:"left"`
	Right  int `json:"right"`
	Top    int `json:"top"`
}
type temp_info struct {
	Bounds             *temp_bound_type `json:"bounds"`
	Checkable          bool             `json:"checkable"`
	Checked            bool             `json:"checked"`
	ChildCount         int              `json:"childCount"`
	ClassName          string           `json:"className"`
	Clickable          bool             `json:"clickable"`
	ContentDescription string           `json:"contentDescription"`
	Enabled            bool             `json:"enabled"`
	Focusable          bool             `json:"focusable"`
	Focused            bool             `json:"focused"`
	LongClickable      bool             `json:"longClickable"`
	PackageName        string           `json:"packageName"`
	ResourceName       string           `json:"resourceName"`
	Scrollable         bool             `json:"scrollable"`
	Selected           bool             `json:"selected"`
	Text               string           `json:"text"`
	VisibleBounds      *temp_bound_type `json:"visibleBounds"`
}

func (u *UiObject) Child(queries ...*UiObjectQuery) *UiObject {
	if reflect.ValueOf(u.query["childOrSibling"]).Len() != 0 {
		return nil
	}
	var q map[string]interface{}
	bytes, _ := json.Marshal(u.query)
	json.Unmarshal(bytes, &q)
	_u := u.d.UiObjectMixIn.UiObject(queries...)
	if _u == nil {
		return nil
	}
	childOrSibling := q["childOrSibling"].([]interface{})
	childOrSibling = append(childOrSibling, "child")
	q["childOrSibling"] = childOrSibling
	childOrSiblingSelector := q["childOrSiblingSelector"].([]interface{})
	childOrSiblingSelector = append(childOrSiblingSelector, _u.query)
	q["childOrSiblingSelector"] = childOrSiblingSelector
	return &UiObject{
		d:     u.d,
		query: q,
	}
}

func (u *UiObject) Sibling(queries ...*UiObjectQuery) *UiObject {
	if reflect.ValueOf(u.query["childOrSibling"]).Len() != 0 {
		return nil
	}
	var q map[string]interface{}
	bytes, _ := json.Marshal(u.query)
	json.Unmarshal(bytes, &q)
	_u := u.d.UiObjectMixIn.UiObject(queries...)
	if _u == nil {
		return nil
	}
	childOrSibling := q["childOrSibling"].([]interface{})
	childOrSibling = append(childOrSibling, "sibling")
	q["childOrSibling"] = childOrSibling
	childOrSiblingSelector := q["childOrSiblingSelector"].([]interface{})
	childOrSiblingSelector = append(childOrSiblingSelector, _u.query)
	q["childOrSiblingSelector"] = childOrSiblingSelector
	return &UiObject{
		d:     u.d,
		query: q,
	}
}

func (u *UiObject) Count() (int, error) {
	c, err := u.d.requestJsonRpc("count", u.query)
	if err != nil {
		return 0, err
	}
	return interface2int(c), nil
}

func (u *UiObject) Index(idx int) *UiObject {
	var q map[string]interface{}
	bytes, _ := json.Marshal(u.query)
	json.Unmarshal(bytes, &q)
	childOrSiblingSelector := q["childOrSiblingSelector"].([]interface{})
	if len(childOrSiblingSelector) != 0 {
		a := childOrSiblingSelector[0].(map[string]interface{})
		a["instance"] = idx
		a["mask"] = interface2int(a["mask"]) | UIOBJECT_FIELDS["instance"]
	} else {
		q["instance"] = idx
		q["mask"] = interface2int(q["mask"]) | UIOBJECT_FIELDS["instance"]
	}
	return &UiObject{
		d:     u.d,
		query: q,
	}
}

func interface2int(val interface{}) int {
	bytes, _ := json.Marshal(val)
	i, _ := strconv.Atoi(string(bytes))
	return i
}

func (u *UiObject) Get() *UiElement {
	raw, err := u.d.requestJsonRpc("objInfo", u.query)
	if err != nil {
		return nil
	}
	bytes, err := json.Marshal(raw)
	if err != nil {
		return nil
	}
	var _info temp_info
	err = json.Unmarshal(bytes, &_info)
	if err != nil {
		return nil
	}
	return &UiElement{
		parent: u,
		info:   _info,
	}
}

func (u *UiObject) Wait(timeout time.Duration) int {
	now := time.Now()
	deadline := now.Add(timeout)
	for time.Now().Before(deadline) {
		c, err := u.Count()
		if err != nil {
			panic(err)
		}
		if c > 0 {
			return c
		}
		time.Sleep(200 * time.Millisecond)
	}
	return -1
}

func (u *UiObject) WaitGone(timeout time.Duration) bool {
	now := time.Now()
	deadline := now.Add(timeout)
	for time.Now().Before(deadline) {
		c, err := u.Count()
		if err != nil {
			panic(err)
		}
		if c == 0 {
			return true
		}
		time.Sleep(200 * time.Millisecond)
	}
	return false
}

func (u *UiElement) Info() *Info {
	_info := u.info
	info := &Info{
		Text:               _info.Text,
		Focusable:          _info.Focusable,
		Enabled:            _info.Enabled,
		Focused:            _info.Focused,
		Scrollable:         _info.Scrollable,
		Selected:           _info.Selected,
		ClassName:          _info.ClassName,
		ContentDescription: _info.ContentDescription,
		LongClickable:      _info.LongClickable,
		PackageName:        _info.PackageName,
		ResourceName:       _info.ResourceName,
		ResourceId:         _info.ResourceName,
		ChildCount:         _info.ChildCount,
	}
	if _info.Bounds != nil {
		info.Bounds = &Bounds{
			LX: _info.Bounds.Left,
			LY: _info.Bounds.Top,
			RX: _info.Bounds.Right,
			RY: _info.Bounds.Bottom,
		}
	}
	return info
}

func (u *UiElement) Bounds() *Bounds {
	return u.Info().Bounds
}

func (u *UiElement) PercentBounds() *PercentBounds {
	bounds := u.Bounds()
	if bounds == nil {
		return nil
	}
	w, h, err := u.parent.d.WindowSize()
	if err != nil {
		return nil
	}
	return &PercentBounds{
		LX: float32(bounds.LX) / float32(w),
		LY: float32(bounds.LY) / float32(h),
		RX: float32(bounds.RX) / float32(w),
		RY: float32(bounds.RY) / float32(h),
	}
}

func (u *UiElement) Rect() *Rect {
	bounds := u.Bounds()
	if bounds == nil {
		return nil
	}
	return &Rect{
		LX:     bounds.LX,
		LY:     bounds.LY,
		Width:  bounds.RX - bounds.LX,
		Height: bounds.RY - bounds.LY,
	}
}

func (u *UiElement) PercentSize() *PercentSize {
	rect := u.Rect()
	if rect == nil {
		return nil
	}
	ww, wh, err := u.parent.d.WindowSize()
	if err != nil {
		return nil
	}
	return &PercentSize{
		Width:  float32(rect.Width) / float32(ww),
		Height: float32(rect.Height) / float32(wh),
	}
}

func (u *UiElement) Text() string {
	return u.info.Text
}

func (u *UiElement) Offset(px, py float32) (int, int, bool) {
	rect := u.Rect()
	if rect == nil {
		return 0, 0, false
	}
	x := int(float32(rect.LX) + float32(rect.Width)*px)
	y := int(float32(rect.LY) + float32(rect.Height)*py)
	return x, y, true
}

func (u *UiElement) Center() (int, int, bool) {
	return u.Offset(0.5, 0.5)
}

func (u *UiElement) Click() bool {
	x, y, ok := u.Center()
	if !ok {
		return false
	}
	err := u.parent.d.Click(float32(x), float32(y))
	return err == nil
}

func (u *UiElement) SwipeInsideList(direction int, scale float32) bool {
	if scale <= 0 || scale >= 1 {
		return false
	}
	bounds := u.Rect()
	if bounds == nil {
		return false
	}
	left := int(float32(bounds.LX) + float32(bounds.Width)*(1-scale)/2.0)
	right := int(float32(bounds.LX) + float32(bounds.Width)*(1+scale)/2.0)
	top := int(float32(bounds.LY) + float32(bounds.Height)*(1-scale)/2.0)
	bottom := int(float32(bounds.LY) + float32(bounds.Height)*(1+scale)/2.0)
	if direction == SWIPE_DIR_LEFT {
		err := u.parent.d.SwipeDefault(float32(right), 0.5, float32(left), 0.5)
		return err == nil
	} else if direction == SWIPE_DIR_RIGHT {
		err := u.parent.d.SwipeDefault(float32(left), 0.5, float32(right), 0.5)
		return err == nil
	} else if direction == SWIPE_DIR_UP {
		err := u.parent.d.SwipeDefault(0.5, float32(bottom), 0.5, float32(top))
		return err == nil
	} else if direction == SWIPE_DIR_DOWN {
		err := u.parent.d.SwipeDefault(0.5, float32(top), 0.5, float32(bottom))
		return err == nil
	} else {
		return false
	}
}

func (u *UiElement) Type(text string) bool {
	ok := u.Click()
	if !ok {
		return false
	}
	err := u.parent.d.SendKeys(text, true)
	return err == nil
}

func (u *UiElement) Screenshot() image.Image {
	rect := u.Rect()
	if rect == nil {
		return nil
	}
	img, _, err := u.parent.d.Screenshot()
	if err != nil {
		return nil
	}
	m := image.NewRGBA(image.Rect(0, 0, rect.Width, rect.Height))
	draw.Draw(m, m.Bounds(), img, image.Point{rect.LX, rect.LY}, draw.Src)
	return m
}

package apis

import (
	"fmt"
	"image"
	"image/draw"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/antchfx/xmlquery"
)

func strict_xpath(xpath string) string {
	if strings.HasPrefix(xpath, "/") {
		return xpath
	} else if strings.HasPrefix(xpath, "@") {
		return fmt.Sprintf(`//*[@resource-id="%s"]`, xpath[1:])
	} else if strings.HasPrefix(xpath, "%") && strings.HasSuffix(xpath, "%") {
		_template := `//*[contains(@text, "{0}") or contains(@content-desc, "{0}")]`
		return strings.ReplaceAll(_template, "{0}", xpath[1:len(xpath)-1])
	} else if strings.HasPrefix(xpath, "%") {
		text := xpath[1:]
		_template := `//*[substring-before(@text, "{0}") or @text="{0}" or substring-before(@content-desc, "{0}") or @content-desc="{0}"]`
		return strings.ReplaceAll(_template, "{0}", text)
	} else if strings.HasSuffix(xpath, "%") {
		text := xpath[:len(xpath)-1]
		_template := `//*[starts-with(@text, "{0}") or starts-with(@content-desc, "{0}")]`
		return strings.ReplaceAll(_template, "{0}", text)
	} else {
		_template := `//*[@text="{0}" or @content-desc="{0}" or @resource-id="{0}"]`
		return strings.ReplaceAll(_template, "{0}", xpath)
	}
}

type XPathMixIn struct {
	d *Device
}

type XPath struct {
	d               *Device
	xpath           string
	source          *xmlquery.Node
	pathType        int
	descendantXpath string
}

const (
	XPATH_TYPE_ORIG = iota
	XPATH_TYPE_CHILD
	XPATH_TYPE_SIBLING
	XPATH_TYPE_DESCENDANT
)

type XMLElement struct {
	parent *XPath
	el     *xmlquery.Node
}

func (xp *XPath) all(useSource bool) ([]*XMLElement, error) {
	// 1. get source
	var _doc *xmlquery.Node
	if xp.source != nil && useSource {
		_doc = xp.source
	} else {
		hierachyTxt, err := xp.d.DumpHierarchyDefault()
		if err != nil {
			return nil, err
		}
		_doc, err = FormatHierachy(hierachyTxt)
		if err != nil {
			return nil, err
		}
	}
	// 2. xpath find
	xpath := strict_xpath(xp.xpath)
	xp.xpath = xpath
	els, err := xmlquery.QueryAll(_doc, xpath)
	if err != nil {
		return nil, err
	}
	if len(els) == 0 {
		return nil, nil
	}
	xmlElements := make([]*XMLElement, 0)
	for _, el := range els {
		xmlElements = append(xmlElements, &XMLElement{
			parent: xp,
			el:     el,
		})
	}
	return xmlElements, nil
}

func (xp *XPathMixIn) XPath(xpath string) *XPath {
	return &XPath{
		d:               xp.d,
		xpath:           xpath,
		source:          nil,
		pathType:        XPATH_TYPE_ORIG,
		descendantXpath: "",
	}
}

func (xp *XPathMixIn) XPath2(xpath string, source *xmlquery.Node) *XPath {
	return &XPath{
		d:               xp.d,
		xpath:           xpath,
		source:          source,
		pathType:        XPATH_TYPE_ORIG,
		descendantXpath: "",
	}
}

func (xp *XPath) All() []*XMLElement {
	els, err := xp.all(true)
	if err != nil {
		panic(err)
	}
	return els
}

func (xp *XPath) First() *XMLElement {
	els, err := xp.all(true)
	if err != nil {
		panic(err)
	}
	if len(els) == 0 {
		return nil
	}
	return els[0]
}

func (xp *XPath) Wait(timeout time.Duration) *XMLElement {
	now := time.Now()
	deadline := now.Add(timeout)
	for time.Now().Before(deadline) {
		els, err := xp.all(false)
		if err != nil {
			panic(err)
		}
		if len(els) > 0 {
			return els[0]
		}
		time.Sleep(200 * time.Millisecond)
	}
	return nil
}

func (xp *XPath) WaitGone(timeout time.Duration) bool {
	now := time.Now()
	deadline := now.Add(timeout)
	for time.Now().Before(deadline) {
		els, err := xp.all(false)
		if err != nil {
			panic(err)
		}
		if len(els) == 0 {
			return true
		}
		time.Sleep(200 * time.Millisecond)
	}
	return false
}

type Bounds struct {
	LX int
	LY int
	RX int
	RY int
}

type PercentBounds struct {
	LX float32
	LY float32
	RX float32
	RY float32
}

type Rect struct {
	LX     int
	LY     int
	Width  int
	Height int
}

type PercentSize struct {
	Width  float32
	Height float32
}

type Info struct {
	Text               string
	Focusable          bool
	Enabled            bool
	Focused            bool
	Scrollable         bool
	Selected           bool
	ClassName          string
	Bounds             *Bounds
	ContentDescription string
	LongClickable      bool
	PackageName        string
	ResourceName       string
	ResourceId         string
	ChildCount         int
}

func findAttribute(attrs []xmlquery.Attr, name string) *xmlquery.Attr {
	for _, a := range attrs {
		if a.Name.Local == name {
			return &a
		}
	}
	return nil
}

func (el *XMLElement) Bounds() *Bounds {
	boundsAttr := findAttribute(el.el.Attr, "bounds")
	if boundsAttr == nil {
		return nil
	}
	str := boundsAttr.Value
	if str == "" {
		return nil
	}
	re := regexp.MustCompile(`^\[(\d+)\,(\d+)\]\[(\d+)\,(\d+)\]$`)
	groups := re.FindStringSubmatch(str)
	if len(groups) == 0 {
		return nil
	}
	lx, err := strconv.Atoi(groups[1])
	if err != nil {
		return nil
	}
	ly, err := strconv.Atoi(groups[2])
	if err != nil {
		return nil
	}
	rx, err := strconv.Atoi(groups[3])
	if err != nil {
		return nil
	}
	ry, err := strconv.Atoi(groups[4])
	if err != nil {
		return nil
	}
	return &Bounds{
		LX: lx,
		LY: ly,
		RX: rx,
		RY: ry,
	}
}

func (el *XMLElement) PercentBounds() *PercentBounds {
	bounds := el.Bounds()
	if bounds == nil {
		return nil
	}
	w, h, err := el.parent.d.WindowSize()
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

func (el *XMLElement) Rect() *Rect {
	bounds := el.Bounds()
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

func (el *XMLElement) PercentSize() *PercentSize {
	rect := el.Rect()
	if rect == nil {
		return nil
	}
	ww, wh, err := el.parent.d.WindowSize()
	if err != nil {
		return nil
	}
	return &PercentSize{
		Width:  float32(rect.Width) / float32(ww),
		Height: float32(rect.Height) / float32(wh),
	}
}

func (el *XMLElement) Text() string {
	textVal := el.Attr("text")
	if textVal != "" {
		return textVal
	}
	contentDescVal := el.Attr("content-desc")
	if contentDescVal != "" {
		return contentDescVal
	}
	return ""
}

func (el *XMLElement) Attr(name string) string {
	a := findAttribute(el.el.Attr, name)
	if a != nil {
		return a.Value
	}
	return ""
}

func (el *XMLElement) Info() *Info {
	text := el.Attr("text")
	focusable := el.Attr("focusable")
	enabled := el.Attr("enabled")
	focused := el.Attr("focused")
	scrollable := el.Attr("scrollable")
	selected := el.Attr("selected")
	className := el.el.Data
	bounds := el.Bounds()
	contentDescription := el.Attr("content-desc")
	longClickable := el.Attr("long-clickable")
	packageName := el.Attr("package")
	resourceName := el.Attr("resource-id")
	resourceId := resourceName
	childCount := 0
	if el.el.FirstChild != nil {
		trimFunc := func(r rune) bool {
			return r == ' ' || r == '\t' || r == '\n'
		}
		for n := el.el.FirstChild; n != el.el.LastChild; n = n.NextSibling {
			if strings.TrimFunc(n.Data, trimFunc) != "" {
				childCount++
			}
		}
		if el.el.FirstChild != el.el.LastChild {
			if strings.TrimFunc(el.el.LastChild.Data, trimFunc) != "" {
				childCount++
			}
		}
	}
	return &Info{
		Text:               text,
		Focusable:          focusable == "true",
		Enabled:            enabled == "true",
		Focused:            focused == "true",
		Scrollable:         scrollable == "true",
		Selected:           selected == "true",
		ClassName:          className,
		Bounds:             bounds,
		ContentDescription: contentDescription,
		LongClickable:      longClickable == "true",
		PackageName:        packageName,
		ResourceName:       resourceName,
		ResourceId:         resourceId,
		ChildCount:         childCount,
	}
}

func (el *XMLElement) Offset(px, py float32) (int, int, bool) {
	rect := el.Rect()
	if rect == nil {
		return 0, 0, false
	}
	x := int(float32(rect.LX) + float32(rect.Width)*px)
	y := int(float32(rect.LY) + float32(rect.Height)*py)
	return x, y, true
}

func (el *XMLElement) Center() (int, int, bool) {
	return el.Offset(0.5, 0.5)
}

func (el *XMLElement) Click() bool {
	x, y, ok := el.Center()
	if !ok {
		return false
	}
	err := el.parent.d.Click(float32(x), float32(y))
	return err == nil
}

const (
	SWIPE_DIR_LEFT = iota + 1
	SWIPE_DIR_RIGHT
	SWIPE_DIR_UP
	SWIPE_DIR_DOWN
)

func (el *XMLElement) SwipeInsideList(direction int, scale float32) bool {
	if scale <= 0 || scale >= 1 {
		return false
	}
	bounds := el.Rect()
	if bounds == nil {
		return false
	}
	left := int(float32(bounds.LX) + float32(bounds.Width)*(1-scale)/2.0)
	right := int(float32(bounds.LX) + float32(bounds.Width)*(1+scale)/2.0)
	top := int(float32(bounds.LY) + float32(bounds.Height)*(1-scale)/2.0)
	bottom := int(float32(bounds.LY) + float32(bounds.Height)*(1+scale)/2.0)
	if direction == SWIPE_DIR_LEFT {
		err := el.parent.d.SwipeDefault(float32(right), 0.5, float32(left), 0.5)
		return err == nil
	} else if direction == SWIPE_DIR_RIGHT {
		err := el.parent.d.SwipeDefault(float32(left), 0.5, float32(right), 0.5)
		return err == nil
	} else if direction == SWIPE_DIR_UP {
		err := el.parent.d.SwipeDefault(0.5, float32(bottom), 0.5, float32(top))
		return err == nil
	} else if direction == SWIPE_DIR_DOWN {
		err := el.parent.d.SwipeDefault(0.5, float32(top), 0.5, float32(bottom))
		return err == nil
	} else {
		return false
	}
}

func (el *XMLElement) Type(text string) bool {
	ok := el.Click()
	if !ok {
		return false
	}
	err := el.parent.d.SendKeys(text, true)
	return err == nil
}

func (el *XMLElement) Children() []*XMLElement {
	if el.el.FirstChild == nil {
		return nil
	}
	_children := make([]*xmlquery.Node, 0)
	trimFunc := func(r rune) bool {
		return r == ' ' || r == '\t' || r == '\n'
	}
	for n := el.el.FirstChild; n != el.el.LastChild; n = n.NextSibling {
		if strings.TrimFunc(n.Data, trimFunc) != "" {
			_children = append(_children, n)
		}
	}
	if el.el.FirstChild != el.el.LastChild {
		if strings.TrimFunc(el.el.LastChild.Data, trimFunc) != "" {
			_children = append(_children, el.el.LastChild)
		}
	}
	if len(_children) == 0 {
		return nil
	}
	result := make([]*XMLElement, 0)
	xpath_parent := *el.parent
	xpath_parent.pathType = XPATH_TYPE_CHILD
	xpath_parent.descendantXpath = ""
	for _, c := range _children {
		result = append(result, &XMLElement{
			parent: &xpath_parent,
			el:     c,
		})
	}
	return result
}

func (el *XMLElement) Siblings() []*XMLElement {
	elem := el.el
	if elem.Parent == nil {
		return nil
	}
	siblings := make([]*xmlquery.Node, 0)
	head := elem.Parent.FirstChild
	tail := elem.Parent.LastChild
	n := head
	trimFunc := func(r rune) bool {
		return r == ' ' || r == '\t' || r == '\n'
	}
	for {
		if n != elem && strings.TrimFunc(n.Data, trimFunc) != "" {
			siblings = append(siblings, n)
		}
		if n != tail {
			n = n.NextSibling
		} else {
			break
		}
	}
	if len(siblings) == 0 {
		return nil
	}
	result := make([]*XMLElement, 0)
	xpath_parent := *el.parent
	xpath_parent.pathType = XPATH_TYPE_SIBLING
	xpath_parent.descendantXpath = ""
	for _, c := range siblings {
		result = append(result, &XMLElement{
			parent: &xpath_parent,
			el:     c,
		})
	}
	return result
}

func (el *XMLElement) Find(xpath string) []*XMLElement {
	elem := el.el
	true_xpath := strict_xpath(xpath)
	nodes, err := xmlquery.QueryAll(elem, true_xpath)
	if err != nil || len(nodes) == 0 || (len(nodes) == 1 && nodes[0] == elem) {
		return nil
	}
	result := make([]*XMLElement, 0)
	xpath_parent := *el.parent
	xpath_parent.pathType = XPATH_TYPE_DESCENDANT
	xpath_parent.descendantXpath = true_xpath
	for _, c := range nodes {
		if c == elem {
			continue
		}
		result = append(result, &XMLElement{
			parent: &xpath_parent,
			el:     c,
		})
	}
	return result
}

func (el *XMLElement) Screenshot() image.Image {
	rect := el.Rect()
	if rect == nil {
		return nil
	}
	img, _, err := el.parent.d.Screenshot()
	if err != nil {
		return nil
	}
	m := image.NewRGBA(image.Rect(0, 0, rect.Width, rect.Height))
	draw.Draw(m, m.Bounds(), img, image.Point{rect.LX, rect.LY}, draw.Src)
	return m
}

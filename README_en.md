[简体中文](./README.md) | English

# GO-MOBILE-AUTOMATION SDK

A full featured Android mobile automation sdk for golang developers

## The purpose

If you are an automation developer, you may find the python/javascript echo system provide the developers great capabilities for manipulating devices and apps. 

For Android automation, some well known tools are:

1. Appium (multi-language primarily javascript/python)
2. uiautomator2 (python)

Our SDK ports the uiautomation2 python library to golang.

Why do we do this? 

1. Easy to deploy - deploy one executable only (at most several dlls) instead of resolving thousands of dependencies like javascript and python. The bigger scope of this family of projects is to provide a mechanism to orchestrate the automation "scripts" to run cross a bunch of platform/systems. Must have a fast & robustic app distribution and deployment mechanism.
2. For using cloud phones (Android) - Some cloud phone providers have very poor quality adb connection. In case the automation process break because of adb connection failure, we would like the process run inside the phone
3. Robust - apps can easily get killed by Android system, but executable not. That why despite the fact that apps like pyto-python3 provides hosting for python but we don't use it.

## Inspired by

Inspired by [OpenAtx](https://github.com/openatx) and the [Uiautomation2](https://github.com/openatx/uiautomator2) python library. We entirely use the openatx drivers, leaving the client side sdk written in golang. This benefits users, because they can use the uiautomator2 tool chain, which is super cool.

## Quick start

There are 4 steps:

1. Setup the Android phone
2. Setup the development environment
3. Start creating a golang project for mobile automation
4. deploy&run

### Setup the Android phone

1. Install the specific version of the app you want to automate.
``` 
$ adb install [package.apk]
```
2. Download atx-agent from [here](https://github.com/openatx/atx-agent/releases) choose the armv7 version unless you run a x86 phone simulator
3. Untar atx-agent and install, follow the intallation insttructions [here](https://github.com/openatx/atx-agent)
```
$ adb push atx-agent /data/local/tmp
$ adb shell chmod 755 /data/local/tmp/atx-agent
# launch atx-agent in daemon mode
$ adb shell /data/local/tmp/atx-agent server -d

# stop already running atx-agent and start daemon
$ adb shell /data/local/tmp/atx-agent server -d --stop
```
4. Download app-uiautomator-test.apk and app-uiautomator.apk from [here](https://github.com/openatx/android-uiautomator-server/releases) and install using adb install
```
$ adb install app-uiautomator-test.apk
$ adb install app-uiautomator.apk
```
5. Grant all priviledges to the app "ATX"
6. Open app "ATX" and click "启动UIAUTOMATOR", click "开启悬浮窗"

## Setup the development environment

1. Install Python3(version 3.6+) from [here](https://www.python.org/downloads/)
2. Install [weditor](https://github.com/alibaba/web-editor)
```
$ pip3 install -U weditor
```
3. You can now open weditor as ui inspector for android applications, by typing
```
$ weditor
```

## Create a golang automation project

1. Create a folder called helloworld
2. Open terminal and type in
```
$ go mod init helloworld
```
3. Add dependency
```
$ go get github.com/fantonglang/go-mobile-automation
```
4. Add the program entry - create main.go file
5. [Here](https://github.com/fantonglang/go-mobile-automation-examples/blob/main/douyin-luo-live/main.go) is an example of main.go file

## Deploy&Run
```
# cross build linux/arm target
$ GOOS=linux GOARCH=arm go build
# deploy - helloworld is the executable name, which is the same with the go module name
$ adb push helloworld /data/local/tmp
$ adb shell chmod 755 /data/local/tmp/helloworld
# run
$ adb shell /data/local/tmp/helloworld
```
If you start a background process, you don't need the phone to connect with your PC/macos. Note that you can transfer your linux shell knowledge to using the adb shell.

## Examples
[This](https://github.com/fantonglang/go-mobile-automation-examples/blob/main/douyin-luo-live) is a working example. Read the comments in main.go carefully. This helps to resolve all dependencies and environment requirements before you start. The comments also give the commands for compilation, deployment, and execution.

# APIS

**[Connect to a device](#connect-to-a-device)**

**[Device APIS](#device-apis)**
  - **[Shell commands](#shell-commands)**
  - **[Retrieve the device info](#retrieve-the-device-info)**
  - **[Clipboard](#clipboard)**
  - **[Key Events](#key-events)**
  - **[Press Key](#press-key)**
  - **[New command timeout](#new-command-timeout)**
  - **[Screenshot](#screenshot)**
  - **[UI Hierarchy](#ui-hierarchy)**
  - **[Touch](#touch)**
  - **[Click](#click)**
  - **[Double Click](#double-click)**
  - **[Long Click](#long-click)**
  - **[Swipe](#swipe)**
  - **[Set Orientation](#set-orientation)**
  - **[Open Quick Settings](#open-quick-settings)**
  - **[Open Url](#open-url)**
  - **[Show float window](#show-float-window)**

**[Input Method](#input-method)**

**[XPATH](#xpath)**
  - **[Finding elements](#finding-elements)**
  - **[Xpath elements API](#xpath-elements-api)**

**[UI Object](#ui-object)**
  - **[construct query](#construct-query)**
  - **[execute ui object query](#execute-ui-object-query)**
  - **[ui object element apis](#ui-object-element-apis)**

## Connect to a device

There are two types of connection: 

1. If the executable is deployed in the phone, use
``` go
package main

import (
	"log"
	"github.com/fantonglang/go-mobile-automation/apis"
)
...
// you don't need to specify device id, because there is no PC connection
d := apis.NewNativeDevice()
```
2. If you debug and deploy in PC/macos, use
``` go
package main

import (
	"log"
	"github.com/fantonglang/go-mobile-automation/apis"
)


//here c574dd45 is the device id, replace it with yours own
d, err := apis.NewHostDevice("c574dd45")
if err != nil {
  log.Println("failed connecting to device")
  return
}
```

Combine 2 code snippets. The following code enables the same piece of code working on the both deployments.
``` go
package main

import (
	"log"
	"runtime"
	"github.com/fantonglang/go-mobile-automation/apis"
)

func getDevice() *apis.Device {
	if runtime.GOARCH == "arm" {
		return apis.NewNativeDevice()
	}
	//here c574dd45 is the device id, replace it with yours own
	_d, err := apis.NewHostDevice("c574dd45")
	if err != nil {
		log.Println("101: failed connecting to device")
		return nil
	}
	return _d
}
...
d := getDevice()
```

Take extra notice if your macos is the ARM architecture. Then you may judge based on GOOS.

## Device APIS
This part showcases how to perform common device operations:

### Shell commands
Example: Force stop douyin(China tiktok) app
```go
d.Shell(`am force-stop com.ss.android.ugc.aweme`)
```
Example: Start douyin app
```go
// You can find the app main activity by using the dumpsys command, hence I didn't implement the uiautomator2 equivalent session API for now.
d.Shell(`am start -n "com.ss.android.ugc.aweme/.main.MainActivity"`)
```

### Retrieve the device info
Get detailed device info
```go
info, err := d.DeviceInfo()
if err != nil {
  log.Println("get device info failed")
  return
}
bytes, err := json.Marshal(info)
if err != nil {
  log.Println("error marshalling")
  return
}
fmt.Println(string(bytes))
```

Below is a possible output:

```json
{
  ...
  "version":"11",
  "serial":"c574dd45",
  ...
  "sdk":30,
  "agentVersion":"0.10.0",
  "display":{"width":1080,"height":2340}
  ...
}
```

Get window size:

```GO
w, h, err := d.WindowSize()
if err != nil {
  log.Println("get window size failed")
  return
}
fmt.Printf("w: %d, h: %d\n", w, h)
// device upright output example: w: 1080, h: 2340
// device horizontal output example: w: 1080, h: 2340
```

### Clipboard
Get or set clipboard content

set clipboard
```go
err := d.SetClipboard("aaa")
if err != nil {
  log.Println("error clipboard")
  return
}
```
get clipboard: This doesn't work in Android > 9.0. Most cloud phone work on lower Android version. I don't mind. 
```go
a, err := d.GetClipboard()
if err != nil {
  log.Println("error clipboard")
  return
}
fmt.Println(a)
```

### Key Events

* Turn on/off screen
```go
err := d.KeyEvent(KEYCODE_POWER) // press power key to turn on/off screen
```
* Home key
```go
err := d.KeyEvent(KEYCODE_HOME)
```

d.KeyEvent is basically the Android "input keyevent " command, please refer to [this doc](https://developer.android.com/reference/android/view/KeyEvent), or if you're behind gfw, [this doc](https://blog.csdn.net/feizhixuan46789/article/details/16801429)

### Press Key
Example: press Home key
```go
err := d.Press("home")
```

supported keys are:
```go
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
```

### New command timeout
How long (in seconds) will wait for a new command from the client before assuming the client quit and ending the uiautomator service

```go
err := d.SetNewCommandTimeout(300) // unit is second
```

### Screenshot
* Screenshot and save - notice Android has readonly file system, this API is only available for host(PC/macos).
```go
err := d.ScreenshotSave("sc.png")
```
* Screenshot and get bytes (preferred, because opencv can accept bytes directly by "cv::imdecode" function)
```go
bytes, err := d.ScreenshotBytes()
```
* Screenshot and get image.Image object
```go
img, format, err := d.Screenshot() // img is the image.Image object, example of format: "jpeg"
```

### UI Hierarchy
* Get hierachy text
```go
content, err := d.DumpHierarchy(false, false) // content is the text
```
* Transform hierachy text to *xmlquery.Node object. (This is useful if you want to do xpath query based on a snapshot - this is a lot faster)
```go
doc, err := FormatHierachy(content) // doc is the *xmlquery.Node object
```

### Touch
Simulate "mouse press down", "mouse hold and move", "mouse up release"
* Get touch object
```go
touch := d.Touch()
```
* Mouse press down at a position
```go
/* press at (relative to the top-left corner)
 *  x: 50% position of the width
 *  y: 60% position of the height
 */ 
err := touch.Down(0.5, 0.6) 
```
* Mouse hold and move
```go
/* then move to (relative to the top-left corner)
 *  x: 50% position of the width
 *  y: 10% position of the height
 */ 
err := touch.Move(0.5, 0.1) 
```
* Mouse up release
```go
/* then mouse up release at (relative to the top-left corner)
 *  x: 50% position of the width
 *  y: 10% position of the height
 */ 
err := touch.Up(0.5, 0.1) 
```

### Click
Click on screen given coordinates

* coordinates using percentage
```go
/* click relative to the top-left corner
 *  x: 48.1% position of the width
 *  y: 24.6% position of the height
 */ 
err := d.Click(0.481, 0.246)
```
* coordinates using absolute pixel values
```go
/* click relative to the top-left corner
 *  at (x: 481, y: 246)
 */ 
err := d.Click(481, 246)
```

### Double Click
```go
err := d.DoubleClickDefault(0.481, 0.246)
```

### Long Click
Mouse click, but there is a certain time interval(0.5s) between mouse down and up 
```go
err := d.LongClickDefault(0.481, 0.246)
```

### Swipe
* Swipe from one point (fx, fy) to another (tx, ty)
```go
var fx, fy, tx, ty float32 = 0.5, 0.5, 0, 0
err := d.SwipeDefault(fx, fy, tx, ty)
```
* Swipe points, you can specify more than 2 points
```go
// swipe from (x=width*0.5, y=height*0.9) to (x=width*0.5,y=height*0.1)
err := d.SwipePoints(0.1, apis.Point4Swipe{0.5, 0.9}, apis.Point4Swipe{0.5, 0.1})
```
* Drag from one point (fx, fy) to another (tx, ty)
```go
var fx, fy, tx, ty float32 = 0.5, 0.5, 0, 0
err := d.DragDefault(fx, fy, tx, ty)
```

### Set Orientation
Accepts 4 orientation parameters:
* "n" - means natural
* "l" - means left
* "u" - means upsidedown
* "r" - means right
```go
err := d.SetOrientation("n")
```

### Open Quick Settings
```go
err := d.OpenQuickSettings()
```

### Open Url
```go
err := d.OpenUrl("https://bing.com")
```

### Show float window
This operation is openatx specific - open a float window to keep the automator app in the front and prevent it from getting killed
```go
err := d.ShowFloatWindow(true)
```


## Input Method
Type text, (you will switch to a special input method)
* Clear Text
```go
err := d.ClearText()
```
* Send Action
```go
err := ime.SendAction(SENDACTION_SEARCH)
```

The following actions are supported:
```go
SENDACTION_GO       = 2
SENDACTION_SEARCH   = 3
SENDACTION_SEND     = 4
SENDACTION_NEXT     = 5
SENDACTION_DONE     = 6
SENDACTION_PREVIOUS = 7
```

* Send Keys - Type text
```go
err := ime.SendKeys("aaa", true)
```

## XPATH
XPATH is the most important way of finding UI Element.

### Finding elements
* Find multiple elements by xpath
```go
els := d.XPath(`//*[@text="your-control-text"]`).All()
for _, el := range els {
  ...
}
```
* Find one element by xpath
```go
el := d.XPath(`//*[@text="your-control-text"]`).First()
```
* Check element exists by xpath
```go
if d.XPath(`//*[@text="your-control-text"]`).First() != nil {
  ...
}
```
* Wait element appear
```go
el := d.XPath(`//*[@text="your-control-text"]`).Wait(time.Minute)
if el == nil {
  log.Println("element doesn't appear within 1 minute")
  return
}
...
```
* Wait element disappear
```go
ok := d.XPath(`//*[@text="your-control-text"]`).WaitGone(time.Minute)
if !ok {
  log.Println("element doesn't disappear within 1 minute")
  return
}
```
* If you want run xpath query based on ui hierachy snaphot,
```go
content, _ := d.DumpHierarchy(false, false) // content is the text
doc, _ := FormatHierachy(content) // doc is the *xmlquery.Node object
...
el := d.XPath2(`//*[@text="your-control-text"]`, doc).First()
```

### Xpath elements API
* Children of Xpath element
```go
// el := d.XPath(`//*[@resource-id="com.taobao.taobao:id/rv_main_container"]`).First()
children := el.Children()
for _, c := range children {
  ...
}
```
* Siblings of Xpath element
```go
// el := d.XPath(`//*[@resource-id="com.taobao.taobao:id/rv_main_container"]/android.widget.FrameLayout[1]`).First()
siblings := el.Siblings()
for _, s := range siblings {
  ...
}
```
* Find descendants based on xpath
```go
// el := d.XPath(`//*[@resource-id="com.taobao.taobao:id/rv_main_container"]`).First()
children := el.Find(`//android.support.v7.widget.RecyclerView`)
for _, c := range children {
  fmt.Println(*c.Info())
}
```
* Find bounding rect

This describes the bounding box surrounding this control
```go
bounds := el.Bounds()
/* bounds has the type of *apis.Bounds:
 * type Bounds struct {
 *	  LX int // top-left-x
 *	  LY int // top-left-y
 *	  RX int // right-bottom-x
 *	  RY int // right-bottom-y
 * }
 */
```
* Find Rect

This also describes the bounding box surrounding this control
```go
rect := el.Rect()
/* rect has the type of *apis.Rect:
 * type Bounds struct {
 *	  LX int      // top-left-x
 *	  LY int      // top-left-y
 *	  Width int   // width of control
 *	  Height int  // height of control
 * }
 */
```
* Get text of control - the text shown in user interface
```go
text := el.Text() // text is string
```
* Get control's info - which is everything, including text and bounding box
```go
info := el.Info()
/* info has the type of *apis.Info
 * type Info struct {
 *     Text               string
 *     Focusable          bool
 *     Enabled            bool
 *     Focused            bool
 *     Scrollable         bool
 *     Selected           bool
 *     ClassName          string
 *     Bounds             *Bounds
 *     ContentDescription string
 *     LongClickable      bool
 *     PackageName        string
 *     ResourceName       string
 *     ResourceId         string
 *     ChildCount         int
 * }
*/
```
* Get center position
```go
x, y, ok := el.Center()
```
* Click
```go
ok := el.Click()
```
* Swipe inside the control - if the control is a (Recycler)List
```go
dir := apis.SWIPE_DIR_LEFT
var scale float32 = 0.8
ok := el.SwipeInsideList(dir, scale)
/* dir can use these 4 values:
 * SWIPE_DIR_LEFT  = 1 // swipe right to left
 * SWIPE_DIR_RIGHT = 2 // swipe left to right
 * SWIPE_DIR_UP    = 3 // swipe bottom to top
 * SWIPE_DIR_DOWN  = 4 // swipe top to bottom

 scale is the percentage of width/height swiped
*/
```
* Type text
```go
ok := el.Type("aaa")
```
* Screenshot - take screenshot of this control
```go
img := el.Screenshot() // img is the image.Image type
```

## UI Object
Find ui elements via attribute matching search. In most platforms including iOS and windows UIA, accessibility api(UI Object) is far more efficient than xpath. For example, windows UIA, to get the full xml structure takes very long time, because fetching element's info involves cross-process COM calls which takes time. But here in Android, this is not the case. Xpath is as fast as accessibility(UI Object) and far more powerful. I would prefer suggest you to use xpath.

### construct query
UI object API doesn't effectively fetch any elements util you call (*UIObject).Get() *UiElement, (*UIObject).Wait(timeout time.Duration) int, or (*UIObject).WaitGone() bool

* Construct query based on attribute values
```go
uo := d.UiObject(apis.NewUiObjectQuery("resourceId", `com.taobao.taobao:id/rv_main_container`)) // this api accepts multiple NewUiObjectQuery in args, with and relationship
```
* Construct sibling query - I find this api's behavior is a bit awkward. Notice that.
```go
c := d.UiObject(apis.NewUiObjectQuery("resourceId", `com.taobao.taobao:id/sv_search_view`)).Sibling(apis.NewUiObjectQuery("className", "android.widget.FrameLayout"))
```
* Construct decendant query
```go
c := d.UiObject(apis.NewUiObjectQuery("resourceId", `com.taobao.taobao:id/rv_main_container`)).Child(apis.NewUiObjectQuery("className", "android.widget.LinearLayout"))
```
* Construct indexed query
```go
c := d.UiObject(
		apis.NewUiObjectQuery("className", `android.support.v7.widget.RecyclerView`)).Index(0)
```

### execute ui object query
* get first ui element
```go
el := d.UiObject(apis.NewUiObjectQuery("resourceId", `com.taobao.taobao:id/rv_main_container`)).Get() // returns an *apis.UiElement
```
* get element count
```go
count := d.UiObject(apis.NewUiObjectQuery("resourceId", `com.taobao.taobao:id/rv_main_container`)).Count()
```
* get the nth element - here in the example - third(with .Index(2))
```go
el := d.UiObject(apis.NewUiObjectQuery("resourceId", `com.taobao.taobao:id/rv_main_container`)).Index(2).Get()
```
* wait element appear
```go
count := d.UiObject(apis.NewUiObjectQuery("resourceId", `com.taobao.taobao:id/rv_main_container`)).Wait(time.Minute)
// if element doesn't appear in 1 minute, it returns -1
```
* wait element disappear
```go
ok := d.UiObject(apis.NewUiObjectQuery("resourceId", `com.taobao.taobao:id/rv_main_container`)).WaitGone(time.Minute)
// returns true if disappear, false otherwise
```

### ui object element apis
* Get info
```go
info := el.Info() // the info's type is same as xpath element's info api
```
* Get bounding rect
```go
bounds := el.Bounds() // the bounds's type is the same as xpath element's bounds api
```
* Get rect
```go
rect := el.Rect() // the rect's type is the same as xpath element's rects api
```
* Get center position
```go
x, y, ok := el.Center()
```
* Click
```go
ok := el.Click()
```
* Get text of control - the text shown in user interface
```go
text := el.Text() // text is string
```
* Swipe inside the control - if the control is a (Recycler)List
```go
dir := apis.SWIPE_DIR_LEFT
var scale float32 = 0.8
ok := el.SwipeInsideList(dir, scale)
/* dir can use these 4 values:
 * SWIPE_DIR_LEFT  = 1 // swipe right to left
 * SWIPE_DIR_RIGHT = 2 // swipe left to right
 * SWIPE_DIR_UP    = 3 // swipe bottom to top
 * SWIPE_DIR_DOWN  = 4 // swipe top to bottom

 scale is the percentage of width/height swiped
*/
```
* Type text
```go
ok := el.Type("aaa")
```
* Screenshot - take screenshot of this control
```go
img := el.Screenshot() // img is the image.Image type
```

# If you want to support author, please donate(wechat), thanks

![image](https://github.com/fantonglang/go-mobile-automation-examples/blob/main/doc/wechat.jpg)

# Contact me if you'd like working with me on the computer vision & speech recognition part.

![image](https://github.com/fantonglang/go-mobile-automation-examples/blob/main/doc/wechat2.jpg)
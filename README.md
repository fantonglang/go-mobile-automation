简体中文 | [English](./README_en.md)

# GO手机自动化SDK

功能齐全的手机自动化Golang SDK，支持编译成二进制可执行文件，适合自动化流程部署到云手机，脱离ADB连接稳定高效执行

## 目的

在（手机）自动化领域，其实python是绝对的主流，其他语言诸如Golang/C#/C++在整个生态中占的份额相对来说比较少。对于安卓手机自动化开发来说，我们起码有下面的工具了：

1. Appium（多语言/全平台自动化支持）
2. Uiautomator2/ATX (python)
3. 各RPA平台，比如云扩科技 (C#)

我们有这么多python库可以用，为什么要再做个Golang自动化SDK呢？

1. 容易部署 - Golang一般只需要部署一个二进制可执行文件，不像python/javascript需要安装一大堆依赖，在墙内如果没有一个稳定可靠的镜像(的确，javascript有CNPM，也可以自己搭镜像)，如果安装依赖挂了就没有然后了。即使下载依赖没有问题，也会有一些依赖版本冲突的问题。C#/JAVA有运行时依赖，你一般不会知道用户会使用哪个版本的.NET SDK/JDK。
2. 部署在云手机上稳定可靠 - 由于云手机使用虚拟化技术，它不像传统机架集群会遇到各式各样的硬件文件 - 比如电池爆膨了（一般2年寿命，工作室机更差），比如电源供电不足手机断电，比如USB线接触不良，等等。云手机相对贵一点但是稳定且省人力成本。但是有些云手机使用外网远程adb连接，如果自动化脚本host在本地的PC/MACOS上远程操控云手机，稳定性非常取决于网络的稳定性。我们当然更希望使用Golang/C++的方式：编译一个二进制可执行文件，推到云手机上定时/按某种触发条件去执行，脱离ADB连接
3. 稳定 - 事实上python能直接跑在手机上。pyto-python3这个app提供了host python脚本在手机上的可能。但是手机应用比较容易被Android系统杀死，而二进制可执行文件一般不会。

所以我们采用Golang做安卓自动化语言是有道理滴～ 

再说说其他自动化工具坑的地方：
1. Appium安装复杂，而且图像识别只支持模版匹配，中间层N多，上层是well known的protocol，要优化性能或者加一些底层支持非常困难。要部署Appium到一台新的主机虚拟机，别提有多烦了
2. UiAutomator2/ATX。这个软件库我是极其推崇，它的API支持比我们的Golang SDK要丰富，毕竟人家才是原版。但是它不提供直接将流程部署到手机的能力。这不能说是缺陷，因为考虑到自动化开发的受众群体，python是绝对绝对的主流
3. 云扩RPA。它其实本身是一个很棒的平台，但是为了做手机自动化安装N多东西不值当。但是如果有别的需求，比如使用低代码做界面层，可视化的流程管理工具，它是一个不错的选择。这个SDK首先会和云扩RPA做集成。

唯一要注意的是Android是只读文件系统，这意味着打日志的话我们可能需要批量收集，在程序里面做批量推送到日志服务器

## 项目启发

项目启发于[OpenAtx](https://github.com/openatx)和[Uiautomation2](https://github.com/openatx/uiautomator2) python库。我事实上在项目中首先大量使用了uiautomator2，在玩儿了云手机之后，因为有需要，自然而然就有了这个Golang库（因为Golang比C++简单，ARM Cortex A编辑不需要NDK）

它事实上是OpenAtx的又一个客户端，但是它可以跨平台执行。它实现了大部分Uiautomator2的API。这么做的好处是，OpenAtx的工具链，包括它的Inspector [weditor](https://github.com/alibaba/web-editor) - 这个真的很好用，而且是web应用不需要额外占我太多存储空间

而且有别于Appium，OpenAtx不需要依赖Session同时在一个App上执行流程，它可以同时操作一个App，点击另一个App的浮层。这是我喜欢并使用OpenAtx的原因。

## Quick start

我分四步讲

1. 设置安卓手机
2. 设置开发环境
3. 创建Golang自动化项目
4. 部署&执行

### 设置安卓手机

1. 安装你要自动化的APP（注意：安装特定的版本 - 因为不同版本同一个元素的属性很可能是不一样的 - 即使同一个版本有些属性其实也会变 - 反映到脚本中就是xpath）
``` 
$ adb install [package.apk]
```
2. 下载 atx-agent: [点这里](https://github.com/openatx/atx-agent/releases) 选择 armv7 除非你部署的手机是x86手机模拟器
3. 解压 atx-agent 并安装, 这里有安装说明: [点这里](https://github.com/openatx/atx-agent) 或者看下面的命令行指令：
```
$ adb push atx-agent /data/local/tmp
$ adb shell chmod 755 /data/local/tmp/atx-agent
# 后台模式执行atx-agent
$ adb shell /data/local/tmp/atx-agent server -d

# 或者后台模式重启atx-agent
$ adb shell /data/local/tmp/atx-agent server -d --stop
```
4. 下载 app-uiautomator-test.apk 和 app-uiautomator.apk [点这里](https://github.com/openatx/android-uiautomator-server/releases) 然后用adb install命令安装
```
$ adb install app-uiautomator-test.apk
$ adb install app-uiautomator.apk
```
5. 给应用 "ATX" 所有权限，并且打开 "ATX" 检查有任何运行时权限请求的话，选择一直许可。
6. 打开应用 "ATX" 点击 "启动UIAUTOMATOR", 点击 "开启悬浮窗"

## 设置开发环境

1. 安装 Python3(版本 3.6+) from [here](https://www.python.org/downloads/)
2. 安装 [weditor](https://github.com/alibaba/web-editor)
```
$ pip3 install -U weditor
```
3. 打开 weditor, 它就是你的UI Inspector，在命令行中输入,
```
$ weditor
```

## 创建Golang自动化项目

1. 创建文件夹 helloworld
2. 打开命令行并输入
```
$ go mod init helloworld
```
3. 添加依赖
```
$ go get github.com/fantonglang/go-mobile-automation
```
4. 添加程序入口 - 创建 main.go 文件
5. [这里](https://github.com/fantonglang/go-mobile-automation-examples/blob/main/douyin-luo-live/main.go) 是示例 main.go 文件

## 部署&执行
```
# 交叉编译(cross build) linux/arm 可执行文件
$ GOOS=linux GOARCH=arm go build
# 部署 - helloworld 是可执行文件名，和Go模块名是相同的
$ adb push helloworld /data/local/tmp
$ adb shell chmod 755 /data/local/tmp/helloworld
# 运行
$ adb shell /data/local/tmp/helloworld
```
如果你用后台模式启动程序，程序启动之后就不需要连接电脑。adb shell命令和普通linux系统是一样的。你可以输入
```
# nohup保证当关闭terminal session的时候，程序不会被杀死，&保证在后台运行程序
$ adb shell nohup /data/local/tmp/helloworld &
```

## 例子
[这里](https://github.com/fantonglang/go-mobile-automation-examples/blob/main/douyin-luo-live) 是一个能跑起来的例子。仔细阅读main函数的注释，里面有手机设置，环境安装，编译，调试，部署，执行的教程。

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

有两种类型的连接：

1. 如果程序是部署在手机上，使用如下代码：
``` go
package main

import (
	"log"
	"github.com/fantonglang/go-mobile-automation/apis"
)
...
// 不需要指定设备ID，因为程序并不需要从电脑通过ADB连接手机
d := apis.NewNativeDevice()
```
2. 如果是在电脑端调试或部署，使用如下代码：
``` go
package main

import (
	"log"
	"github.com/fantonglang/go-mobile-automation/apis"
)


//这里 c574dd45 是设备ID, 你可以从adb devices指令中获取到它, 把它替换成你自己的手机设备ID
d, err := apis.NewHostDevice("c574dd45")
if err != nil {
  log.Println("failed connecting to device")
  return
}
```

结合上述两个代码片段，下面的代码既能在电脑端(假定是x86架构)调试的时候工作，又能在Android设备中部署之后工作。它判断运行时如果是ARM，使用手机的方式连接；反之使用电脑的方式连接。特别注意苹果MAC最近的ARM架构电脑，如果这种情况，最好判断一下系统(GOOS)
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
	//这里 c574dd45 是设备ID, 你可以从adb devices指令中获取到它, 把它替换成你自己的手机设备ID
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

## Device APIS
这部分展示如何执行常见的设备操作

### Shell commands
示例: 强制停止抖音APP
```go
d.Shell(`am force-stop com.ss.android.ugc.aweme`)
```
示例: 打开抖音APP
```go
// 使用dumpsys命令，你可以找到APP的启动Activity。所以目前并不实现uiautomator2的启动APP API，以及session API。
d.Shell(`am start -n "com.ss.android.ugc.aweme/.main.MainActivity"`)
```

### Retrieve the device info
获取详细设备信息
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

下面是可能的输出：

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

获取窗口大小：

```GO
w, h, err := d.WindowSize()
if err != nil {
  log.Println("get window size failed")
  return
}
fmt.Printf("w: %d, h: %d\n", w, h)
// 设备竖屏时的输出示例: w: 1080, h: 2340
// 设备横屏时的输出示例: w: 1080, h: 2340
```

### Clipboard
获取和设置剪切板内容

设置剪切板内容
```go
err := d.SetClipboard("aaa")
if err != nil {
  log.Println("error clipboard")
  return
}
```
获取剪切板内容: 在Android大于9.0这个API并不工作. 但是大多数云手机使用低版本Android(比如7.0)，所以我并不care
```go
a, err := d.GetClipboard()
if err != nil {
  log.Println("error clipboard")
  return
}
fmt.Println(a)
```

### Key Events

* 打开/关闭屏幕
```go
err := d.KeyEvent(KEYCODE_POWER) // 打开关闭屏幕都是按电源键
```
* Home键
```go
err := d.KeyEvent(KEYCODE_HOME)
```

d.KeyEvent 是在调用 Android 的 "input keyevent " 命令, 请参照 [这个文档](https://developer.android.com/reference/android/view/KeyEvent), 或者你如果没有翻墙工具, 参照 [这个文档](https://blog.csdn.net/feizhixuan46789/article/details/16801429)

### Press Key
示例: 按Home键
```go
err := d.Press("home")
```

Press API支持下面的按键:
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
设置Uiautomator服务的超时时间

```go
err := d.SetNewCommandTimeout(300) // 单位秒
```

### Screenshot
* 截屏并保存文件 - 注意 Android 使用只读文件系统, 这个API只有在电脑端有效.
```go
err := d.ScreenshotSave("sc.png")
```
* 截屏并获取[]byte字节流 (推荐, 因为opencv使用cv::imdecode函数能直接读取字节流)
```go
bytes, err := d.ScreenshotBytes()
```
* 截屏并获取image.Image对象 （如果需要不同的图片编码，比如jpeg/png，使用image.Image可以帮你做到转换）
```go
img, format, err := d.Screenshot() // img 是 image.Image 对象, format的示例: "jpeg"
```

### UI Hierarchy
* 获取UI结构的XML文本
```go
content, err := d.DumpHierarchy(false, false) // content 是文本
```
* 将UI结构的XML文本转换成 *xmlquery.Node 对象. (我们如果要基于snaphot执行xpath查询，那 *xmlquery.Node 对象是非常有用的 - 它不涉及Uiautomator调用，所以速度会非常快，适合广告页面识别关闭的场景)
```go
doc, err := FormatHierachy(content) // doc 是 *xmlquery.Node 对象
```

### Touch
模拟“手指按下触屏”，“手指按住触屏拖动”，以及“手指离开触屏”
* 获取 touch 对象
```go
touch := d.Touch()
```
* 手指按下触屏 - 在某个位置
```go
/* （相对于屏幕左上角）手指按下触屏，位置为
 *  x: 50% 宽度坐标
 *  y: 60% 高度坐标
 */ 
err := touch.Down(0.5, 0.6) 
```
* 手指按住触屏拖动 - 到某个位置
```go
/* （相对于屏幕左上角）然后拖动到下述位置
 *  x: 50% 宽度坐标
 *  y: 10% 高度坐标
 */ 
err := touch.Move(0.5, 0.1) 
```
* 手指离开触屏
```go
/* （相对于屏幕左上角）然后手指离开触屏，位置为
 *  x: 50% 宽度坐标
 *  y: 10% 高度坐标
 */ 
err := touch.Up(0.5, 0.1) 
```

### Click
在指定坐标位置点击屏幕

* 使用百分比 - 如果x,y的任意一个值在 [0,1) 的范围内，它就是表示百分比，反之如果在[1, maxScreenWidth/maxScreenHeight]的范围内，它就是绝对坐标值
```go
/* 相对于屏幕左上角点击
 *  x: 48.1% 宽度坐标
 *  y: 24.6% 高度坐标
 */ 
err := d.Click(0.481, 0.246)
```
* 使用绝对坐标 - 如果x,y的任意一个值在 [0,1) 的范围内，它就是表示百分比，反之如果在[1, maxScreenWidth/maxScreenHeight]的范围内，它就是绝对坐标值
```go
// 相对于屏幕左上角点击(x: 481, y: 246) 
err := d.Click(481, 246)
```

### Double Click
双击
```go
err := d.DoubleClickDefault(0.481, 0.246)
```

### Long Click
长按点击，相当于按下和松开之间隔了一个给定时间，默认0.5s。函数名后面带Default的一般有一个非Default版本，提供更多参数选择。
```go
err := d.LongClickDefault(0.481, 0.246)
```

### Swipe
滑动
* 从起始点 (fx, fy) 滑动到终点 (tx, ty)
```go
var fx, fy, tx, ty float32 = 0.5, 0.5, 0, 0
err := d.SwipeDefault(fx, fy, tx, ty)
```
* 多点滑动, 参数可以指定多个apis.Point4Swipe对象，代表滑动途径的坐标点
```go
// 滑动途经 起点(x=width*0.5, y=height*0.9) 到 终点(x=width*0.5,y=height*0.1)，滑动总时长0.1秒
err := d.SwipePoints(0.1, apis.Point4Swipe{0.5, 0.9}, apis.Point4Swipe{0.5, 0.1})
```
* 从起点(fx, fy) 拖动到 终点(tx, ty)
```go
var fx, fy, tx, ty float32 = 0.5, 0.5, 0, 0
err := d.DragDefault(fx, fy, tx, ty)
```

### Set Orientation
设置屏幕方向，接受下面这四种参数:
* "n" - 代表正常的竖屏
* "l" - 代表朝左的横屏
* "u" - 代表倒过来的竖屏
* "r" - 代表朝右的横屏
```go
err := d.SetOrientation("n")
```

### Open Quick Settings
打开[快速设置菜单](https://www.lifewire.com/quick-settings-menu-android-4121299)
```go
err := d.OpenQuickSettings()
```

### Open Url
打开浏览器输入URL，打开网页
```go
err := d.OpenUrl("https://bing.com")
```

### Show float window
显示弹窗。这个操作是openatx特有的 - 它打开一个浮层来做应用保活。在这个SDK的产品级部署上面，我们在每次流程开始时都会用d.Shell函数打开ATX应用，再使用自研的分辨率无关的控件图像识别找到“启动UIAUTOMATOR”按钮并点击，再打开ATX应用用同样的方法找到“开启悬浮窗”按钮并点击。以避免在流程执行的过程中，由于uiautomator已经被杀死而重新唤起它 - 虽然这个操作会自动发生，但是有时候会很慢。
```go
err := d.ShowFloatWindow(true)
```


## Input Method
用来输入文字，(会使用一个特殊的输入法)
* 清除文字
```go
err := d.ClearText()
```
* 发送动作
```go
err := ime.SendAction(SENDACTION_SEARCH)
```

支持下面的动作：
```go
SENDACTION_GO       = 2
SENDACTION_SEARCH   = 3
SENDACTION_SEND     = 4
SENDACTION_NEXT     = 5
SENDACTION_DONE     = 6
SENDACTION_PREVIOUS = 7
```

* 发送按键 - 输入文字（包括中文以及其他Unicode）
```go
err := ime.SendKeys("aaa", true)
```

## XPATH
在Android自动化中，XPATH是最重要的寻找UI元素的方法。快速又强大

### Finding elements
* 通过Xpath查找多个元素
```go
els := d.XPath(`//*[@text="your-control-text"]`).All()
for _, el := range els {
  ...
}
```
* 通过Xpath查找一个元素（首个或者nil）
```go
el := d.XPath(`//*[@text="your-control-text"]`).First()
```
* 通过Xpath检查元素是否存在, el.First()函数当元素不存在的时候返回nil
```go
if d.XPath(`//*[@text="your-control-text"]`).First() != nil {
  ...
}
```
* 等待元素出现
```go
el := d.XPath(`//*[@text="your-control-text"]`).Wait(time.Minute)
if el == nil {
  log.Println("element doesn't appear within 1 minute")
  return
}
...
```
* 等待元素消失
```go
ok := d.XPath(`//*[@text="your-control-text"]`).WaitGone(time.Minute)
if !ok {
  log.Println("element doesn't disappear within 1 minute")
  return
}
```
* 如果我们基于UI结构的snapshot来做Xpath查询, 这样，基于一次uiautomator的调用，我们可以执行多次Xpath查询，这在需要广告识别的场景非常有效
```go
content, _ := d.DumpHierarchy(false, false) // content 是文本
doc, _ := FormatHierachy(content) // doc 是 *xmlquery.Node 对象
...
el := d.XPath2(`//*[@text="your-control-text"]`, doc).First()
```

### Xpath elements API
* Xpath元素的子元素（子级后代）
```go
// el := d.XPath(`//*[@resource-id="com.taobao.taobao:id/rv_main_container"]`).First()
children := el.Children()
for _, c := range children {
  ...
}
```
* Xpath元素的兄弟元素
```go
// el := d.XPath(`//*[@resource-id="com.taobao.taobao:id/rv_main_container"]/android.widget.FrameLayout[1]`).First()
siblings := el.Siblings()
for _, s := range siblings {
  ...
}
```
* 通过Xpath查找Xpath元素的后代元素（它基于查找Xpath元素时获取的UI结构snapshot，所以不会涉及Uiautomator调用）
```go
// el := d.XPath(`//*[@resource-id="com.taobao.taobao:id/rv_main_container"]`).First()
children := el.Find(`//android.support.v7.widget.RecyclerView`)
for _, c := range children {
  fmt.Println(*c.Info())
}
```
* 获取 bounding rect

获取控件的边框，注释中解释了返回值类型
```go
bounds := el.Bounds()
/* bounds 的类型为 *apis.Bounds:
 * type Bounds struct {
 *	  LX int // 左上角x
 *	  LY int // 左上角y
 *	  RX int // 右下角x
 *	  RY int // 右下角y
 * }
 */
```
* 获取 Rect

也是获取控件的边框，注释中解释了返回值类型
```go
rect := el.Rect()
/* rect 的类型为 *apis.Rect:
 * type Bounds struct {
 *	  LX int      // 左上角x
 *	  LY int      // 左上角y
 *	  Width int   // 控件宽度
 *	  Height int  // 控件高度
 * }
 */
```
* 获取控件在UI中显示的文本
```go
text := el.Text() // text 是文本
```
* 获取控件的所有信息 - 包括文本和边框，注释中解释了返回值类型
```go
info := el.Info()
/* info 的类型是 *apis.Info
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
* 获取元素中心点位置
```go
x, y, ok := el.Center()
```
* 点击
```go
ok := el.Click()
```
* 在控件中滑动 - 如果控件是(Recycler)List
```go
dir := apis.SWIPE_DIR_LEFT
var scale float32 = 0.8
ok := el.SwipeInsideList(dir, scale)
/* 滑动方向dir 有四种值(int):
 * SWIPE_DIR_LEFT  = 1 // 从右向左滑动
 * SWIPE_DIR_RIGHT = 2 // 从左向右滑动
 * SWIPE_DIR_UP    = 3 // 从下到上滑动
 * SWIPE_DIR_DOWN  = 4 // 从上到下滑动

 scale 是滑动距离的比例，相对于在此滑动方向下，宽度或者高度。比如，再这个例子中，我们是从右到左横向滑动，那么scale = 0.8就意味着滑动80%的宽度距离
*/
```
* 输入文字
```go
ok := el.Type("aaa")
```
* 截图 - 控件截图
```go
img := el.Screenshot() // img 是 image.Image 类型对象
```

## UI Object
通过属性匹配方式查找UI元素. 在大多数平台，包括iOS和Windows UIA, 这种典型的辅助功能API(UI Object)远比Xpath更快. 举例来说, windows UIA, 获取桌面根元素下的XML UI结构异常的慢, 因为获取元素的信息会有大量的夸进程COM调用，它是很慢的. 但是在安卓获取XML UI结构非常的快，基本和UI Object方式一样快. 而且Xpath还更强大，支持更有表现力的查询。

与其你使用UI Object方式查找元素，我会更建议你使用上一节的Xpath方式。我给两种方式提供了类似的元素操作API，比如Info, Click, Wait等。

有些API在UI Object中看起来不那么自然，比如Count，那是因为获取所有元素在UI Object可能会牵涉非常多次的Uiautomator调用。与其让用户因为一个不那么好的查询条件等太长时间，不如把底层暴露给用户，让用户自己决定实现。

### construct query
UI Object API事实上在调用这几个API前是不调用Uiautomator的: 
  1. (*UIObject).Get() *UiElement, 
  2. (*UIObject).Wait(timeout time.Duration) int, 
  3. (*UIObject).WaitGone() bool

* 基于属性值构造UI Object查询
```go
uo := d.UiObject(apis.NewUiObjectQuery("resourceId", `com.taobao.taobao:id/rv_main_container`)) // 这个 api 接受多个 NewUiObjectQuery 参数, AND关系
```
* 构造兄弟元素UI Object查询 - 这个API的Uiautomator返回有点怪，稍微注意一下
```go
c := d.UiObject(apis.NewUiObjectQuery("resourceId", `com.taobao.taobao:id/sv_search_view`)).Sibling(apis.NewUiObjectQuery("className", "android.widget.FrameLayout"))
```
* 构造后代元素UI Object查询 - 为了和uiautomation2 python库尽可能保持API一致，我使用Child的方法名称，而不是Descendant
```go
c := d.UiObject(apis.NewUiObjectQuery("resourceId", `com.taobao.taobao:id/rv_main_container`)).Child(apis.NewUiObjectQuery("className", "android.widget.LinearLayout"))
```
* 构建指数UI Object查询 - 一个查询条件可能返回多个元素，当你通过Count API知道一共有多少个匹配的元素时，你可以指定小于Count值的Index（从0开始），用于指定获取哪个元素
```go
c := d.UiObject(
		apis.NewUiObjectQuery("className", `android.support.v7.widget.RecyclerView`)).Index(0)
```

### execute ui object query
执行UI Object查询
* 获取第一个UI元素
```go
el := d.UiObject(apis.NewUiObjectQuery("resourceId", `com.taobao.taobao:id/rv_main_container`)).Get() // 返回 *apis.UiElement 对象
```
* 获取匹配查询的元素数量
```go
count := d.UiObject(apis.NewUiObjectQuery("resourceId", `com.taobao.taobao:id/rv_main_container`)).Count()
```
* 获取第N个元素 - 代码示例片段中 - 是第三个( .Index(2)，因为index值是从0开始的)
```go
el := d.UiObject(apis.NewUiObjectQuery("resourceId", `com.taobao.taobao:id/rv_main_container`)).Index(2).Get()
```
* 等待元素出现
```go
count := d.UiObject(apis.NewUiObjectQuery("resourceId", `com.taobao.taobao:id/rv_main_container`)).Wait(time.Minute)
// 返回匹配的元素的数量，如果在参数给定的时间内不出现匹配元素，返回 -1
```
* 等待元素消失
```go
ok := d.UiObject(apis.NewUiObjectQuery("resourceId", `com.taobao.taobao:id/rv_main_container`)).WaitGone(time.Minute)
// 如果消失，返回true，反之false
```

### ui object element apis
UI Object元素API
* 获取控件的所有信息 - 包括文本和边框，注释中解释了返回值类型
```go
info := el.Info() // info的类型与xpath元素的同名API相同
```
* 获取 bounding rect

获取控件的边框，注释中解释了返回值类型
```go
bounds := el.Bounds() // bounds的类型与xpath元素的同名API相同
```
* 获取 Rect

也是获取控件的边框，注释中解释了返回值类型
```go
rect := el.Rect() // rect的类型与xpath元素的同名API相同
```
* 获取元素中心点位置
```go
x, y, ok := el.Center()
```
* 点击
```go
ok := el.Click()
```
* 获取控件在UI中显示的文本
```go
text := el.Text() // text 是文本
```
* 在控件中滑动 - 如果控件是(Recycler)List
```go
dir := apis.SWIPE_DIR_LEFT
var scale float32 = 0.8
ok := el.SwipeInsideList(dir, scale)
/* 滑动方向dir 有四种值(int):
 * SWIPE_DIR_LEFT  = 1 // 从右向左滑动
 * SWIPE_DIR_RIGHT = 2 // 从左向右滑动
 * SWIPE_DIR_UP    = 3 // 从下到上滑动
 * SWIPE_DIR_DOWN  = 4 // 从上到下滑动

 scale 是滑动距离的比例，相对于在此滑动方向下，宽度或者高度。比如，再这个例子中，我们是从右到左横向滑动，那么scale = 0.8就意味着滑动80%的宽度距离
*/
```
* 输入文字
```go
ok := el.Type("aaa")
```
* 截图 - 控件截图
```go
img := el.Screenshot() // img 是 image.Image 类型对象
```


**如果想支持作者，左边是微信打赏码；如果想和作者交朋友或者一起做好玩的编程事情，右边是微信加好友的二维码**

| ![image](https://github.com/fantonglang/go-mobile-automation-examples/blob/main/doc/wechat.jpg) | ![image](https://github.com/fantonglang/go-mobile-automation-examples/blob/main/doc/wechat2.jpg) |
| -- | -- |
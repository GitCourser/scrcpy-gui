# scrcpy GUI

scrcpy 无线连接的跨平台图形界面。界面由 [Viva](https://github.com/vivaui/viva) 原生绘制。

## 功能

- 选择或拖入 scrcpy 目录，目录内包含 `adb` 与 `scrcpy` 即可
- 无线连接：IP、端口、启动时检查上次端口、端口范围扫描
- 扫描设置：卡片右上角的图标按钮手动扫描；并发数、超时、屏蔽端口；多个开放端口会复制到剪贴板
- 无线调试配对：配对端口与配对码
- 启动选项：置顶、熄屏且禁止休眠、退出后锁屏、录屏、H.265、禁用音频转发、禁用鼠标悬停、仅录音
- 视频参数：热键、高度、码率、缓冲
- 参数自动保存到 `scrcpy.cfg`
- 跟随系统的明暗主题：系统切换外观时界面立即切换

## 环境要求

- Go 1.26 或更高版本
- viva 的平台依赖：
  - Linux：X11、XCB、Pango、Cairo；推荐安装 Vulkan 驱动
  - Windows：Direct3D 11 与 DirectWrite
  - macOS：Xcode 命令行工具
- `adb` 与 `scrcpy` 放在选定目录，或加入 `PATH`

## 构建与运行

```sh
go build -o scrcpy-gui .
./scrcpy-gui
```

Windows 下生成 `scrcpy-gui.exe`。`-H=windowsgui` 让程序以图形子系统运行，双击时不会附带终端窗口：

```sh
go build -ldflags="-H=windowsgui" -o scrcpy-gui.exe .
```

### 最小体积构建

`-trimpath` 不写入构建机器的绝对路径，`-ldflags="-s -w"` 去掉符号表与调试信息，产物体积约减少三成：

```sh
go build -trimpath -ldflags="-s -w" -o scrcpy-gui .
```

Windows 下同样适用，并保留 GUI 子系统：

```sh
go build -trimpath -ldflags="-s -w -H=windowsgui" -o scrcpy-gui.exe .
```

### 在 Linux 上编译 Windows 版

viva 的 Windows 后端是纯 Go 实现，交叉编译无需安装 MinGW 等额外工具链：

```sh
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -trimpath -ldflags="-s -w -H=windowsgui" -o scrcpy-gui.exe .
```

ARM64 设备将 `GOARCH` 改为 `arm64`：

```sh
CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build -trimpath -ldflags="-s -w -H=windowsgui" -o scrcpy-gui-arm64.exe .
```

运行测试：

```sh
go test ./...
```

## 使用

1. 选择 scrcpy 目录。留空时会在程序所在目录和 `PATH` 中查找 `adb` 与 `scrcpy`。
2. 填写 IP 与端口。勾选“扫描”后，启动时会先检查上次使用的端口；端口未开放时按设定范围扫描。
3. 首次无线调试时勾选“配对”，填写手机显示的配对端口和配对码。
4. 点击“启动”，或在窗口任意位置按回车。程序完成配对与连接后启动 scrcpy，并关闭自身窗口；scrcpy 继续运行。

## 配置

`scrcpy.cfg` 优先读取工作目录，其次读取可执行文件所在目录；不存在时写入工作目录。文件键名如下：

| 键                                                | 含义                              |
| ------------------------------------------------- | --------------------------------- |
| `strPath`                                         | scrcpy 目录                       |
| `strIP` / `strPort`                               | 设备地址                          |
| `isScan`                                          | 启动时扫描                        |
| `strPortRange`                                    | 扫描端口范围，例如 `36000-44000`  |
| `strUnPort`                                       | 屏蔽端口，逗号或空格分隔          |
| `numThread` / `numTimeout`                        | 扫描并发数与超时                  |
| `isTop` / `isSw` / `isPoweroff`                   | 置顶、熄屏、退出后锁屏            |
| `isHEVC` / `isNoAudio` / `isNoMousehover`         | H.265、禁用音频转发、禁用鼠标悬停 |
| `isRecord` / `isRecordAudio`                      | 录屏、仅录音                      |
| `strKey` / `strSize` / `strBitrate` / `strBuffer` | 热键、高度、码率、缓冲            |

默认值不会写入文件。

## 行为说明

- 仅录音开启时只生成 `.opus` 录音，并去掉“熄屏”和“退出后锁屏”参数；录屏与仅录音同时勾选时以仅录音为准。
- 勾选“录屏”后标签变为红色，提示录像已开启。
- 录像与录音保存在工作目录的 `rec/` 下，文件名使用启动时间。
- 启动前会检测是否已有 `scrcpy` 进程在运行，避免重复启动。
- 从终端启动时程序会断开控制台，不占用终端窗口；scrcpy 在启动器退出后继续运行。
- 扫描使用 Go 协程并发连接，结果按端口升序返回。

## 代码结构

| 文件                                      | 职责                                  |
| ----------------------------------------- | ------------------------------------- |
| `main.go`                                 | 应用入口、窗口创建、跟随系统明暗主题  |
| `ui.go`                                   | viva 页面、控件绑定、后台 Effect 接入 |
| `config.go`                               | `scrcpy.cfg` 读取、保存与路径解析     |
| `scan.go`                                 | 单端口与范围扫描                      |
| `launch.go`                               | 参数校验、adb/scrcpy 命令构建与启动   |
| `process_windows.go` / `process_other.go` | 平台相关的子进程窗口处理              |

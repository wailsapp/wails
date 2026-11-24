# Wails3 原生鸿蒙笔记：Android / iOS 对照、ArkWeb 与权限边界

> 范围：本机 Wails3 v3.1.0 派生分支 `openharmony`，基线 `f23bee819`。本文区分源码事实、目标设计和设备验收；不把已有 server HAP、x86_64 模拟器或编译探针当成原生后端通过。OpenHarmony 与商用 HarmonyOS 并非同一套发布矩阵：华为文档用于找 API 与规则，目标 OpenHarmony SDK 的头文件/类型声明以及目标设备行为才是选型依据。

**首验目标：OpenHarmony，使用 Windows DevEco Studio 配套模拟器。** HDC 已连接 `127.0.0.1:5555`，设备报告 `emulator 7.0.0.107`、`x86_64`、API 26；已安装的 `com.wails.hello` 编译 SDK 为 OpenHarmony `26.0.0.105`。现有对照包为 `com.homecast.music`、`com.demo.hmdemo`、`com.wails.hello`，未卸载。arm64 真机与商用 HarmonyOS 不计入本次首验。

## 1. Android 和 iOS 为什么可以

| 责任 | Android | iOS |
|---|---|---|
| 谁拥有 UI、启动和销毁 | Java `MainActivity` 持有 WebView，`WailsBridge` 用 `System.loadLibrary("wails")` 加载 Go 原生库；Activity 管生命周期 | `main.m` 的 `UIApplicationMain` 占用真正的系统主线程；AppDelegate 先建窗口，再后台启动 `WailsIOSMain`；UIKit 持有 WKWebView |
| Go 怎么操作 UI | `application_android.go` 缓存 `JavaVM`；需要时通过 `AttachCurrentThread` 拿当前线程的 `JNIEnv`，用 `Call*Method` 调 Java。能在当前线程调用 Java 方法，不代表能在任意线程操作 WebView：`executeJavaScript` 仍走 `runOnUiThread` | Go 经 cgo 调 Objective-C；`mainthread_ios.go` 将回调投递到 `dispatch_get_main_queue()`，`webview_window_ios.m` 在主队列操作 WKWebView。个别需要同步返回的路径先判断是否已在主线程，否则同步调度；不能套在互相等待的调用链上 |
| 页面如何拿资源（生产模式） | `https://wails.localhost/` 是 WebView 内的入口；`shouldInterceptRequest` / `WebViewAssetLoader` 将请求导给 Go。不是监听这个域名对应的 HTTP 端口 | WKWebView 注册 `WKURLSchemeHandler`，拦截 `wails://`；Go 收到 `ServeAssetRequest` 后把工作交给后台处理，由原生层完成或取消请求 |
| 页面如何调 Go、Go 如何回送 | `addJavascriptInterface` 暴露 `wails`；runtime POST 请求体不经 `shouldInterceptRequest`，走 Java 接口 → JNI → `nativeHandleRuntimeCall`。Go 回送可由 Java 在 UI 线程执行 JS | `WKScriptMessageHandler` 接收 `external`/`wails` 消息 → `HandleJSMessage`；Go 经 Objective-C 在主队列调用 `evaluateJavaScript` |

本地源码：`v3/pkg/application/{application_android.go,application_ios.go,application_ios_delegate.m,mainthread_ios.go,webview_window_ios.m}`；Android 官方构建模板：`v3/internal/commands/build_assets/android/app/src/main/java/com/wails/app/{MainActivity.java,WailsBridge.java}`。

**共同点是宿主接管 UI，Go 承载业务和 Wails runtime，宿主用两条通路分别解决资源与消息。** 原生支持不要求用系统控件重写 Web UI，也不要求 `fetch` 原样充当绑定传输。Android 给了接口覆盖参考；iOS 的主队列、冷启动顺序和请求取消更值得借鉴。两者都不能证明鸿蒙可以直接复制它们的线程或 WebView API。

## 2. 鸿蒙该怎么做：只验证一个 Web 窗口

- **宿主层**：Stage 模型 `UIAbility` 管运行，`WindowStage` 管窗口，页面持有 ArkWeb `WebviewController`。启动拆成两个信号：Go runtime ready、WebView/controller ready；两者齐备且会话仍有效才允许绑定。UI 回调只交接工作，不等待 Go 启动或完成业务。
- **Native 层**：优先试用由 OHOS SDK/Clang 构建的薄 NAPI shim，持有 `napi_env` 对应的引用、TSFN 和队列；与 Go 库用小而明确的 C ABI 交换字节、ID 与状态。不在 Go API 中传 `napi_value`，不让 NAPI 对象跨所属环境使用。把 shim 作为设计选择，不当成已经定位到的 cgo 修复。
- **Go 层**：复用 asset handler、bindings、events。命令投递与有返回值调用分开：`runOnMainThread` 投递即可，不等回包；有结果的调用用 request ID、超时/取消，只有不会阻塞 ArkTS 消费队列的后台路径才可等待。ArkTS → Go 的同步入口只做校验、复制、入队等短操作，耗时服务用异步结果。
- **页面 ↔ 宿主**：候选为 ArkWeb `javaScriptProxy` / `registerJavaScriptProxy` 接收消息，`runJavaScript` 回送结果与事件。代理暴露最小接口；在目标 SDK 核对来源限制，验证可信页面、外部导航、iframe 和重载。外链交系统或隔离 WebView，不给远程页面复用拥有 Go 权限的桥。
- **资源 ↔ 宿主**：优先试 `onInterceptRequest` 响应受控应用 origin 的 GET，再评估 `WebSchemeHandler` 是否需要且受目标 API 支持。最小验收含 HTML、JS 模块、CSS、二进制、MIME/状态码、找不到资源、刷新和取消；绑定 RPC 单独走消息，不预设能拦截 POST body。候选 origin 的安全上下文、同源、CSP、拦截覆盖和是否产生真实外网请求均须实测。若最小实验不成立，再决定换 scheme 还是把本地端口保留为明确的兼容模式。

### Native 桥必须守住的线程和内存规则

`napi_call_threadsafe_function` 是向所属 JS 环境投递；`napi_tsfn_blocking` 是入队行为，不等 JS 执行。建议有界队列、非阻塞投递，调用方能区分 not-ready、queue-full、closing、timeout。关闭时先拒新请求、取消未完成请求；无论队列仍有数据还是收到迟到回包，都只释放/完成一次。`onWindowStageDestroy` 后 WebView 句柄失效；晚到的回包不可命中新建页面的会话。会话级 ID 与请求 ID 的宽度要端到端一致。

本机代码观察（未运行桥）：`bridge_openharmony.go` 的 C 端在入队失败等路径释放 `name/args`，Go 失败分支又释放，须先明确所有权；请求 ID 从 Go `uint64` 缩成 NAPI `uint32`，回收又读 `int32`；`mainthread_openharmony.go` 用会等待结果的 `ohBridgeVoidInt` 做 `dispatchOnMainThread`，如果调度方正是待消费的 ArkTS 线程，会卡到超时。这些与下文已定位的 cgo 前导码问题不同。

## 3. 鸿蒙平台规则：不能搬 AndroidManifest 或把权限藏在 Go 里

| 场景 | 平台规则与适配契约 |
|---|---|
| 权限声明、授权 | `module.json5` 的 `requestPermissions` 声明应用权限；`system_grant` 和 `user_grant` 分开，后者还要由合适的前台上下文按需 `requestPermissionsFromUser`，按目标 SDK 填写 `reason`、`usedScene`。还要核对 APL/受限范围：写了配置不等于普通应用能拿到。Go / C 同属应用沙箱，不能凭 `.so` 绕过权限；模板默认不替业务索要权限 |
| 联网和明文 | `ohos.permission.INTERNET` 是应用联网的系统授权类权限，不弹运行时授权框；原生桥本身不因此要开端口。如果使用 loopback，应用联网权限、ArkWeb 明文/混合内容策略、Network Kit 和 Go `net/http` 是不同层，分别测。旧笔记“localhost 默认放行”和“加 `networkSecurityConfig` 即可”不能当跨设备规则；不要默认开放 LAN 或全局 cleartext |
| 麦克风、摄像头、定位 | 按能力查目标权限表，通常涉及用户授权；Web 页面还要处理 ArkWeb 的来源与对应网页权限回调。应用授权不是网页授权。处理拒绝、不再询问、设置撤销，并将具体状态回传 Go/页面。不在启动时批量弹窗 |
| 通知、后台播放 | 业务使用通知时走 Notification Kit 的授权/开关流程，不套通用权限弹窗。`onBackground` 不是保活承诺，Go goroutine 也不会阻止系统挂起或回收。后台音频要按目标系统核对 AVSession、`audioPlayback` 长时任务及其声明/停止约束，真机测锁屏和停止；不作为 Wails 模板默认能力 |
| 文件、沙箱 | 应用私有数据从 Ability Context 取得 `filesDir/cacheDir`。系统 Picker 返回 URI 及对应授权，不能当任意路径 `os.Open`；需要让 Go 长期使用时，按 API 规则用受控句柄或复制到沙箱，明确关闭、授权期限和大小。不要为了导入一个文件就申请广泛存储权限 |
| 剪贴板 | 读写分开；OpenHarmony 文档记录 API 12 起剪贴板读取增加管控。按设备版本和具体 API 选择系统用户操作/粘贴控件或合规授权，不把读取当任意后台同步调用。不可用就明确返回能力或权限错误 |
| 生命周期 | `onCreate`/`onDestroy` 管 Ability 状态，`onWindowStageCreate`/`onWindowStageDestroy` 管窗口与 WebView 会话，`onForeground`/`onBackground` 管可见状态。后台不等于销毁；销毁回调不能作为唯一持久化机会。首次启动、回到前台、窗口销毁重建和冷启动都须测 |
| 工具链、打包及签名 | 固定 Go port、DevEco/OHOS SDK、NDK 与目标 API。分别验 x86_64 模拟器和 arm64-v8a 真机；检查 HAP 包内 Go 库及 shim、ABI、依赖库名/`DT_NEEDED`，用目标签名配置安装冷启动。已有 Demo 表明 `IMPORTED` 库在当时的配置里没有自动进 HAP；以产物为准，不以 CMake 链接成功为准 |

桌面托盘、任意窗口定位等不是阶段验收项。以后按系统版本、形态和公开 API 做能力矩阵；不因手机上没测通就断言 PC/平板上也不存在，不把不支持伪装成成功。

## 4. 里程碑与接手顺序

| 里程碑 | 过关证据 | 当前 |
|---|---|---|
| Server 基线（独立） | HelloWails HAP 加载页面、资产及 WebSocket 事件 | 资料记录 x86_64 模拟器通过；本文未重跑 |
| M1：Native 构建 | `pkg/application` 用 `GOOS=openharmony GOARCH=amd64 CGO_ENABLED=1`、不带 `server` 编译通过；最小 Go 应用通过 `c-shared` 严格链接 | **通过**；`llvm-readelf -d` 确认依赖 `libace_napi.z.so` 和 `libc.so`；同架构纯 Go `-tags server` 包编译也通过 |
| M2：最小 Native HAP | 不依赖 Wails 应用，ArkTS → Go 入队 → Go 后台 → TSFN → ArkTS 返回；验证多请求、超时、关闭会话后的迟到结果与事件循环响应 | **通过**：OpenHarmony x86_64 模拟器首次运行与强制停止后重新启动均输出 `M2 PASS` |
| M3：原生 HelloWails | 单窗口、本地资源、`Greet` binding、Go → 页面事件；无 `server` tag、无本地 HTTP 监听；新进程启动与前后台回切 | **x86_64 模拟器通过**；两次新进程 PASS，同一进程回切后页面计数继续更新。首次 ArkWeb renderer 有一次连接超时，见下文 |
| M4：ArkWeb 文档会话 | 同进程内同一 ArkWeb 重新加载，或销毁并重建 ArkWeb 组件；新文档复用旧请求 ID，回包不得串到新页面 | **通过**；修复前两种路径都收到 `reply=old`，修复后 `M4 PASS refresh/rebuild` |
| M5：宿主重建 | 同一进程销毁并重建 WindowStage；第二 Stage 重新拦截资源、绑定与事件，复用 ID 的回包指向新会话 | **x86_64 模拟器通过**；`smoke.py` 三用例全过，stage PID `25798` 同进程销毁重建后 `M4 PASS stage` |
| M6：通用壳模板 | 把 M3 测试工程里已验证的 ArkTS/C++ 壳提炼为 `build_assets/openharmony/` 模板；模板生成的 HAP 必须跑通同一条资源/绑定/事件链路 | **x86_64 模拟器通过**；模板 HAP PID `31792` 与冷重启 `1730` 两次 `M3 PASS` |

**M1 实测**：同一工具链、独立临时目录，原桥文件运行 `go tool cgo` 退出 2；只去掉 `*/` 与 `import "C"` 间的空行，退出 0。空行让 Go 解析器没有把 C 前导码注释附到 `import "C"`，不是 NAPI 类型本身不兼容。改回工作树重跑 `go build ./pkg/application` 已越过 cgo 识别，下一批报错包括 `setupSignalHandler` 重复、缺 `fatalHandler` / `logPlatformInfo` / `platformEnvironment` / `newGlobalShortcutImpl` / `newMenuImpl` 等。

**M1 收尾**：排除 native 下重复的桌面信号实现，保留 `server` 行为；补齐平台信息、日志和 Go 内部回调，沿用无桌面菜单/托盘及全局快捷键的能力边界。`bridge_openharmony.go` 显式链接 SDK 的 `-lace_napi.z`。复现环境：`go1.27.4-hmos-devel`，OHOS Clang `--target=x86_64-linux-ohos --sysroot=<SDK>/native/sysroot`。分别执行 `go build ./pkg/application`、`CGO_ENABLED=0 go build -tags server ./pkg/application`；独立临时模块通过本地 `replace` 引入 fork，入口调用 `application.New`，执行 `go build -mod=mod -buildmode=c-shared -ldflags=-extldflags=-Wl,--no-undefined -o probe.so .`，三项均退出 0。该探针只验证编译/链接，不交付 HAP，也没调用运行时桥。

**M2 实测记录**：独立项目 `C:\Users\a1\wails-oh-m2`，复用 `hm-demo` 的 ArkTS/CMake 打包结构，包名 `com.wails.nativebridgeprobe`，不声明 `INTERNET`、不启动 HTTP。Go 后台根据请求 ID 计算结果，经有界 TSFN 回传；页面通过 `WailsM2` 日志报告结果。首次 ArkTS 构建报 `arkts-no-any-unknown`，给 NAPI 调用结果加显式类型后，`hvigorw.bat assembleHap --no-daemon` 成功，HDC 安装 `harmony/entry/build/default/outputs/default/entry-default-unsigned.hap` 成功。`aa start -b com.wails.nativebridgeprobe -a EntryAbility` 返回 `10106102`，开发者模式下锁定屏幕不能自动解锁。解锁后直接重跑启动，查看 `hdc -t 127.0.0.1:5555 shell hilog -x -T WailsM2`，不用重装。探针只验证桥设计；不能替代 fork 自带桥、ArkWeb 或真实页面/Ability 销毁重建的验收。

**M2 运行结果**：解锁后 `aa start` 成功，随后 `aa force-stop com.wails.nativebridgeprobe` 并再次启动。`hilog -x -T WailsM2` 分别在 PID `26997` 和 `32110` 记录 `M2 PASS`，两次均为 `sent=true complete=true pending=true closed=true cancelled=true deadline=true`；ArkTS 计时器分别触发 13 和 14 次。本次关闭的是探针的 TSFN 会话，强停重启验证的是新进程冷启动，并未验证同一进程内 Ability 或 ArkWeb 的销毁重建。

**M3 复现**：测试工程 `v3/tests/openharmony-hellowails/`，包名 `com.wails.nativehello`。使用 M1 的 x86_64 OpenHarmony Go/Clang，在 `go/` 下运行 `go build -mod=mod -buildmode=c-shared -ldflags=-extldflags=-Wl,--no-undefined -o ../harmony/entry/src/main/cpp/libs/x86_64/libhomecast_go.so .`；将同目录生成的 `libhomecast_go.h` 复制到 `harmony/entry/src/main/cpp/include/`。在 `harmony/` 配好本机 SDK 路径，运行 `hvigorw.bat assembleHap --no-daemon`；确认 `entry/build/default/outputs/default/entry-default-unsigned.hap` 包含 `libs/x86_64/libentry.so` 和 `libs/x86_64/libhomecast_go.so`。用 `hdc -t 127.0.0.1:5555 install <HAP>`、`hdc -t 127.0.0.1:5555 shell aa start -b com.wails.nativehello -a EntryAbility`，再查 `hdc -t 127.0.0.1:5555 shell hilog -x -T WailsM3`。成功信号是 `M3 PASS: local HTML/JS/CSS, Greet and Go event`，不是 `Local-origin frontend loaded`。

**M3 结果**：OpenHarmony `7.0.0.107`、x86_64、API 26；PID `24381`、`28147` 分别通过。ArkWeb 拦截 `https://wails.localhost/`、`/style.css`、`/main.js`、`/wails/runtime.js`，直接调用 Go asset handler。真实 Wails runtime 通过 ArkWeb proxy 调用 `Greet`，Go 通过 TSFN 和 `runJavaScript` 回送绑定结果及事件。Home → `aa start` 后仍为 PID `24381`，页面显示 `Hello OpenHarmony!`，两张截图的事件计数从 `135` 到 `137`。`aa force-stop` 后新 PID `28147` 再次 PASS。不带 `server` 标签、不启动本地 HTTP 服务；设备 `netstat -lntp` 快照没有 loopback TCP 监听，但 shell 无权确认所有套接字归属。纯 Go `CGO_ENABLED=0 go build -tags server ./pkg/application` 交叉编译也通过。

**ArkWeb 边界**：此镜像上，未声明 `ohos.permission.INTERNET` 的 HTTPS 虚拟 origin 实验没有进入资源拦截；声明后能拦截，应用仍不启动 HTTP 服务。权限不是无网络活动的保证：仅验证应用资源和绑定不依赖网络。宿主拒绝其他 origin 的资源和顶层导航；未做完整网络隔离验收。声明权限后第一次运行返回了 HTML 200，但 ArkWeb renderer 报 `channel connect timeout(6s)` 并出现白屏；强停后连续两个新进程通过。没有修复该镜像首次 renderer 超时，也未验证同进程销毁重建、arm64、商业 HarmonyOS、多窗口或完整移动 API。

**M4：同进程内的 ArkWeb 文档会话。** 在原生 HelloWails 里让旧页面先发一个延迟 Go 绑定，再做同一 ArkWeb 内重新加载，或销毁并创建新的 ArkWeb 组件和 controller。新文档复用 `lifecycle-reused-id`。修复前两种路径都实际收到 `reply=old`，PID 分别为 `29831`、`13214`。修复后，宿主在导航或组件销毁时清除当前文档会话，新文档在装好监听器后登记自己的会话；Go 回包脚本在执行时再次核对会话，避免同一请求 ID 被新页面的回调消费。旧 Go 调用可以继续执行，只丢弃回包；不是业务取消或撤销。

**M4 一键验收**：沿用 M3 的 OpenHarmony Go、x86_64 `CC`、本机 DevEco SDK 配置和解锁的模拟器。从 `v3/tests/openharmony-hellowails/` 执行 `python smoke.py --go <OpenHarmony-go.exe> --hvigor <hvigorw.bat>`（HDC 在 PATH；默认目标 `127.0.0.1:5555`）。脚本只重新安装和启动 `com.wails.nativehello`，判定每个用例前后 PID 不变、文档会话不同、请求 ID 相同，同时检查新 `Greet` 和至少两个递增事件。一次完整运行输出 `M4 PASS refresh`（PID `15978`）和 `M4 PASS rebuild`（PID `17594`）。脚本不重试失败的 ArkWeb 启动。范围限单窗口的文档和 ArkWeb 组件，不包含整个 UIAbility/WindowStage 销毁重建、renderer 恢复或真机兼容性。

**M5 arm64 构建通过，真机未验。** 以 `eece2a12f` 的原生 HelloWails 源码，在单独副本将 `harmony/entry/build-profile.json5` 的 `abiFilters` 设成 `["arm64-v8a"]`。沿用 M3 命令，将 Go 架构改为 `GOARCH=arm64`，Clang 改为 `--target=aarch64-linux-ohos`，Go 库输出到 `cpp/libs/arm64-v8a/libhomecast_go.so`，生成头仍复制到 `cpp/include/`。`c-shared` 严格链接及 `assembleHap --no-daemon` 均成功。检查未签名 HAP，`libs/arm64-v8a/` 下的 `libentry.so`、`libhomecast_go.so` 和 `libc++_shared.so` 的 ELF `e_machine` 均为 `183`（AArch64）。未安装到真机，不能推断鸿蒙 3.0 兼容性。

**M5 完整宿主重建未通过。** 在独立测试包 `com.wails.nativehostprobe` 声明一个不创建 ArkWeb 的 `KeeperAbility`，先启动它，再启动 `EntryAbility`；让 `EntryAbility` 调用 `context.terminateSelf()`、记录 `onWindowStageDestroy` / `onDestroy`，随后重新启动 `EntryAbility`。现有实现在 PID `22185` 记录销毁后进程退出，不能据此认定具体退出原因。仅把 TSFN 从页面提升到 ArkTS 模块的实验保住了 PID `29404`，且记录同 PID 的 Stage destroy → Stage create；但第二次创建时报 `Unable to initialize the Wails native dispatcher`，未恢复绑定或事件。这个实验未并入主线。

**M5 宿主交接已并入主线并在模拟器通过。** 页面级调度改为进程级：ArkTS 模块持有 `BridgeState`（`opened` 标志和当前 `dispatch` 回调），存进 `AppStorage`；模块级 `dispatchBridge` 始终自己处理 `runOnMainThread`，无页面时对其余命令回空结果；`aboutToAppear` 只在首次 `native.open`，之后仅重新挂接 `dispatch`；`aboutToDisappear` 只解绑仍属于本页的回调，不再 `native.close`。为验收新增 `KeeperAbility`（不建 ArkWeb、只保进程）和 `--ps wailsLifecycle stage [--pb closeStage true]` 模式。`smoke.py` 三次启动全过：refresh PID `24279`、rebuild PID `25374`、stage PID `25798`；第二轮独立运行 refresh `26137`、rebuild `29081`、stage `29492`，三项同样全过。stage 用例内：旧 Stage 载入 `?probe=stage&phase=old`，4 秒后 `terminateSelf`，同一 PID 记录 `Stage destroy`/`Ability destroy`；新 Stage 重建页面，重新拦截资源，再次 `M3 PASS`，复用 `lifecycle-reused-id` 的回包解析为新会话结果 `M4 PASS stage`，旧回包未串入。旧的只提升 TSFN 到模块层实验（PID `29404`）失败于二次 `native.open`；仅提升位置不够，必须把打开状态放进 `AppStorage` 才能跨 Stage 保留。

**首次 ArkWeb 冷启动易冻结。** 安装后第一次启动可能触发系统 `appfreeze`（`THREAD_BLOCK_6S`，主线程阻塞于 `libarkweb_engine.so` 初始化），进程被系统杀，`pidof` 变空、页面不出现。本机表现：16:14/16:15 两次首启失败（PID `21108`、`21185`，有 faultlog），预热后连续三个新进程（`24279`/`25374`/`25798`）全过。这与 M3 记录的首次 renderer 超时同源，不是桥的回归；失败后重跑一次即可，smoke 不自动重试，失败会如实报出。

**M6：通用壳模板。** M3–M5 的宿主一直住在 `v3/tests/openharmony-hellowails/` 里，协议是测试工程自己定的（`libentry.so` 导出 `open/result/close/start/ready/asset/request/message/runOnMainThread`，Go 侧主包导出 `WailsBridgeOpen/Result/Close` 等）。产品化要把它变成 `v3/internal/commands/build_assets/openharmony/` 模板，让 `wails3 build` 能直接生成鸿蒙工程。第一版模板曾另起一套协议（`libwails.so` + `init/bridgeResult/handleMessage/serveAsset`），与设备上跑通的 M3 协议不兼容，已推倒重写；**模板必须复用已验证协议，不要另造**。

模板内容：`main_openharmony.go`（经 Go build overlay 注入应用主包，`init()` 调 `application.RegisterOpenharmonyMain(main)`，并导出 C ABI）、`app/src/main/cpp/napi_init.cpp`（把 Go 导出注册为 NAPI 模块）、`app/src/main/cpp/CMakeLists.txt`（`IMPORTED` + `IMPORTED_NO_SONAME` 引入 Go 库，POST_BUILD 复制进 HAP）、`app/src/main/ets/`（`pages/Index.ets`、`bridge/WailsBridge.ets`、`common/WailsContext.ets`、`entryability/EntryAbility.ets`）、`app/src/main/module.json5`、`app/src/main/cpp/types/libentry/`。`OpenHarmonyOverlayGen` 与 Android/iOS 的 overlay 同构，接线在 `manifest_pipeline.go` 的 `assets` 分支 `case "openharmony"`。

为此在 `pkg/application` 补了两个导出：`WailsOhNativeStartMain`（启动 `RegisterOpenharmonyMain` 注册的 main；此前 `openharmonyMainFunc` 只注册、从未被调用）和 `WailsOhNativeResolveRuntimeCall`（把异步 runtime call 的回包送回发起请求的 ArkWeb 文档，按文档会话校验，旧文档的回包不会命中新文档）。

**M6 编译与打包实测**：Go 侧 `go build -mod=mod -overlay=<overlay.json> -buildmode=c-shared -ldflags=-extldflags=-Wl,--no-undefined` 退出 0；`llvm-nm -D` 实测导出 `WailsBridgeOpen/WailsBridgeResult/WailsBridgeClose/WailsStart/WailsReady/WailsAsset/WailsMime/WailsRequest/WailsMessage/WailsRunMain` 与 `WailsOhNativeInit/WailsOhNativeStartMain/WailsOhNativeOnBridgeResult/WailsOhClose` 全部存在。`hvigorw.bat assembleHap --no-daemon` 构建成功；未签名 HAP 的 `module.json` 只含 `EntryAbility`，`libs/x86_64/` 下有 `libentry.so`、`libwails_go.so`、`libc++_shared.so`。

**M6 ArkTS 编译踩坑（可复用）**：首次 `CompileArkTS` 报 6 个错，同属三类。其一，`arkts-no-untyped-obj-literals`：对象字面量必须对应显式声明的 interface，并用 `as` 标注（如 `replyJSON(callId, { dark: false } as DarkModeResult)`）。其二，`arkts-no-any-unknown`：`import { asset } from 'libentry.so'` 的具名导入拿不到类型，`asset()` 退化为 `any`；必须写成 `import * as native from 'libentry.so'` 再标注 `const served: native.Asset | null = native.asset(path)`。其三，`cpp/types/libentry/` 下除了 `Index.d.ts` 还必须有 `oh-package.json5`（`"types": "./Index.d.ts"`），否则 d.ts 不生效，导入报 "module for 'libentry.so' is not verified"。三条都是 ArkTS 静态规则，与桥的实现无关。

**M6 设备验收**：在独立副本套用模板（`build_assets/openharmony/` 全量覆盖 M3 测试工程的 `ets/`、`cpp/`、`module.json5`，Go 侧改用 overlay 注入模板 `main_openharmony.go`），重新编译 Go 库并打包 HAP 后装到模拟器。首次启动复现历史同型的 `appfreeze`（`appfreeze-com.wails.nativehello-20020067-20261124161738949.log`，进程被系统杀、`pidof` 变空），与 M5 记录的首启冻结同源，非桥回归；重跑一次即过。第二次启动 PID `31792`：`Stage create` → `Runtime ready; creating local-origin ArkWeb` → `Controller attached` → 依次 `Intercept` `/`、`/style.css`、`/main.js`、`/wails/runtime.js` → `Local-origin frontend loaded` → `M3 PASS: local HTML/JS/CSS, Greet and Go event`。资源来自 Go asset handler，绑定经真实 Wails runtime，Go 事件回到页面，全程无本地 HTTP 监听。`aa force-stop` 后冷重启新 PID `1730` 再次 `M3 PASS`。脚本 `v3/tests/openharmony-hellowails/verify_template.py` 一键复验同样通过（PID `29337`，`TEMPLATE PASS`）；脚本对 hdc 在 Windows 上把 `C:/` 绝对路径拼到 cwd 的问题改用切到 HAP 目录只传文件名，并对 force-stop 后过早起 Ability 导致进程存活但 onCreate 不执行的情况加了 3 秒稳定等待与最多 3 次重试。设备 `netstat -lntp` 快照没有 loopback TCP 监听（shell 无权确认所有套接字归属，结论限于未观察到本地监听）。首启冻结、arm64 真机与鸿蒙 3.0 兼容仍未处理。

**商用 HarmonyOS 仅编译验证（未安装、未运行）。** 为确认模板也能产出商用版本 HAP，在独立副本把 `harmony/build-profile.json5` 的 `runtimeOS` 从 `OpenHarmony` 改为 `HarmonyOS`，其余配置不变（本机 DevEco 安装的是 HarmonyOS 26.0.0.105 SDK，apiVersion 26，`sdk-pkg.json` 明示 displayName 为 HarmonyOS 26.0.0；`sdk/default/hms` 提供商用 ETS/native API 声明，`sdk/default/openharmony/native/llvm` 提供共用编译器与 sysroot）。商用模式需设置环境变量 `DEVECO_SDK_HOME=<DevEco Studio>/sdk`，否则 hvigor 报 `Invalid value of DEVECO_SDK_HOME`（OH 模式用 local.properties 的 sdk.dir 即可，不需要该变量）。`hvigorw.bat assembleHap --no-daemon` BUILD SUCCESSFUL，34 任务全执行；日志出现 `Will skip sign 'hos_hap'`，`hos_hap` 为商用 HarmonyOS 包格式标识（OH 模式为 `hap`）。产物 `entry-default-unsigned.hap` 含 `libs/x86_64/{libentry,libwails_go,libc++_shared}.so`、bundle `com.wails.nativehello`、ability `EntryAbility`；`pack.info` 记录 apiVersion compatible/target 均为 26、releaseType Release。**本次仅编译打包，未安装、未在真机或模拟器运行**，不推断鸿蒙 3.0/4.x/5.x 的运行兼容性。

**平台命名口径确认与新增 HarmonyOS 目标（2026-11-24）**：“真 Harmony”正式名称是 **HarmonyOS**（本机 `<DevEco>/sdk/sdk-pkg.json`：`HarmonyOS 26.0.0`，apiVersion 26，version 26.0.0.105；商用 API 声明在 `sdk/default/hms`），**不是 HarmonyOS NEXT**；NEXT 是公开文档里 HarmonyOS 5.0.0 / API 12 时期的早期称呼，全仓与资料包均无该词。起初代码只有一个平台值 `openharmony`（Go fork 的 GOOS 也只支持 `openharmony`，clang 目标为 `x86_64-linux-ohos` / `aarch64-linux-ohos`；`go1.27.4-hmos` / `go-hmos` / `ohos-go` / `hmgo` 只是路径与生态项目名）。**本轮已把商用 HarmonyOS 做成独立编译目标**：平台标识 `harmonyos`（amd64/arm64）、包格式 `hos_hap`，Go 侧继续用 `GOOS=openharmony` 的同一套编译器与 c-shared 流程，差异落在 hvigor 的 `runtimeOS: HarmonyOS` 与包格式；不新增 GOOS，不引入第二套资源目录。默认 targets 不含 harmonyos，需显式 `-targets harmonyos/<arch> -formats hos_hap`；签名沿用 openharmony signingConfigs。商用包仍只承诺构建产出与内容校验，商用真机运行未验证（见边界段）。

**端到端验收（2026-11-24）**：在 `C:\Users\a1\oh-hos-test`（标准 `wails.yaml` + 普通 Wails `main.go` + `frontend/`）执行 `wails3 build -targets harmonyos/amd64 -formats hos_hap`（`WAILS_EXP_USE_WAKE=1`、`DEVECO_SDK_HOME=<DevEco>/sdk`、PATH 需含定制 Go、DevEco `jbr/bin` 与 node）5/5 阶段成功：产物 `bin/hose2e.hap` 11,643,640 字节 ≈ 11.1 MiB，`module.json` 的 `compileSdkType` 为 `HarmonyOS`、`bundleName` 为 `com.wails.hose2e`，`libs/x86_64/` 含 `libentry.so`、`libwails_go.so`、`libc++_shared.so`；嵌入模板的 `build-profile.json5` 仍是 `runtimeOS: OpenHarmony`，商用值只作用于 staged 副本。同工程回归 `openharmony/amd64 -formats hap` 亦 5/5 成功，产物 `compileSdkType` 为 `OpenHarmony`。`wails3 harmonyos devices/deploy/logs` 在 Wake 模式下可用，`harmonyos devices` 见到 `127.0.0.1:5555`。商用 HAP 未装机、未在真机运行。

## 5. 产品化管线端到端实测（2026-11-24）

M6 模板接进 Wake/YAML 构建管线后，用单条 `wails3 build` 命令端到端产出可装机 HAP，并在 DevEco x86_64 模拟器复现 M3 PASS。测试项目在 `C:\Users\a1\oh-e2e-test`（标准 `wails.yaml` + `main.go` + `frontend/`，前端把 M3 的 `index.html` / `main.js` 复制过来以复现同一条 binding/event 链路）。

**精确命令（PowerShell）**：

```powershell
$env:WAILS_EXP_USE_WAKE = '1'
$env:DEVECO_SDK_HOME = 'C:\Program Files\Huawei\DevEco Studio/sdk'
$env:PATH = 'C:\Users\a1\sdk\go1.27.4-hmos\toolchains\81bad8b625b9c845961010ea2af776a3cb426b4a-windows-amd64\bin;' + $env:PATH
wails3 build -targets openharmony/amd64 -formats hap
```

`wails3` 是 fork 构建的 `wails3-oh.exe`（`go build -o wails3-oh.exe ./v3/cmd/wails3`）。`WAILS_EXP_USE_WAKE=1` 打开 Wake/YAML 路由；`DEVECO_SDK_HOME` 缺失报 `hap packaging for openharmony/amd64 requires an OpenHarmony SDK (DevEco)`；`PATH` 前置的定制 Go 为 `go1.27.4-hmos-devel windows/amd64`，stock Go 1.26.5 报 `buildmode c-shared not supported`。`OHOS_BASE_SDK_HOME` 由管线从 `DEVECO_SDK_HOME` 自动注入给 hvigor，无需手工设置。管线还需把 c-shared 生成的头复制到模板 `cpp/include/`，并按 8.3 短路径/整体加引号处理含空格的 CC/sysroot（`quoted.Split` 只认整 token 引号）。

**产物**：`bin/ohe2e.hap`，**11,715,510 字节 ≈ 11.2 MiB**；`python zipfile` 读 `module.json` 得 bundleName `com.wails.ohe2e`、ability `EntryAbility`，`libs/x86_64/` 下为 `libentry.so`、`libwails_go.so`、`libc++_shared.so`。模板产物路径为 `bin/<name>.hap`。

**装机、启动与日志（hilog 标签为 `Wails`）**：

```bash
cd C:/Users/a1/oh-e2e-test
hdc -t 127.0.0.1:5555 install bin/ohe2e.hap
hdc -t 127.0.0.1:5555 shell aa start -b com.wails.ohe2e -a EntryAbility
hdc -t 127.0.0.1:5555 shell hilog -x -T Wails
```

绝对路径 `install C:/.../ohe2e.hap` 会被 hdc 误解析，必须在 HAP 目录用相对路径。模板 ArkTS 用 `hilog.info(0x0, 'Wails', ...)`，与 M3 测试工程的 `WailsM3` 标签不同；`-T WailsM3` 查不到产品化 HAP 的 PASS。实测 PASS 文本：

```
M3 PASS: local HTML/JS/CSS, Greet and Go event
```

同批日志含 `Runtime ready; creating local-origin ArkWeb`、`Controller attached`、`Intercept` `/`、`/style.css`、`/main.js`、`/wails/runtime.js`、`Local-origin frontend loaded`。设备为 OpenHarmony `7.0.0.107`、x86_64、API 26；资源走 Go asset handler，`Greet` 走真实 Wails runtime，Go 事件回页面，无本地 HTTP 监听。首启仍可能 `appfreeze`（`THREAD_BLOCK_6S`），重跑一次即过，与 M5/M6 同源、非桥回归。

**OH 签名 stage 已实现（未端到端验证）。** `v3/internal/commands/openharmony_sign.go` 的 `signOpenHarmony` 从 `signing.openharmony` 读字段，映射为：`certificate` → `.p12` keystore、`entitlements` → `.cer` 签名证书、`provisioning_profile` → `.p7b` profile、`identity`/`key_alias` → keystore 内别名（两者取一）、`credential` → 保存 keystore 密码的**环境变量名**（不是密码本身）。`wails.yaml` 示例：

```yaml
openharmony:
  signing:
    certificate: signing/openharmony.p12
    entitlements: signing/openharmony.cer
    provisioning_profile: signing/openharmony.p7b
    identity: wails
    credential: OHOS_SIGNING_PASSWORD
```

签名时设 `$env:OHOS_SIGNING_PASSWORD = '<p12 password>'`。实现方式是往生成的 hvigor 工程的 staging 副本 `build-profile.json5` 注入 `signingConfigs` 段，再重跑 `hvigorw assembleHap` 产出 signed HAP（产物 `<input>.signed`）；还要求 `DEVECO_SDK_HOME` 已设，且能从 `<DevEco Studio>/tools/hvigor/bin/hvigorw(.bat)` 找到 hvigorw。**验证边界**：本机无 DevEco 签名凭据，只验证注入函数单测 `TestInjectOpenHarmonySigningConfig` PASS；真实签名出包、装机与运行均未端到端验证，上面的 M3 PASS 来自 unsigned HAP。

**边界**：仅 DevEco 模拟器 x86_64 端到端实测；arm64 真机、商用 HarmonyOS、多窗口、完整移动 API 未验证（arm64 与 `runtimeOS: HarmonyOS` 仅构建/编译通过）。仍缺：legacy（无 `wails.yaml`）Taskfile、`openharmony:dev`/dev-run、deploy/device run、真实签名出包与安装、CI 构建与模拟器验收任务（签名边界详见上段）。站点文档：`docs/mpress/content/guides/mobile/openharmony.md`。


## 6. 桌面 Wails Demo 手感回归（2026-11-24）

用本机 `C:\Users\a1\Desktop\wails3-harmony-demo\Harmony_Demo`（包名 `com.wails.harmonydemo`，原生 ArkWeb + Go backend）在同一台 DevEco x86_64 模拟器回归桌面 demo 的交互手感：8 项生成 binding（Greet / Reverse / Add / Fib / Ping / Version / Delay / Echo）、本地 PNG 与 CSS 资源、前端事件全部通过，日志为 `DEMO PASS: bindings, events, local assets and CSS`（PID 3785）。再 `aa force-stop` 冷启动，新 PID 9713 重复同一 PASS。无本地 HTTP 监听：`netstat -an` 的 3 个 TCP LISTEN 在 force-stop 前后完全一致，属于模拟器系统进程。首启仍有已知的 `appfreeze`，重跑一次即过。

**关键修复：前端桥必须在 npm runtime 模块求值前就绪。** `@wailsio/runtime` 在模块顶层通过 `window.wails?.invoke` 捕获消息通道，而 ES module 的 import 先于 `main.js` 模块体执行。原来在模块体里赋值 `window.wails = {invoke: ...}` 太晚，`System.invoke("wails:runtime:ready")` 落入浏览器警告分支；Go 侧 `runtimeLoaded` 一直为 false，`WebviewWindow.ExecJS` 只把回包脚本塞进 `pendingJS`，binding 调用 5 秒超时。修法是在 `frontend/index.html` 的 `<script type="module">` 之前加普通脚本：`window.wails = {invoke: m => window.wailsOpenHarmony.invoke(m)}`。M3 测试工程用动态 `import('/wails/runtime.js')`，所以之前没有暴露这个时序问题。

**验证边界**：只在这台 DevEco x86_64 模拟器与该 demo 工程验证；桌面 Wails 其他 demo、arm64 真机、商用 HarmonyOS 未验。诊断期间在 fork 源码临时加的 hilog 打点已全部删除，工作树相对已提交的 M3 内容无新增改动。

## 参考资料

- 鸿蒙 [UIAbility 生命周期](https://developer.huawei.com/consumer/cn/doc/harmonyos-guides/uiability-lifecycle)、[Node-API 线程安全开发](https://developer.huawei.com/consumer/cn/doc/harmonyos-guides/use-napi-thread-safety)、[Node-API 参考](https://developer.huawei.com/consumer/cn/doc/harmonyos-references/napi)；[Node.js TSFN 语义参考](https://nodejs.org/api/n-api.html#asynchronous-thread-safe-function-calls)（不是鸿蒙支持矩阵）。
- ArkWeb [请求拦截](https://developer.huawei.com/consumer/en/doc/harmonyos-guides/web-scheme-handler)、[页面调用应用侧函数](https://developer.huawei.com/consumer/cn/doc/harmonyos-guides/web-in-page-app-function-invoking)、[WebviewController](https://developer.huawei.com/consumer/en/doc/harmonyos-references/arkts-apis-webview-webviewcontroller)、[组件安全](https://developer.huawei.com/consumer/cn/doc/best-practices/bpta-arkweb-component-security)。
- 权限 [声明权限](https://developer.huawei.com/consumer/cn/doc/harmonyos-guides/declare-permissions)、[system_grant](https://developer.huawei.com/consumer/cn/doc/harmonyos-guides/permissions-for-all)、[user_grant](https://developer.huawei.com/consumer/cn/doc/doccenter-capabilities/permissions-for-all-user)、[HTTP 访问](https://developer.huawei.com/consumer/cn/doc/harmonyos-guides/http-request)。
- [文件 Picker](https://developer.huawei.com/consumer/cn/doc/harmonyos-references/js-apis-file-picker)、[OpenHarmony 剪贴板读取权限](https://gitee.com/openharmony/docs/blob/master/zh-cn/application-dev/basic-services/pasteboard/get-pastedata-permission-guidelines.md)、[后台播放](https://developer.huawei.com/consumer/cn/doc/harmonyos-guides/avsession-background-scene)、[长时任务](https://developer.huawei.com/consumer/cn/doc/harmonyos-guides/continuous-task)。

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

**M5 宿主交接已并入主线并在模拟器通过。** 页面级调度改为进程级：ArkTS 模块持有 `BridgeState`（`opened` 标志和当前 `dispatch` 回调），存进 `AppStorage`；模块级 `dispatchBridge` 始终自己处理 `runOnMainThread`，无页面时对其余命令回空结果；`aboutToAppear` 只在首次 `native.open`，之后仅重新挂接 `dispatch`；`aboutToDisappear` 只解绑仍属于本页的回调，不再 `native.close`。为验收新增 `KeeperAbility`（不建 ArkWeb、只保进程）和 `--ps wailsLifecycle stage [--pb closeStage true]` 模式。`smoke.py` 三次启动全过：refresh PID `24279`、rebuild PID `25374`、stage PID `25798`。stage 用例内：旧 Stage 载入 `?probe=stage&phase=old`，4 秒后 `terminateSelf`，同一 PID `25798` 记录 `Stage destroy`/`Ability destroy`；新 Stage 重建页面，重新拦截资源，再次 `M3 PASS`，复用 `lifecycle-reused-id` 的回包解析为新会话结果 `M4 PASS stage`，旧回包未串入。旧的只提升 TSFN 到模块层实验（PID `29404`）失败于二次 `native.open`；仅提升位置不够，必须把打开状态放进 `AppStorage` 才能跨 Stage 保留。

**首次 ArkWeb 冷启动易冻结。** 安装后第一次启动可能触发系统 `appfreeze`（`THREAD_BLOCK_6S`，主线程阻塞于 `libarkweb_engine.so` 初始化），进程被系统杀，`pidof` 变空、页面不出现。本机表现：16:14/16:15 两次首启失败（PID `21108`、`21185`，有 faultlog），预热后连续三个新进程（`24279`/`25374`/`25798`）全过。这与 M3 记录的首次 renderer 超时同源，不是桥的回归；失败后重跑一次即可，smoke 不自动重试，失败会如实报出。

## 参考资料

- 鸿蒙 [UIAbility 生命周期](https://developer.huawei.com/consumer/cn/doc/harmonyos-guides/uiability-lifecycle)、[Node-API 线程安全开发](https://developer.huawei.com/consumer/cn/doc/harmonyos-guides/use-napi-thread-safety)、[Node-API 参考](https://developer.huawei.com/consumer/cn/doc/harmonyos-references/napi)；[Node.js TSFN 语义参考](https://nodejs.org/api/n-api.html#asynchronous-thread-safe-function-calls)（不是鸿蒙支持矩阵）。
- ArkWeb [请求拦截](https://developer.huawei.com/consumer/en/doc/harmonyos-guides/web-scheme-handler)、[页面调用应用侧函数](https://developer.huawei.com/consumer/cn/doc/harmonyos-guides/web-in-page-app-function-invoking)、[WebviewController](https://developer.huawei.com/consumer/en/doc/harmonyos-references/arkts-apis-webview-webviewcontroller)、[组件安全](https://developer.huawei.com/consumer/cn/doc/best-practices/bpta-arkweb-component-security)。
- 权限 [声明权限](https://developer.huawei.com/consumer/cn/doc/harmonyos-guides/declare-permissions)、[system_grant](https://developer.huawei.com/consumer/cn/doc/harmonyos-guides/permissions-for-all)、[user_grant](https://developer.huawei.com/consumer/cn/doc/doccenter-capabilities/permissions-for-all-user)、[HTTP 访问](https://developer.huawei.com/consumer/cn/doc/harmonyos-guides/http-request)。
- [文件 Picker](https://developer.huawei.com/consumer/cn/doc/harmonyos-references/js-apis-file-picker)、[OpenHarmony 剪贴板读取权限](https://gitee.com/openharmony/docs/blob/master/zh-cn/application-dev/basic-services/pasteboard/get-pastedata-permission-guidelines.md)、[后台播放](https://developer.huawei.com/consumer/cn/doc/harmonyos-guides/avsession-background-scene)、[长时任务](https://developer.huawei.com/consumer/cn/doc/harmonyos-guides/continuous-task)。

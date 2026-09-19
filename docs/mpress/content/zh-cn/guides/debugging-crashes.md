---
title: "调试崩溃"
description: "收集 Wails 应用程序的诊断报告和小型转储"
slug: "guides/debugging-crashes"
sourcePath: "guides/debugging-crashes.md"
---

生产环境出现问题时，日志往往不足以定位原因。本指南介绍 `github.com/wailsapp/wails/v3/pkg/debug` 包，它能生成结构化诊断报告，并在 Windows 上生成可用 WinDbg 或 Visual Studio 打开的小型转储。

这个包刻意保持精简，提供两个入口：

- `debug.Dump(...)` — 将当前进程的 Windows 小型转储写入文件。
- `debug.Report(...)` — 收集详细的诊断快照，可选择包含小型转储。

两者都使用调用用户的令牌运行，无需提升权限、`SeDebugPrivilege` 或 `OpenProcess`。转储当前进程使用 `GetCurrentProcess` 伪句柄，从而绕过对自身进程访问的 ACL 检查。

## 快速参考

```go
import "github.com/wailsapp/wails/v3/pkg/debug"

// 1. Just a minidump (Windows only).
path, err := debug.Dump()

// 2. Dump to a specific path.
path, err := debug.Dump(debug.WithPath("C:\\crashes\\my.dmp"))

// 3. Full-memory minidump (large, gigabytes).
path, err := debug.Dump(debug.WithFullMemory())

// 4. Full diagnostic report, no minidump.
r, err := debug.Report()

// 5. Report + minidump.
r, err := debug.Report(debug.WithDump())

// 6. Report + minidump at a specific path, full-memory.
r, err := debug.Report(
    debug.WithDumpPath("C:\\crashes\\my.dmp"),
    debug.WithDumpFullMemory(),
)
```

`debug.Report` 始终返回 `*CrashReport`，即使部分操作失败也不例外。调用方仍能获得出错前成功收集的信息，因此即使 `err != nil`，报告也不会是 `nil`。

## `CrashReport` 结构

```go
type CrashReport struct {
    Timestamp   time.Time          // when the report was generated
    System      SystemInfo         // OS, hardware, CPU, GPU from doctor
    Build       BuildInfo          // Go version, buildmode, compiler, CGO flag
    Crash       *CrashInfo         // process, memory, modules, env vars
    Diagnostics []DiagnosticResult // doctor's health-check results
    DumpPath    string             // set only if WithDump() was passed
}
```

使用 `encoding/json` 序列化报告以发送到崩溃报告服务；如果只需要部分信息，也可以直接读取字段。

## 与 `PanicHandler` 集成

最实用的模式是在 Wails 捕获 panic 时立即收集报告及可选的小型转储，再将合并后的数据交给自己的报告处理流程：

```go
app := application.New(application.Options{
    PanicHandler: func(pd *application.PanicDetails) {
        // Always: capture a diagnostic snapshot plus a minidump.
        report, reportErr := debug.Report(debug.WithDump())
        if reportErr != nil {
            log.Printf("debug.Report: %v", reportErr)
        }

        // Now you have:
        //   pd.Error, pd.StackTrace, pd.FullStackTrace  — from wails
        //   report.DumpPath                              — minidump (Windows)
        //   report.System / report.Crash / report.Build — context
        //
        // Ship it off (Sentry, S3, support ticket, local log...).
        mycrashservice.Upload(pd, report)
    },
})
```

Wails **不会**自动调用 `debug.Report`。崩溃报告通常涉及用户同意或个人身份信息脱敏，因此应由你的处理函数决定。关于 Wails 何时调用 `PanicHandler`，以及如何覆盖用户自行创建的 goroutine，请参阅 [panic 处理指南](/guides/panic-handling/)。

## 手动调试

没有发生 panic 时，`debug.Dump` 和 `debug.Report` 也很有用：

- 卡死检测：看门狗触发时转储进程，以便在 WinDbg 中准确查看每个线程阻塞在哪里。
- 异常状态报告：例如 goroutine 数量持续增长时，收集报告，让应用继续运行，稍后再调查。
- 按需支持资料：将“报告错误”菜单项连接到 `debug.Report(debug.WithDump())`，把结果附加到用户的支持工单。

## Windows 小型转储实践

小型转储默认采用信息丰富但体积紧凑的标志组合：

- `MiniDumpNormal`
- `MiniDumpWithThreadInfo`
- `MiniDumpWithHandleData`
- `MiniDumpWithUnloadedModules`

生成的转储约为 5–50 MB，保留全部线程堆栈、句柄表和模块列表，足以重建转储时每个 goroutine 的位置。

`debug.WithFullMemory()`，或用于 `Report` 的 `debug.WithDumpFullMemory()`，会添加 `MiniDumpWithFullMemory`，捕获整个进程地址空间。这适合调查某个特定指针为何被破坏，但文件会达到数百 MB 甚至 GB。请仅在深入调查时使用。

可在 WinDbg 中使用 `windbg -z C:\path\to\your.dmp` 打开生成的 `.dmp`，或在 Visual Studio 中选择 File → Open Crash Dump。发布构建通常会剥离符号；要获得有用的堆栈跟踪，应搭配同一次构建的 `.pdb`，或包含调试信息的 Go 二进制文件。

## 非 Windows 平台

在非 Windows 平台上，`debug.Dump` 返回“未实现”错误。Linux 和 macOS 使用不同的事后分析工具链，例如通过 `GOTRACEBACK=crash` 生成核心转储、使用 `rr` 或平台专用崩溃报告工具，无法直接对应 `MiniDumpWriteDump` API。`debug.Report` 仍然可用，除了转储文件之外的信息都能获取。

你可以检查 `Dump` 返回的错误，以便平稳降级：

```go
path, err := debug.Dump()
if err != nil {
    log.Printf("minidump unavailable on %s: %v", runtime.GOOS, err)
    // Continue with just debug.Report() or pd.StackTrace from wails.
}
```

## 安全注意事项

- **不需要提升权限。** 这个包刻意避免使用凭据窃取工具所采用的 `SeDebugPrivilege` 和 `OpenProcess` / 远程 PID 访问方式。它只转储调用进程。
- **转储文件包含内存中的全部内容。** 包括正在使用的凭据、身份验证令牌、用户数据和加密密钥。应像对待完整内存映像一样保护转储：保护传输过程、清理存储中的敏感内容，并考虑在上传给第三方之前进行脱敏。
- **默认路径是 `os.TempDir()`。** 在 Windows 上，它解析为 `%LOCALAPPDATA%\Temp`，属于当前用户，并非所有人都能读取，但文件仍会持续保留。如果要保留转储，请将其移到安全位置。

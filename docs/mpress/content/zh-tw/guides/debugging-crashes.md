---
title: "偵錯當機"
description: "收集 Wails 應用程式的診斷報告與小型傾印"
slug: "guides/debugging-crashes"
sourcePath: "guides/debugging-crashes.md"
---

正式環境發生問題時，記錄往往不足以找出原因。本指南介紹 `github.com/wailsapp/wails/v3/pkg/debug` 套件，它可產生結構化診斷報告，並在 Windows 上產生可用 WinDbg 或 Visual Studio 開啟的小型傾印。

這個套件刻意保持精簡，提供兩個進入點：

- `debug.Dump(...)` — 將目前處理程序的 Windows 小型傾印寫入檔案。
- `debug.Report(...)` — 收集詳細診斷快照，並可選擇包含小型傾印。

兩者都使用呼叫使用者的權杖執行，不需要提升權限、`SeDebugPrivilege` 或 `OpenProcess`。傾印目前處理程序時使用 `GetCurrentProcess` 虛擬控制代碼，可略過存取自身處理程序時的 ACL 檢查。

## 快速參考

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

`debug.Report` 一律傳回 `*CrashReport`，即使部分操作失敗也不例外。呼叫端仍可取得出錯前成功收集的資訊，因此即使 `err != nil`，報告也不會是 `nil`。

## `CrashReport` 結構

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

使用 `encoding/json` 序列化報告，以傳送至當機回報服務；若只需要部分資訊，也可直接讀取欄位。

## 與 `PanicHandler` 整合

最實用的模式是在 Wails 捕捉到 panic 時立即收集報告及選用的小型傾印，再把合併的資料交給自己的回報處理流程：

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

Wails **不會**自動呼叫 `debug.Report`。當機回報通常涉及使用者同意或個人識別資訊遮蔽，因此應由你的處理函式決定。關於 Wails 何時呼叫 `PanicHandler`，以及如何涵蓋使用者自行建立的 goroutine，請參閱 [panic 處理指南](/guides/panic-handling/)。

## 手動偵錯

沒有發生 panic 時，`debug.Dump` 和 `debug.Report` 也很有用：

- 停滯偵測：監控計時器觸發時傾印處理程序，以便在 WinDbg 中確認每個執行緒被什麼阻塞。
- 異常狀態回報：例如 goroutine 數量持續增加時，收集報告，讓應用程式繼續執行，之後再調查。
- 按需支援資料：將「回報錯誤」選單項目連接至 `debug.Report(debug.WithDump())`，並將結果附加到使用者的支援單。

## Windows 小型傾印實務

小型傾印預設採用資訊豐富且體積精簡的旗標組合：

- `MiniDumpNormal`
- `MiniDumpWithThreadInfo`
- `MiniDumpWithHandleData`
- `MiniDumpWithUnloadedModules`

產生的傾印約為 5–50 MB，保留所有執行緒堆疊、控制代碼資料表及模組清單，足以重建傾印當下每個 goroutine 的位置。

`debug.WithFullMemory()`，或用於 `Report` 的 `debug.WithDumpFullMemory()`，會加入 `MiniDumpWithFullMemory`，擷取整個處理程序的位址空間。這適合調查特定指標為何損毀，但檔案會達到數百 MB 甚至 GB。請僅在深入調查時使用。

在 WinDbg 中使用 `windbg -z C:\path\to\your.dmp` 開啟產生的 `.dmp`，或在 Visual Studio 中選擇 File → Open Crash Dump。發行組建通常已移除符號；若要取得有用的堆疊追蹤，請搭配同一次組建的 `.pdb`，或含有偵錯資訊的 Go 二進位檔。

## 非 Windows 平台

在非 Windows 平台上，`debug.Dump` 會傳回「未實作」錯誤。Linux 和 macOS 使用不同的事後分析工具鏈，例如透過 `GOTRACEBACK=crash` 產生核心傾印、使用 `rr` 或平台專用當機回報工具，無法直接對應到 `MiniDumpWriteDump` API。`debug.Report` 仍可使用，除了傾印檔以外的資訊都能取得。

你可以檢查 `Dump` 傳回的錯誤，以便平順地降級處理：

```go
path, err := debug.Dump()
if err != nil {
    log.Printf("minidump unavailable on %s: %v", runtime.GOOS, err)
    // Continue with just debug.Report() or pd.StackTrace from wails.
}
```

## 安全注意事項

- **不需要提升權限。** 此套件刻意避免使用憑證竊取工具採用的 `SeDebugPrivilege` 和 `OpenProcess` / 遠端 PID 存取方式。它只傾印呼叫端處理程序。
- **傾印檔包含記憶體中的所有內容。** 包括使用中的憑證、驗證權杖、使用者資料及加密金鑰。應將傾印視為完整記憶體映像：保護傳輸過程、清理儲存中的敏感內容，並考慮在上傳給第三方之前遮蔽敏感資訊。
- **預設路徑是 `os.TempDir()`。** 在 Windows 上會解析為 `%LOCALAPPDATA%\Temp`。這是個別使用者的路徑，並非所有人都可讀取，但檔案仍會持續保留。若要保留傾印，請移至安全位置。

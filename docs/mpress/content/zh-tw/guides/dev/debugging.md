---
title: "偵錯"
description: "調查問題並分析應用程式效能"
slug: "guides/dev/debugging"
sourcePath: "guides/dev/debugging.md"
---

本指南示範如何使用各種工具檢查並調查 Wails 應用程式中可能存在的效能問題，包括使用

- [`runtime/trace`](https://pkg.go.dev/runtime/trace)建立效能圖表，並在瀏覽器中檢查這些圖表

## 建立效能追蹤

@steps
### 準備應用程式以進行追蹤
請確認程式進入點附近有類似以下內容的程式碼

```go
// Create the file to store our trace data within
traceFile, err := os.Create("trace.out")
if err != nil {
  log.Fatalf("trace.out could not be created: %v", err)
}

// Start the trace
if err := trace.Start(traceFile); err != nil {
  _ = traceFile.Close()
  log.Fatalf("trace.start could not start: %v", err)
}

// Trace cleanup on exit. Alternatively,
defer func() {
  trace.Stop()
  _ = traceFile.Close()
}()

...Start your wails app here...
```

輸出會儲存至目前工作路徑中的`trace.out`，其中包含涵蓋應用程式完整執行期間的指標。預設追蹤資料相當原始；強烈建議進一步閱讀 trace 的相關資料，以加入更多情境資訊，並限制實際記錄的內容。例如：`WithRegion, NewTask, Log`

### 執行應用程式以建立追蹤資料
應用程式執行時會持續產生追蹤資料，直到應用程式結束為止；請在應用程式中執行一些操作，完成後再結束應用程式

### 安裝視覺化輔助工具
部分檢視需要[Graphviz](https://graphviz.org/)。您可以執行`dot -V`來確認是否已安裝

```bash
# macOS
brew install graphviz

# Ubuntu/Debian
sudo apt-get install graphviz

# Arch
sudo pacman -S graphviz
```

### 將追蹤資料視覺化
取得追蹤資料後，即可啟動網頁介面

```bash
go tool trace trace.out
```

這應該會在您的預設瀏覽器中開啟追蹤事件檢視器首頁（建議使用 Chromium 核心的瀏覽器）。

對初次使用者而言，`Syscall profile`畫面可能最實用。此畫面會詳細呈現程式究竟在執行哪些作業，以及各項作業的耗時

@end

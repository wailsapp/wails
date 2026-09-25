---
title: "macOS 上的 Syso 檔案"
description: "疑難排解 macOS 上的 Syso 檔案建置錯誤"
slug: "troubleshooting/mac-syso"
sourcePath: "troubleshooting/mac-syso.md"
---

## 問題

嘗試在 macOS 上建置 Wails 應用程式時，建置會失敗，並顯示類似以下的錯誤：

```
Error: Users/runner/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.23.1.darwin-arm64/pkg/tool/darwin_arm64/link: running clang failed: exit status 1
ld: unknown file type in '/private/var/folders/ml/x_tvfgn50_s7p67dm1ypcqqm0000gn/T/go-link-774134794/000000.o'
clang: error: linker command failed with exit code 1 (use -v to see invocation)
```

## 為什麼會發生這種情況？

針對 Windows 建置時，Wails 會為應用程式圖示、視窗圖示和選單圖示產生 `.syso` 檔案。在 Windows 上建置應用程式時需要這些檔案。 這些 `.syso` 檔案位於專案的根目錄中，可能會在針對 macOS 建置時造成問題。

## 解決方案

若要修正此問題，可以從專案的根目錄移除 syso 檔案；如果是自行產生這些檔案，請將它們命名為 `wails_windows_<arch>.syso`，例如 `wails_windows_arm64.syso`。

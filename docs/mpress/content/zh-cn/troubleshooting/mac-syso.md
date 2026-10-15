---
title: "macOS 上的 Syso 文件"
description: "排查 macOS 上的 Syso 文件构建错误"
slug: "troubleshooting/mac-syso"
sourcePath: "troubleshooting/mac-syso.md"
---

## 问题

尝试在 macOS 上构建 Wails 应用程序时，构建会失败，并显示类似以下内容的错误：

```
Error: Users/runner/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.23.1.darwin-arm64/pkg/tool/darwin_arm64/link: running clang failed: exit status 1
ld: unknown file type in '/private/var/folders/ml/x_tvfgn50_s7p67dm1ypcqqm0000gn/T/go-link-774134794/000000.o'
clang: error: linker command failed with exit code 1 (use -v to see invocation)
```

## 为什么会出现这种情况？

为 Windows 构建时，Wails 会为应用程序图标、窗口图标和菜单图标生成 `.syso` 文件。在 Windows 上构建应用程序需要这些文件。 这些 `.syso` 文件位于项目的根目录中，可能会导致在 macOS 上构建时出现问题。

## 解决方案

要解决此问题，可以从项目的根目录中删除 syso 文件；如果这些文件由你自行生成，请将其命名为 `wails_windows_<arch>.syso`，例如 `wails_windows_arm64.syso`。

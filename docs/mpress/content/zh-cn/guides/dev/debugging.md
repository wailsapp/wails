---
title: "调试"
description: "调查问题并分析应用性能"
slug: "guides/dev/debugging"
sourcePath: "guides/dev/debugging.md"
---

本指南介绍用于检查和调查 Wails 应用中潜在性能问题的各种工具，具体使用

- [`runtime/trace`](https://pkg.go.dev/runtime/trace)创建性能图并在浏览器中检查这些图

## 创建性能跟踪

@steps
### 为应用启用跟踪
在靠近程序入口点的位置，确保存在类似下面的代码

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

输出将保存到工作路径中的`trace.out`，其中的指标涵盖应用的整个运行周期。默认跟踪数据相当原始；强烈建议进一步阅读有关 trace 的资料，以添加更多上下文信息，并限制实际记录的内容。例如：`WithRegion, NewTask, Log`

### 运行应用以创建跟踪数据
应用运行时将持续生成跟踪数据，直至退出。因此，请在应用中执行一些操作，并在完成后退出

### 安装可视化辅助工具
某些视图需要[Graphviz](https://graphviz.org/)。可以运行`dot -V`来验证是否已安装

```bash
# macOS
brew install graphviz

# Ubuntu/Debian
sudo apt-get install graphviz

# Arch
sudo pacman -S graphviz
```

### 可视化跟踪数据
获得跟踪数据后，即可启动 Web 界面

```bash
go tool trace trace.out
```

这应该会使用默认浏览器打开跟踪事件查看器的主页（建议使用基于 Chrome 的浏览器）。

对于首次使用的用户，`Syscall profile`屏幕可能最为实用。它会详细列出程序正在执行的操作及其耗时

@end

---
title: "디버깅"
description: "문제 조사 및 앱 성능 프로파일링"
slug: "guides/dev/debugging"
sourcePath: "guides/dev/debugging.md"
---

이 가이드에서는 다음 도구를 사용하여 Wails 앱에서 발생할 수 있는 성능 문제를 점검하고 조사하는 다양한 방법을 설명합니다.

- [`runtime/trace`](https://pkg.go.dev/runtime/trace)을 사용하여 성능 그래프를 생성하고 브라우저에서 살펴봅니다.

## 성능 트레이스 생성

@steps
### 트레이싱을 위한 앱 준비
프로그램 진입점과 가까운 위치에 다음과 같은 코드가 있는지 확인하세요.

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

앱의 전체 실행 시간에 해당하는 메트릭이 작업 경로의 `trace.out`에 저장됩니다. 기본 트레이스는 상당히 가공되지 않은 상태입니다. 트레이스 관련 자료를 더 살펴보고 상황별 정보를 추가하며 실제로 기록하는 범위를 제한하는 것이 좋습니다. 예: `WithRegion, NewTask, Log`

### 앱을 실행하여 트레이스 생성
앱을 실행하면 종료 시점까지 연속 트레이스가 생성됩니다. 앱에서 몇 가지 작업을 수행한 후 완료되면 종료하세요.

### 시각화 도우미 설치
일부 보기에는 [Graphviz](https://graphviz.org/)가 필요합니다. `dot -V`을 실행하여 설치 여부를 확인할 수 있습니다.

```bash
# macOS
brew install graphviz

# Ubuntu/Debian
sudo apt-get install graphviz

# Arch
sudo pacman -S graphviz
```

### 트레이스 시각화
트레이스 데이터를 사용할 수 있으면 웹 인터페이스를 시작할 수 있습니다.

```bash
go tool trace trace.out
```

기본 브라우저에서 트레이스 이벤트 뷰어 홈페이지가 열립니다. Chrome 기반 브라우저를 사용하는 것이 좋습니다.

처음 사용하는 경우 `Syscall profile` 화면이 가장 유용할 것입니다. 이 화면에서는 프로그램이 정확히 무엇을 하는지와 각 작업의 소요 시간을 세부적으로 확인할 수 있습니다.

@end

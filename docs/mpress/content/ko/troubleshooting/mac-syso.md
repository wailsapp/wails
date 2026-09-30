---
title: "macOS의 Syso 파일"
description: "macOS에서 발생하는 Syso 파일 빌드 오류 해결"
slug: "troubleshooting/mac-syso"
sourcePath: "troubleshooting/mac-syso.md"
---

## 문제

macOS에서 Wails 애플리케이션을 빌드하려고 하면 다음과 유사한 오류와 함께 빌드가 실패합니다:

```
Error: Users/runner/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.23.1.darwin-arm64/pkg/tool/darwin_arm64/link: running clang failed: exit status 1
ld: unknown file type in '/private/var/folders/ml/x_tvfgn50_s7p67dm1ypcqqm0000gn/T/go-link-774134794/000000.o'
clang: error: linker command failed with exit code 1 (use -v to see invocation)
```

## 이 문제가 발생하는 이유는 무엇인가요?

Windows용으로 빌드할 때 Wails는 애플리케이션 아이콘, 창 아이콘 및 메뉴 아이콘에 사용할 `.syso` 파일을 생성합니다. 이러한 파일은 Windows에서 애플리케이션을 빌드하는 데 필요합니다. 이러한 `.syso` 파일은 프로젝트의 루트 디렉터리에 있으며 macOS용으로 빌드할 때 문제를 일으킬 수 있습니다.

## 해결 방법

이 문제를 해결하려면 프로젝트의 루트 디렉터리에서 syso 파일을 제거하세요. 또는 직접 생성하는 경우 `wails_windows_arm64.syso`처럼 파일 이름을 `wails_windows_<arch>.syso`으로 지정하세요.

---
title: "Файлы Syso в macOS"
description: "Устранение ошибок сборки из-за файлов Syso в macOS"
slug: "troubleshooting/mac-syso"
sourcePath: "troubleshooting/mac-syso.md"
---

## Проблема

При попытке собрать приложение Wails в macOS сборка завершается ошибкой, подобной следующей:

```
Error: Users/runner/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.23.1.darwin-arm64/pkg/tool/darwin_arm64/link: running clang failed: exit status 1
ld: unknown file type in '/private/var/folders/ml/x_tvfgn50_s7p67dm1ypcqqm0000gn/T/go-link-774134794/000000.o'
clang: error: linker command failed with exit code 1 (use -v to see invocation)
```

## Почему это происходит?

При сборке для Windows Wails создаёт файлы `.syso` для значка приложения, значка окна и значка меню. Эти файлы необходимы для сборки приложения в Windows. Эти файлы `.syso` находятся в корневом каталоге проекта и могут вызывать проблемы при сборке для macOS.

## Решение

Чтобы устранить эту проблему, можно удалить файлы syso из корневого каталога проекта или, если вы создаёте их самостоятельно, присвоить им имена `wails_windows_<arch>.syso`, например `wails_windows_arm64.syso`.

---
title: "Syso-Dateien unter macOS"
description: "Fehler beim Erstellen aufgrund von Syso-Dateien unter macOS beheben"
slug: "troubleshooting/mac-syso"
sourcePath: "troubleshooting/mac-syso.md"
---

## Problem

Beim Versuch, eine Wails-Anwendung unter macOS zu erstellen, schlägt der Build mit einem Fehler ähnlich dem folgenden fehl:

```
Error: Users/runner/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.23.1.darwin-arm64/pkg/tool/darwin_arm64/link: running clang failed: exit status 1
ld: unknown file type in '/private/var/folders/ml/x_tvfgn50_s7p67dm1ypcqqm0000gn/T/go-link-774134794/000000.o'
clang: error: linker command failed with exit code 1 (use -v to see invocation)
```

## Warum tritt dieses Problem auf?

Beim Erstellen für Windows generiert Wails `.syso`-Dateien für das Anwendungssymbol, das Fenstersymbol und das Menüsymbol. Diese Dateien sind erforderlich, damit die Anwendung unter Windows erstellt werden kann. Diese `.syso`-Dateien befinden sich im Stammverzeichnis des Projekts und können beim Erstellen für macOS Probleme verursachen.

## Lösung

Um dieses Problem zu beheben, entfernen Sie die syso-Dateien aus dem Stammverzeichnis des Projekts. Wenn Sie sie selbst generieren, benennen Sie sie alternativ `wails_windows_<arch>.syso`, z. B. `wails_windows_arm64.syso`.

---
title: "Fichiers syso sous macOS"
description: "Résoudre les erreurs de compilation liées aux fichiers syso sous macOS"
slug: "troubleshooting/mac-syso"
sourcePath: "troubleshooting/mac-syso.md"
---

## Problème

Lorsque vous tentez de compiler une application Wails sous macOS, la compilation échoue avec une erreur semblable à la suivante :

```
Error: Users/runner/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.23.1.darwin-arm64/pkg/tool/darwin_arm64/link: running clang failed: exit status 1
ld: unknown file type in '/private/var/folders/ml/x_tvfgn50_s7p67dm1ypcqqm0000gn/T/go-link-774134794/000000.o'
clang: error: linker command failed with exit code 1 (use -v to see invocation)
```

## Pourquoi cela se produit-il ?

Lors de la compilation pour Windows, Wails génère des fichiers `.syso` pour l’icône de l’application, celle de la fenêtre et celle du menu. Ces fichiers sont nécessaires pour compiler l’application sous Windows. Ces fichiers `.syso` se trouvent dans le répertoire racine du projet et peuvent provoquer des problèmes lors de la compilation pour macOS.

## Solution

Pour résoudre ce problème, vous pouvez supprimer les fichiers syso du répertoire racine du projet ou, si vous les générez vous-même, les nommer `wails_windows_<arch>.syso`, par exemple `wails_windows_arm64.syso`.

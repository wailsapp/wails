---
title: "Arquivos Syso no macOS"
description: "Solucione erros de compilação causados por arquivos Syso no macOS"
slug: "troubleshooting/mac-syso"
sourcePath: "troubleshooting/mac-syso.md"
---

## Problema

Ao tentar compilar um aplicativo Wails no macOS, a compilação falha com um erro semelhante ao seguinte:

```
Error: Users/runner/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.23.1.darwin-arm64/pkg/tool/darwin_arm64/link: running clang failed: exit status 1
ld: unknown file type in '/private/var/folders/ml/x_tvfgn50_s7p67dm1ypcqqm0000gn/T/go-link-774134794/000000.o'
clang: error: linker command failed with exit code 1 (use -v to see invocation)
```

## Por que isso acontece?

Ao compilar para Windows, o Wails gera arquivos `.syso` para o ícone do aplicativo, o ícone da janela e o ícone do menu. Esses arquivos são necessários para compilar o aplicativo no Windows. Esses arquivos `.syso` ficam no diretório raiz do projeto e podem causar problemas durante a compilação para macOS.

## Solução

Para corrigir esse problema, remova os arquivos syso do diretório raiz do projeto ou, se você mesmo os gerar, nomeie-os como `wails_windows_<arch>.syso`, por exemplo, `wails_windows_arm64.syso`.

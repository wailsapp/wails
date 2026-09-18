---
title: "Atualizar a partir de uma versão alpha da v3"
description: "Migrar um projeto Wails v3 alpha existente para uma versão beta fixa"
slug: "migration/alpha-to-beta"
sourcePath: "migration/alpha-to-beta.md"
---

Este guia é destinado a projetos v3 alpha existentes. Para Wails v2, use o [guia de v2 para v3](/migration/v2-to-v3/).

## Antes de atualizar

Faça um commit ou uma cópia de segurança do projeto. Leia o [registro de alterações](/changelog/) entre sua versão alpha e a beta escolhida: o código-fonte, as APIs ou a configuração de compilação podem precisar de alterações. Verifique a [política de compatibilidade para desktop](/status/) e os requisitos da sua plataforma.

Os comandos abaixo usam a versão publicada `v3.0.0-beta.23` como exemplo de versão exata, não como recomendação para acompanhar sempre a versão mais recente. Se escolher outra versão, verifique as versões da CLI, do módulo Go e do runtime npm e atualize os comandos em conjunto. Neste exemplo, a versão npm corresponde à versão Go sem o prefixo `v`.

## 1. Atualizar a CLI

```sh
go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0-beta.23
wails3 version
```

Verifique se `wails3 version` informa a versão instalada. Um executável antigo que aparece antes no `PATH` pode ocultar a nova CLI.

## 2. Atualizar o módulo Go

Execute na raiz do projeto. Revise as alterações nas dependências; não atualize indiscriminadamente módulos não relacionados.

```sh
go get github.com/wailsapp/wails/v3@v3.0.0-beta.23
go mod tidy
```

## 3. Atualizar o runtime frontend

Para projetos que usam npm e um diretório `frontend`:

```sh
cd frontend
npm install --save-exact @wailsio/runtime@3.0.0-beta.23
cd ..
```

Mantenha o arquivo de travamento e revise suas alterações. Se o frontend usa outro gerenciador de pacotes ou diretório, adapte esta etapa mantendo uma versão exata do runtime.

## 4. Gerar novamente, compilar e testar

Na raiz do projeto, gere novamente as vinculações a partir dos serviços Go e compile:

```sh
wails3 generate bindings
wails3 build
```

Execute o aplicativo compilado e teste seus fluxos de trabalho em cada plataforma suportada para a qual você distribui. Revise e faça um commit conjunto do código-fonte, das vinculações geradas, dos arquivos do módulo e das alterações no arquivo de travamento do frontend.

## Se a atualização falhar

Verifique a CLI no `PATH`, a versão do módulo com `go list -m github.com/wailsapp/wails/v3` e o runtime instalado com `npm --prefix frontend ls @wailsio/runtime`. Gere novamente as vinculações após resolver incompatibilidades de versão. Não presuma que toda versão alpha pode ser atualizada sem alterações no código.

Se o problema persistir, [relate uma issue reproduzível](https://github.com/wailsapp/wails/issues/new/choose) com as versões antiga e nova, o erro exato e a saída de `wails3 doctor`. Siga a [política de segurança](https://github.com/wailsapp/wails/blob/master/SECURITY.md) para vulnerabilidades.

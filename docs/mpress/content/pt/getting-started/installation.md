---
title: "Instalação"
description: "Instale o Wails e configure seu ambiente de desenvolvimento"
slug: "getting-started/installation"
sourcePath: "getting-started/installation.md"
---

## Plataformas compatíveis

- Windows AMD64/ARM64
- macOS 10.15+ AMD64 (pode ser implantado no macOS 10.13+)
- macOS 11.0+ ARM64
- Ubuntu 24.04 AMD64/ARM64 (outras distribuições Linux também podem funcionar!)

## Dependências

O Wails tem várias dependências comuns que são necessárias antes da instalação.

@note{type="tip"}
Depois de instalar a CLI do Wails, você pode executar `wails3 setup` para verificar automaticamente essas dependências e obter ajuda para instalá-las.

@end

@tabs
[Go (pelo menos 1.24)]
Baixe o Go na [página de downloads do Go](https://go.dev/dl/).

Siga as [instruções oficiais de instalação do Go](https://go.dev/doc/install). Verifique também se a variável de ambiente `PATH` inclui o caminho para o diretório `~/go/bin`. Reinicie o terminal e faça as seguintes verificações:

- Verifique se o Go está instalado corretamente: `go version`
- Verifique se `~/go/bin` está na variável PATH
  - Mac/Linux: `echo $PATH | grep go/bin`
  - Windows: `$env:PATH -split ';' | Where-Object { $_ -like '*\go\bin' }`


[npm (opcional)]
Embora o Wails não exija a instalação do npm, a maioria dos modelos incluídos precisa dele.

Baixe o instalador mais recente do Node na [página de downloads do Node](https://nodejs.org/en/download/). É recomendável usar a versão mais recente, pois geralmente é a que testamos.

Execute `npm --version` para verificar.

@note{type="info"}
Se preferir outro gerenciador de pacotes ao npm, você pode usá-lo. Será necessário atualizar os Taskfiles do projeto para utilizá-lo.

@end

@end

## Dependências específicas da plataforma

Você também precisará instalar as dependências específicas da plataforma:

@tabs{sync-key="platform"}
[Mac]
O Wails exige que as ferramentas de linha de comando do xcode estejam instaladas. Para instalá-las, execute:

```sh
xcode-select --install
```

[Windows]
O Wails exige que o [WebView2 Runtime](https://developer.microsoft.com/en-us/microsoft-edge/webview2/) esteja instalado. Ele já estará instalado em quase todas as instalações do Windows. Você pode verificar isso usando o comando `wails doctor`.

[Linux]
O Linux requer as ferramentas de compilação `gcc` padrão, além de `gtk4` e `webkitgtk-6.0`. Após a instalação, execute <code>wails3 doctor</code> para ver como instalar as dependências. A pilha legada GTK3/WebKit2GTK 4.1 continuará disponível por meio de `-tags gtk3` (consulte [Empacotamento no Linux — suporte legado ao GTK3](/guides/build/linux/#legacy-gtk3-support)) até a v3.1. Se sua distribuição ou seu gerenciador de pacotes não for compatível, avise-nos no discord.

@end

## Instalação

Para instalar a CLI do Wails usando Go Modules, execute os seguintes comandos:

```shell
go install -v github.com/wailsapp/wails/v3/cmd/wails3@latest
```

Se quiser instalar a versão de desenvolvimento mais recente, execute os seguintes comandos:

```shell
git clone https://github.com/wailsapp/wails.git
cd wails
cd v3/cmd/wails3
go install
```

Ao usar a versão de desenvolvimento, todos os projetos gerados usarão a diretiva [replace](https://go.dev/ref/mod#go-mod-file-replace) do Go para garantir que utilizem a versão de desenvolvimento do Wails.

## Próximas etapas

Depois de instalar a CLI, execute o assistente de configuração para configurar seu ambiente de desenvolvimento:

```shell
wails3 setup
```

@note{type="caution" title="Experimental"}
O assistente de configuração é novo e foi testado principalmente no Linux. Se encontrar problemas, [relate-os](https://github.com/wailsapp/wails/issues/4904) e siga as etapas de instalação manual de dependências abaixo.

@end

O assistente de configuração fará o seguinte:

- Verificará as dependências da plataforma e ajudará a instalá-las
- Configurará os padrões do projeto (informações do autor e prefixo do ID do pacote)
- Opcionalmente, configurará o Docker para compilações multiplataforma
- Configure a assinatura de código (se necessário)

Consulte o [guia de configuração](/getting-started/setup/) para obter mais detalhes.

## Instalação manual de dependências

Se preferir instalar as dependências manualmente ou se o assistente de configuração não funcionar em seu sistema, siga as instruções específicas da plataforma acima e depois execute:

```shell
wails3 doctor
```

Isso verificará se as dependências corretas estão instaladas e informará o que está faltando.

## O comando `wails3` parece estar ausente?

Se o sistema informar que o comando `wails3` está ausente, verifique o seguinte:

- Certifique-se de ter seguido corretamente o **guia de instalação do Go** acima e de que o diretório `go/bin` esteja na variável de ambiente `PATH`.
- Feche e reabra os terminais atuais para que reconheçam a nova variável `PATH`.

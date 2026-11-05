---
title: "Instalação"
description: "Instale o Wails e configure seu ambiente de desenvolvimento"
slug: "getting-started/installation"
sourcePath: "getting-started/installation.md"
---

## Plataformas compatíveis

- Windows AMD64/ARM64
- macOS AMD64/ARM64. As tarefas de build geradas têm como destino o macOS 12.0 ou posterior; esse é o destino de implantação do build padrão, não uma afirmação de que todas as configurações de build personalizadas são compatíveis com sistemas mais antigos.
- Ubuntu 24.04 AMD64/ARM64 (outras distribuições Linux também podem funcionar!)

## Dependências

O Wails tem várias dependências comuns que precisam estar disponíveis antes da instalação.

@note{type="tip"}
Depois de instalar a CLI do Wails, você pode executar `wails3 setup` para verificar essas dependências automaticamente e obter ajuda para instalá-las.

@end

@tabs
[Go (no mínimo 1.26)]
O Wails v3 seleciona uma cadeia de ferramentas do Go 1.26 ou posterior. Um inicializador do Go mais antigo pode baixá-la automaticamente com `GOTOOLCHAIN=auto`; caso contrário, instale-a explicitamente. A declaração de uma versão mínima do Go inferior no módulo inicial não substitui o requisito da dependência do Wails.

Baixe o Go na [página de downloads do Go](https://go.dev/dl/).

Siga as [instruções oficiais de instalação do Go](https://go.dev/doc/install). Verifique também se a variável de ambiente `PATH` inclui o caminho para o diretório `~/go/bin`. Reinicie o terminal e faça as seguintes verificações:

- Verifique se o Go está instalado corretamente: `go version`
- Verifique se `~/go/bin` está na variável PATH
  - Mac/Linux: `echo $PATH | grep go/bin`
  - Windows: `$env:PATH -split ';' | Where-Object { $_ -like '*\go\bin' }`


[npm (opcional)]
Embora o Wails não exija a instalação do npm, a maioria dos modelos incluídos precisa dele.

Baixe o instalador mais recente do Node na [página de downloads do Node](https://nodejs.org/en/download/). É recomendável usar a versão mais recente, pois normalmente é com ela que realizamos os testes.

Execute `npm --version` para verificar.

@note{type="info"}
Se preferir outro gerenciador de pacotes em vez do npm, você pode usá-lo. Será necessário atualizar os Taskfiles do projeto para utilizá-lo.

@end

@end

## Dependências específicas da plataforma

Você também precisará instalar as dependências específicas da plataforma:

@tabs{sync-key="platform"}
[Mac]
O Wails exige a instalação das ferramentas de linha de comando do Xcode. Para instalá-las, execute:

```sh
xcode-select --install
```

[Windows]
O Wails exige a instalação do [WebView2 Runtime](https://developer.microsoft.com/en-us/microsoft-edge/webview2/). Ele já estará instalado em quase todas as instalações do Windows. Você pode verificar isso usando o comando `wails doctor`.

[Linux]
No Linux, são necessárias as ferramentas de build padrão `gcc`, além de `gtk4` e `webkitgtk-6.0`. Após a instalação, execute <code>wails3 doctor</code> para ver como instalar as dependências. A pilha legada GTK3/WebKit2GTK 4.1 continuará disponível por meio de `-tags gtk3` (consulte [Empacotamento para Linux — compatibilidade legada com GTK3](/guides/build/linux/#legacy-gtk3-support)) até v3.1. Se sua distribuição ou seu gerenciador de pacotes não for compatível, avise-nos no Discord.

@end

## Instalação

Para instalar a CLI do Wails usando Go Modules, execute os seguintes comandos:

```shell
go install -v github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0
```

Para instalar a versão de desenvolvimento mais recente, execute os seguintes comandos:

```shell
git clone https://github.com/wailsapp/wails.git
cd wails
cd v3/cmd/wails3
go install
```

Ao usar a versão de desenvolvimento, todos os projetos gerados usarão a diretiva [replace](https://go.dev/ref/mod#go-mod-file-replace) do Go para garantir que utilizem a versão de desenvolvimento do Wails.

## Próximas etapas

Depois de instalar a CLI, execute o assistente de configuração para preparar seu ambiente de desenvolvimento:

```shell
wails3 setup
```

@note{type="caution" title="Experimental"}
O assistente de configuração é novo e foi testado principalmente no Linux. Se você encontrar problemas, [relate-os](https://github.com/wailsapp/wails/issues/4904) e siga as etapas de instalação manual das dependências abaixo.

@end

O assistente de configuração irá:

- Verificar as dependências da plataforma e ajudar a instalá-las
- Configurar os padrões do projeto (informações do autor e prefixo do ID do pacote)
- Configurar opcionalmente o Docker para builds multiplataforma
- Configurar a assinatura de código (se necessário)

Consulte o [Guia de configuração](/getting-started/setup/) para obter mais detalhes.

## Instalação manual das dependências

Se preferir instalar as dependências manualmente ou se o assistente de configuração não funcionar em seu sistema, siga as instruções específicas da plataforma acima e depois execute:

```shell
wails3 doctor
```

Isso verificará se as dependências corretas estão instaladas e informará o que está faltando.

## O comando `wails3` parece estar ausente?

Se o sistema informar que o comando `wails3` está ausente, verifique o seguinte:

- Confirme que você seguiu corretamente o **guia de instalação do Go** acima e que o diretório `go/bin` está incluído na variável de ambiente `PATH`.
- Feche e reabra os terminais atuais para que eles reconheçam a nova variável `PATH`.

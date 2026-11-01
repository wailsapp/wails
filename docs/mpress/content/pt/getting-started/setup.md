---
title: "Configuração"
slug: "getting-started/setup"
sourcePath: "getting-started/setup.md"
---

@note{type="caution" title="Experimental"}
O assistente de configuração é novo e foi testado principalmente no Linux. Se você encontrar problemas, [relate-os](https://github.com/wailsapp/wails/issues/4904) e siga as [etapas de instalação manual](/getting-started/installation/#platform-specific-dependencies).

@end

## Início rápido

```bash
# 1. Install the Wails CLI
go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0

# 2. Run the setup wizard
wails3 setup
```

O assistente é aberto no navegador e orienta você durante as verificações de dependências, a definição dos padrões do projeto e a configuração opcional de compilação multiplataforma. As instruções acima usam o fluxo de configuração padrão; o fluxo experimental para projetos Wake é descrito abaixo.

Depois disso, você estará pronto para criar seu primeiro projeto:

```bash
wails3 init -n myapp -t vanilla
cd myapp && wails3 dev
```

## O que ele faz

- **Verifica as dependências** — verifica o Go, o npm e as ferramentas da plataforma
- **Configura os padrões** — informações do autor, prefixo do ID do pacote e modelos preferidos
- **Compilações multiplataforma** — configuração opcional do Docker para compilar em qualquer host
- **Assinatura de código** — configuração opcional para macOS, Windows e Linux

A configuração é salva em `$XDG_CONFIG_HOME/wails/defaults.yaml`, com `~/.config/wails/defaults.yaml` como alternativa, e é usada por `wails3 init`. Ela é independente do `build/config.yml` de um projeto.

## Configurar um projeto Wake

Para criar um primeiro aplicativo completo, siga o [tutorial de introdução ao Wake](/experimental/wake/).
Ele orienta você na criação do projeto e em cada tela de configuração, explica o que esperar
do recurso experimental e aborda a migração de uma compilação existente.

Execute a configuração no diretório do projeto com o [experimento Wake ativado](/experimental/yaml-builds/#enable-yaml-builds):

```bash
export WAILS_EXP_USE_WAKE=1
wails3 setup
```

No PowerShell:

```powershell
$env:WAILS_EXP_USE_WAKE = '1'
wails3 setup
```

Isso abre um fluxo voltado para o projeto:

1. **Projeto** — insira os metadados do projeto quando não houver um manifesto nativo.
2. **Aplicação** — escolha GUI, CLI ou servidor, um backend de compilação ou executáveis já compilados. Go puro é uma opção explicitamente experimental.
3. **Destinos** — selecione os sistemas operacionais e as arquiteturas. iOS e Android continuam experimentais; defina as configurações de assinatura deles no YAML.
4. **Pacotes** — escolha os formatos de pacote e se eles devem ser assinados.
5. **Ferramentas** — revise os requisitos do fluxo de compilação e empacotamento selecionado.
6. **Assinatura** — configure a assinatura quando ela estiver selecionada, usando os controles de assinatura existentes.
7. **Salvar** — grave a configuração do projeto e mostre o comando a ser executado.

Os requisitos de ferramentas dependem das suas escolhas. A compilação de aplicações de desktop em Go puro não requer Docker nem um compilador C. As compilações cruzadas com CGo exigem um compilador compatível e as bibliotecas da plataforma de destino, ou uma cadeia de ferramentas em contêiner compatível. O NSIS ainda requer `makensis`; AppImages com GUI compiladas em um host diferente exigem uma árvore de runtime do Linux fornecida. As identidades de assinatura exportadas podem usar os assinadores em Go do Wails, enquanto as identidades do chaveiro nativo ou do repositório de certificados continuam dependendo das ferramentas da plataforma. Consulte os [requisitos de empacotamento e assinatura](/experimental/yaml-builds/#package-and-sign-for-distribution).

Executáveis já compilados ignoram a compilação do Go, a instalação e a compilação do frontend, os bindings e o UPX. Eles ainda exigem as ferramentas de empacotamento e assinatura selecionadas, além de binários de entrada válidos. Consulte os [caminhos de entrada dos binários pré-compilados](/experimental/yaml-builds/#use-prebuilt-desktop-executables). Essa escolha se aplica ao comando exibido após salvar; ela não é uma configuração permanente de compilação no YAML.

Salvar grava `wails.yaml`, atualiza o tipo do aplicativo, a ativação do frontend e as configurações selecionadas do backend, além de substituir a seleção específica de destino/pacote/assinatura de `profiles.setup`. Os outros perfis, comandos do frontend, executáveis complementares, hooks e configurações não relacionadas do projeto são preservados. As configurações do backend se aplicam a todo o projeto; portanto, os outros perfis que as herdam também usam o backend selecionado. Os novos metadados do projeto permanecem no assistente até você selecionar Salvar. A ausência de ferramentas é informada para que você possa instalá-las antes da compilação; salvar uma configuração válida não instala ferramentas nem executa uma compilação.

As compilações assinadas reutilizam as configurações de assinatura ausentes dos seus padrões no nível do usuário, incluindo arquivos de certificado fora do projeto. As configurações explícitas do projeto têm precedência. O assistente de configuração não copia chaves privadas, senhas nem caminhos globais de credenciais para `wails.yaml`. Para chaves do Windows e do Linux configuradas globalmente, forneça a senha por meio do repositório de credenciais existente ou de `WAILS_WINDOWS_CERT_PASSWORD` / `WAILS_PGP_PASSWORD`. Para uma identidade Apple criptografada, defina `darwin.signing.credential` como o nome de uma variável de ambiente que contenha a senha correspondente.

O perfil salvo é executado com `wails3 build setup`; a opção de executável já compilado acrescenta `--skip-build`.

Mantenha o Wake ativado ao executar o comando exibido. Um `wails.yaml` ativo nunca usa Taskfiles como alternativa quando a flag não está definida. O assistente de configuração não converte Taskfiles existentes nem metadados legados: use o [fluxo de migração revisado](/experimental/yaml-builds/#migrate-an-existing-project) antes de configurar esses projetos com o Wake.

## Subcomandos

```bash
wails3 signing      # Configure code signing
wails3 entitlements # Configure macOS entitlements
```

## Está com problemas?

1. Execute `wails3 doctor` para diagnosticar problemas
2. Siga as [etapas de instalação manual](/getting-started/installation/#platform-specific-dependencies)
3. [Relate o problema](https://github.com/wailsapp/wails/issues/4904) e inclua a saída de `wails3 doctor`

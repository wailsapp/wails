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

O assistente é aberto no navegador e orienta você nas verificações de dependências, nas configurações globais padrão dos projetos, nas ferramentas multiplataforma opcionais e nas credenciais de assinatura. Essas preferências podem ser reutilizadas em vários projetos.

Depois disso, você estará pronto para criar seu primeiro projeto:

```bash
wails3 init -n myapp -t vanilla
cd myapp && wails3 dev
```

## O que ele faz

- **Verifica as dependências** - Verifica o Go, o npm e as ferramentas da plataforma
- **Configura os padrões** - Informações do autor, prefixo do ID do pacote e modelos preferidos
- **Builds multiplataforma** - Configuração opcional do Docker para fazer builds a partir de qualquer host
- **Assinatura de código** - Configuração opcional para macOS, Windows e Linux

A configuração é salva em `$XDG_CONFIG_HOME/wails/defaults.yaml`, com fallback para `~/.config/wails/defaults.yaml`, e usada por `wails3 init`. Ela é separada do `build/config.yml` de um projeto.

## Criar um projeto Wake

Com o experimento do Wake ativado, `wails3 init -ui` coleta as opções de compilação específicas do projeto e as salva em `wails.yaml` quando a criação é concluída com êxito. Siga o [tutorial de introdução ao Wake](/experimental/wake/) para escolher um back-end de compilação, os destinos e os formatos de pacote durante a criação de um projeto para desktop. `wails3 setup` mantém sua função de armazenar as preferências globais.

Os projetos Wake existentes são configurados em `wails.yaml`. Use o [guia de build com YAML](/experimental/yaml-builds/) para alterar seus perfis e suas configurações de empacotamento. Antes de usar o Wake, os projetos Taskfile existentes precisam passar pela [migração unidirecional revisada](/experimental/wake/#migrate-an-existing-project); a configuração não converte nem executa seus Taskfiles.

## Reutilizar preferências de assinatura

Os builds assinados do Wake reutilizam as configurações de assinatura ausentes dos padrões do usuário, inclusive arquivos de certificado localizados fora do projeto. As configurações explícitas do projeto têm precedência. Chaves privadas, senhas e caminhos globais de credenciais não são copiados para `wails.yaml`.

Para chaves do Windows e do Linux configuradas globalmente, forneça a senha por meio do armazenamento de credenciais existente ou de `WAILS_WINDOWS_CERT_PASSWORD` / `WAILS_PGP_PASSWORD`. Para uma identidade Apple criptografada, defina `darwin.signing.credential` como o nome de uma variável de ambiente que contenha sua senha.

## Subcomandos

```bash
wails3 signing      # Configure code signing
wails3 entitlements # Configure macOS entitlements
```

## Está com problemas?

1. Execute `wails3 doctor` para diagnosticar problemas
2. Siga as [etapas de instalação manual](/getting-started/installation/#platform-specific-dependencies)
3. [Relate o problema](https://github.com/wailsapp/wails/issues/4904) incluindo a saída de `wails3 doctor`

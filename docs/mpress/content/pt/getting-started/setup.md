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

O assistente é aberto no navegador e orienta você na verificação das dependências, na definição dos padrões do projeto e na configuração opcional de compilação multiplataforma.

Depois disso, você estará pronto para criar seu primeiro projeto:

```bash
wails3 init -n myapp -t vanilla
cd myapp && wails3 dev
```

## O que ele faz

- **Verifica as dependências** — Verifica o Go, o npm e as ferramentas da plataforma
- **Configura os padrões** — Informações do autor, prefixo do ID do pacote e modelos preferidos
- **Compilações multiplataforma** — Configuração opcional do Docker para compilar em qualquer sistema host
- **Assinatura de código** — Configuração opcional para macOS, Windows e Linux

A configuração é salva em `$XDG_CONFIG_HOME/wails/defaults.yaml`, com `~/.config/wails/defaults.yaml` como alternativa, e é usada por `wails3 init`. Ela é separada do `build/config.yml` de um projeto.

## Subcomandos

```bash
wails3 signing      # Configure code signing
wails3 entitlements # Configure macOS entitlements
```

## Está com problemas?

1. Execute `wails3 doctor` para diagnosticar problemas
2. Siga as [etapas de instalação manual](/getting-started/installation/#platform-specific-dependencies)
3. [Relate o problema](https://github.com/wailsapp/wails/issues/4904) incluindo a saída de `wails3 doctor`

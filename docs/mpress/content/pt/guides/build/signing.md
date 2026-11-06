---
title: "Assinatura de código"
description: "Guia para assinar seus aplicativos Wails em todas as plataformas"
slug: "guides/build/signing"
sourcePath: "guides/build/signing.md"
---

## Como assinar seu aplicativo

Este guia explica como assinar seus aplicativos Wails para macOS, Windows e Linux. O Wails v3 fornece ferramentas de CLI integradas para assinatura de código, notarização e gerenciamento de chaves PGP.

- **macOS** — Assine e notarize seus aplicativos macOS
- **Windows** — Assine seus executáveis e pacotes do Windows
- **Linux** — Assine pacotes DEB e RPM com chaves PGP

## Matriz de assinatura multiplataforma

Esta matriz mostra o que você pode assinar a partir de cada plataforma de origem:

| Formato de destino | A partir do Windows | A partir do macOS | A partir do Linux |
| --- | :---: | :---: | :---: |
| EXE/MSI do Windows | ✅ | ✅ | ✅ |
| Pacote .app do macOS | ❌ | ✅ | ❌ |
| Notarização do macOS | ❌ | ✅ | ❌ |
| DEB do Linux | ✅ | ✅ | ✅ |
| RPM do Linux | ✅ | ✅ | ✅ |

@note{type="tip"}
Os pacotes do Windows e do Linux podem ser assinados a partir de **qualquer plataforma**. A assinatura do macOS requer um Mac devido aos requisitos das ferramentas da Apple.

@end

### Back-ends de assinatura

O Wails seleciona automaticamente o melhor back-end de assinatura disponível:

| Plataforma | Back-end nativo | Back-end multiplataforma |
| --- | --- | --- |
| Windows | `signtool.exe` (Windows SDK) | Integrado |
| macOS | `codesign` (Xcode) | Não disponível |
| Linux | N/A | Integrado |

Quando executado na plataforma nativa, o Wails usa as ferramentas nativas para garantir a máxima compatibilidade. Na compilação cruzada, ele usa o suporte integrado à assinatura.

## Início rápido

A maneira mais rápida de configurar a assinatura é usar o assistente de configuração:

```bash
wails3 setup
```

Isso grava uma configuração de assinatura compartilhada em `~/.config/wails/defaults.yaml`, que é **respeitada no momento da assinatura em todas as plataformas** (consulte [Precedência de configuração](#precedncia-de-configurao)). A etapa de assinatura:

- Detecta quais ferramentas de assinatura estão instaladas para a plataforma de destino — **a partir de qualquer sistema host** — e mostra o comando de instalação para seu sistema operacional (por exemplo, `brew install gnupg`, `sudo apt install osslsigncode`, `winget install GnuPG.Gpg4win`).
- No macOS, lista os certificados Developer ID do seu chaveiro.
- Para Linux, lista suas chaves GPG e pode **gerar e exportar** uma nova.
- Para Windows, pode **gerar um certificado autoassinado** (para testes) por meio do OpenSSL.
- **Armazena senhas com segurança no chaveiro do sistema** (não nos Taskfiles).

@note{type="tip"}
As senhas são armazenadas no repositório de credenciais nativo do sistema (Chaves do macOS, Gerenciador de Credenciais do Windows ou Secret Service do Linux). Isso significa que sua configuração de assinatura permanece segura e funciona em todos os seus projetos Wails.

@end

## Precedência de configuração

Quando você executa uma tarefa de assinatura, cada opção é determinada nesta ordem (a primeira correspondência prevalece):

1. Uma opção explícita passada para `wails3 tool sign` (por exemplo, `--pgp-key`, `--certificate`, `--identity`).
2. A **variável correspondente no Taskfile do projeto** (`PGP_KEY`, `SIGN_CERTIFICATE`/`SIGN_THUMBPRINT`, `SIGN_IDENTITY`, …).
3. A **configuração global** em `~/.config/wails/defaults.yaml` (gravada por `wails3 setup`).

Isso significa que as variáveis do Taskfile são **substituições opcionais**: se uma variável não estiver definida, será usada a chave, o certificado ou a identidade configurada globalmente. Se nenhuma das três opções fornecer um valor, o comando de assinatura exibirá um erro claro, explicando como configurá-lo.

## Configuração por projeto

Para configurar a assinatura de apenas um projeto, em vez de fazê-lo globalmente, execute o assistente específico do projeto no diretório dele — o assistente grava o `vars` nos arquivos `build/<platform>/Taskfile.yml` desse projeto:

```bash
wails3 signing                                   # all detected platforms
wails3 signing --platform windows --platform linux
```

### Configuração manual

Como alternativa, você pode editar manualmente os Taskfiles específicos de cada plataforma. Edite a seção `vars` no início de cada arquivo:

@tabs
[macOS]
Edite `build/darwin/Taskfile.yml`:

```yaml
vars:
  SIGN_IDENTITY: "Developer ID Application: Your Company (TEAMID)"
  KEYCHAIN_PROFILE: "my-notarize-profile"
  # ENTITLEMENTS: "build/darwin/entitlements.plist"
```

Em seguida, execute:

```bash
wails3 task darwin:sign           # Sign only
wails3 task darwin:sign:notarize  # Sign and notarize
```

[Windows]
Edite `build/windows/Taskfile.yml`:

```yaml
vars:
  SIGN_CERTIFICATE: "path/to/certificate.pfx"
  # Or use thumbprint instead:
  # SIGN_THUMBPRINT: "certificate-thumbprint"
  # TIMESTAMP_SERVER: "http://timestamp.digicert.com"
```

A senha é recuperada do chaveiro do sistema (execute `wails3 signing` para configurá-la).

Em seguida, execute:

```bash
wails3 task windows:sign           # Sign executable
wails3 task windows:sign:installer # Sign NSIS installer
```

[Linux]
Edite `build/linux/Taskfile.yml`:

```yaml
vars:
  PGP_KEY: "path/to/signing-key.asc"
  # SIGN_ROLE: "builder"  # Options: origin, maint, archive, builder
```

A senha é recuperada do chaveiro do sistema (execute `wails3 signing` para configurá-la).

Em seguida, execute:

```bash
wails3 task linux:sign:deb       # Sign DEB package
wails3 task linux:sign:rpm       # Sign RPM package
wails3 task linux:sign:packages  # Sign all packages
```

@end

Você também pode inspecionar o estado da assinatura diretamente no sistema:

```bash
# List available macOS code-signing identities
security find-identity -v -p codesigning

# List PGP keys (Linux package signing)
gpg --list-keys
```

Para configurar interativamente a assinatura de todas as plataformas, use o assistente:

```bash
wails3 signing
```

## Assinatura de código no macOS

### Pré-requisitos

- Conta do Apple Developer (US$ 99/ano)
- Certificado Developer ID Application
- Xcode Command Line Tools instaladas

### Identidades de assinatura

Verifique as identidades de assinatura disponíveis:

```bash
security find-identity -v -p codesigning
```

Saída:

```
Found 2 signing identities:

  Developer ID Application: Your Company (ABCD1234) [valid]
    Hash: ABC123DEF456...

  Apple Development: your@email.com (XYZ789) [valid]
    Hash: DEF789ABC123...
```

@note{type="tip"}
Para distribuição fora da App Store, você precisa de um certificado **Developer ID Application**.

@end

### Configuração

Edite `build/darwin/Taskfile.yml` e defina as variáveis de assinatura:

```yaml
vars:
  SIGN_IDENTITY: "Developer ID Application: Your Company (TEAMID)"
  KEYCHAIN_PROFILE: "my-notarize-profile"
  ENTITLEMENTS: "build/darwin/entitlements.plist"
```

| Variável | Obrigatória | Descrição |
| --- | --- | --- |
| `SIGN_IDENTITY` | Sim | Seu Developer ID (por exemplo, "Developer ID Application: Sua Empresa (TEAMID)") |
| `KEYCHAIN_PROFILE` | Para notarização | Nome do perfil do chaveiro com as credenciais armazenadas |
| `ENTITLEMENTS` | Não | Caminho para o arquivo de direitos |

Em seguida, execute:

```bash
wails3 task darwin:sign           # Build, package, and sign
wails3 task darwin:sign:notarize  # Build, package, sign, and notarize
```

### Direitos

Os direitos controlam quais recursos seu aplicativo pode acessar. Aplicativos Wails geralmente precisam de direitos diferentes para desenvolvimento e produção:

- **Desenvolvimento**: requer direitos de JIT, memória não assinada e depuração
- **Produção**: direitos mínimos (apenas acesso à rede)

Use o assistente de configuração interativo para gerar ambos os arquivos:

```bash
wails3 entitlements
```

Isso cria:

- `build/darwin/entitlements.dev.plist` — Para builds de desenvolvimento
- `build/darwin/entitlements.plist` — Para builds de produção/assinados

**Predefinições disponíveis:**\

| Predefinição | Descrição |
| --- | --- |
| Desenvolvimento | JIT, memória não assinada, depuração, rede |
| Produção | Somente rede (mínimo e mais seguro) |
| Ambos | Cria os arquivos de desenvolvimento e produção (recomendado) |
| App Store | Sandbox habilitado com acesso à rede e a arquivos |
| Personalizado | Escolha direitos individuais |

@note{type="note"}
A tarefa `run` no Taskfile do Darwin usa `entitlements.dev.plist` automaticamente. As tarefas `sign` usam `entitlements.plist` para builds de produção.

@end

Em seguida, defina `ENTITLEMENTS` nas variáveis do seu Taskfile para apontar para o arquivo apropriado.

### Notarização

A Apple exige que todos os aplicativos distribuídos sejam notarizados.

@steps
### **Armazene suas credenciais no chaveiro** (configuração única). Execute `wails3 signing` (o comando solicita esses valores e chama `notarytool` internamente) ou chame `notarytool` diretamente:
```bash
xcrun notarytool store-credentials "my-notarize-profile" \
  --apple-id "your@email.com" \
  --team-id "ABCD1234" \
  --password "app-specific-password"
```

### **Defina KEYCHAIN_PROFILE no seu Taskfile** para corresponder ao nome do perfil acima.
### **Assine e notarize seu aplicativo**:
```bash
wails3 task darwin:sign:notarize
```

### **Verifique a notarização**:
```bash
spctl --assess --verbose=2 bin/MyApp.app
```

@end

@note{type="note"}
A notarização normalmente leva 1-2 minutos. O tíquete é anexado automaticamente ao seu aplicativo.

@end

## Assinatura de código no Windows

### Pré-requisitos

- Certificado de assinatura de código (da DigiCert, Sectigo etc.)
- Para assinatura nativa no Windows: Windows SDK instalado (para `signtool.exe`)
- Para assinatura multiplataforma no macOS/Linux: [`osslsigncode`](https://github.com/mtrojnar/osslsigncode) (a etapa de assinatura `wails3 setup` mostra o comando de instalação específico do sistema host)

### Como gerar um certificado autoassinado (teste)

Se você só precisa testar o pipeline de assinatura, execute `wails3 setup`, abra a guia **Windows** da etapa de assinatura e escolha **Gerar um certificado autoassinado**. Isso usa o OpenSSL para criar um `.pfx` de assinatura de código e registra o respectivo caminho na configuração global.

@note{type="caution"}
Certificados autoassinados destinam-se **somente a testes e distribuição interna** — eles acionam avisos do SmartScreen para os usuários finais. Versões públicas precisam de um certificado emitido por uma AC confiável.

@end

### Configuração

Edite `build/windows/Taskfile.yml` e defina as variáveis de assinatura:

```yaml
vars:
  SIGN_CERTIFICATE: "path/to/certificate.pfx"
  # Or use thumbprint instead:
  # SIGN_THUMBPRINT: "certificate-thumbprint"
  # TIMESTAMP_SERVER: "http://timestamp.digicert.com"
```

| Variável | Obrigatório | Descrição |
| --- | --- | --- |
| `SIGN_CERTIFICATE` | Substituição | Caminho para o arquivo de certificado .pfx/.p12 (usa a configuração global se não for definido) |
| `SIGN_THUMBPRINT` | Substituição | Impressão digital do certificado no repositório de certificados do Windows (alternativa a `SIGN_CERTIFICATE`) |
| `TIMESTAMP_SERVER` | Não | URL do servidor de carimbo de data/hora (padrão: http://timestamp.digicert.com) |

@note{type="note"}
Essas variáveis são **substituições** opcionais. Se não forem definidas, será usado o certificado configurado globalmente por meio de `wails3 setup`. Consulte [Precedência de configuração](#precedncia-de-configurao).

@end

@note{type="note"}
A senha do certificado é armazenada no chaveiro do sistema, não no Taskfile. Execute `wails3 signing` para configurá-la ou defina a variável de ambiente `WAILS_WINDOWS_CERT_PASSWORD` na CI.

@end

Em seguida, execute:

```bash
wails3 task windows:sign           # Build and sign executable
wails3 task windows:sign:installer # Build and sign NSIS installer
```

### Assinatura multiplataforma

Os executáveis do Windows podem ser assinados em qualquer plataforma. A mesma configuração e os mesmos comandos do Taskfile funcionam no macOS e no Linux.

### Formatos do Windows compatíveis

| Formato | Extensão | Observações |
| --- | --- | --- |
| Executáveis | .exe | Assinatura PE padrão |
| Instaladores | .msi | Pacotes do Windows Installer |
| Pacotes de aplicativos | .msix, .appx | Aplicativos modernos do Windows |

## Assinatura de pacotes Linux

Os pacotes Linux (DEB e RPM) são assinados com chaves PGP/GPG. Diferentemente da assinatura de código do Windows e do macOS, a assinatura de pacotes Linux comprova que o pacote veio de uma fonte confiável, e não que o sistema operacional confia no código.

### Pré-requisitos

- Par de chaves PGP (pode ser gerado com o Wails)

### Como gerar uma chave PGP

A maneira mais fácil é usar o assistente de configuração: execute `wails3 setup`, abra a guia **Linux** da etapa de assinatura e escolha **Criar uma nova chave GPG**. O assistente:

- gera uma chave RSA 4096 no seu chaveiro GPG (deixe a frase secreta em branco para obter uma chave adequada à CI e que possa ser usada sem supervisão),
- **exporta-a para `~/.wails/signing/<keyid>.asc`** (o arquivo usado para assinar um build) e
- registra tanto o ID da chave quanto o caminho exportado em `~/.config/wails/defaults.yaml`, para que ela seja usada automaticamente no momento da assinatura.

@note{type="note"}
Se uma chave tiver sido configurada **somente pelo ID** (por exemplo, se tiver sido criada antes de o assistente exportar chaves automaticamente), na próxima vez que você acessar a guia Linux, ela será exportada para um arquivo e o caminho será preenchido — nenhuma etapa manual será necessária. Os builds são assinados com um *arquivo* de chave, portanto o que importa é o caminho.

@end

Você também pode fazer isso manualmente com `gpg`:

```bash
# Interactive — the wizard will prompt for name, email, key size and expiry.
gpg --full-generate-key

# Export the key pair to ASCII-armoured files for the Taskfile to consume.
gpg --armor --export-secret-keys "your@email.com" > signing-key.asc
gpg --armor --export "your@email.com" > signing-key.pub.asc
```

Configurações recomendadas: RSA de 4096 bits, validade de 1 ano e proteção com uma senha forte.

@note{type="caution"}
Mantenha sua chave privada segura! Armazene-a de forma criptografada e faça um backup seguro.

@end

### Configuração

Edite `build/linux/Taskfile.yml` e defina as variáveis de assinatura:

```yaml
vars:
  PGP_KEY: "path/to/signing-key.asc"
  # SIGN_ROLE: "builder"  # Options: origin, maint, archive, builder
```

| Variável | Obrigatória | Descrição |
| --- | --- | --- |
| `PGP_KEY` | Substituição | Caminho para o arquivo exportado da chave privada PGP (se não for definido, usa a chave configurada globalmente por meio de `wails3 setup`) |
| `SIGN_ROLE` | Não | Função de assinatura de DEB (padrão: builder) |

@note{type="note"}
A senha da chave PGP é armazenada no chaveiro do sistema, não no Taskfile. Execute `wails3 signing` para configurá-la ou defina a variável de ambiente `WAILS_PGP_PASSWORD` no CI.

@end

Em seguida, execute:

```bash
wails3 task linux:sign:deb       # Build and sign DEB package
wails3 task linux:sign:rpm       # Build and sign RPM package
wails3 task linux:sign:packages  # Build and sign all packages
```

### Funções de assinatura de DEB

Para pacotes DEB, você pode especificar a função de assinatura por meio de `SIGN_ROLE`:

- `origin`: assinatura da origem do pacote
- `maint`: assinatura do mantenedor do pacote
- `archive`: assinatura do mantenedor do repositório
- `builder`: assinatura do responsável pela compilação do pacote (padrão)

### Assinatura multiplataforma

Os pacotes Linux podem ser assinados em qualquer plataforma. A mesma configuração e os mesmos comandos do Taskfile funcionam no Windows e no macOS.

### Visualização das informações da chave

```bash
gpg --show-keys signing-key.asc
```

Saída:

```
pub   rsa4096 2024-01-15 [SC] [expires: 2025-01-15]
      1234 5678 90AB CDEF 1234 5678 90AB CDEF 1234 5678
uid                      Your Name <your@email.com>
```

### Verificação de pacotes Linux

```bash
# Verify DEB signature
dpkg-sig --verify myapp_1.0.0_amd64.deb

# Verify RPM signature
rpm --checksig myapp-1.0.0.x86_64.rpm
```

### Distribuição da sua chave pública

Os usuários precisam da sua chave pública para verificar os pacotes:

```bash
# Export public key for distribution
gpg --armor --export "your@email.com" > myapp-signing.pub.asc

# Users can import it:
# For DEB (apt):
sudo apt-key add myapp-signing.pub.asc
# Or for modern apt:
sudo cp myapp-signing.pub.asc /etc/apt/trusted.gpg.d/

# For RPM:
sudo rpm --import myapp-signing.pub.asc
```

## Integração com o GitHub Actions

Em ambientes de CI, as senhas são fornecidas por meio de variáveis de ambiente, em vez do chaveiro do sistema:

| Variável de ambiente | Descrição |
| --- | --- |
| `WAILS_WINDOWS_CERT_PASSWORD` | Senha do certificado do Windows |
| `WAILS_PGP_PASSWORD` | Senha da chave PGP para pacotes Linux |

Você também pode passar as variáveis do Taskfile diretamente:

```bash
wails3 task darwin:sign SIGN_IDENTITY="$SIGN_IDENTITY" KEYCHAIN_PROFILE="$KEYCHAIN_PROFILE"
```

### Fluxo de trabalho do macOS

```yaml
name: Build and Sign macOS

on:
  push:
    tags: ['v*']

jobs:
  build:
    runs-on: macos-latest
    steps:
      - uses: actions/checkout@v4

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.23'

      - name: Install Wails
        run: go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0

      - name: Import Certificate
        env:
          CERTIFICATE_BASE64: ${{ secrets.MACOS_CERTIFICATE }}
          CERTIFICATE_PASSWORD: ${{ secrets.MACOS_CERTIFICATE_PASSWORD }}
        run: |
          echo $CERTIFICATE_BASE64 | base64 --decode > certificate.p12
          security create-keychain -p "" build.keychain
          security default-keychain -s build.keychain
          security unlock-keychain -p "" build.keychain
          security import certificate.p12 -k build.keychain -P "$CERTIFICATE_PASSWORD" -T /usr/bin/codesign
          security set-key-partition-list -S apple-tool:,apple:,codesign: -s -k "" build.keychain

      - name: Store Notarization Credentials
        env:
          APPLE_ID: ${{ secrets.APPLE_ID }}
          APPLE_TEAM_ID: ${{ secrets.APPLE_TEAM_ID }}
          APPLE_APP_PASSWORD: ${{ secrets.APPLE_APP_PASSWORD }}
        run: |
          xcrun notarytool store-credentials "notarize-profile" \
            --apple-id "$APPLE_ID" \
            --team-id "$APPLE_TEAM_ID" \
            --password "$APPLE_APP_PASSWORD"

      - name: Build, Sign, and Notarize
        env:
          SIGN_IDENTITY: ${{ secrets.MACOS_SIGN_IDENTITY }}
        run: |
          wails3 task darwin:sign:notarize \
            SIGN_IDENTITY="$SIGN_IDENTITY" \
            KEYCHAIN_PROFILE="notarize-profile"

      - name: Upload Artifact
        uses: actions/upload-artifact@v4
        with:
          name: MyApp-macOS
          path: bin/*.app
```

### Fluxo de trabalho do Windows

```yaml
name: Build and Sign Windows

on:
  push:
    tags: ['v*']

jobs:
  build:
    runs-on: windows-latest
    steps:
      - uses: actions/checkout@v4

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.23'

      - name: Install Wails
        run: go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0

      - name: Import Certificate
        env:
          CERTIFICATE_BASE64: ${{ secrets.WINDOWS_CERTIFICATE }}
        run: |
          $certBytes = [Convert]::FromBase64String($env:CERTIFICATE_BASE64)
          [IO.File]::WriteAllBytes("certificate.pfx", $certBytes)

      - name: Build and Sign
        env:
          WAILS_WINDOWS_CERT_PASSWORD: ${{ secrets.WINDOWS_CERTIFICATE_PASSWORD }}
        run: |
          wails3 task windows:sign SIGN_CERTIFICATE=certificate.pfx

      - name: Upload Artifact
        uses: actions/upload-artifact@v4
        with:
          name: MyApp-Windows
          path: bin/*.exe
```

### Fluxo de trabalho multiplataforma (executor Linux)

Assine pacotes do Windows e do Linux usando um único executor Linux:

```yaml
name: Build and Sign (Cross-Platform)

on:
  push:
    tags: ['v*']

jobs:
  build-and-sign:
    runs-on: ubuntu-latest

    steps:
      - uses: actions/checkout@v4

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.23'

      - name: Install Wails
        run: go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0

      - name: Install Build Dependencies
        run: |
          sudo apt-get update
          sudo apt-get install -y nsis rpm

      # Import certificates
      - name: Import Certificates
        env:
          WINDOWS_CERT_BASE64: ${{ secrets.WINDOWS_CERTIFICATE }}
          PGP_KEY_BASE64: ${{ secrets.PGP_PRIVATE_KEY }}
        run: |
          echo "$WINDOWS_CERT_BASE64" | base64 -d > certificate.pfx
          echo "$PGP_KEY_BASE64" | base64 -d > signing-key.asc

      # Build and sign Windows
      - name: Build and Sign Windows
        env:
          WAILS_WINDOWS_CERT_PASSWORD: ${{ secrets.WINDOWS_CERTIFICATE_PASSWORD }}
        run: |
          wails3 task windows:sign SIGN_CERTIFICATE=certificate.pfx

      # Build and sign Linux packages
      - name: Build and Sign Linux Packages
        env:
          WAILS_PGP_PASSWORD: ${{ secrets.PGP_PASSWORD }}
        run: |
          wails3 task linux:sign:packages PGP_KEY=signing-key.asc

      # Cleanup secrets
      - name: Cleanup
        if: always()
        run: rm -f certificate.pfx signing-key.asc

      - name: Upload Artifacts
        uses: actions/upload-artifact@v4
        with:
          name: signed-binaries
          path: |
            bin/*.exe
            bin/*.deb
            bin/*.rpm
```

@note{type="note"}
Usar um executor Linux para assinatura multiplataforma simplifica o CI/CD ao eliminar a necessidade de executores Windows separados. O macOS ainda exige um executor macOS para a assinatura devido aos requisitos das ferramentas nativas da Apple.

@end

## Referência da CLI

### wails3 signing

Assistente interativo para configurar a assinatura do seu projeto.

```bash
wails3 signing [flags]

Flags:
  --platform    Platform to configure (darwin, windows, linux). Repeatable.
                If omitted, auto-detects which platforms to configure from the build directory.
```

O assistente orienta você nestas etapas:

- **macOS**: seleção de um certificado Developer ID e configuração das credenciais de notarização (chama `xcrun notarytool store-credentials`).
- **Windows**: escolha entre um arquivo de certificado ou uma impressão digital e definição da senha e do servidor de carimbo de data/hora.
- **Linux**: uso de uma chave PGP existente ou geração de uma nova (chama `gpg`) e configuração da função de assinatura.

### wails3 entitlements

Assistente interativo para configurar os direitos do macOS.

```bash
wails3 entitlements [flags]

Flags:
  --output    Output path for entitlements.plist (default: build/darwin/entitlements.plist)
```

**Predefinições:**

- **Desenvolvimento**: cria `entitlements.dev.plist` com JIT, depuração e rede
- **Produção**: cria `entitlements.plist` com o mínimo de direitos
- **Ambos**: cria os dois arquivos (recomendado)
- **App Store**: cria direitos com sandbox para a Mac App Store
- **Personalizado**: escolha direitos individuais e o arquivo de destino

### wails3 sign

Assina binários e pacotes para a plataforma atual ou especificada. Este é um comando wrapper que chama a tarefa de assinatura específica da plataforma apropriada.

```bash
wails3 sign
wails3 sign GOOS=darwin
wails3 sign GOOS=windows
wails3 sign GOOS=linux
```

Isso executa a tarefa `<platform>:sign` correspondente, que usa a configuração de assinatura do seu Taskfile.

### wails3 tool sign

Comando de baixo nível para assinar diretamente um arquivo específico. Usado internamente pelos Taskfiles.

```bash
wails3 tool sign [flags]
```

**Opções comuns:**\

| Opção | Descrição |
| --- | --- |
| `--input` | Caminho para o arquivo a ser assinado |
| `--output` | Caminho de saída (opcional; por padrão, sobrescreve no local) |
| `--verbose` | Ativa a saída detalhada |

**Opções do Windows/macOS:**\

| Opção | Descrição |
| --- | --- |
| `--certificate` | Caminho para o certificado PKCS#12 (.pfx/.p12) |
| `--password` | Senha do certificado |
| `--timestamp` | URL do servidor de carimbo de data/hora |

**Opções específicas do macOS:**\

| Opção | Descrição |
| --- | --- |
| `--identity` | Identidade de assinatura (use '-' para assinatura ad hoc) |
| `--entitlements` | Caminho para o plist de direitos |
| `--hardened-runtime` | Ativa o runtime reforçado (padrão: true) |
| `--notarize` | Envia para notarização |
| `--keychain-profile` | Perfil do chaveiro para notarização |

**Opções específicas do Windows:**\

| Opção | Descrição |
| --- | --- |
| `--thumbprint` | Impressão digital do certificado no repositório do Windows |

**Opções específicas do Linux:**\

| Opção | Descrição |
| --- | --- |
| `--pgp-key` | Caminho para a chave privada PGP |
| `--pgp-password` | Senha da chave PGP |
| `--role` | Função de assinatura DEB (origin/maint/archive/builder) |

### Como inspecionar o estado da assinatura (ferramentas nativas)

`wails3 signing` configura as credenciais de assinatura. Para inspecionar o estado da assinatura, use diretamente as ferramentas nativas:

| Tarefa | Comando |
| --- | --- |
| Listar identidades de assinatura de código do macOS | `security find-identity -v -p codesigning` |
| Armazenar credenciais de notarização | `xcrun notarytool store-credentials "<profile>" --apple-id … --team-id … --password …` |
| Inspecionar um arquivo de chave PGP | `gpg --show-keys <key.asc>` |
| Gerar um par de chaves PGP | `gpg --full-generate-key` |
| Exportar uma chave pública | `gpg --armor --export <email>` |

## Solução de problemas

### Problemas no macOS

**"Nenhum certificado de ID de desenvolvedor encontrado"**

- Verifique se o certificado está instalado nas Chaves do macOS
- Use `security find-identity -v -p codesigning` para verificar se ele não expirou
- Verifique se você tem um certificado "Developer ID Application" (e não apenas "Apple Development")

**"Falha na notarização"**

- Verifique o log de notarização: `xcrun notarytool log <submission-id> --keychain-profile <profile>`
- Verifique se o hardened runtime está habilitado
- Verifique se o aplicativo não inclui binários não assinados

**"Falha no codesign"**

- Verifique se as Chaves estão desbloqueadas: `security unlock-keychain`
- Verifique as permissões de arquivo do pacote do aplicativo

### Problemas no Windows

**"Certificado não encontrado"**

- Verifique se o caminho do certificado está correto
- Verifique a senha do certificado
- Verifique se o certificado é válido (não expirou nem foi revogado)

**"Erro no servidor de carimbo de data/hora"**

- Tente usar outro servidor de carimbo de data/hora:
  - `http://timestamp.digicert.com`
  - `http://timestamp.sectigo.com`
  - `http://timestamp.comodoca.com`


### Problemas no Linux

**"Chave PGP inválida"**

- Verifique se o arquivo de chave está no formato ASCII-armored
- Use `gpg --show-keys <key.asc>` para verificar se a chave não expirou
- Verifique se a senha está correta

**"Falha na verificação da assinatura"**

- Verifique se a chave pública foi importada corretamente
- Verifique se o pacote não foi modificado após a assinatura

## Recursos adicionais

### Documentação oficial

- [Guia de assinatura de código da Apple](https://developer.apple.com/support/code-signing/)
- [Documentação de notarização da Apple](https://developer.apple.com/documentation/security/notarizing_macos_software_before_distribution)
- [Assinatura de código da Microsoft](https://docs.microsoft.com/en-us/windows-hardware/drivers/dashboard/get-a-code-signing-certificate)
- [Assinatura de pacotes Debian](https://wiki.debian.org/SecureApt)
- [Assinatura de pacotes RPM](https://rpm-software-management.github.io/rpm/manual/signatures.html)

---
title: "Empacotamento para Windows"
description: "Empacote seu aplicativo Wails para distribuição no Windows"
slug: "guides/build/windows"
sourcePath: "guides/build/windows.md"
---

## Instalador NSIS

O formato de empacotamento padrão cria um instalador NSIS:

```bash
wails3 package GOOS=windows
```

Isso executa `wails3 task windows:package`, que:

1. Compila o aplicativo
2. Gera o bootstrapper do WebView2
3. Cria um instalador NSIS

Saída: `build/windows/nsis/<AppName>-installer.exe`

### Pacote MSIX

Para distribuição pela Microsoft Store ou implantação moderna no Windows:

```bash
wails3 package GOOS=windows FORMAT=msix
```

Saída: `bin/<AppName>-<arch>.msix`

@note{type="note"}
O MSIX requer `makeappx.exe` (Windows SDK) ou as ferramentas MSIX autônomas. O Taskfile do Windows disponibiliza a tarefa de instalação como `wails3 task install:msix:tools`.

@end

## Personalização do instalador

A configuração do NSIS está em `build/windows/nsis/project.nsi`. Edite esse arquivo para personalizar:

- Interface do instalador e identidade visual
- Diretório de instalação
- Atalhos do menu Iniciar e da área de trabalho
- Associações de arquivos
- Contrato de licença

Os metadados do aplicativo vêm de `build/windows/info.json`:

```json
{
  "fixed": {
    "file_version": "1.0.0"
  },
  "info": {
    "0000": {
      "ProductVersion": "1.0.0",
      "CompanyName": "My Company",
      "FileDescription": "My Application",
      "ProductName": "MyApp"
    }
  }
}
```

## Assinatura de código

Assine o executável e o instalador para evitar avisos do SmartScreen:

```bash
# Using the wrapper (auto-detects platform)
wails3 sign GOOS=windows

# Or using tasks directly
wails3 task windows:sign
wails3 task windows:sign:installer
```

Configure a assinatura em `build/windows/Taskfile.yml`:

```yaml
vars:
  SIGN_CERTIFICATE: "path/to/certificate.pfx"
  # Or use thumbprint for certificates in Windows store
  SIGN_THUMBPRINT: "certificate-thumbprint"
  TIMESTAMP_SERVER: "http://timestamp.digicert.com"
```

Armazene a senha do certificado com segurança:

```bash
wails3 setup signing
```

Consulte [Assinatura de aplicativos](/guides/build/signing/) para obter detalhes.

## Compilação para ARM

```bash
wails3 build GOOS=windows GOARCH=arm64
wails3 package GOOS=windows GOARCH=arm64
```

## Solução de problemas

### makensis não encontrado

Instale o NSIS:

```bash
# Windows
winget install NSIS.NSIS

# Or download from https://nsis.sourceforge.io/
```

### Aviso do SmartScreen

Seu executável não está assinado. Consulte [Assinatura de código](#assinatura-de-cdigo) acima.

### WebView2 ausente

O instalador inclui um bootstrapper do WebView2 que baixa o runtime quando necessário. Se você precisar de uma instalação offline, baixe o Evergreen Standalone Installer da Microsoft.

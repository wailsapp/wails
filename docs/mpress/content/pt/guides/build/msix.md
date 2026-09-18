---
title: "Empacotamento MSIX"
description: "Empacotamento do seu aplicativo Wails v3 como um pacote MSIX"
slug: "guides/build/msix"
sourcePath: "guides/build/msix.md"
---

MSIX é o formato moderno de empacotamento de aplicativos para Windows. O Wails pode gerar um pacote MSIX como parte da compilação para Windows.

As instruções para empacotamento MSIX estão documentadas no guia [Empacotamento para Windows](/guides/build/windows/#msix-package).

## Empacotamento direto pela CLI

Execute a ferramenta MSIX no Windows, inclusive na CI. O backend padrão usa `MakeAppx.exe`; a assinatura também exige `signtool.exe`. Ambos fazem parte do SDK do Windows. O assistente de instalação abre a Microsoft Store e, se necessário, a página de download do SDK; conclua a instalação antes de criar o pacote.

Defina a identidade em `build/config.yml`, depois compile e empacote o executável:

```yaml
info:
  companyName: "Example Corp"
  productName: "MyApp"
  productIdentifier: "com.example.myapp"
  description: "MyApp"
  version: "1.0.0"
```

```powershell
wails3 tool msix-install-tools
wails3 build GOOS=windows
wails3 tool msix --executable bin/myapp.exe --name myapp.exe
```

O comando direto grava `MyApp.msix` no diretório atual. Já a [tarefa de empacotamento do Windows](/guides/build/windows/#msix-package) fornece seu próprio caminho de saída.

## Opções da CLI

| Opção | Significado |
| --- | --- |
| `--config` | Arquivo de configuração; padrão: `build/config.yml`. |
| `--executable`, `--name` | Executável existente e seu nome de arquivo dentro do pacote; ambos são obrigatórios. |
| `--out` | Arquivo de saída; padrão: `<ProductName>.msix`. |
| `--arch` | Arquitetura do pacote: `x64` (padrão), `x86`, `arm`, `arm64`, `x86a64` ou `neutral`. Os aliases Go `amd64` e `386` são aceitos. Use a arquitetura do executável. |
| `--publisher` | Identidade do publicador; padrão: `CN=<companyName>`. |
| `--cert`, `--cert-password` | Caminho do certificado PFX e senha para assinatura. |
| `--use-makeappx` | Usar o empacotador padrão do SDK do Windows. |
| `--use-msix-tool` | Selecionar explicitamente `MsixPackagingTool.exe`, que deve estar no `PATH`. |

## Assinatura e CI

Para distribuição fora da Store, assine com um certificado confiável na máquina de destino. O campo Subject do certificado deve corresponder exatamente a `--publisher`. O backend MakeAppx invoca o SignTool com SHA256 quando `--cert` é fornecido. Consulte o [guia de assinatura da Microsoft](https://learn.microsoft.com/en-us/windows/msix/package/sign-msix-package-guide).

Esta etapa de workflow do Windows pressupõe que Wails e o SDK estão instalados e que uma etapa anterior disponibilizou com segurança o arquivo PFX em `CERT_PATH`. Um segredo com um caminho não transfere o certificado por si só:

```yaml
- name: MSIX
  if: runner.os == 'Windows'
  shell: pwsh
  run: |
    wails3 build GOOS=windows
    wails3 tool msix --executable bin/myapp.exe --name myapp.exe --publisher "$env:MSIX_PUBLISHER" --cert "$env:CERT_PATH" --cert-password "$env:CERT_PASSWORD"
  env:
    MSIX_PUBLISHER: ${{ vars.MSIX_PUBLISHER }}
    CERT_PATH: ${{ secrets.WINDOWS_CERT_PATH }}
    CERT_PASSWORD: ${{ secrets.WINDOWS_CERT_PASSWORD }}
```

## Associações de arquivos e recursos

Adicione extensões sem o ponto inicial em `build/config.yml`; o manifesto gerado acrescenta o ponto:

```yaml
fileAssociations:
  - ext: myext
    name: MyApp Document
    description: MyApp Document
    iconName: fileicon
```

Trate a abertura de arquivos em tempo de execução conforme descrito em [Associações de arquivos](/guides/file-associations/). Atualmente, o backend MakeAppx copia apenas o executável e gera imagens transparentes provisórias. Ele não importa arquivos de `Assets/` do projeto nem converte ícones de `iconName`. Use um fluxo de empacotamento personalizado para recursos visuais próprios ou DLLs adicionais.

| Recurso gerado | Tamanho (pixels) |
| --- | --- |
| `Square150x150Logo.png` | 150×150 |
| `Square44x44Logo.png` | 44×44 |
| `Wide310x150Logo.png` | 310×150 |
| `StoreLogo.png` | 50×50 |
| `SplashScreen.png` | 620×300 |
| `FileIcon.png` | 44×44 |

`FileIcon.png` é gerado somente quando há associações de arquivos configuradas. Esses arquivos ficam no diretório `Assets/` do pacote.

## Envio para a Store e solução de problemas

Reserve seu aplicativo no [portal Partner Center](https://partner.microsoft.com/dashboard) e use a identidade do pacote e o publicador fornecidos por ele ao preparar o envio. A Store assina pacotes MSIX durante o envio; você não precisa comprar um certificado de assinatura para essa forma de distribuição. Consulte os [requisitos de pacotes da Microsoft](https://learn.microsoft.com/en-us/windows/apps/publish/publish-your-app/msix/app-package-requirements).

Se `MakeAppx.exe` ou `signtool.exe` não for encontrado, instale ou repare o SDK do Windows. O Wails pesquisa no `PATH` e nos locais padrão do SDK. Para falhas de assinatura, verifique o campo Subject do certificado, sua validade e a confiança na máquina de destino; consulte a [solução de problemas de MSIX](https://learn.microsoft.com/en-us/windows/msix/msix-troubleshooting-guide).

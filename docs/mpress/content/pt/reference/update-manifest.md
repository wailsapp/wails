---
title: "Protocolo de manifesto de atualização"
description: "O protocolo JSON aberto que os aplicativos Wails usam para descobrir e verificar atualizações do próprio aplicativo, podendo ser servido por qualquer hospedagem de arquivos estáticos ou servidor de atualização dinâmico."
slug: "reference/update-manifest"
sourcePath: "reference/update-manifest.md"
---

O protocolo de manifesto de atualização do Wails é um contrato JSON pequeno e aberto entre um aplicativo Wails e uma fonte de atualização. Qualquer serviço capaz de disponibilizar um arquivo JSON por HTTPS pode fornecer atualizações do Wails: um bucket do S3, o GitHub Pages, uma CDN ou um servidor de atualização dinâmico que condicione as versões a uma licença.

O lado do cliente é fornecido com o framework como o provedor `endpoint` (`github.com/wailsapp/wails/v3/pkg/updater/providers/endpoint`). Esta página é a referência do formato de comunicação para quem implementa o lado do servidor.

## Objetivos de design

1. **Compatível com hospedagem estática.** Um único arquivo de manifesto por canal, contendo o artefato de cada plataforma, constitui uma implementação completa. Nenhum código de servidor é necessário.
2. **Compatível com servidores dinâmicos.** O cliente envia `platform`, `arch`, `version` e `channel` em cada verificação, permitindo que um servidor responda com exatamente um artefato, aplique regras de licenciamento ou retorne `204 No Content` quando o solicitante já estiver atualizado.
3. **Verificação em primeiro lugar.** O manifesto contém somas de verificação e assinaturas para cada artefato, e o atualizador do Wails as verifica usando uma chave pública incorporada ao binário do aplicativo durante a compilação. A fonte de atualização nunca escolhe sua própria raiz de confiança.

## A solicitação

O cliente envia uma solicitação `GET` à URL de manifesto configurada, com `Accept: application/json` e todos os cabeçalhos configurados pelo aplicativo (por exemplo, `Authorization: License <key>`).

A URL pode conter espaços reservados, que o cliente substitui em cada verificação:

| Espaço reservado | Substituído por |
| --- | --- |
| `{{platform}}` | O sistema operacional em execução, como um valor `GOOS` do Go (`darwin`, `windows`, `linux`) |
| `{{arch}}` | A arquitetura em execução, como um valor `GOARCH` do Go (`amd64`, `arm64`, ...) |
| `{{version}}` | A versão instalada no momento |
| `{{channel}}` | O canal de lançamento configurado, quando definido |

Qualquer um dos quatro valores que não for usado por um espaço reservado será acrescentado como um parâmetro de consulta com o mesmo nome (`channel` somente quando configurado). Portanto, estas duas configurações são válidas e equivalentes:

```text
# Dynamic server: reads query parameters
https://updates.example.com/check
  -> GET /check?platform=darwin&arch=arm64&version=1.0.0&channel=stable

# Static host: one manifest per platform/arch/channel path
https://cdn.example.com/updates/{{platform}}/{{arch}}/{{channel}}.json
  -> GET /updates/darwin/arm64/stable.json?version=1.0.0
```

As hospedagens estáticas simplesmente ignoram os parâmetros de consulta recebidos.

## A resposta

| Status | Significado |
| --- | --- |
| `200 OK` | Um manifesto vem em seguida. O cliente decide se ele representa uma atualização. |
| `204 No Content` | O servidor comparou as versões e o solicitante está atualizado. |
| `404 Not Found` | Nada foi publicado (tratado da mesma forma que estar atualizado). |
| Qualquer outro | Um erro. O atualizador passa para o próximo provedor configurado. |

Um corpo `200` é um documento de manifesto:

```json
{
  "schemaVersion": 1,
  "version": "2.1.0",
  "channel": "stable",
  "name": "Summer Release",
  "notes": "## What's new\n\n- Faster startup\n- New themes",
  "publishedAt": "2026-07-03T10:00:00Z",
  "artifacts": [
    {
      "url": "MyApp-2.1.0-darwin-arm64.zip",
      "platform": "darwin",
      "arch": "arm64",
      "filetype": "zip",
      "size": 8388608,
      "digestAlgo": "sha512",
      "digest": "base64-encoded digest bytes",
      "signatureAlgo": "ed25519ph",
      "signature": "base64-encoded signature bytes"
    },
    {
      "url": "MyApp-2.1.0-windows-amd64.zip",
      "platform": "windows",
      "arch": "amd64",
      "filetype": "zip",
      "size": 9437184,
      "digestAlgo": "sha512",
      "digest": "...",
      "signatureAlgo": "ed25519ph",
      "signature": "..."
    }
  ]
}
```

### Campos de nível superior

| Campo | Tipo | Obrigatório | Observações |
| --- | --- | --- | --- |
| `schemaVersion` | int | não | Versão do protocolo. Quando omitida, significa `1`. Os clientes rejeitam valores mais recentes do que aqueles que conseguem interpretar. |
| `version` | string | **sim** | SemVer 2.0.0, com ou sem `v` no início. |
| `channel` | string | não | Informativo. Um cliente configurado para outro canal trata o manifesto como se não houvesse atualização. |
| `name` | string | não | Título da versão legível por pessoas, exibido na janela de atualização. |
| `notes` | string | não | Notas da versão em Markdown, renderizadas na janela de atualização. |
| `publishedAt` | string | não | Carimbo de data e hora no formato RFC 3339. |
| `artifacts` | array | **sim** | Uma entrada por artefato disponível para download. A ordem expressa a preferência do publicador. |
| `metadata` | object | não | Dados de chave/valor em formato livre, repassados ao aplicativo. |

Os clientes ignoram campos desconhecidos, portanto os servidores podem adicionar seus próprios campos sem causar incompatibilidades. As adições específicas do servidor devem ficar em `metadata`.

### Campos do artefato

| Campo | Tipo | Obrigatório | Observações |
| --- | --- | --- | --- |
| `url` | string | **sim** | Absoluta ou relativa à URL do manifesto. Somente `http(s)`. |
| `platform` | string | não | Valor de `GOOS` do Go. Apelidos comuns (`macos`, `win`, ...) são aceitos. Um valor vazio corresponde a todas as plataformas. |
| `arch` | string | não | Valor de `GOARCH` do Go. Apelidos comuns (`x86_64`, `aarch64`, ...) são aceitos. Um valor vazio corresponde a todas as arquiteturas. |
| `filename` | string | não | O padrão é o último segmento do caminho de `url`. |
| `filetype` | string | não | O padrão é a extensão do nome do arquivo. |
| `size` | int | não | Bytes, usados para indicar o progresso do download. |
| `digestAlgo` / `digest` | string / base64 | não | `sha256` ou `sha512`. |
| `signatureAlgo` / `signature` | string / base64 | não | `ed25519`, `ed25519ph` ou `ecdsa-p256`. `signatureAlgo` é obrigatório sempre que `signature` estiver presente. Consulte o [guia do atualizador](/guides/updater/#cryptographic-verification) para saber o que cada algoritmo assina. |

O cliente seleciona o **primeiro** artefato cujos valores de `platform` e `arch` correspondam ao sistema em execução. Valores em Base64 são aceitos com ou sem preenchimento.

### Comparação de versões

A decisão sobre o manifesto representar ou não uma atualização sempre é tomada no cliente, segundo a precedência do SemVer 2.0.0: o valor de `version` do manifesto deve ser estritamente mais recente que a versão instalada. Isso torna a hospedagem estática trivialmente correta (o manifesto sempre descreve a versão mais recente, e os clientes atualizados simplesmente não fazem nada), enquanto `204` continua disponível para que servidores dinâmicos economizem largura de banda.

## Verificação e confiança

As somas de verificação e assinaturas são incluídas no manifesto, mas a raiz de confiança não: as assinaturas são verificadas com a chave pública que a aplicação fixou por meio de `updater.Config.PublicKey` durante a compilação. Uma fonte de atualização comprometida ou substituída não pode fornecer sua própria chave. Um artefato que contenha uma assinatura quando a aplicação não tiver uma chave fixada falha de forma segura, assim como uma assinatura sem um `signatureAlgo` declarado ou que não possa ser decodificada: os clientes nunca recorrem silenciosamente à verificação apenas por resumo criptográfico.

Artefatos que contêm apenas o resumo criptográfico são instalados após sua verificação, o que protege contra corrupção, mas depende de TLS e da própria integridade do host para resistir a adulterações. Distribua assinaturas para tudo que seja sensível à segurança.

Assinar um artefato com o esquema `ed25519ph` do framework requer apenas algumas linhas de Go:

```go
digest := sha512.Sum512(artifactBytes)
sig, _ := privateKey.Sign(nil, digest[:], &ed25519.Options{Hash: crypto.SHA512})
manifest.Artifacts[i].DigestAlgo = "sha512"
manifest.Artifacts[i].Digest = base64.StdEncoding.EncodeToString(digest[:])
manifest.Artifacts[i].SignatureAlgo = "ed25519ph"
manifest.Artifacts[i].Signature = base64.StdEncoding.EncodeToString(sig)
```

Na prática, raramente é necessário escrever esse código: a CLI faz isso por você.

## Publicação com a CLI wails3

O grupo de comandos `wails3 updater` abrange todo o pipeline de publicação. Uma versão é publicada com três comandos:

```bash
# Once per application: create the signing keypair.
wails3 updater genkey
# updater.key      keep secret (CI secret store), signs every release
# updater.key.pub  embed in the app and pass as updater.Config.PublicKey

# Per release: digest, sign and describe every artifact in one manifest.
wails3 updater manifest -version 2.1.0 -channel stable \
    -key updater.key -notes-file notes.md \
    -url-prefix "https://cdn.example.com/myapp/2.1.0" \
    bin/updates/

# Before uploading: re-verify the files exactly as a shipped app would.
wails3 updater verify -manifest manifest.json -publickey updater.key.pub
```

`manifest` aceita arquivos ou diretórios (materiais de chave, `.json`, arquivos auxiliares de soma de verificação e notas são ignorados automaticamente), processa cada artefato em fluxo com SHA-512, assina o resumo com Ed25519ph quando `-key` é fornecido e infere `platform` e `arch` a partir de nomes de arquivo convencionais, como `MyApp-2.1.0-darwin-arm64.zip` (apelidos comuns como `macOS`, `win64`, `x86_64` e `aarch64` são reconhecidos; um aviso é exibido para tudo que não puder ser inferido, que então corresponderá a todas as plataformas). Omita `-url-prefix` para gerar URLs relativas e envie o manifesto junto aos artefatos.

`verify` termina com um código diferente de zero diante de qualquer divergência, o que o torna uma etapa natural de validação da CI entre a compilação e a publicação. Para servidores que montam os próprios manifestos, `wails3 updater sign -key updater.key <files...>` imprime os campos `digest`/`signature` de cada arquivo como JSON pronto para ser incorporado ao seu próprio documento.

## Autenticação

A autenticação é responsabilidade do servidor; o protocolo apenas transporta cabeçalhos. O cliente reenvia os cabeçalhos configurados em cada solicitação do manifesto. Nos downloads de artefatos, o cabeçalho `Authorization` só é enviado quando a URL do artefato está no mesmo host que o manifesto e não rebaixa de `https` para `http`. Ele é removido em qualquer redirecionamento entre origens ou que faça esse rebaixamento, de modo que as credenciais nunca vazem para uma CDN ou um armazenamento de objetos e nunca sejam transmitidas como texto simples.

Exemplo condicionado a uma licença, que se integra naturalmente a serviços de licenciamento hospedados:

```go
ep, _ := endpoint.New(endpoint.Config{
    URL:     "https://updates.example.com/check",
    Headers: map[string]string{"Authorization": "License " + licenseKey},
})
```

## Configuração do cliente

Consulte o [guia do atualizador](/guides/updater/#providers) para ver a referência completa de `endpoint.Config` e como o provedor se encaixa em cadeias de fallback ao lado dos provedores GitHub, keygen.sh e AppCast.

---
title: "Empacotamento para macOS"
description: "Empacote seu aplicativo Wails para distribuição no macOS"
slug: "guides/build/macos"
sourcePath: "guides/build/macos.md"
---

## APIs privadas do macOS

Por padrão, o Wails v3 usa APIs públicas do macOS. Para habilitar recursos que exigem APIs não documentadas da Apple, compile seu aplicativo com a única tag de compilação do Go `private_mac_apis`:

```bash
wails3 build -tags private_mac_apis
EXTRA_TAGS=private_mac_apis wails3 dev
wails3 package GOOS=darwin EXTRA_TAGS=private_mac_apis
```

Para uma compilação direta com o Go, use `go build -tags private_mac_apis .` (ou `-tags production,private_mac_apis` para produção). Aplicativos existentes que dependem de comportamentos privados devem adicionar essa tag para mantê-los. A tag se aplica apenas a compilações de desktop para macOS.

Consulte [APIs privadas do macOS](/guides/build/private-macos-apis/) para ver a lista completa dos recursos e valores de opções afetados, os comportamentos alternativos exatos das compilações públicas, os mapeamentos de estilo do Liquid Glass e as combinações de compilação do inspetor. Sem a tag, as operações exclusivas das APIs privadas não realizam nenhuma ação; a API pública do Go permanece inalterada.

## Pacote do aplicativo

Empacote seu aplicativo como um pacote `.app` padrão do macOS:

```bash
wails3 package GOOS=darwin
```

Isso cria `bin/<AppName>.app`, que contém:

- O binário compilado em `Contents/MacOS/`
- O ícone do aplicativo em `Contents/Resources/` (proveniente de `icons.icns` ou, quando presente, de um catálogo de recursos `Assets.car`)
- `Info.plist` com os metadados do aplicativo

## Recursos do pacote

`Contents/Resources/` é o local padrão para arquivos somente leitura distribuídos com um aplicativo macOS. Use-o para modelos maiores, dados iniciais, mídia, pacotes de idiomas ou outros conteúdos que devam ser abertos sob demanda, em vez de compilados no executável Go com `embed`.

O Wails já coloca o ícone do aplicativo nesse diretório. Para adicionar seus próprios arquivos, coloque-os em um diretório de origem, como `build/resources/`, e adicione uma etapa de cópia à tarefa `create:app:bundle` em `build/darwin/Taskfile.yml`:

```yaml
tasks:
  create:app:bundle:
    cmds:
      # Existing bundle creation commands...
      - |
          if [ -d build/resources ]; then
            cp -R build/resources/. "{{.BIN_DIR}}/{{.APP_NAME}}.app/Contents/Resources/"
          fi
```

Se você usa a tarefa `darwin:run` do Taskfile, adicione o comando equivalente à tarefa `run` correspondente, tendo `{{.BIN_DIR}}/{{.APP_NAME}}.dev.app/Contents/Resources/` como destino.

### Leitura de recursos no Go

Importe o pacote da plataforma macOS:

```go
import (
	"io/fs"

	"github.com/wailsapp/wails/v3/pkg/mac"
)
```

Para arquivos pequenos, use `LoadResource`:

```go
func loadSplash() ([]byte, error) {
	return mac.LoadResource("images/splash.png")
}
```

Para arquivos maiores, use `ResourceFS`. Essa função retorna um `io/fs.FS` cuja raiz é `Contents/Resources`, permitindo que os chamadores abram e transmitam um recurso sem antes carregá-lo por completo em um slice de bytes do Go:

```go
func openCatalogue() (fs.File, error) {
	resources, err := mac.ResourceFS()
	if err != nil {
		return nil, err
	}

	return resources.Open("catalogue/defaults.json")
}
```

Os nomes dos recursos são caminhos separados por barras, relativos a `Contents/Resources`. `ResourceFS` e `LoadResource` retornam `mac.ErrNotInAppBundle`, a menos que o executável seja executado a partir de `.app/Contents/MacOS`.

Trate os recursos do pacote como imutáveis. Alterar arquivos dentro de um aplicativo assinado invalida sua assinatura de código; em vez disso, armazene dados baixados, gerados ou editáveis pelo usuário no diretório Application Support do usuário.

### Binário universal

Compile para Macs com Apple Silicon e Intel:

```bash
wails3 task darwin:package:universal
```

Isso cria um único `.app` que é executado nativamente nas duas arquiteturas. Binários universais podem ser compilados em qualquer plataforma — no Linux e no Windows, `wails3 tool lipo` é usado automaticamente.

## Personalização do pacote

Edite `build/darwin/Info.plist` para personalizar:

- Identificador do pacote (`CFBundleIdentifier`)
- Nome e versão do aplicativo
- Versão mínima do macOS
- Associações de arquivos
- Esquemas de URL

O ícone do aplicativo é gerado a partir dos recursos no diretório `build/`. Use a tarefa `generate:icons`:

```bash
wails3 task common:generate:icons
```

Isso usa `build/appicon.png` para gerar `darwin/icons.icns` e `windows/icon.ico`. No macOS, você também pode fornecer `build/appicon.icon` (formato do Icon Composer): a tarefa passa `-iconcomposerinput appicon.icon -macassetdir darwin`, que gera `Assets.car` e `darwin/icons.icns` a partir do arquivo `.icon` (essa etapa é ignorada em plataformas que não sejam macOS). Quando `Assets.car` estiver presente, execute a tarefa `update:build-assets` para que `Info.plist` e `CFBundleIconName` sejam atualizados de acordo:

```bash
wails3 task common:update:build-assets
```

Para executar manualmente o comando de ícone a partir do diretório `build/`:

```bash
cd build
wails3 generate icons -input appicon.png -macfilename darwin/icons.icns -windowsfilename windows/icon.ico -iconcomposerinput appicon.icon -macassetdir darwin
```

## Assinatura de código

Assine seu aplicativo para distribuição:

```bash
# Using the wrapper (auto-detects platform)
wails3 sign GOOS=darwin

# Or using the task directly
wails3 task darwin:sign
```

Configure a assinatura em `build/darwin/Taskfile.yml`:

```yaml
vars:
  SIGN_IDENTITY: "Developer ID Application: Your Company (TEAMID)"
  ENTITLEMENTS: "build/darwin/entitlements.plist"
```

### Notarização

Para aplicativos distribuídos fora da Mac App Store, a Apple exige notarização:

```bash
wails3 task darwin:sign:notarize
```

Primeiro, armazene suas credenciais. Execute o assistente interativo (`wails3 setup signing`) ou chame `notarytool` diretamente:

```bash
xcrun notarytool store-credentials "my-notarize-profile" \
  --apple-id "you@email.com" \
  --team-id "TEAMID" \
  --password "app-specific-password"
```

Configure em `build/darwin/Taskfile.yml`:

```yaml
vars:
  SIGN_IDENTITY: "Developer ID Application: Your Company (TEAMID)"
  KEYCHAIN_PROFILE: "my-notarize-profile"
```

Consulte [Assinatura de aplicativos](/guides/build/signing/) para obter detalhes.

## Instalador DMG

O modelo fornecido com o Wails 3 disponibiliza `wails3 task darwin:package:dmg`. Primeiro, ele cria o `.app` e depois gera um DMG estilizado com a biblioteca DMG. Por padrão, o DMG usa um plano de fundo em degradê com a identidade visual do Wails, contendo o símbolo do dragão vermelho e o logotipo textual WAILS.

```bash
wails3 task darwin:package:dmg
```

A tarefa de nível mais baixo `darwin:create:dmg` cria um DMG a partir de um pacote `.app` existente e pode ser configurada diretamente no Taskfile:

```yaml
vars:
  # These are the template defaults; override them when needed.
  DMG_BACKGROUND: build/darwin/dmg-background.png
  DMG_VOLUME_ICON: build/darwin/icons.icns
  DMG_FILE_ICON: build/darwin/dmg-file-icon.icns
  DMG_WINDOW_WIDTH: 540
  DMG_WINDOW_HEIGHT: 380
  DMG_FILES: "Install.command=build/darwin/Install.command,README.txt=README.md"
```

### Layout padrão

O DMG gerado contém:

- O pacote do aplicativo à esquerda
- Um link para `Applications` à direita
- Uma janela do Finder com 540 × 380 pixels
- Ícones com 96 pontos e rótulos abaixo de cada ícone
- O plano de fundo com a identidade visual do Wails proveniente de `build/darwin/dmg-background.png`

Os ícones do aplicativo e de `Applications` são posicionados em relação às dimensões configuradas da janela. Assim, ao alterar `DMG_WINDOW_WIDTH` ou `DMG_WINDOW_HEIGHT`, o espaçamento proporcional do layout padrão de dois ícones é mantido. Para obter o melhor resultado, use uma imagem de fundo com as mesmas dimensões em pixels da janela do Finder.

### Substituição dos recursos do DMG

Os arquivos gerados em `build/darwin/` são recursos comuns do projeto e podem ser substituídos:

- `DMG_BACKGROUND` controla a imagem exibida atrás do conteúdo da janela do Finder.
- `DMG_VOLUME_ICON` controla o ícone exibido para o volume montado.
- `DMG_FILE_ICON` controla o ícone exibido no Finder para o arquivo `.dmg` resultante.

O ícone do volume e o ícone do arquivo DMG são recursos separados. Substituir o ícone do aplicativo não substitui automaticamente nenhum deles.

### Como adicionar arquivos extras

Use `DMG_FILES` para incluir scripts do instalador, notas de versão, licenças ou outros recursos junto ao aplicativo. O valor é uma lista de pares `name=path` separados por vírgulas:

```yaml
vars:
  DMG_FILES: "Install.command=build/darwin/Install.command,README.txt=README.md"
```

O nome antes de `=` é o nome de arquivo exibido dentro do DMG. O caminho depois de `=` corresponde ao arquivo de origem no projeto. Os espaços em branco no início e no fim são ignorados.

Cada nome exibido deve ser único. Os arquivos extras não podem substituir entradas já criadas pelo empacotador, incluindo o pacote do aplicativo ou a entrada `Applications`. Em caso de conflito de nomes, o empacotamento falha com um erro, em vez de produzir um DMG corrompido.

@note{type="note"}
A criação de DMGs é compatível apenas com o macOS, pois utiliza as ferramentas de imagem de disco e do Finder do macOS. Os pacotes `.app` podem ser compilados de forma cruzada em outro sistema, mas o DMG final deve ser produzido em um Mac.

@end

## Solução de problemas

### "O app está danificado e não pode ser aberto"

O app não está assinado. Assine-o com um certificado Developer ID ou, alternativamente, os usuários podem ignorar o Gatekeeper:

```bash
xattr -cr /path/to/YourApp.app
```

### Falha na notarização

Problemas comuns:

- **Credenciais inválidas**: execute `xcrun notarytool store-credentials` novamente (ou `wails3 setup signing`)
- **Runtime reforçado obrigatório**: verifique se os entitlements incluem `com.apple.security.cs.allow-unsigned-executable-memory`, caso necessário
- **Carimbo de data e hora ausente**: o processo de assinatura deve incluir automaticamente um carimbo de data e hora

### O app compilado de forma cruzada não é executado

Os binários do macOS compilados de forma cruzada não são assinados. Transfira-os para um Mac e assine-os antes de testar:

```bash
codesign --force --deep --sign - YourApp.app
```

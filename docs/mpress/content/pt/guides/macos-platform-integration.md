---
title: "Integração com a plataforma macOS"
description: "Janelas de documentos, recursos adicionais do Dock e dos menus, painéis nativos, itens de status, feedback, área de transferência com conteúdo avançado, arraste para fora, permissões, energia e ciclo de vida no macOS"
slug: "guides/macos-platform-integration"
sourcePath: "guides/macos-platform-integration.md"
---

Plataformas relevantes: macOS

O Wails v3 oferece à sua aplicação os comportamentos que usuários do macOS esperam de um aplicativo nativo: janelas de documentos com ícones proxy e disposição em cascata, menus com símbolos e indicadores, menu do Dock com progresso, alertas e painéis nativos, itens de status removíveis, respostas táteis e fala, área de transferência com conteúdo avançado e arraste para fora, informações do sistema sobre permissões, energia e localidade, além de integração com o menu Serviços, Handoff, AppleScript e Quick Look. Tudo é controlado em Go pelo pacote `application`.

O mesmo código compila no Windows e no Linux. Os setters armazenam seus valores, as consultas retornam valores zero e as operações que exigem macOS retornam um erro documentado, como `ErrMacOnly`, `ErrDialogNotSupported` ou `ErrClipboardNotSupported`. As [notas sobre plataformas](#platform-notes) abaixo descrevem o comportamento de cada área fora do macOS.

Para elementos nativos da janela (barras de ferramentas, barras laterais, inspetores, acessórios e abas de janela), consulte o guia [Elementos nativos da janela no macOS](/guides/macos-native-chrome).

## Janelas de documentos

Uma janela de documento mostra na barra de título o arquivo que representa, sinaliza alterações não salvas com um ponto no botão de fechar e abre novas janelas em cascata. Esses recursos são métodos de `WebviewWindow` e opções de `MacWindow`.

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:           "Report.md",
    URL:             "/editor",
    InitialPosition: application.WindowCascade,
    Mac: application.MacWindow{
        FrameAutosaveName: "editor.main",
        TitleBar: application.MacTitleBar{
            WindowButtonsOffset: &application.Point{X: 12, Y: 8},
        },
    },
})

window.SetRepresentedFile("/Users/me/Documents/Report.md")
window.SetSubtitle("Documents")
window.SetDocumentEdited(true)
```

- `SetRepresentedFile` mostra o ícone proxy do arquivo na barra de título. O usuário pode arrastar o ícone para outro aplicativo ou clicar nele com a tecla Command pressionada para revelar o caminho. Passe `""` para removê-lo. `RepresentedFile` lê o valor de volta.
- `SetDocumentEdited` mostra o ponto de alterações não salvas no botão de fechar e esmaece o ícone proxy. `IsDocumentEdited` lê o estado de volta.
- `SetSubtitle` mostra uma segunda linha abaixo do título no macOS 11 e posteriores.
- `InitialPosition: application.WindowCascade` posiciona a janela abaixo e à direita da última janela em cascata, como na abertura de novos documentos. `CascadeFrom(other)` faz o mesmo para uma janela existente e atualiza o ponto de cascata para janelas posteriores.
- `Mac.FrameAutosaveName` restaura a posição e o tamanho salvos antes que a janela seja mostrada pela primeira vez e continua a salvá-los enquanto ela se move. Uma geometria restaurada tem precedência sobre `X`, `Y`, `Width`, `Height` e `InitialPosition`. `SetFrameAutosaveName` muda o nome em uma janela ativa.
- `MacTitleBar.WindowButtonsOffset` desloca os botões de fechar, minimizar e ampliar por um número de pontos. `SetWindowButtonsOffset` e `ResetWindowButtonsOffset` alteram esse deslocamento em tempo de execução.

Os três setters podem ser chamados antes da criação da janela nativa; os valores são aplicados quando ela é criada.

### Solicitações de atenção

`RequestAttention` faz o ícone do Dock saltar enquanto a aplicação está em segundo plano. Uma solicitação informativa faz o ícone saltar uma vez. Uma solicitação crítica mantém o ícone saltando até que o usuário ative a aplicação ou você cancele a solicitação.

```go
request := window.RequestAttention(true)

// Once the work that needed attention is done:
request.Cancel()
```

`Flash` continua sendo a forma multiplataforma de solicitar atenção uma vez.

### Impressão e exportação

`PrintWithOptions` imprime a WebView com configurações de página explícitas. O valor zero mostra o painel de impressão com as configurações de impressão compartilhadas. `Print` mantém seu comportamento histórico (orientação paisagem e margens de 30 pontos).

`ExportPDF` renderiza a página em um documento PDF, e `Snapshot` a captura como PNG. Ambos aguardam o WebKit; portanto, chame-os de uma goroutine, nunca da thread da aplicação. Chamá-los nessa thread retorna `ErrMacExportOnMainThread`.

```go
err := window.PrintWithOptions(application.PrintOptions{
    Orientation: application.PrintOrientationPortrait,
    Margins:     application.PrintMargins{Top: 36, Left: 36, Bottom: 36, Right: 36},
    Silent:      false,
})
if err != nil {
    log.Println("print:", err)
}

go func() {
    pdf, err := window.ExportPDF(application.PDFExportOptions{})
    if err != nil {
        log.Println("export:", err)
        return
    }
    if err := os.WriteFile("report.pdf", pdf, 0o644); err != nil {
        log.Println("write:", err)
    }

    png, err := window.Snapshot(application.SnapshotOptions{Width: 800})
    if err != nil {
        log.Println("snapshot:", err)
        return
    }
    if err := os.WriteFile("preview.png", png, 0o644); err != nil {
        log.Println("write:", err)
    }
}()
```

`PrintOptions` também aceita `PrinterName`, `PaperName` (um nome PostScript como `"iso-a4"`) e `Scale`. `PDFExportOptions` e `SnapshotOptions` aceitam um `Rect` opcional para limitar a captura e um `Timeout` cujo padrão é `DefaultMacExportTimeout` (30 segundos).

## Janelas modais anexadas

Uma janela modal anexada é uma segunda janela presa à parte superior da janela pai, como os painéis de Salvar. Qualquer `WebviewWindow` pode ser apresentada como janela modal anexada a outra com `PresentSheet` e encerrada com `EndSheet` e um código de resposta recebido pelos callbacks `OnSheetEnd` da janela. Crie a janela modal com `Hidden` definido para que ela não apareça brevemente na tela antes de ser anexada; o AppKit a oculta novamente quando ela termina, permitindo apresentar a mesma janela repetidas vezes.

```go
sheet := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:  "Rename",
    URL:    "/rename",
    Width:  420,
    Height: 180,
    Hidden: true,
})
sheet.OnSheetEnd(func(code int) {
    if code == application.MacSheetResponseOK {
        log.Println("renamed")
    }
})

// From a menu item or a bound method:
if err := window.PresentSheet(sheet); err != nil {
    log.Println(err)
}

// From the sheet's own page, through a bound method:
sheet.EndSheet(application.MacSheetResponseOK)
```

`PresentCriticalSheet` mostra a janela modal à frente de qualquer janela modal comum já anexada, em vez de colocá-la na fila atrás dela. `PresentNativeSheet` faz o mesmo com uma `NativeWindow`. `IsSheet`, `SheetParent`, `AttachedSheet`, `AttachedNativeSheet` e `HasAttachedSheet` descrevem o estado atual. Fechar uma janela modal anexada em vez de encerrá-la produz `MacSheetResponseStop`.

## Popovers

Um `MacPopover` é um `NSPopover`: um painel temporário ancorado a um retângulo em uma janela, a um item da barra de ferramentas ou a um item de status na barra de menus. Seu conteúdo é uma faixa de controles nativos `MacAccessory`, do mesmo tipo usado para acessórios da barra de título no guia [Elementos nativos da janela no macOS](/guides/macos-native-chrome). Adicione todos os controles antes da primeira exibição.

```go
strip := application.NewMacAccessory(application.MacAccessoryLayoutBottom)
strip.AddLabel("Sort by")
strip.AddSegmented([]string{"Date", "Title"}, 0)
strip.AddFlexibleSpace()
strip.AddButton("Apply").OnClick(func(*application.Context) {
    // apply the sort
})

popover := application.NewMacPopover(application.MacPopoverOptions{
    Width:    320,
    Behavior: application.MacPopoverBehaviorTransient,
    Content:  strip,
})
popover.OnClose(func() {
    log.Println("popover closed")
})

// Anchored to a rectangle in the page, in window content coordinates:
err := popover.ShowRelativeTo(application.Rect{X: 20, Y: 60, Width: 120, Height: 28}, window, application.MacRectEdgeMaxY)
if err != nil {
    log.Println(err)
}
```

`MacToolbarItem.ShowPopover` e `SystemTray.ShowPopover` ancoram o mesmo popover a um item da barra de ferramentas ou de status. `MacPopoverBehaviorTransient` fecha o popover com qualquer clique fora dele; `Semitransient` o fecha somente com cliques na janela que o apresenta; o comportamento padrão o mantém aberto até `Close`. `MacRectEdge` escolhe o lado em que o popover aparece. `SetContentSize` e `SetBehavior` ajustam um popover ativo, e `Destroy` libera o popover nativo e a faixa de conteúdo para uso em outro lugar.

## Restauração de estado

O macOS reabre as janelas de uma aplicação após uma falha, um encerramento forçado ou uma reinicialização, e também após um encerramento normal quando "Fechar janelas ao encerrar um aplicativo" está desativado nos Ajustes do Sistema. Atribua um `Mac.RestorationID` à janela, armazene o necessário para recriá-la com `SetRestorationData` e registre `app.Window.OnRestore` para reconstruí-la na próxima inicialização.

```go
app.Window.OnRestore(func(id string, state application.RestorationState) application.Window {
    if id != "editor" {
        return nil
    }
    path := state.Get("path")
    restored := app.Window.NewWithOptions(application.WebviewWindowOptions{
        Title: filepath.Base(path),
        URL:   "/editor?path=" + path,
        Mac:   application.MacWindow{RestorationID: id},
    })
    restored.SetRestorationData(state.Data)
    return restored
})

window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "Report.md",
    URL:   "/editor?path=/Users/me/Documents/Report.md",
    Mac:   application.MacWindow{RestorationID: "editor"},
})
window.SetRestorationData(map[string]string{"path": "/Users/me/Documents/Report.md"})
```

Somente janelas visíveis quando a aplicação termina são salvas. `RestorationState.Data` é um mapa de strings; limite seu conteúdo a identificadores, caminhos e posições. `InteractionState` retorna a lista de navegação anterior e posterior da WebView e as posições de rolagem como um bloco opaco (macOS 12+) que `RestoreInteractionState` aplica a uma janela recriada, normalmente armazenado com codificação base64 nos dados de restauração. `SetRestorationID` e `RestorationID` alteram e leem o identificador em uma janela ativa.

## Opções de apresentação

`MacPresentationOptions` espelha `NSApplication.presentationOptions`: uma máscara de bits que oculta o Dock ou a barra de menus e desabilita a troca de processos, o encerramento forçado, o logout ou o comando Ocultar enquanto a aplicação está ativa. Defina `Mac.PresentationOptions` nas opções da aplicação para aplicá-la na inicialização ou altere-a em tempo de execução com `SetPresentationOptions`. Combinações inválidas são rejeitadas antes de chegar ao AppKit, com um erro que encapsula `ErrMacPresentationOptionsInvalid`.

```go
kiosk := application.MacPresentationHideDock |
    application.MacPresentationHideMenuBar |
    application.MacPresentationDisableProcessSwitching |
    application.MacPresentationDisableForceQuit

if err := app.SetPresentationOptions(kiosk); err != nil {
    log.Println(err)
}
log.Println("presentation:", app.PresentationOptions())

// Restore the standard Dock and menu bar:
_ = app.SetPresentationOptions(application.MacPresentationDefault)
```

Ocultar a barra de menus (`HideMenuBar` ou `AutoHideMenuBar`) exige uma das opções do Dock, e `AutoHideToolbar` exige tanto `FullScreen` quanto `AutoHideMenuBar`. `Validate` informa a primeira regra violada por um valor, e `Has` verifica sinalizadores individuais.

## Menus e Dock

Os itens de menu passam a ter SF Symbols, indicadores, cabeçalhos de seção, paletas de cores, estado de seleção misto, alternativas e recuo. Todos são métodos de `MenuItem` e `Menu`; portanto, funcionam no menu da aplicação, em menus de contexto, menus da bandeja e no menu do Dock.

```go
menu := app.NewMenu()
menu.AddRole(application.AppMenu)

fileMenu := menu.AddSubmenu("File")
fileMenu.Add("Open...").SetSymbol("folder").SetAccelerator("CmdOrCtrl+o").
    OnClick(func(*application.Context) {
        // open the file, then note it in the recent list:
        app.Menu.AddRecentDocument("/Users/me/Documents/Report.md")
    })
fileMenu.AddRole(application.OpenRecent)

view := menu.AddSubmenu("View")
view.AddSectionHeader("Mailboxes")
inbox := view.Add("Inbox").SetSymbol("tray.full").SetBadge(3)
inbox.OnClick(func(*application.Context) {
    inbox.ClearBadge()
})
view.Add("Updates").SetBadgeText("New")

view.AddSeparator()
wrap := view.AddCheckbox("Wrap lines", false).SetMixed()
view.Add("Reset").SetIndentationLevel(1).OnClick(func(*application.Context) {
    wrap.SetMixed()
})

view.AddSeparator()
view.Add("Close Tab").SetAccelerator("CmdOrCtrl+w").OnClick(func(*application.Context) {})
view.Add("Close All Tabs").SetAccelerator("CmdOrCtrl+OptionOrAlt+w").SetAlternate(true).
    OnClick(func(*application.Context) {})

colours := []application.RGBA{
    application.NewRGB(255, 59, 48),
    application.NewRGB(255, 149, 0),
    application.NewRGB(52, 199, 89),
}
view.AddPalette([]string{"tag.fill"}, colours, 0, func(_ *application.Context, index int) {
    log.Println("tag colour", index)
}).SetLabel("Tag colour")

app.Menu.Set(menu)
```

- `SetSymbol` mostra um SF Symbol ao lado do título (macOS 11+). Ele substitui uma imagem definida com `SetBitmap`.
- `SetBadge` mostra uma contagem, e `SetBadgeText` mostra uma string curta após o título (macOS 14+). `ClearBadge` remove o indicador; `BadgeCount` e `BadgeText` leem seu valor.
- `AddSectionHeader` adiciona um cabeçalho não interativo (macOS 14+). Em versões anteriores, ele é um item desabilitado com o mesmo título.
- `AddPalette` adiciona uma linha de amostras de cores baseada no menu de paletas de `NSMenu` (macOS 14+). Passe um símbolo para cada amostra, um por cor, ou um slice vazio para círculos preenchidos. Sem um rótulo, a paleta aparece diretamente no menu pai; `SetLabel` a apresenta como um submenu com título. `PaletteSelected` retorna o índice selecionado.
- `SetMixed` coloca uma caixa de seleção no estado misto, representado por um traço. Um clique a ativa por completo, como faz o AppKit.
- `SetAlternate(true)` mostra o item no lugar daquele acima enquanto a tecla modificadora diferente estiver pressionada. Os dois itens devem compartilhar uma tecla e diferir nas teclas modificadoras.
- `SetIndentationLevel` recua o título em até 15 níveis.

### Abrir Recentes

`fileMenu.AddRole(application.OpenRecent)` adiciona o submenu padrão Abrir Recentes. No macOS, `NSDocumentController` o preenche sempre que ele é aberto e inclui um item Limpar Menu. Adicione arquivos com `app.Menu.AddRecentDocument`, liste-os com `RecentDocuments` e esvazie a lista com `ClearRecentDocuments`. A lista persiste entre reinicializações.

A escolha de um arquivo recente chega como o mesmo evento de um arquivo aberto pelo Finder:

```go
app.Event.OnApplicationEvent(events.Common.ApplicationOpenedWithFile, func(event *application.ApplicationEvent) {
    path := event.Context().Filename()
    log.Println("open", path)
})
```

### Menu do Dock

`app.Menu.SetDockMenu` instala um menu estático exibido com um clique direito no ícone do Dock. `OnDockMenu` cria um menu sob demanda sempre que ele está prestes a ser mostrado, a escolha adequada quando os itens refletem um estado variável. Uma função de criação tem precedência sobre o menu estático.

```go
dockMenu := application.NewMenu()
dockMenu.Add("New Document").SetSymbol("doc.badge.plus").OnClick(func(*application.Context) {
    // create a document
})
app.Menu.SetDockMenu(dockMenu)

app.Menu.OnDockMenu(func() *application.Menu {
    dynamic := application.NewMenu()
    dynamic.AddSectionHeader("Recent")
    for _, path := range app.Menu.RecentDocuments() {
        path := path
        dynamic.Add(filepath.Base(path)).OnClick(func(*application.Context) {
            app.Menu.OpenRecentDocument(path)
        })
    }
    return dynamic
})
```

### Progresso no Dock

O serviço do Dock desenha uma barra de progresso sobre o ícone do Dock, além do suporte existente a indicadores. Registre `dock.New()` como serviço e chame `SetProgress` com uma fração entre 0 e 1.

```go
import "github.com/wailsapp/wails/v3/pkg/services/dock"

dockService := dock.New()

app := application.New(application.Options{
    Name: "Exporter",
    Services: []application.Service{
        application.NewService(dockService),
    },
})

app.Event.On("export:progress", func(event *application.CustomEvent) {
    if fraction, ok := event.Data.(float64); ok {
        _ = dockService.SetProgress(fraction)
    }
})
app.Event.On("export:done", func(*application.CustomEvent) {
    _ = dockService.ClearProgress()
})
```

`GetProgress` retorna a fração atual ou `nil` quando nenhuma barra é mostrada.

## Caixas de diálogo

As caixas de diálogo de mensagem, abertura e salvamento aceitam opções do macOS, e o gerenciador de diálogos passa a oferecer uma solicitação de texto, além dos painéis de cor e fonte do sistema.

### Alertas

`SetSuppression` adiciona uma caixa de seleção "Não mostrar esta mensagem novamente", e `SetHelp` mostra o botão de ajuda. Leia o estado da caixa com `Suppressed` em um callback de botão ou registre `OnSuppression` para recebê-lo primeiro.

```go
dialog := app.Dialog.Warning().
    SetTitle("Delete 3 items?").
    SetMessage("The items will be moved to the Bin.").
    SetSuppression("Do not warn me again").
    SetHelp(func() {
        log.Println("help requested")
    }).
    AttachToWindow(window)

dialog.AddButton("Delete").SetAsDefault().OnClick(func() {
    if dialog.Suppressed() {
        // remember not to ask again
    }
})
dialog.AddButton("Cancel").SetAsCancel()
dialog.Show()
```

### Solicitação de texto

`Prompt` mostra um alerta com um campo de texto e bloqueia até que ele seja fechado; portanto, chame-o de uma goroutine ou de um método vinculado. `Secure` transforma o campo em um campo de senha, e `Window` apresenta o alerta como uma janela modal anexada.

```go
go func() {
    value, ok, err := app.Dialog.Prompt(application.PromptOptions{
        Title:        "Name this document",
        Message:      "The name is used for the exported file.",
        Placeholder:  "Untitled",
        DefaultValue: "Quarterly report",
        OKLabel:      "Rename",
        Window:       window,
    })
    if err != nil || !ok {
        return
    }
    log.Println("renamed to", value)
}()
```

### Painéis de arquivos

`AddContentType` filtra por identificador de tipo uniforme nas caixas de diálogo de abertura e salvamento. Ele funciona junto com `AddFilter`: assim, `"public.image"` corresponde a todos os tipos de imagem conhecidos pelo sistema, enquanto um filtro ainda captura PDFs pela extensão.

```go
paths, err := app.Dialog.OpenFile().
    SetTitle("Choose images").
    AddFilter("PDF", "*.pdf").
    AddContentType("public.image").
    AttachToWindow(window).
    PromptForMultipleSelection()
if err == nil {
    log.Println(paths)
}
```

Os painéis de salvamento passam a ter um menu suspenso Formato, um rótulo personalizado para o campo de nome e etiquetas do Finder. `SetFormats` troca o tipo permitido e a extensão do campo de nome quando o usuário muda a opção do menu, e `SelectedFormat` informa a escolha final.

```go
formats := []application.DialogFormat{
    {Label: "PNG image", Extension: "png", UTI: "public.png"},
    {Label: "PDF document", Extension: "pdf", UTI: "com.adobe.pdf"},
}
save := app.Dialog.SaveFile().
    SetFilename("Quarterly report.png").
    SetNameFieldLabel("Export As:").
    SetTags([]string{"Reports", "Draft"}).
    AttachToWindow(window)
save.SetFormats(formats, 0, nil)

path, err := save.PromptForSingleSelection()
if err == nil && path != "" {
    log.Println("export", path, "as", formats[save.SelectedFormat()].Label)
}
```

### Painéis de cor e fonte

`PickColor` e `PickFont` abrem os painéis compartilhados do sistema e bloqueiam até que o painel seja fechado. `OnChange` entrega cada seleção enquanto o painel está aberto, permitindo que a página mostre uma prévia em tempo real. Somente um painel de cada tipo pode estar aberto por vez; uma segunda chamada retorna `ErrDialogInProgress`.

```go
go func() {
    colour, changed, err := app.Dialog.PickColor(application.ColorPickerOptions{
        Initial:    application.NewRGB(52, 120, 246),
        ShowsAlpha: true,
        Title:      "Accent colour",
        OnChange: func(colour application.RGBA) {
            app.Event.Emit("accent:preview", colour)
        },
    })
    if err == nil && changed {
        log.Println("accent", colour)
    }

    font, changed, err := app.Dialog.PickFont(application.FontPickerOptions{
        Family: "Helvetica Neue",
        Size:   18,
    })
    if err == nil && changed {
        log.Println(font.Family, font.Face, font.PostScriptName, font.Size)
    }
}()
```

## Itens de status e feedback

### Itens de status

Um item da bandeja do sistema no macOS é um `NSStatusItem`. Ele pode ser desenhado a partir de um SF Symbol, conter uma dica de ferramenta e ser removido pelo usuário da mesma forma que os itens integrados.

```go
tray := app.SystemTray.New()
tray.SetSymbol("waveform.circle").SetSymbolConfiguration(0, application.MacSymbolWeightMedium)
tray.SetTooltip("Recorder. Command-drag to remove.")
tray.SetRemovable(true, "com.example.recorder.tray")
tray.OnVisibilityChange(func(visible bool) {
    log.Println("status item visible:", visible)
})

trayMenu := app.Menu.New()
trayMenu.Add("Show").OnClick(func(*application.Context) {
    window.Show().Focus()
})
tray.SetMenu(trayMenu)
tray.Run()
```

- `SetSymbol` renderiza o símbolo como uma imagem modelo para acompanhar a aparência da barra de menus (macOS 11+). `SetSymbolConfiguration` define o tamanho em pontos e a espessura.
- `SetTooltip` define o texto exibido ao passar o cursor. `Tooltip` lê o valor de volta.
- `SetRemovable(true, name)` permite ao usuário arrastar o item para fora da barra de menus com a tecla Command pressionada. Forneça um nome de salvamento automático estável para que o macOS se lembre da remoção entre execuções. `Show` ou `SetVisible(true)` o traz de volta.
- `IsVisible` lê `NSStatusItem.visible`; portanto, retorna falso depois que o usuário remove o item. `OnVisibilityChange` informa todas as alterações.

### Respostas táteis

`app.Haptics.Perform` reproduz um padrão em um trackpad Force Touch ou Magic Trackpad enquanto a aplicação está ativa.

```go
app.Haptics.Perform(application.HapticAlignment)
```

Os tipos são `HapticGeneric`, `HapticAlignment` (um item encaixado na posição) e `HapticLevelChange` (um ponto de resistência ou estágio de clique). `IsSupported` informa se a plataforma consegue produzir feedback.

### Sons

`app.Sound` reproduz o som de alerta, um som do sistema identificado por nome ou um arquivo de áudio.

```go
app.Sound.Beep()

if err := app.Sound.Play("Glass"); err != nil {
    log.Println(err)
}

for _, name := range app.Sound.SystemSounds() {
    log.Println(name)
}
```

`Play` aceita um nome de `SystemSounds` ou um caminho absoluto para qualquer arquivo que o Core Audio consiga decodificar. `PlayData` reproduz um arquivo de áudio completo da memória.

### Fala

`app.Speech.Speak` coloca texto na fila para a voz do sistema e retorna um `Utterance`. As falas são reproduzidas uma após a outra; `Stop` descarta uma, e `StopAll` limpa a fila. `Voices` lista as vozes instaladas com seus identificadores e idiomas.

```go
utterance, err := app.Speech.Speak("Export finished", application.SpeechOptions{
    Voice: "com.apple.voice.compact.en-GB.Daniel",
    Rate:  0.5,
})
if err != nil {
    log.Println(err)
    return
}
utterance.OnFinished(func() {
    log.Println("done, stopped:", utterance.WasStopped())
})
```

`Recognize` transcreve o áudio do microfone padrão com `SFSpeechRecognizer`. A primeira chamada solicita permissões de microfone e reconhecimento de fala e bloqueia até que o usuário responda; portanto, chame-a de uma goroutine. Transcrições parciais chegam por `OnPartial`; `Stop` encerra a captura e retorna o texto final.

```go
go func() {
    session, err := app.Speech.Recognize(application.RecognitionOptions{
        Locale: "en-US",
        OnPartial: func(text string) {
            app.Event.Emit("dictation:partial", text)
        },
    })
    if err != nil {
        if errors.Is(err, application.ErrSpeechRecognitionDenied) {
            _ = app.Permissions.OpenSystemSettings(application.PermissionKindMicrophone)
        }
        log.Println(err)
        return
    }
    time.Sleep(5 * time.Second)
    text, err := session.Stop()
    if err != nil {
        log.Println(err)
        return
    }
    app.Event.Emit("dictation:final", text)
}()
```

O reconhecimento exige uma aplicação empacotada cujo `Info.plist` declare `NSSpeechRecognitionUsageDescription` e `NSMicrophoneUsageDescription`. Sem essas chaves, o macOS recusa o acesso, e `Recognize` retorna `ErrSpeechRecognitionUsageDescription`.

## Área de transferência e arraste

### Área de transferência com conteúdo avançado

`app.Clipboard` lê e grava imagens, referências a arquivos, HTML, RTF e dados brutos sob qualquer identificador de tipo uniforme, além de texto simples. `Types` lista o que há na área de transferência, e `OnChange` informa alterações feitas por qualquer aplicação.

```go
if err := app.Clipboard.SetHTML("<p>Rich <b>HTML</b></p>", "Rich HTML"); err != nil {
    log.Println(err)
}
_ = app.Clipboard.SetFiles([]string{"/Users/me/Documents/Report.md"})
_ = app.Clipboard.SetData("com.example.record", []byte(`{"id":42}`))

stop := app.Clipboard.OnChange(func() {
    log.Println("clipboard changed", app.Clipboard.ChangeCount(), app.Clipboard.Types())
    if files, err := app.Clipboard.Files(); err == nil && len(files) > 0 {
        log.Println("files:", files)
    }
})
defer stop()
```

`SetImage` e `Image` trabalham com bytes PNG; imagens copiadas como TIFF por outros aplicativos são convertidas automaticamente. Não há notificação do sistema para alterações na área de transferência; por isso, `OnChange` consulta a contagem de alterações a cada 500 ms enquanto existir pelo menos um observador.

### Arraste para fora

`StartDrag` inicia um arraste do sistema a partir da janela, como se o usuário tivesse pegado os itens no Finder. Ele oferece arquivos existentes, promessas de arquivos cujo conteúdo só é produzido quando um destino aceita a soltura, ou texto simples. Inicie-o durante um gesto do mouse: vincule um método Go e chame-o no manipulador `mousedown` ou `pointerdown` da página, no elemento arrastável, com o atributo HTML `draggable` definido como `false` para que o WebKit não inicie seu próprio arraste.

```go
// DragService is bound to the page and called from a mousedown handler.
type DragService struct {
    window *application.WebviewWindow
}

func (s *DragService) DragExport() error {
    return s.window.StartDrag(application.DragItems{
        Promises: []application.DragPromise{{
            Filename: "export.csv",
            Data: func() ([]byte, error) {
                return []byte("id,name\n1,Wails\n"), nil
            },
        }},
        Operations: application.DragOperationCopy,
    })
}
```

```go
window.OnDragEnd(func(operation application.DragOperation) {
    log.Println("drag ended:", operation)
})
```

Chamar `StartDrag` fora de um gesto retorna `ErrDragOutNoGesture`. `DragItems.Image` e `ImageOffset` definem a imagem sob o cursor.

### Itens soltos por outros aplicativos

A soltura de arquivos continua usando o evento `WindowFilesDropped`. Para aceitar texto, URLs ou imagens arrastados de outros aplicativos, liste os tipos em `DropTypes` e registre `OnDrop`. Esses itens são entregues ao Go, em vez de aos manipuladores de soltura HTML5 da própria página.

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "Inbox",
    URL:   "/",
    DropTypes: []application.DropType{
        application.DropFiles,
        application.DropText,
        application.DropURLs,
        application.DropImages,
    },
})

window.OnWindowEvent(events.Common.WindowFilesDropped, func(event *application.WindowEvent) {
    log.Println("files:", event.Context().DroppedFiles())
})

window.OnDrop(func(_ *application.Context, data application.DropData) {
    log.Println("text:", data.Text, "urls:", data.URLs, "images:", len(data.Images), "at", data.X, data.Y)
})
```

## Sistema

### Permissões

`app.Permissions` informa e solicita permissões de privacidade do sistema (câmera, microfone, gravação de tela, acessibilidade, localização, notificações, monitoramento de entrada e acesso total ao disco). `Status` nunca solicita permissão. `Request` solicita permissões cujo estado ainda não foi determinado e bloqueia até que o usuário responda; portanto, chame-o de uma goroutine. `OpenSystemSettings` abre o painel correspondente em Privacidade e Segurança.

```go
go func() {
    status := app.Permissions.Status(application.PermissionKindCamera)
    if status == application.PermissionStatusNotDetermined {
        status, _ = app.Permissions.Request(application.PermissionKindCamera)
    }
    if status == application.PermissionStatusDenied {
        _ = app.Permissions.OpenSystemSettings(application.PermissionKindCamera)
    }
}()
```

O acesso total ao disco não pode ser solicitado e retorna `ErrPermissionNotRequestable`; oriente o usuário a abrir o painel de ajustes. Uma solicitação que nunca recebe resposta retorna `ErrPermissionRequestTimeout`, o que no macOS geralmente significa que falta a chave de descrição de uso correspondente em `Info.plist`.

A opção de janela `Permissions` agora é respeitada no macOS. Ela determina como são tratadas as solicitações de `getUserMedia` feitas pela página: `PermissionAllow` dispensa a solicitação da própria WebView, `PermissionDeny` recusa sem perguntar e `PermissionDefault` mostra a solicitação. A solicitação TCC no nível do sistema ainda aparece na primeira vez que a câmera ou o microfone é usado.

```go
app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "Meeting",
    URL:   "/",
    Permissions: map[application.PermissionType]application.Permission{
        application.PermissionMicrophone: application.PermissionAllow,
        application.PermissionCamera:     application.PermissionDefault,
    },
})
```

### Energia

`app.Power.PreventSleep` mantém o sistema ativo e, com `Display`, também a tela, até que a função de liberação retornada seja chamada. As retenções são contabilizadas, de modo que várias partes da aplicação podem manter o sistema ativo ao mesmo tempo. O motivo é mostrado no Monitor de Atividade.

```go
release, err := app.Power.PreventSleep("Exporting video", application.PreventSleepOptions{Display: true})
if err != nil {
    log.Println(err)
}
defer release()

state := app.Power.State()
if state.LowPowerMode || state.ThermalState >= application.ThermalStateSerious {
    // trim background work
}
log.Println("battery", state.BatteryLevel, "charging", state.Charging, "on battery", state.OnBattery)
```

As alterações chegam como `events.Mac.ApplicationDidChangePowerState` (alternância do Modo de Pouca Energia) e `events.Mac.ApplicationDidChangeThermalState`.

### Ciclo de vida

O macOS pode encerrar instantaneamente uma aplicação ociosa durante o logout ou desligamento quando ela aceita esse comportamento com `NSSupportsSuddenTermination`, e pode encerrar uma aplicação ociosa sem janelas quando ela aceita esse comportamento com `NSSupportsAutomaticTermination`. `app.Lifecycle.HoldTermination` suspende ambos durante uma seção crítica, como o salvamento de um arquivo.

```go
release := app.Lifecycle.HoldTermination("Saving document")
defer release()
// write the file
```

`SetSuddenTerminationEnabled` alterna o encerramento repentino em tempo de execução; `SuddenTerminationEnabled` informa o estado atual, cujo valor inicial vem da chave em `Info.plist`.

### Ambiente

`app.Env` passa a oferecer três consultas sobre as configurações do usuário.

```go
a11y := app.Env.Accessibility()
if a11y.ReduceMotion || a11y.ReduceTransparency {
    app.Event.Emit("theme:calm", true)
}

layout := app.Env.KeyboardLayout()
log.Println(layout.ID, layout.Name, layout.Languages)

locale := app.Env.Locale()
log.Println(locale.Identifier, locale.Language, locale.Region, locale.Preferred)
```

- `Accessibility` reflete Reduzir Movimento, Reduzir Transparência, Aumentar Contraste, Diferenciar sem Cor, Inverter Cores, VoiceOver e Controle Assistivo.
- `KeyboardLayout` retorna a fonte de entrada ativa com seu identificador, nome localizado e idiomas.
- `Locale` retorna a localidade selecionada pelo AppKit para a aplicação, além da lista completa e ordenada `Preferred` do usuário. `Identifier` reflete apenas um idioma declarado pelo pacote em `CFBundleLocalizations`; use `Preferred` para escolher um idioma por conta própria.

### Eventos

Estes eventos da aplicação são novos. Cada um é entregue por `app.Event.OnApplicationEvent`; consulte o gerenciador correspondente para obter o valor atualizado.

| Evento | Acionado quando | Ler com |
|-------|------------|-----------|
| `events.Mac.ApplicationDidChangePowerState` | o Modo de Pouca Energia é alternado | `app.Power.State()` |
| `events.Mac.ApplicationDidChangeThermalState` | a pressão térmica muda | `app.Power.State()` |
| `events.Common.AccessibilitySettingsChanged` | uma configuração de exibição de acessibilidade muda | `app.Env.Accessibility()` |
| `events.Mac.ApplicationDidChangeAccessibilitySettings` | ocorre a forma macOS da mesma alteração | `app.Env.Accessibility()` |
| `events.Mac.ApplicationDidChangeKeyboardLayout` | a fonte de entrada muda | `app.Env.KeyboardLayout()` |
| `events.Mac.ApplicationDidChangeLocale` | a localidade muda | `app.Env.Locale()` |

```go
app.Event.OnApplicationEvent(events.Mac.ApplicationDidChangeThermalState, func(*application.ApplicationEvent) {
    app.Event.Emit("system:power", app.Power.State())
})
app.Event.OnApplicationEvent(events.Common.AccessibilitySettingsChanged, func(*application.ApplicationEvent) {
    app.Event.Emit("system:accessibility", app.Env.Accessibility())
})
```

## Integração

### Menu Serviços

`app.ServicesProvider.Register` adiciona uma entrada ao submenu Serviços que todos os aplicativos macOS mostram para texto ou arquivos selecionados. O manipulador recebe a área de transferência como `ServiceRequest` e retorna um `ServiceResponse` para gravar o resultado; uma resposta vazia deixa a seleção intacta. Mantenha o manipulador rápido, pois o AppKit o aguarda na thread principal.

```go
err := app.ServicesProvider.Register(application.ServiceDefinition{
    Name:          "summarise",
    MenuTitle:     "Summarise with Notes",
    SendTypes:     []string{"public.utf8-plain-text"},
    ReturnTypes:   []string{"public.utf8-plain-text"},
    KeyEquivalent: "S",
    Handler: func(_ *application.Context, request application.ServiceRequest) (application.ServiceResponse, error) {
        return application.ServiceResponse{Text: "Summary: " + request.Text}, nil
    },
})
if err != nil {
    log.Println(err)
}

// Paste this inside the top-level <dict> of Info.plist.
log.Println(app.ServicesProvider.InfoPlistXML())
```

`Name` é a mensagem enviada pelo AppKit e deve ser um identificador simples. `SendTypes` e `ReturnTypes` são tipos da área de transferência; um serviço precisa de pelo menos um deles. O registro, por si só, não torna o serviço visível: o `Info.plist` do pacote deve declará-lo em `NSServices`. `InfoPlistXML` retorna esse bloco pronto para colar, e `InfoPlistEntries` retorna os mesmos dados como mapas para um serializador de plist. `NSPortName` nessas entradas é o `Name` da aplicação, que deve corresponder a `CFBundleName`.

Projetos compilados com a CLI do Wails podem declarar os mesmos serviços uma vez em `build/config.yml` e deixar que o empacotamento gere o bloco `NSServices`:

```yaml
services:
  - name: SummariseText
    menuTitle: Summarise with My Product
    sendTypes:
      - public.utf8-plain-text
    returnTypes:
      - public.utf8-plain-text
    keyEquivalent: S
```

Cada entrada deve corresponder a uma `ServiceDefinition` registrada em Go com o mesmo `Name`. Execute `pbs -update` após instalar uma nova compilação para que o menu Serviços reconheça a alteração sem exigir logout.

### Handoff e atividades do usuário

`app.Activity.Publish` torna uma `NSUserActivity` atual para que o usuário possa continuá-la em outro dispositivo, encontrá-la no Spotlight ou receber uma sugestão da Siri. A `PublishedActivity` retornada pode ser atualizada conforme o estado muda e invalidada quando o documento é fechado. Publicar uma nova atividade substitui a anterior.

```go
activity, err := app.Activity.Publish(application.UserActivity{
    Type:               "com.example.notes.editing",
    Title:              "Editing Quarterly report",
    UserInfo:           map[string]any{"note": "quarterly-report"},
    WebpageURL:         "https://example.com/notes/quarterly-report",
    EligibleForHandoff: true,
    EligibleForSearch:  true,
    Keywords:           []string{"report", "quarterly"},
})
if err != nil {
    log.Println(err)
    return
}
if err := activity.Update(map[string]any{"note": "quarterly-report", "cursor": 120}); err != nil {
    log.Println(err)
}
// When the document closes:
activity.Invalidate()
```

As atividades recebidas chegam por `OnContinue`. Links universais têm o tipo `UserActivityTypeBrowsingWeb`, com a página em `WebpageURL`, e também são entregues como `events.Common.ApplicationLaunchedWithUrl`, permitindo que a aplicação use um único fluxo de tratamento de URLs. `OnWillContinue`, `OnFailed` e `OnUpdated` cobrem o restante do delegate.

```go
app.Activity.OnContinue(func(_ *application.Context, incoming application.UserActivity) bool {
    if incoming.Type == application.UserActivityTypeBrowsingWeb {
        log.Println("universal link", incoming.WebpageURL)
        return true
    }
    note, _ := incoming.UserInfo["note"].(string)
    log.Println("continue editing", note)
    return true
})
```

Todos os tipos de atividade devem constar em `NSUserActivityTypes` no `Info.plist`. Links universais também exigem o entitlement `com.apple.developer.associated-domains` com uma entrada `applinks:example.com` e o arquivo `apple-app-site-association` correspondente nesse domínio.

### Apple Events

`app.AppleEvents.Handle` registra um manipulador para uma classe e um ID de evento, permitindo que AppleScript, Shortcuts e outros aplicativos controlem a aplicação. Os códigos são strings de quatro caracteres. O parâmetro direto é decodificado em um valor Go (`string`, `[]string` de caminhos de arquivos, `int64`, `float64`, `bool`, `[]any` ou `AppleEventRawData`), e o `Result` da resposta aceita os mesmos tipos. Os manipuladores executam em suas próprias goroutines enquanto o evento está suspenso.

```go
err := app.AppleEvents.Handle("WAIL", "note", func(_ *application.Context, event application.AppleEvent) (application.AppleEventReply, error) {
    text, _ := event.DirectObject.(string)
    return application.AppleEventReply{Result: "noted: " + text}, nil
})
if err != nil {
    log.Println(err)
}

go func() {
    result, err := app.AppleEvents.Send("com.apple.finder", "misc", "actv", nil)
    log.Println(result, err)
}()

sdef := app.AppleEvents.ScriptingDefinition()
if err := os.WriteFile("Notes.sdef", []byte(sdef), 0o644); err != nil {
    log.Println(err)
}
```

Um script pode chamar o manipulador imediatamente com a sintaxe bruta de eventos:

```applescript
tell application id "com.example.notes" to «event WAILnote» "hello"
```

`ScriptingDefinition` gera um `.sdef` mínimo que dá um nome de comando a cada manipulador. Inclua-o em `Contents/Resources` e faça `Info.plist` apontar para ele com `NSAppleScriptEnabled` e `OSAScriptingDefinition`; o Editor de Scripts então o mostra em Arquivo > Abrir Dicionário. O Wails já trata o evento Get URL para esquemas de URL personalizados; um manipulador para `"GURL"`/`"GURL"` é encadeado a esse tratamento, enquanto um manipulador para `"aevt"`/`"odoc"` substitui a entrega integrada de Open Documents. `Send` direciona a chamada a uma aplicação em execução pelo identificador do pacote, exige `NSAppleEventsUsageDescription` em uma aplicação empacotada e bloqueia a goroutine chamadora até a resposta chegar.

### Quick Look

`app.QuickLook.Preview` abre o painel compartilhado do Quick Look para um ou mais arquivos; com vários caminhos, o painel mostra setas para alternar entre eles. `Thumbnail` renderiza um arquivo por meio dos provedores de miniaturas do sistema e retorna um PNG; portanto, funciona com documentos, imagens, PDFs e vídeos.

```go
if err := app.QuickLook.Preview([]string{"/Users/me/Documents/Report.pdf"}); err != nil {
    log.Println(err)
}

go func() {
    png, err := app.QuickLook.Thumbnail("/Users/me/Documents/Report.pdf", application.ThumbnailOptions{
        Width: 256,
        Scale: 2,
    })
    if err != nil {
        log.Println(err)
        return
    }
    if err := os.WriteFile("thumbnail.png", png, 0o644); err != nil {
        log.Println(err)
    }
}()
```

Os caminhos devem ser absolutos e existir. `ClosePreview` e `IsPreviewOpen` gerenciam o painel. `ThumbnailOptions.IconMode` desenha a borda de documento no estilo do Finder, e `Scale: 2` produz uma imagem Retina. `Thumbnail` bloqueia a goroutine chamadora; portanto, chame-o de uma goroutine ou de um método vinculado.

### Funções auxiliares do espaço de trabalho

`app.Browser` passa a oferecer três funções auxiliares para `NSWorkspace`. `OpenWith` abre um arquivo com um aplicativo específico, indicado pelo identificador ou caminho do pacote. `ApplicationsForFile` lista os aplicativos instalados que podem abrir um arquivo, com o manipulador padrão primeiro. `ActivateApplication` traz uma aplicação em execução para a frente.

```go
apps := app.Browser.ApplicationsForFile("/Users/me/Documents/Report.md")
for _, info := range apps {
    log.Println(info.Name, info.BundleID, info.Path)
}

if err := app.Browser.OpenWith("/Users/me/Documents/Report.md", "com.apple.TextEdit"); err != nil {
    log.Println(err)
}

if err := app.Browser.ActivateApplication("com.apple.TextEdit"); err != nil {
    log.Println(err)
}
```

### Spotlight

`app.Spotlight.Index` adiciona conteúdo da aplicação ao índice de pesquisa do sistema por meio do Core Spotlight. Cada `SearchableItem` tem um `ID`, um `Title` e, opcionalmente, um `Domain` para remoção em massa, uma `Description`, `Keywords`, um `ContentType`, uma miniatura PNG, um `URL` de link direto e um prazo de expiração. `OnOpen` é chamado quando o usuário escolhe um dos itens no Spotlight ou seleciona "Buscar no Aplicativo" com uma consulta.

```go
err := app.Spotlight.Index([]application.SearchableItem{{
    ID:          "note:quarterly-report",
    Domain:      "notes",
    Title:       "Quarterly report",
    Description: "Draft for the board meeting",
    Keywords:    []string{"finance", "q3"},
    ContentType: "public.plain-text",
    URL:         "notes://open/quarterly-report",
}})
if err != nil {
    log.Println(err)
}

stop := app.Spotlight.OnOpen(func(_ *application.Context, id string, query string) {
    if query != "" {
        log.Println("search in app:", query)
        return
    }
    log.Println("open item", id)
})
defer stop()

// Later:
_ = app.Spotlight.Delete([]string{"note:quarterly-report"})
_ = app.Spotlight.DeleteDomain("notes")
```

`IsAvailable` informa se o índice aceita itens. A indexação exige uma aplicação empacotada: itens indexados por um binário `go run` sem pacote nunca aparecem no Spotlight. `DeleteAll` remove tudo o que a aplicação indexou.

## Em breve

Um exemplo `mac-windows-extra` que aborda janelas modais anexadas, popovers, opções de apresentação e restauração de estado está sendo adicionado e será vinculado aqui quando estiver disponível.

## Requisitos de versão

Tudo nesta página é exclusivo do macOS, e a API Go é idêntica em todas as plataformas. O Wails oferece suporte ao macOS 10.13 e posteriores; recursos que exigem uma versão mais recente têm o comportamento alternativo indicado.

| Recurso | macOS mínimo | Comportamento em versões anteriores |
|---------|---------------|-------------------------------|
| Estado das permissões de câmera e microfone | 10.14 | informado como autorizado (versões anteriores não restringem dispositivos de captura) |
| `Speech.Speak` e `Voices` | 10.14 | `ErrSpeechNotSupported` |
| Permissões de gravação de tela e monitoramento de entrada | 10.15 | informadas como autorizadas |
| `Speech.Recognize` | 10.15 | `ErrSpeechRecognitionNotSupported` |
| `QuickLook.Thumbnail` | 10.15 | `ErrQuickLookNotSupported` |
| `SetSubtitle` | 11 | ignorado com um registro de depuração |
| `ExportPDF` | 11 | `ErrMacExportUnsupported` |
| `SetSymbol` em itens de menu e de status | 11 | nenhuma imagem é mostrada |
| `AddContentType`, `SetFormats` por UTI | 11 | os mesmos identificadores são aplicados pela API legada de tipos de arquivo permitidos |
| Ícones de promessas de arquivos em `StartDrag` | 11 | um ícone genérico de documento |
| `PowerState.LowPowerMode` e seu evento | 12 | sempre falso; o evento nunca é acionado |
| `InteractionState` e `RestoreInteractionState` | 12 | `ErrMacInteractionStateUnsupported` |
| Painel de Notificações em `OpenSystemSettings` | 13 | o painel de preferências de Notificações mais antigo é aberto |
| Indicadores de menu, cabeçalhos de seção, paletas | 14 | os indicadores não são mostrados; os cabeçalhos são itens desabilitados; as paletas ficam ocultas |
| `MacToolbarItem.ShowPopover` em itens sem uma visualização personalizada | 14 | `ErrMacPopoverAnchorUnavailable` |

Todos os demais recursos não exigem versão além do mínimo do Wails.

## Chaves do Info.plist

Vários recursos dependem de chaves no `Info.plist` da aplicação. As descrições de uso são mostradas ao usuário na solicitação de permissão; sem a chave, o macOS nunca mostra a solicitação e ela expira.

| Chave | Necessária para |
|-----|-----------|
| `NSCameraUsageDescription` | `Permissions.Request(PermissionKindCamera)`, acesso à câmera pela página |
| `NSMicrophoneUsageDescription` | `Permissions.Request(PermissionKindMicrophone)`, acesso ao microfone pela página, `Speech.Recognize` |
| `NSSpeechRecognitionUsageDescription` | `Speech.Recognize` |
| `NSLocationUsageDescription` | `Permissions.Request(PermissionKindLocation)` |
| `NSSupportsSuddenTermination` | o estado inicial de `Lifecycle.SuddenTerminationEnabled`; `HoldTermination` o suspende |
| `NSSupportsAutomaticTermination` | permite ao macOS encerrar a aplicação ociosa; `HoldTermination` o suspende |
| `CFBundleLocalizations` | quais idiomas `Env.Locale().Identifier` pode informar |
| `NSServices` | uma entrada por serviço registrado com `app.ServicesProvider`; gere o bloco com `InfoPlistXML` |
| `NSUserActivityTypes` | cada `UserActivity.Type` publicado ou continuado por `app.Activity` |
| `NSAppleScriptEnabled` e `OSAScriptingDefinition` | marcam a aplicação como controlável por script e identificam o `.sdef` gerado por `app.AppleEvents.ScriptingDefinition` |
| `NSAppleEventsUsageDescription` | `app.AppleEvents.Send` para outros aplicativos |

A permissão de notificações e o reconhecimento de fala também exigem que a aplicação seja executada como um pacote com identificador de pacote; um binário `go run` sem pacote informa `PermissionStatusUnsupported` para notificações.

<a id="platform-notes"></a>

## Notas sobre plataformas

Todas as APIs desta página compilam no Windows e no Linux. Fora do macOS:

- Recursos adicionais de janela: `SetRepresentedFile`, `SetDocumentEdited`, `SetSubtitle`, `CascadeFrom`, `SetFrameAutosaveName` e `SetWindowButtonsOffset` não têm efeito. `WindowCascade` se comporta como `WindowCentered`. `RequestAttention` retorna um identificador cujo `Cancel` não tem efeito. `PrintWithOptions` chama `Print`. `ExportPDF` e `Snapshot` retornam `ErrMacOnly`.
- Menus: símbolos, indicadores, estado misto, alternativas e recuo são armazenados, mas não desenhados. Cabeçalhos de seção são itens desabilitados. Paletas ficam ocultas. O submenu Abrir Recentes é preenchido com a lista de recentes em Go quando o menu é criado. Os menus do Dock nunca são mostrados.
- Caixas de diálogo: `SetSuppression`, `SetHelp`, `SetNameFieldLabel` e `SetTags` são ignorados. `AddContentType` e `SetFormats` se tornam filtros de extensão quando há uma tradução conhecida. `Prompt`, `PickColor` e `PickFont` retornam `ErrDialogNotSupported`.
- Itens de status: `SetSymbol`, `SetRemovable` e `OnVisibilityChange` não têm efeito. `IsVisible` reflete a última chamada a `Show` ou `Hide`.
- Feedback: as respostas táteis chegam ao iOS e ao Android; no Windows e no Linux, não têm efeito. `Sound.Beep` e `Sound.Play` funcionam no Windows com arquivos WAV e aliases do Registro; em outros sistemas, `Play` retorna `ErrSoundNotSupported`. A fala retorna `ErrSpeechNotSupported` e `ErrSpeechRecognitionNotSupported`.
- Área de transferência: os métodos de conteúdo avançado retornam `ErrClipboardNotSupported`, `Types` está vazio, `ChangeCount` é 0 e `OnChange` nunca é acionado.
- Arraste: `StartDrag` retorna `ErrDragOutUnsupported`. Tipos de itens soltos que não são arquivos não são entregues; a soltura de arquivos continua funcionando por `WindowFilesDropped`.
- Sistema: `Permissions.Status` informa `PermissionStatusUnsupported`, e `Request` retorna `ErrPermissionsUnsupported`. `PreventSleep` retorna `ErrPreventSleepUnsupported` com uma função de liberação sem efeito. `HoldTermination` retorna uma função de liberação sem efeito. `Accessibility` tem todos os campos falsos, `KeyboardLayout` é o valor zero e `Locale` é derivada de `LC_ALL`, `LC_MESSAGES` e `LANG`. A opção de janela `Permissions` é multiplataforma.
- Integração: `ServicesProvider.Register` retorna `ErrServicesUnsupported`, e `Activity.Publish` retorna `ErrActivityUnsupported`; `InfoPlistXML`, `InfoPlistEntries` e os manipuladores de atividade continuam funcionando. `Handle` e `Send` de `AppleEvents` retornam `ErrAppleEventsNotSupported`; `ScriptingDefinition` é gerado em todas as plataformas. `QuickLook.Preview` e `Thumbnail` retornam `ErrQuickLookNotSupported`. Os métodos de indexação de `Spotlight` retornam `ErrSpotlightNotSupported`, e `OnOpen` nunca é acionado. `Browser.OpenWith` inicia o executável indicado, passando o caminho como argumento; `ApplicationsForFile` está vazio, e `ActivateApplication` retorna `ErrApplicationNotRunning`.
- Opções de apresentação: `SetPresentationOptions` retorna `ErrMacOnly`, e `PresentationOptions` é `MacPresentationDefault`.
- Janelas modais anexadas: `PresentSheet`, `PresentCriticalSheet` e `PresentNativeSheet` retornam `ErrMacSheetUnsupported`; `EndSheet` não tem efeito, e os métodos de consulta informam que não há janela modal anexada.
- Popovers: `NewMacPopover` funciona, os métodos de exibição retornam `ErrMacPopoverUnsupported`, e `IsShown` é falso.
- Restauração de estado: `SetRestorationID` e `SetRestorationData` não têm efeito, `OnRestore` nunca é chamado, e `InteractionState` e `RestoreInteractionState` retornam `ErrMacOnly`.

## Exemplos

Cada exemplo é uma aplicação completa e executável:

- [`v3/examples/mac-toolbar`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-toolbar): `windowextras.go` integra ao editor de notas o ponto de alterações não salvas, o subtítulo, a exportação de PDF, as opções de impressão e as janelas em cascata.
- [`v3/examples/mac-menus-dock`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-menus-dock): símbolos, indicadores, cabeçalhos de seção, estado misto, alternativas, uma paleta, Abrir Recentes, menu dinâmico do Dock e progresso no Dock.
- [`v3/examples/mac-dialogs`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-dialogs): caixa de seleção de supressão e botões de ajuda, solicitações de texto, tipos de conteúdo, menu suspenso Formato, etiquetas do Finder e painéis de cor e fonte.
- [`v3/examples/mac-feedback`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-feedback): um item de status removível com SF Symbol, respostas táteis, sons do sistema, síntese de fala e reconhecimento de fala.
- [`v3/examples/mac-clipboard-drag`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-clipboard-drag): área de transferência com conteúdo avançado e acompanhamento de alterações, arraste para fora com promessas de arquivos e soltura de texto, URLs e imagens.
- [`v3/examples/mac-system`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-system): permissões, prevenção de suspensão, retenções de encerramento, estado de energia, acessibilidade, layout do teclado e localidade com atualizações em tempo real.
- [`v3/examples/mac-integration`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-integration): uma entrada no menu Serviços, uma atividade Handoff com manipuladores de continuação e as funções auxiliares do espaço de trabalho em `app.Browser`.
- [`v3/examples/mac-search-preview`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-search-preview): indexação do Spotlight com `OnOpen`, prévias e miniaturas do Quick Look e um Apple Event personalizado com sua definição de script.

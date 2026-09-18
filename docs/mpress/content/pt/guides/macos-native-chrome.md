---
title: "Elementos nativos da janela no macOS"
description: "Crie barras de ferramentas, barras laterais, listas de conteúdo, inspetores, acessórios e abas de janela nativos do AppKit em torno da sua janela Wails"
slug: "guides/macos-native-chrome"
sourcePath: "guides/macos-native-chrome.md"
---

Plataformas relevantes: macOS

O Wails v3 pode envolver sua WebView com elementos de janela reais do AppKit. A barra de ferramentas, a barra lateral, a lista de conteúdo, o inspetor e as faixas da barra de título são controles nativos criados em Go. Nada disso é HTML; assim, eles recebem os materiais, o tratamento de teclado, as animações e a persistência do AppKit, enquanto seu frontend permanece focado no conteúdo.

Todas as APIs desta página estão no pacote `application` e têm o prefixo `Mac`. O mesmo código compila no Windows e no Linux: construtores e setters funcionam em todas as plataformas, chamadas de anexação não têm efeito ou retornam um erro, e a janela mantém sua WebView única habitual.

## Estrutura da janela

Uma janela completa tem estas partes, da borda inicial à borda final:

| Parte | Tipo | Classe do AppKit |
|------|------|--------------|
| Barra de ferramentas | `MacToolbar` | `NSToolbar` |
| Barra lateral | `MacSidebar` | Lista de origem `NSOutlineView` em um item de divisão da barra lateral |
| Lista de conteúdo | `MacContentList` | `NSTableView` em um item de divisão da lista de conteúdo |
| Conteúdo principal | sua WebView ou `MacTextEditor` | `WKWebView` ou `NSTextView` |
| Inspetor | `MacInspector` | controles nativos de propriedades em um item de divisão do inspetor |
| Acessórios | `MacAccessory` | `NSTitlebarAccessoryViewController` ou `NSSplitViewItemAccessoryViewController` |

Os painéis são organizados por um `MacSplitView`, que é um `NSSplitViewController`. Crie primeiro as partes, adicione-as à visualização dividida da borda inicial à final, anexe a visualização dividida à janela e então anexe a barra de ferramentas. Você pode fazer isso com `SetSplitView` e `SetToolbar` ou em uma única etapa por meio das opções de janela `Mac.SplitView` e `Mac.Toolbar`.

Uma configuração típica combina esses elementos com as seguintes opções de janela para que o conteúdo possa rolar sob uma barra de ferramentas unificada:

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Mac: application.MacWindow{
        ContentLayout: application.MacContentLayoutEdgeToEdge,
        TitleBar: application.MacTitleBar{
            FullSizeContent:      true,
            HideToolbarSeparator: true,
            ToolbarStyle:         application.MacToolbarStyleUnified,
        },
    },
})
```

## Barra de ferramentas

`NewMacToolbar` cria uma `NSToolbar`. Adicione itens com os métodos `Add`, encadeie setters e callbacks nos identificadores retornados e anexe a barra de ferramentas com `SetToolbar`. Os identificadores são gerados automaticamente.

```go
toolbar := application.NewMacToolbar().
    SetDisplayMode(application.MacToolbarDisplayModeIconOnly)

// Standard AppKit items. They have no handle because AppKit owns them.
toolbar.AddSidebarToggle()
toolbar.AddSidebarTrackingSeparator()

toolbar.AddButton("New").
    SetSymbol("square.and.pencil").
    SetTooltip("Create a note").
    SetBordered(true).
    OnClick(func(*application.Context) {
        // create a note
    })

toolbar.AddSearch("Search").
    SetSearchPlaceholder("Search notes").
    SetSearchIncremental(true).
    OnSearch(func(_ *application.Context, query string) {
        filterNotes(query)
    })

toolbar.AddFlexibleSpace()

mode := toolbar.AddGroup("Mode", application.ToolbarGroupSelectOne)
mode.AddButton("Write").SetSymbol("pencil").OnClick(func(*application.Context) {
    mode.SetSelectedIndex(0)
})
mode.AddButton("Preview").SetSymbol("doc.richtext").OnClick(func(*application.Context) {
    mode.SetSelectedIndex(1)
})

actions := application.NewMenu()
actions.Add("Export...").OnClick(func(*application.Context) {})
toolbar.AddMenu("Actions", actions).SetSymbol("ellipsis.circle")

window.SetToolbar(toolbar)
```

Os tipos de item são:

- `AddButton` adiciona um botão. Todo botão precisa de um `OnClick` antes que a barra de ferramentas seja anexada; caso contrário, `SetToolbar` informa um erro e mantém a barra de ferramentas anterior.
- `AddSearch` adiciona um `NSSearchToolbarItem`. `OnSearch` é obrigatório. Use `SetSearchPlaceholder`, `SetSearchIncremental`, `SetSearchRecentsKey` para um menu persistente de pesquisas recentes e `SetSearchMenu` para um menu personalizado atrás da lupa.
- `AddShare` adiciona o item de compartilhamento do sistema (veja abaixo).
- `AddGroup` adiciona um `NSToolbarItemGroup` segmentado. Adicione membros com o `AddButton` do grupo e escolha `ToolbarGroupSelectOne`, `ToolbarGroupMomentary` ou `ToolbarGroupSelectAny`.
- `AddMenu` adiciona um `NSMenuToolbarItem` suspenso controlado por um `Menu` comum. `SetShowsIndicator(false)` oculta a seta.
- `AddSpace` e `AddFlexibleSpace` adicionam os espaçadores padrão.
- `AddSidebarToggle` e `AddSidebarTrackingSeparator` adicionam os itens de barra lateral do AppKit. O separador mantém tudo o que vem antes dele alinhado acima do divisor da barra lateral; por isso, a janela deve ter uma visualização dividida com um painel de barra lateral.
- `AddInspectorToggle` e `AddInspectorTrackingSeparator` fazem o mesmo para um painel de inspetor.

`SetDisplayMode` escolhe entre `MacToolbarDisplayModeIconAndLabel` (o padrão), `MacToolbarDisplayModeIconOnly`, `MacToolbarDisplayModeLabelOnly` e `MacToolbarDisplayModeDefault`.

### Atualizações em tempo real

Todo identificador também permite atualizações em tempo real. Setters aplicados após a anexação atualizam o item nativo na thread da aplicação.

```go
save := toolbar.AddButton("Save").SetSymbol("checkmark.circle")
save.OnClick(func(*application.Context) {
    save.SetBadgeCount(0).SetProminent(false)
})

// Later, when the document changes:
save.SetBadgeCount(1).SetProminent(true)
save.SetVisibilityPriority(application.MacToolbarVisibilityPriorityHigh)

toolbar.Move(save, 0)
toolbar.Remove(save)
```

Outros setters úteis são `SetLabel`, `SetTooltip`, `SetEnabled`, `SetHidden`, `SetTintColor` e `SetNavigational`, que mantém um item na borda inicial, como o Safari faz com os botões de voltar e avançar. `SetVisibilityPriority` determina quais itens vão primeiro para o menu de itens excedentes quando a janela se estreita.

### Personalização pelo usuário

`SetCustomizable` habilita a janela padrão "Personalizar Barra de Ferramentas..." e salva a organização do usuário sob a chave fornecida. Dê a cada item uma `SetPersistenceKey` estável para que a organização salva sobreviva à reinicialização e use `SetInDefaultSet(false)` para itens que só devem aparecer depois que o usuário os adicionar. Chame `SetCustomizable` antes de anexar a barra de ferramentas.

```go
toolbar := application.NewMacToolbar().SetCustomizable("myapp.main-toolbar")

newNote := toolbar.AddButton("New").
    SetPersistenceKey("new").
    OnClick(func(*application.Context) {})

// Offered in the palette but hidden until the user adds it.
toolbar.AddButton("Archive").
    SetPersistenceKey("archive").
    SetInDefaultSet(false).
    OnClick(func(*application.Context) {})

toolbar.SetCenteredItems(newNote)

// From a menu item, for example:
toolbar.RunCustomizationPalette()
```

### Compartilhamento

`AddShare` retorna um `MacToolbarShareItem`. Ele permanece desabilitado até que um `MacShareProvider` anuncie pelo menos uma representação. O Wails solicita os bytes ao provedor somente quando um serviço de compartilhamento os pede; assim, exportações grandes são geradas sob demanda. `MacShareProviderFunc` adapta duas funções para formar um provedor; aplicações com estado podem implementar a interface diretamente.

```go
share := toolbar.AddShare("Share")
share.SetProvider(application.MacShareProviderFunc{
    Available: []application.MacShareRepresentation{
        {ContentType: application.MacShareTypePDF},
        {ContentType: application.MacShareTypePlainText},
    },
    Load: func(request application.MacShareRequest) ([]byte, error) {
        if request.ContentType == application.MacShareTypePDF {
            return renderPDF()
        }
        return []byte(currentText()), nil
    },
}).SetSubject("Saturday, slowly").SetSuggestedName("Note")

share.OnShared(func(_ *application.Context, service string) {
    // service is the localised name of the sharing service
})
share.OnShareError(func(_ *application.Context, service string, err error) {
    window.Error("share via %s failed: %s", service, err)
})
```

Os tipos de conteúdo comuns são `MacShareTypePlainText`, `MacShareTypeHTML`, `MacShareTypePDF`, `MacShareTypePNG` e `MacShareTypeJPEG`. Qualquer outra string UTI também é aceita.

### Anexar e desanexar

Uma barra de ferramentas pertence a uma janela por vez. `WebviewWindow.SetToolbar` e `NativeWindow.SetToolbar` aceitam uma barra de ferramentas antes ou depois da criação da janela; passar `nil` a remove e a libera para uso em outro lugar. `SetToolbar` em uma `WebviewWindow` informa problemas de validação por meio de `Window.Error`, enquanto a versão de `NativeWindow` retorna o erro. A opção de janela `Mac.Toolbar` anexa uma barra de ferramentas na criação; ela é aplicada depois de `Mac.SplitView` para que um separador de acompanhamento encontre a barra lateral com a qual deve se alinhar.

## Visualização dividida

`MacSplitView` organiza os painéis. Adicione-os da borda inicial à final. `AddSidebar`, `AddContentList` e `AddInspector` recebem o modelo nativo que hospedam e retornam um `MacSplitPane` para controlar tamanho e recolhimento. `AddPrimaryContent` posiciona a WebView existente da janela e retorna um `MacSplitWebviewPane`, que acrescenta `SetContentLayout`.

```go
split := application.NewMacSplitView().SetAutosaveName("myapp.main-window")

sidebarPane := split.AddSidebar(sidebar).
    SetMinimumThickness(200).
    SetMaximumThickness(320).
    SetCollapsible(true)

split.AddContentList(list).
    SetMinimumThickness(240).
    SetCollapsible(true)

split.AddPrimaryContent().
    SetContentLayout(application.MacContentLayoutEdgeToEdge)

split.AddInspector(inspector).
    SetPreferredThicknessFraction(0.25).
    SetHoldingPriority(300).
    SetCollapsible(true).
    SetCanCollapseFromWindowResize(false).
    SetCollapsed(true)

sidebarPane.OnCollapsedChange(func(_ *application.Context, collapsed bool) {
    // fired for toggles, gestures, menu items and SetCollapsed alike
})

window.SetSplitView(split)

// At runtime:
sidebarPane.Toggle()
```

`SetSplitView` funciona antes ou depois da criação da janela nativa. Se chamado antes, o layout fica na fila e é instalado quando a janela é criada. Se chamado depois, por exemplo, de um callback de menu ou de bandeja em uma aplicação em execução, o layout é instalado imediatamente: a WebView existente da janela se torna o painel principal, a barra de ferramentas atual é anexada novamente para alinhar seus separadores de acompanhamento e todos os acessórios pendentes são anexados.

Uma janela criada após `app.Run` pode ser configurada em uma única chamada com as opções `Mac.SplitView` e `Mac.Toolbar`:

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "Notes",
    URL:   "/",
    Mac: application.MacWindow{
        SplitView: split,
        Toolbar:   toolbar,
    },
})
```

As regras de um layout são:

- Um layout precisa de pelo menos dois painéis e exatamente um painel principal (`AddPrimaryContent` para uma `WebviewWindow`, `AddTextEditor` para uma `NativeWindow`).
- No máximo uma lista de conteúdo, posicionada depois da barra lateral e antes do painel principal.
- A estrutura dos painéis fica fixa após a anexação da visualização dividida. As configurações dos painéis, o estado de recolhimento e o conteúdo da barra lateral, da lista e do inspetor ainda podem mudar a qualquer momento.
- Um layout instalado não pode ser substituído. Uma segunda chamada a `SetSplitView` na mesma janela informa `ErrMacSplitViewAlreadyInstalled` (por meio de `Window.Error` em uma `WebviewWindow`, como valor de retorno em uma `NativeWindow`) e não altera a janela. Passar `nil` antes da instalação limpa um layout pendente.
- Uma barra lateral, lista, inspetor ou visualização dividida pertence a uma janela por vez.

Os setters de painel são `SetMinimumThickness`, `SetMaximumThickness`, `SetPreferredThicknessFraction`, `SetHoldingPriority`, `SetCollapsible`, `SetCanCollapseFromWindowResize`, `SetCollapsed`, `Toggle`, `IsCollapsed` e `OnCollapsedChange`. `SetAutosaveName` persiste as posições dos divisores entre execuções.

`SetContentLayout` no painel principal escolhe entre `MacContentLayoutBelowToolbar` e `MacContentLayoutEdgeToEdge`. `MacContentLayoutAutomatic` herda `MacWindow.ContentLayout`, que, por sua vez, segue `TitleBar.FullSizeContent`. A disposição de borda a borda permite que o AppKit aplique seu efeito de borda de rolagem sob a barra de ferramentas no macOS 26 e posteriores.

## Barra lateral

`MacSidebar` é uma lista de origem nativa. Ela contém linhas raiz, seções e linhas aninhadas em qualquer profundidade. Guarde os identificadores retornados para atualizar as linhas depois.

```go
sidebar := application.NewMacSidebar()

// A root row above the sections.
sidebar.AddItem("All Notes").
    SetSymbol("tray.full").
    SetBadge(12).
    OnClick(func(*application.Context) {})

notes := sidebar.AddSection("Notes")
draft := notes.AddItem("Saturday, slowly").
    SetSymbol("doc.text").
    SetAccessorySymbol("pin.fill").
    SetTooltip("Field notes").
    SetEditable(true).
    OnRename(func(_ *application.Context, label string) {
        // the row label is already updated
    })

// Nested rows with a tinted symbol.
tags := sidebar.AddSection("Tags")
tint := application.NewRGB(0, 122, 255)
work := tags.AddItem("Work").SetSymbol("briefcase").SetTintColor(&tint)
work.AddItem("Meetings").SetSymbol("tag")
work.SetExpanded(true)

sidebar.SetSelectedItem(draft)
```

Os setters de linha são `SetLabel`, `SetSymbol`, `SetTooltip`, `SetEnabled`, `SetHidden`, `SetBadge`, `SetAccessorySymbol`, `SetTintColor`, `SetEditable` e `SetExpanded`. `OnClick` é acionado quando o AppKit seleciona a linha, `OnExpandedChange` quando o usuário abre ou fecha suas linhas aninhadas e `OnRename` depois que uma renomeação embutida é confirmada.

### Seleção

A seleção única é o padrão. `SetSelectedItem` seleciona uma linha sem acionar seu `OnClick`. Com a seleção múltipla habilitada, `OnClick` ainda é acionado para a linha clicada e `OnSelectionChange` informa o conjunto completo.

```go
sidebar.SetAllowsMultipleSelection(true)
sidebar.OnSelectionChange(func(_ *application.Context, items []*application.MacSidebarItem) {
    for _, item := range items {
        _ = item.Section()
    }
})
```

### Menus de contexto

Um clique com o botão direito procura um menu nesta ordem: o `SetContextMenu` da própria linha, depois o callback `OnContextMenu` da barra lateral e, por fim, o `SetContextMenu` de reserva da barra lateral. O callback roda na thread da aplicação enquanto o AppKit aguarda; portanto, mantenha-o rápido.

```go
sidebar.OnContextMenu(func(_ *application.Context, item *application.MacSidebarItem) *application.Menu {
    if item == nil {
        return nil // fall back to the menu set with SetContextMenu
    }
    menu := application.NewMenu()
    menu.Add("Delete").OnClick(func(*application.Context) { item.Remove() })
    return menu
})

empty := application.NewMenu()
empty.Add("New Note").OnClick(func(*application.Context) {})
sidebar.SetContextMenu(empty)
```

### Reordenação por arrastar

`SetReorderable` permite ao usuário arrastar linhas dentro de uma seção, entre seções e para a raiz ou a partir dela. O modelo Go é atualizado antes de `OnMove` ser acionado.

```go
sidebar.SetReorderable(true)
sidebar.OnMove(func(_ *application.Context, item *application.MacSidebarItem, section *application.MacSidebarSection, index int) {
    // section is nil when the row was dropped at the sidebar root
})
```

### Remoção

Remover uma linha também remove suas linhas aninhadas. Depois disso, os identificadores ficam inativos.

```go
notes.Remove(draft)
sidebar.RemoveSection(tags)
```

## Lista de conteúdo

`MacContentList` é a coluna central do Finder, do Mail e dos navegadores de documentos. Sem colunas, ela mostra linhas detalhadas: título, subtítulo, símbolo inicial, detalhe final e indicador de contagem.

```go
list := application.NewMacContentList().
    SetStyle(application.MacContentListStyleInset).
    SetEmptyText("No notes match")

row := list.AddRow("Saturday, slowly").
    SetSubtitle("A slow day is still a day well spent.").
    SetDetail("Yesterday").
    SetSymbol("doc.text").
    SetBadge(2)

list.OnSelectionChange(func(_ *application.Context, rows []*application.MacContentListRow) {})
list.OnActivate(func(_ *application.Context, row *application.MacContentListRow) {
    // double-click or Return
})
list.SetSelectedRow(row)
```

Os estilos são `MacContentListStyleAutomatic`, `MacContentListStyleInset`, `MacContentListStyleSourceList`, `MacContentListStylePlain` e `MacContentListStyleFullWidth`. `SetRowHeight`, `SetAlternatingRowBackgrounds`, `SetHeaderVisible` e `SetAllowsMultipleSelection` completam as opções de apresentação. As linhas permitem `SetHidden`, `SetEnabled`, `SetTooltip`, `InsertRow` em uma posição, `Remove` e `RemoveAll`.

### Colunas

`SetColumns` muda para o modo de tabela. As linhas então mostram seus valores de `SetCells`, um por coluna. Colunas marcadas como `Sortable` mostram um indicador de ordenação; `SetSortable` habilita cliques nos cabeçalhos. Sem um callback `OnSort`, a lista se ordena pelo texto da coluna com `SortBy`.

```go
table := application.NewMacContentList().
    SetColumns(
        application.MacContentListColumn{Title: "Name", Width: 220, Sortable: true},
        application.MacContentListColumn{Title: "Size", Width: 80, Alignment: application.MacContentListAlignTrailing},
        application.MacContentListColumn{Title: "Modified", Sortable: true},
    ).
    SetSortable(true).
    SetAlternatingRowBackgrounds(true)

table.AddRow("notes.txt").SetCells("notes.txt", "4 KB", "Today")
table.AddRow("ideas.txt").SetCells("ideas.txt", "1 KB", "Yesterday")

// Without OnSort the list sorts itself by the column text.
table.OnSort(func(_ *application.Context, column int, ascending bool) {
    table.SortRows(func(a, b *application.MacContentListRow) bool {
        return (a.Cells()[column] < b.Cells()[column]) == ascending
    })
})
```

### Menus de contexto

Os menus de contexto são resolvidos na mesma ordem da barra lateral: o `SetContextMenu` da linha, depois `OnContextMenu` e, por fim, o menu de reserva da lista.

```go
list.OnContextMenu(func(_ *application.Context, row *application.MacContentListRow) *application.Menu {
    if row == nil {
        return nil
    }
    menu := application.NewMenu()
    menu.Add("Delete").OnClick(func(*application.Context) { row.Remove() })
    return menu
})
```

## Inspetor

`MacInspector` é um painel de propriedades na borda final, construído com controles nativos agrupados em seções. Cada método `Add` retorna um identificador `MacInspectorControl` com setters e callbacks específicos do tipo.

```go
inspector := application.NewMacInspector()

document := inspector.AddSection("Document")
document.AddTextField("Title", "Saturday, slowly").
    OnTextChange(func(_ *application.Context, value string) {})
document.AddPopup("Category", []string{"Personal", "Work"}, 0).
    OnSelectionChange(func(_ *application.Context, index int, value string) {})
document.AddCheckbox("Pinned", false).
    OnToggle(func(_ *application.Context, checked bool) {})

appearance := inspector.AddSection("Appearance").SetCollapsible(true)
priority := appearance.AddSlider("Priority", 0, 5, 0)
priority.OnValueChange(func(_ *application.Context, value float64) {})
appearance.AddStepper("Indent", 0, 8, 1, 2)
appearance.AddSegmented("Align", []string{"Left", "Centre", "Right"}, 0)
appearance.AddColorWell("Tint", application.NewRGB(0, 122, 255)).
    OnColorChange(func(_ *application.Context, colour application.RGBA) {})
appearance.AddDatePicker("Due", time.Time{}).
    OnDateChange(func(_ *application.Context, t time.Time) {})
appearance.AddButton("Reset").OnClick(func(*application.Context) {
    priority.SetFloatValue(0)
})

statistics := inspector.AddSection("Statistics")
words := statistics.AddLabel("Words", "0")

// Programmatic setters never fire the callbacks.
words.SetValue("128")
```

Os tipos de controle e seus setters:

| Controle | Setters | Callback |
|---------|---------|----------|
| `AddLabel` | `SetValue` | nenhum |
| `AddTextField` | `SetValue` | `OnTextChange` |
| `AddCheckbox` | `SetChecked` | `OnToggle` |
| `AddPopup`, `AddSegmented` | `SetOptions`, `SetSelectedIndex` | `OnSelectionChange` |
| `AddSlider`, `AddStepper` | `SetFloatValue`, `SetRange`, `SetStep` | `OnValueChange` |
| `AddColorWell` | `SetColor` | `OnColorChange` |
| `AddDatePicker` | `SetDate` | `OnDateChange` |
| `AddButton` | `SetLabel` | `OnClick` |

`SetLabel`, `SetTooltip`, `SetEnabled` e `SetHidden` se aplicam a todos os tipos. As seções podem ser recolhíveis, e tanto as seções quanto os controles podem ser movidos ou removidos a qualquer momento.

```go
appearance.SetCollapsed(true)
inspector.MoveSection(statistics, 0)
appearance.Remove(priority)
inspector.RemoveSection(statistics)
```

## Acessórios

`MacAccessory` é uma faixa de controles nativos. Escolha um layout ao criá-lo: `MacAccessoryLayoutLeading` fica ao lado dos botões da janela, `MacAccessoryLayoutTrailing` na borda final da barra de título, `MacAccessoryLayoutBottom` ocupa toda a largura abaixo da barra de título e da barra de ferramentas, e `MacAccessoryLayoutTop` fica no topo de um painel da visualização dividida. Adicione todos os controles antes de anexar o acessório.

```go
// Next to the window buttons.
folders := application.NewMacAccessory(application.MacAccessoryLayoutLeading)
folders.AddSegmented([]string{"Inbox", "Starred"}, 0).
    SetSegmentSymbols("tray", "star").
    OnSelectionChange(func(_ *application.Context, index int, label string) {})

// At the trailing edge of the titlebar.
tools := application.NewMacAccessory(application.MacAccessoryLayoutTrailing)
tools.AddSearch("Search mail").
    SetIncremental(true).
    SetWidth(200).
    OnSearch(func(_ *application.Context, query string) {})
tools.AddSymbolButton("square.and.pencil").
    SetTooltip("Compose").
    OnClick(func(*application.Context) {})

// A full-width strip beneath the titlebar and toolbar.
status := application.NewMacAccessory(application.MacAccessoryLayoutBottom).
    SetHeight(30).
    SetPreferredScrollEdgeEffectStyle(application.MacScrollEdgeEffectStyleSoft)
label := status.AddLabel("Saved").SetSymbol("checkmark.circle")
status.AddFlexibleSpace()
status.AddButton("Mark All Read").OnClick(func(*application.Context) {})

for _, accessory := range []*application.MacAccessory{folders, tools, status} {
    if err := window.AddTitlebarAccessory(accessory); err != nil {
        window.Error("titlebar accessory: %s", err)
    }
}

// Live updates and lifecycle.
label.SetText("Edited").SetSymbol("pencil.circle")
status.SetHidden(true)
tools.Remove()
```

Os controles são `AddSearch`, `AddSegmented`, `AddButton`, `AddSymbolButton`, `AddMenuButton`, `AddLabel`, `AddFlexibleSpace` e, para integrações nativas, `AddNativeView`. `AddTitlebarAccessory` existe tanto em `WebviewWindow` quanto em `NativeWindow`; chamadas feitas antes da criação da janela ficam na fila e são aplicadas quando ela é criada. `Remove` desanexa um acessório para que ele possa ser anexado novamente em outro lugar, e `SetHidden` o recolhe no local.

### Acessórios de painel

No macOS 26 e posteriores, um acessório pode ficar na parte superior ou inferior de um painel dividido, onde o Finder mantém o campo de filtro da barra lateral. Crie o acessório com o layout `Top` ou `Bottom` e anexe-o ao painel com `AddTopAccessory` ou `AddBottomAccessory`.

```go
filter := application.NewMacAccessory(application.MacAccessoryLayoutTop)
filter.AddSearch("Filter").
    SetIncremental(true).
    OnSearch(func(_ *application.Context, query string) { filterNotes(query) })

if err := sidebarPane.AddTopAccessory(filter); err != nil {
    window.Error("sidebar filter: %s", err)
}
```

`SetPreferredScrollEdgeEffectStyle` escolhe o tratamento `Automatic`, `Soft` ou `Hard` do AppKit para conteúdo que rola atrás do acessório. Estilos explícitos exigem macOS 26.1; em versões anteriores, a solicitação é informada pelo manipulador de erros da janela e o estilo automático continua em vigor.

## Juntando tudo

Este programa cria uma janela de três painéis com barra lateral, WebView e inspetor, além de uma barra de ferramentas que acompanha os dois divisores.

```go title="main.go"
package main

import (
	"embed"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed assets
var assets embed.FS

func main() {
	app := application.New(application.Options{
		Name:   "Notes",
		Assets: application.AssetOptions{Handler: application.AssetFileServerFS(assets)},
	})

	window := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "Notes",
		Width:  1100,
		Height: 700,
		URL:    "/",
		Mac: application.MacWindow{
			ContentLayout: application.MacContentLayoutEdgeToEdge,
			TitleBar: application.MacTitleBar{
				FullSizeContent:      true,
				HideToolbarSeparator: true,
				ToolbarStyle:         application.MacToolbarStyleUnified,
			},
		},
	})

	sidebar := application.NewMacSidebar()
	notes := sidebar.AddSection("Notes")
	notes.AddItem("Saturday, slowly").SetSymbol("doc.text").OnClick(func(*application.Context) {
		app.Event.Emit("note:selected", "saturday")
	})

	inspector := application.NewMacInspector()
	inspector.AddSection("Document").AddTextField("Title", "Saturday, slowly").
		OnTextChange(func(_ *application.Context, value string) {
			app.Event.Emit("note:title", value)
		})

	split := application.NewMacSplitView().SetAutosaveName("notes.main")
	split.AddSidebar(sidebar).SetMinimumThickness(200).SetCollapsible(true)
	split.AddPrimaryContent()
	split.AddInspector(inspector).SetMinimumThickness(240).SetCollapsible(true)
	window.SetSplitView(split)

	toolbar := application.NewMacToolbar().SetDisplayMode(application.MacToolbarDisplayModeIconOnly)
	toolbar.AddSidebarToggle()
	toolbar.AddSidebarTrackingSeparator()
	toolbar.AddButton("New").SetSymbol("square.and.pencil").SetBordered(true).
		OnClick(func(*application.Context) {
			notes.AddItem("Untitled").SetSymbol("doc.text")
		})
	toolbar.AddFlexibleSpace()
	toolbar.AddInspectorTrackingSeparator()
	toolbar.AddInspectorToggle()
	window.SetToolbar(toolbar)

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
```

Os elementos nativos da janela e o frontend se comunicam pelos eventos e serviços habituais do Wails. A barra lateral emite `note:selected`, o inspetor emite `note:title` e a página escuta por meio de `Events.On` do runtime.

## Janelas nativas

@note{type="caution" title="Experimental"}
`NativeWindow`, `NativeWindowManager` e `MacTextEditor` são experimentais na v3. A API é deliberadamente pequena e pode mudar quando a API comum de janelas for reformulada para a v4.
@end

Uma `NativeWindow` não tem WebView. Seu conteúdo principal é um `MacTextEditor`, um `NSTextView` dentro de um `NSScrollView`, e ela aceita os mesmos tipos de barra de ferramentas, visualização dividida e acessórios que uma `WebviewWindow`. Crie uma com `app.NativeWindow.New` ou `app.NativeWindow.NewWithOptions` e encontre-a depois com `Get` ou `GetByID`.

Uma janela nativa só é criada quando tem conteúdo. Forneça uma visualização dividida cujo painel principal tenha sido adicionado com `AddTextEditor`, seja por `NativeWindowOptions.SplitView` ou chamando `SetSplitView`. Antes de `app.Run`, o layout fica na fila; em uma aplicação em execução, `SetSplitView` cria e mostra a janela imediatamente (a menos que `Hidden` esteja definido) e retorna qualquer erro de criação. Uma janela sem layout permanece com a criação adiada, e `Run` retorna `ErrNativeWindowContentRequired`; um layout sem editor de texto é rejeitado com `ErrNativeWindowEditorRequired`. Use `NativeWindowOptions.Toolbar` e `NativeWindowOptions.SplitView` em vez dos campos `Mac.Toolbar` e `Mac.SplitView`, que uma janela nativa ignora.

```go title="main.go"
package main

import (
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
)

func main() {
	app := application.New(application.Options{
		Name:       "Native Notes",
		NativeOnly: true,
	})

	editor := application.NewMacTextEditor()
	editor.OnChange(func(*application.Context) {
		// mark the document dirty; call editor.Text() only when saving
	})

	sidebar := application.NewMacSidebar()
	sidebar.AddSection("Files").AddItem("README.txt").
		SetSymbol("doc.plaintext").
		OnClick(func(*application.Context) {
			editor.SetText("Hello from AppKit")
		})

	split := application.NewMacSplitView().SetAutosaveName("native-notes.main")
	split.AddSidebar(sidebar).SetMinimumThickness(200).SetCollapsible(true)
	split.AddTextEditor(editor).SetMinimumThickness(400)

	window := app.NativeWindow.NewWithOptions(application.NativeWindowOptions{
		Title:  "Native Notes",
		Width:  900,
		Height: 600,
		Mac: application.MacWindow{
			TitleBar: application.MacTitleBar{
				FullSizeContent: true,
				ToolbarStyle:    application.MacToolbarStyleUnified,
			},
		},
	})
	if err := window.SetSplitView(split); err != nil {
		log.Fatal(err)
	}

	toolbar := application.NewMacToolbar()
	toolbar.AddSidebarToggle()
	toolbar.AddSidebarTrackingSeparator()
	toolbar.AddFlexibleSpace()
	toolbar.AddButton("Save").SetSymbol("square.and.arrow.down").SetBordered(true).
		OnClick(func(*application.Context) {
			log.Printf("%d bytes", len(editor.Text()))
		})
	if err := window.SetToolbar(toolbar); err != nil {
		log.Fatal(err)
	}

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
```

A mesma janela pode ser criada em uma única chamada a partir de uma aplicação em execução, passando os elementos da janela como opções:

```go
window := app.NativeWindow.NewWithOptions(application.NativeWindowOptions{
    Title:     "Native Notes",
    SplitView: split,
    Toolbar:   toolbar,
})
```

`MacTextEditor` oferece `SetText`, `Text`, `SetEditable`, `OnChange` e `Focus`. Chamadas programáticas a `SetText` não acionam `OnChange`; portanto, carregar um arquivo nunca o marca como modificado. `Text` lê o documento inteiro do AppKit; chame-o quando precisar do conteúdo, em vez de fazê-lo a cada alteração.

Duas opções mantêm enxuta uma aplicação exclusivamente nativa:

- `NativeOnly: true` em `application.Options` dispensa, em tempo de execução, o transporte do frontend e o servidor de recursos. Não crie uma `WebviewWindow` quando essa opção estiver definida.
- A tag de compilação `wails_native` exclui inteiramente do binário o código da WebView, do frontend e do atualizador e define `NativeOnly` automaticamente:

```sh
go build -tags wails_native .
```

O suporte a instância única também é excluído de uma compilação `wails_native`; adicione a tag `wails_single_instance` quando precisar dele.

## Abas de janela

O macOS pode agrupar janelas em abas. O modo de abas é fixado quando a janela é criada; portanto, defina `Mac.TabbingMode` como `MacWindowTabbingModePreferred` ou `MacWindowTabbingModeAutomatic` em cada janela que deva participar. `MacWindowTabbingModeDisallowed` (e o padrão quando não definido) mantém a janela fora dos grupos de abas.

```go
first := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "Document 1",
    URL:   "/",
    Mac: application.MacWindow{
        TabbingMode: application.MacWindowTabbingModePreferred,
    },
})
```

As operações de abas precisam de janelas nativas ativas; portanto, chame-as de um manipulador de menu, de um método de serviço ou de outro código executado após o início de `app.Run`. `TabGroup` retorna um identificador `MacWindowTabGroup` que resolve o grupo nativo a cada chamada; é seguro usar um identificador nulo, e todos os seus métodos retornam seu valor zero.

```go
second := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "Document 2",
    URL:   "/",
    Mac: application.MacWindow{
        TabbingMode: application.MacWindowTabbingModeAutomatic,
    },
})
if err := first.AddTab(second, application.MacTabOrderAbove); err != nil {
    app.Logger.Error("add tab", "error", err)
}
second.SetTabTitle("Draft")

if group := first.TabGroup(); group != nil {
    group.SelectNext()
    group.ToggleTabBar()
    for _, member := range group.Windows() {
        app.Logger.Info("tab", "name", member.Name())
    }
}

first.MoveTabToNewWindow()
first.MergeAllWindows()
```

`MacWindowTabGroup` fornece `Identifier`, `Count`, `Windows`, `NativeWindows`, `SelectedWindow`, `Select`, `SelectNative`, `SelectNext`, `SelectPrevious`, `IsTabBarVisible`, `ToggleTabBar`, `IsOverviewVisible` e `ToggleTabOverview`. `app.Window.TabGroups` lista todos os grupos que contêm uma `WebviewWindow`. `AddNativeTab` adiciona uma `NativeWindow` ao grupo de uma janela WebView, e `SetTabTooltip` define o texto exibido ao passar o cursor sobre uma aba. Fora do macOS, os métodos de abas retornam `ErrMacWindowTabsUnsupported`.

## Requisitos de versão

Tudo nesta página é exclusivo do macOS. Os recursos têm comportamento alternativo em versões anteriores, conforme indicado abaixo; a API Go é idêntica em todas as plataformas.

| Recurso | macOS mínimo | Comportamento em versões anteriores |
|---------|---------------|-------------------------------|
| Abas de janela | 10.12 | indisponíveis |
| Grupos da barra de ferramentas, itens com borda, `AddMenu` | 10.15 | os itens de menu são omitidos; os grupos usam a apresentação legada |
| SF Symbols (`SetSymbol` em itens da barra de ferramentas, barra lateral, lista e acessórios) | 11 | nenhuma imagem é mostrada |
| `AddSearch` como `NSSearchToolbarItem`, `SetNavigational`, separador de acompanhamento da barra lateral | 11 | a pesquisa usa um campo de pesquisa comum; o separador é omitido |
| Funções de divisão de inspetor e lista de conteúdo | 11 | os mesmos painéis são hospedados em itens de divisão comuns |
| Estilos de lista de conteúdo | 11 | ignorados |
| `SetCenteredItems` | 13 | ignorado |
| Botão de alternância e separador de acompanhamento do inspetor | 14 | o Wails fornece um botão de alternância nativo; o separador é omitido |
| `SetBadgeCount`, `SetProminent`, `SetTintColor` em itens da barra de ferramentas | 26 | armazenados e aplicados quando disponíveis |
| Acessórios de painel (`AddTopAccessory`, `AddBottomAccessory`) | 26 | `ErrMacSplitItemAccessoryUnavailable` |
| Estilos explícitos do efeito de borda de rolagem | 26.1 | informados pelo manipulador de erros da janela; o estilo automático permanece |

Acessórios da barra de título, visualizações divididas, barras laterais, inspetores e o editor de texto não têm requisito de versão além do mínimo do Wails.

## Exemplos

Cada exemplo é uma aplicação completa e executável:

- [`v3/examples/mac-toolbar`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-toolbar): um editor de notas com barra de ferramentas, barra lateral, lista de conteúdo, inspetor, provedor de compartilhamento, personalização da barra de ferramentas e acessório de painel.
- [`v3/examples/mac-titlebar-accessory`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-titlebar-accessory): acessórios nas posições inicial, final e inferior da barra de título que controlam uma visualização de caixa de correio.
- [`v3/examples/mac-window-tabs`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-window-tabs): modos de abas, `AddTab`, `AddNativeTab` e a API de grupos de abas a partir de um menu.
- [`v3/examples/mac-native-editor`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-native-editor): um editor de texto `NativeWindow` sem WebView, compilado com a tag `wails_native`.

## Saiba mais

Os elementos da janela desta página são metade de uma aplicação nativa para macOS. A outra metade é o comportamento da aplicação: janelas de documentos com ícones proxy e disposição em cascata, menus com símbolos e indicadores, menu do Dock com progresso, alertas e painéis nativos, itens de status removíveis, respostas táteis e fala, área de transferência com conteúdo avançado e arraste para fora, além de informações do sistema sobre permissões, energia e localidade. Essas APIs são abordadas no guia [Integração com a plataforma macOS](/guides/macos-platform-integration).

---
title: "Menüreferenz"
description: "Vollständige Referenz zu Menüeintragstypen, Eigenschaften und Methoden"
slug: "features/menus/reference"
sourcePath: "features/menus/reference.md"
---

## Menüreferenz

Vollständige Referenz zu Menüeintragstypen, Eigenschaften und dynamischem Verhalten. Erstellen Sie professionelle, reaktionsfähige Menüs mit Kontrollkästchen, Optionsgruppen, Trennlinien und dynamischen Aktualisierungen.

## Menüeintragstypen

### Reguläre Menüeinträge

Der häufigste Typ – zeigt Text an und löst eine Aktion aus:

```go
menuItem := menu.Add("Click Me")
menuItem.OnClick(func(ctx *application.Context) {
    fmt.Println("Menu item clicked!")
})
```

**Verwendung:** Befehle, Aktionen, Öffnen von Fenstern

### Kontrollkästchen

Umschaltbare Menüeinträge mit aktiviertem oder deaktiviertem Auswahlzustand:

```go
checkbox := menu.AddCheckbox("Enable Feature", true)  // true = initially checked
checkbox.OnClick(func(ctx *application.Context) {
    isChecked := ctx.ClickedMenuItem().Checked()
    fmt.Printf("Feature is now: %v\n", isChecked)
})
```

**Verwendung:** Boolesche Einstellungen, Funktionsumschalter, Ansichtsoptionen

**Wichtig:** Der Auswahlzustand wird beim Anklicken automatisch umgeschaltet.

### Optionsgruppen

Sich gegenseitig ausschließende Optionen – nur eine kann ausgewählt werden:

```go
menu.AddRadio("Small", true)   // true = initially selected
menu.AddRadio("Medium", false)
menu.AddRadio("Large", false)
```

**Verwendung:** Sich gegenseitig ausschließende Auswahlmöglichkeiten (Größe, Theme, Modus)

**Funktionsweise der Gruppierung:**

- Benachbarte Optionsfelder bilden automatisch eine Gruppe
- Bei Auswahl eines Eintrags werden die anderen Einträge der Gruppe abgewählt
- Trennen Sie Gruppen durch eine Trennlinie oder einen regulären Menüeintrag

**Beispiel mit mehreren Gruppen:**

```go
// Group 1: Size
menu.AddRadio("Small", true)
menu.AddRadio("Medium", false)
menu.AddRadio("Large", false)

menu.AddSeparator()

// Group 2: Theme
menu.AddRadio("Light", true)
menu.AddRadio("Dark", false)
```

### Untermenüs

Verschachtelte Menüstrukturen zur Organisation:

```go
submenu := menu.AddSubmenu("More Options")
submenu.Add("Submenu Item 1").OnClick(func(ctx *application.Context) {
    // Handle click
})
submenu.Add("Submenu Item 2")
```

**Verwendung:** Gruppieren zusammengehöriger Einträge, Reduzieren von Unübersichtlichkeit

**Verschachtelungslimit:** Die meisten Plattformen unterstützen 2-3 Ebenen. Vermeiden Sie eine tiefere Verschachtelung.

### Trennlinien

Visuelle Trennlinien zwischen Menüeinträgen:

```go
menu.Add("Item 1")
menu.AddSeparator()
menu.Add("Item 2")
```

**Verwendung:** Visuelles Gruppieren zusammengehöriger Einträge

**Bewährte Vorgehensweise:** Beginnen oder beenden Sie Menüs nicht mit Trennlinien.

## Eigenschaften von Menüeinträgen

### Beschriftung

Der für den Menüeintrag angezeigte Text:

```go
menuItem := menu.Add("Initial Label")
menuItem.SetLabel("New Label")

// Get current label
label := menuItem.Label()
```

**Dynamische Beschriftungen:**

```go
updateMenuItem := menu.Add("Check for Updates")
updateMenuItem.OnClick(func(ctx *application.Context) {
    updateMenuItem.SetLabel("Checking...")
    menu.Update()  // Important on Windows!
    
    // Perform update check
    checkForUpdates()
    
    updateMenuItem.SetLabel("Check for Updates")
    menu.Update()
})
```

### Aktivierungszustand

Legen Sie fest, ob mit dem Menüeintrag interagiert werden kann:

```go
menuItem := menu.Add("Save")
menuItem.SetEnabled(false)  // Greyed out, can't click

// Enable it later
menuItem.SetEnabled(true)
menu.Update()  // Important: Call this after changing enabled state!

// Check current state
isEnabled := menuItem.Enabled()
```

@note{type="caution" title="Menüverhalten unter Windows"}
Unter Windows müssen Menüs neu aufgebaut werden, wenn sich ihr Zustand ändert. **Rufen Sie nach dem Aktivieren oder Deaktivieren von Menüeinträgen immer `menu.Update()` auf**, insbesondere wenn der Eintrag im deaktivierten Zustand erstellt wurde.

**Grund:** Windows-Menüs werden bei Aktualisierungen vollständig neu aufgebaut. Wenn Sie `Update()` nicht aufrufen, werden Klick-Handler nicht ordnungsgemäß ausgelöst.

@end

**Beispiel: Dynamisches Aktivieren und Deaktivieren**

```go
var hasSelection bool

cutMenuItem := menu.Add("Cut")
cutMenuItem.SetEnabled(false)  // Initially disabled

copyMenuItem := menu.Add("Copy")
copyMenuItem.SetEnabled(false)

// When selection changes
func onSelectionChanged(selected bool) {
    hasSelection = selected
    cutMenuItem.SetEnabled(hasSelection)
    copyMenuItem.SetEnabled(hasSelection)
    menu.Update()  // Critical on Windows!
}
```

**Gängiges Muster: Bedingtes Aktivieren**

```go
saveMenuItem := menu.Add("Save")

func updateSaveMenuItem() {
    canSave := hasUnsavedChanges() && !isSaving()
    saveMenuItem.SetEnabled(canSave)
    menu.Update()
}

// Call whenever state changes
onDocumentChanged(func() {
    updateSaveMenuItem()
})
```

### Auswahlzustand

Steuern Sie bei Kontrollkästchen und Optionsfeldern den Auswahlzustand oder fragen Sie ihn ab:

```go
checkbox := menu.AddCheckbox("Feature", false)
checkbox.SetChecked(true)
menu.Update()

// Query state
isChecked := checkbox.Checked()
```

**Automatisches Umschalten:** Kontrollkästchen werden beim Anklicken automatisch umgeschaltet. Sie müssen `SetChecked()` nicht im Klick-Handler aufrufen.

**Manuelle Steuerung:**

```go
checkbox := menu.AddCheckbox("Auto-save", false)

// Sync with external state
func syncAutoSave(enabled bool) {
    checkbox.SetChecked(enabled)
    menu.Update()
}
```

### Tastenkombinationen

Fügen Sie Menüeinträgen Tastenkombinationen hinzu:

```go
saveMenuItem := menu.Add("Save")
saveMenuItem.SetAccelerator("CmdOrCtrl+S")

quitMenuItem := menu.Add("Quit")
quitMenuItem.SetAccelerator("CmdOrCtrl+Q")
```

**Format für Tastenkombinationen:**

- `CmdOrCtrl` – Cmd unter macOS, Ctrl unter Windows/Linux
- `Shift`, `Alt`, `Option` – Modifikatortasten
- `A-Z`, `0-9` – Buchstaben- und Zifferntasten
- `F1-F12` – Funktionstasten
- `Enter`, `Space`, `Backspace` usw. – Sondertasten

**Beispiele:**

```go
"CmdOrCtrl+S"           // Save
"CmdOrCtrl+Shift+S"     // Save As
"CmdOrCtrl+W"           // Close Window
"CmdOrCtrl+Q"           // Quit
"F5"                    // Refresh
"CmdOrCtrl+,"           // Preferences (macOS convention)
"Alt+F4"                // Close (Windows convention)
```

**Plattformspezifische Tastenkombinationen:**

```go
if runtime.GOOS == "darwin" {
    prefsMenuItem.SetAccelerator("Cmd+,")
} else {
    prefsMenuItem.SetAccelerator("Ctrl+P")
}
```

### Tooltip

Fügen Sie Menüeinträgen einen beim Darüberfahren angezeigten Text hinzu (die Plattformunterstützung variiert):

```go
menuItem := menu.Add("Advanced Options")
menuItem.SetTooltip("Configure advanced settings")
```

**Plattformunterstützung:**

- **Windows:** ✅ Unterstützt
- **macOS:** ❌ Nicht unterstützt (Tooltips sind für Menüs nicht üblich)
- **Linux:** ⚠️ Abhängig von der Desktop-Umgebung

### Sichtbarkeitszustand

Blenden Sie Menüeinträge aus, ohne sie zu entfernen:

```go
debugMenuItem := menu.Add("Debug Mode")
debugMenuItem.SetHidden(true)  // Hidden

// Show in debug builds
if isDebugBuild {
    debugMenuItem.SetHidden(false)
    menu.Update()
}
```

**Verwendung:** Debug-Optionen, Feature-Flags, bedingte Funktionen

## Ereignisbehandlung

### OnClick-Handler

Code ausführen, wenn auf den Menüeintrag geklickt wird:

```go
menuItem := menu.Add("Click Me")
menuItem.OnClick(func(ctx *application.Context) {
    // Handle click
    fmt.Println("Clicked!")
})
```

**Der Kontext stellt Folgendes bereit:**

- `ctx.ClickedMenuItem()` – Der angeklickte Menüeintrag
- Fensterkontext (bei Aufruf über das Fenstermenü)
- Anwendungskontext

**Beispiel: Im Handler auf den Menüeintrag zugreifen**

```go
checkbox := menu.AddCheckbox("Feature", false)
checkbox.OnClick(func(ctx *application.Context) {
    item := ctx.ClickedMenuItem()
    isChecked := item.Checked()
    fmt.Printf("Feature is now: %v\n", isChecked)
})
```

### Mehrere Handler

Sie können mehrere Handler festlegen (der zuletzt festgelegte wird verwendet):

```go
menuItem := menu.Add("Action")
menuItem.OnClick(func(ctx *application.Context) {
    fmt.Println("First handler")
})

// This replaces the first handler
menuItem.OnClick(func(ctx *application.Context) {
    fmt.Println("Second handler - this one runs")
})
```

**Empfehlung:** Legen Sie den Handler nur einmal fest und verwenden Sie darin bei Bedarf bedingte Logik.

## Dynamische Menüs

### Menüeinträge aktualisieren

**Die goldene Regel:** Rufen Sie nach jeder Änderung des Menüzustands `menu.Update()` auf.

```go
// ✅ Correct
menuItem.SetEnabled(true)
menu.Update()

// ❌ Wrong (especially on Windows)
menuItem.SetEnabled(true)
// Forgot to call Update() - click handlers may not work!
```

**Warum das wichtig ist:**

- **Windows:** Menüs werden beim Aktualisieren neu erstellt
- **macOS/Linux:** Weniger wichtig, aber dennoch empfohlen
- **Klick-Handler:** Werden ohne Update() nicht ordnungsgemäß ausgelöst

### Menüs neu erstellen

Erstellen Sie bei umfangreichen Änderungen das gesamte Menü neu:

```go
func rebuildFileMenu() {
    menu := app.Menu.New()
    
    menu.Add("New").OnClick(handleNew)
    menu.Add("Open").OnClick(handleOpen)
    
    if hasRecentFiles() {
        recentMenu := menu.AddSubmenu("Open Recent")
        for _, file := range getRecentFiles() {
            recentMenu.Add(file).OnClick(func(ctx *application.Context) {
                openFile(file)
            })
        }
    }
    
    menu.AddSeparator()
    menu.Add("Quit").OnClick(handleQuit)
    
    // Set the new menu
    window.SetMenu(menu)
}
```

**Wann das Menü neu erstellt werden sollte:**

- Die Liste der zuletzt verwendeten Dateien ändert sich
- Plugin-Menüs ändern sich
- Umfangreiche Zustandsübergänge

**Wann eine Aktualisierung ausreicht:**

- Einträge aktivieren oder deaktivieren
- Beschriftungen ändern
- Kontrollkästchen umschalten

### Kontextsensitive Menüs

Menüs an den Anwendungszustand anpassen:

```go
func updateEditMenu() {
    cutMenuItem.SetEnabled(hasSelection())
    copyMenuItem.SetEnabled(hasSelection())
    pasteMenuItem.SetEnabled(hasClipboardContent())
    undoMenuItem.SetEnabled(canUndo())
    redoMenuItem.SetEnabled(canRedo())
    menu.Update()
}

// Call whenever state changes
onSelectionChanged(updateEditMenu)
onClipboardChanged(updateEditMenu)
onUndoStackChanged(updateEditMenu)
```

## Plattformunterschiede

### Position der Menüleiste

| Plattform | Position | Hinweise |
| --- | --- | --- |
| **macOS** | Am oberen Bildschirmrand | Globale Menüleiste |
| **Windows** | Am oberen Fensterrand | Fensterspezifisches Menü |
| **Linux** | Am oberen Fensterrand | Fensterspezifisch (normalerweise) |

### Standardmenüs

**macOS:**

- Verfügt über ein Anwendungsmenü (mit dem Namen der App)
- „Einstellungen“ im Anwendungsmenü
- „Beenden“ im Anwendungsmenü

**Windows/Linux:**

- Kein Anwendungsmenü
- „Einstellungen“ im Menü „Bearbeiten“ oder „Extras“
- „Beenden“ im Menü „Datei“

**Beispiel: Plattformgerechte Struktur**

```go
menu := app.Menu.New()

// macOS gets Application menu
if runtime.GOOS == "darwin" {
    menu.AddRole(application.AppMenu)
}

// File menu
fileMenu := menu.AddSubmenu("File")
fileMenu.Add("New")
fileMenu.Add("Open")

// Preferences location varies
if runtime.GOOS == "darwin" {
    // On macOS, preferences are in Application menu (added by AppMenu role)
} else {
    // On Windows/Linux, add to Edit or Tools menu
    editMenu := menu.AddSubmenu("Edit")
    editMenu.Add("Preferences")
}
```

### Konventionen für Tastenkürzel

**macOS:**

- `Cmd+` für die meisten Tastenkürzel
- `Cmd+,` für „Einstellungen“
- `Cmd+Q` für „Beenden“

**Windows:**

- `Ctrl+` für die meisten Tastenkürzel
- `Ctrl+P` oder `Ctrl+,` für „Einstellungen“
- `Alt+F4` für „Beenden“ (oder `Ctrl+Q`)

**Linux:**

- Folgt im Allgemeinen den Windows-Konventionen
- Die Desktop-Umgebung kann diese überschreiben

## Bewährte Vorgehensweisen

### ✅ Empfohlen

- **menu.Update() aufrufen**, nachdem Sie den Menüzustand geändert haben (insbesondere unter Windows)
- **Optionsgruppen verwenden** für Optionen, die sich gegenseitig ausschließen
- **Kontrollkästchen verwenden** für Funktionen, die ein- und ausgeschaltet werden können
- **Fügen Sie häufig verwendeten Aktionen Tastenkürzel hinzu**
- **Gruppieren Sie zusammengehörige Einträge** mit Trennlinien
- **Testen Sie auf allen Plattformen** – das Verhalten variiert

### ❌ Vermeiden

- **Vergessen Sie menu.Update() nicht** – Click-Handler funktionieren sonst nicht ordnungsgemäß
- **Verschachteln Sie Menüs nicht zu tief** – maximal 2-3 Ebenen
- **Beginnen oder beenden Sie Menüs nicht mit Trennlinien** – das wirkt unprofessionell
- **Verwenden Sie unter macOS keine Tooltips** – sie werden nicht unterstützt
- **Codieren Sie plattformspezifische Tastenkürzel nicht fest** – verwenden Sie `CmdOrCtrl`

## Fehlerbehebung

### Menüeinträge reagieren nicht

**Symptom:** Click-Handler werden nicht ausgelöst

**Ursache:** Nach dem Aktivieren des Eintrags wurde `menu.Update()` nicht aufgerufen

**Lösung:**

```go
menuItem.SetEnabled(true)
menu.Update()  // Add this!
```

### Menüeinträge sind ausgegraut

**Symptom:** Menüeinträge lassen sich nicht anklicken

**Ursache:** Die Einträge sind deaktiviert

**Lösung:**

```go
menuItem.SetEnabled(true)
menu.Update()
```

### Tastenkürzel funktionieren nicht

**Symptom:** Tastenkürzel lösen keine Menüeinträge aus

**Ursachen:**

1. Das Format des Tastenkürzels ist falsch
2. Konflikt mit Systemtastenkürzeln
3. Das Fenster hat nicht den Fokus

**Lösung:**

```go
// Check format
menuItem.SetAccelerator("CmdOrCtrl+S")  // ✅ Correct
menuItem.SetAccelerator("Ctrl+S")       // ❌ Wrong (macOS uses Cmd)

// Avoid conflicts
// ❌ Cmd+H (Hide Window on macOS - system shortcut)
// ✅ Cmd+Shift+H (Custom shortcut)
```

## Nächste Schritte

- [Anwendungsmenüs](/features/menus/application/) – Menüleisten für Anwendungen erstellen
- [Kontextmenüs](/features/menus/context/) – Kontextmenüs für Rechtsklicks
- [Infobereichsmenüs](/features/menus/systray/) – Menüs für den Infobereich beziehungsweise die Menüleiste
- [Menümuster](/guides/menus/) – gängige Menümuster und bewährte Vorgehensweisen

---

**Fragen?** Fragen Sie auf [Discord](https://discord.gg/JDdSxwjhGf) nach oder sehen Sie sich die [Menübeispiele](https://github.com/wailsapp/wails/tree/master/v3/examples/menu) an.
